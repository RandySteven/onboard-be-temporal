package onboarding_workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	temporal_client "github.com/RandySteven/go-cook/temporal"
	"github.com/RandySteven/onboard-be/entities/models"
	"github.com/RandySteven/onboard-be/entities/payloads/requests"
	"github.com/RandySteven/onboard-be/entities/payloads/responses"
	"github.com/RandySteven/onboard-be/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func onboardingRegisterRequestMapping(request *requests.RegisterRequest) (*models.Onboarding, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &models.Onboarding{
		Name:           utils.MergeName(request.FirstName, request.LastName),
		Email:          request.Email,
		Password:       string(hashed),
		PhoneNumber:    request.PhoneNumber,
		Address:        request.Address,
		RegisterAs:     request.RegisterAs.ToString(),
		Status:         "PENDING",
		AdditionalInfo: utils.AdditionalInfoOrEmpty(request.AdditionalInfo),
	}, nil
}

func mappingUserInfo(additionalInfo json.RawMessage) *models.User {
	userRequest := &requests.RegisterUserRequestData{}
	if err := utils.MapJSONToInterface(additionalInfo, userRequest); err != nil {
		return nil
	}
	return &models.User{
		AdditionalInfo: utils.AdditionalInfoOrEmpty(additionalInfo),
	}
}

func activationSecret() []byte {
	if key := os.Getenv("ACTIVATION_JWT_KEY"); key != "" {
		return []byte(key)
	}
	return []byte("activation_hash")
}

type DecodeActivationTokenObject struct {
	jwt.RegisteredClaims
	UserID       uint64 `json:"user_id"`
	OnboardingID uint64 `json:"onboarding_id"`
	WorkflowID   string `json:"workflow_id"`
	RunID        string `json:"run_id"`
}

func activationTokenGenerate(userID, onboardingID uint64, workflowID, runID string) string {
	claims := &DecodeActivationTokenObject{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "ActivationRequest",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
		UserID:       userID,
		OnboardingID: onboardingID,
		WorkflowID:   workflowID,
		RunID:        runID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(activationSecret())
	if err != nil {
		return ""
	}
	return signedToken
}

func decodeActivationToken(tokenString string) (*DecodeActivationTokenObject, error) {
	claims := &DecodeActivationTokenObject{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return activationSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func (o *onboardingWorkflow) waitForRegisterResponse(ctx context.Context, workflowID, runID string) (*responses.RegisterResponse, error) {
	waitCtx := ctx
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		waitCtx, cancel = context.WithTimeout(ctx, registerClientWaitTimeout)
		defer cancel()
	}

	ticker := time.NewTicker(registerClientPollInterval)
	defer ticker.Stop()

	for {
		status, statusErr := o.workflow.GetWorkflowStatus(waitCtx, workflowID, runID)
		if statusErr == nil {
			switch status {
			case temporal_client.StatusFailed, temporal_client.StatusRejected:
				return nil, fmt.Errorf("onboarding workflow %s", status)
			}
		}

		raw, err := o.temporal.QueryWorkflow(waitCtx, workflowID, runID, queryRegisterResponse)
		if err == nil {
			if resp := decodeRegisterQueryResult(raw); registerResponseReady(resp) {
				return resp, nil
			}
		}

		select {
		case <-waitCtx.Done():
			return nil, waitCtx.Err()
		case <-ticker.C:
		}
	}
}

func decodeRegisterQueryResult(raw interface{}) *responses.RegisterResponse {
	if raw == nil {
		return nil
	}
	if resp, ok := raw.(*responses.RegisterResponse); ok {
		return resp
	}
	if resp, ok := raw.(responses.RegisterResponse); ok {
		return &resp
	}

	bytes, err := json.Marshal(raw)
	if err != nil {
		return nil
	}
	var resp responses.RegisterResponse
	if err := json.Unmarshal(bytes, &resp); err != nil {
		return nil
	}
	return &resp
}

func registerResponseReady(resp *responses.RegisterResponse) bool {
	return resp != nil && (resp.ActivationToken != "" || resp.OnboardingID != 0)
}