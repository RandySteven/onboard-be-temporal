package onboarding_workflow

import (
	"encoding/json"

	temporal_client "github.com/RandySteven/go-cook/temporal"
	"github.com/RandySteven/onboard-be/entities/models"
	"github.com/RandySteven/onboard-be/entities/payloads/requests"
	"github.com/RandySteven/onboard-be/entities/payloads/responses"
)

type ExecutionData struct {
	Request    *requests.RegisterRequest   `json:"request"`
	Response   *responses.RegisterResponse `json:"response"`
	WorkflowID string                      `json:"workflow_id"`
	RunID      string                      `json:"run_id"`
	Status     string                      `json:"status"`

	Onboarding      *models.Onboarding `json:"onboarding"`
	User            *models.User       `json:"user"`
	ActivationToken string             `json:"activation_token"`
	CurrActivity    string             `json:"curr_activity"`
	NextActivity    string             `json:"next_activity"`
}

func (e *ExecutionData) GetActivity() string {
	return e.CurrActivity
}

func (e *ExecutionData) SetActivity(activityName string) {
	e.CurrActivity = activityName
}

func (e *ExecutionData) GetStatus() string {
	return e.Status
}

func (e *ExecutionData) SetStatus(status string) {
	e.Status = status
}

func (e *ExecutionData) ApplyCorrection(payload json.RawMessage) error {
	var asRequest requests.RegisterRequest
	if err := json.Unmarshal(payload, &asRequest); err == nil && isRegisterRequestPatch(&asRequest) {
		e.Request = &asRequest
		return nil
	}

	var next ExecutionData
	if err := json.Unmarshal(payload, &next); err != nil {
		return err
	}
	if next.Request != nil {
		e.Request = next.Request
	}
	if next.ActivationToken != "" {
		e.ActivationToken = next.ActivationToken
	}
	return nil
}

func isRegisterRequestPatch(req *requests.RegisterRequest) bool {
	return req.Email != "" || req.Username != "" || req.FirstName != "" || req.PhoneNumber != ""
}

func (e *ExecutionData) JSONString() (string, error) {
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

var (
	_ temporal_client.ExecutionWorkflow = (*ExecutionData)(nil)
	_ temporal_client.StatusReporter    = (*ExecutionData)(nil)
	_ temporal_client.CorrectionApplier = (*ExecutionData)(nil)
)
