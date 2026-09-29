package onboarding_workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	temporal_client "github.com/RandySteven/go-cook/temporal"
	"github.com/RandySteven/onboard-be/apperror"
	"github.com/RandySteven/onboard-be/entities/payloads/requests"
	"github.com/RandySteven/onboard-be/entities/payloads/responses"
	"github.com/RandySteven/onboard-be/enums"
	"github.com/RandySteven/onboard-be/repositories"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	onboardingWorkflowExecution = "onboarding_workflow_execution"

	persistOnboardingRequestActivity = "persist_onboarding_request_activity"
	registerUserActivity             = "register_user_activity"
	updateOnboardingStatusActivity   = "update_onboarding_status_activity"

	sgOnboardingCorrection = temporal_client.DefaultCorrectionSignal
	sgActivatedUser        = "activated_user_signal"

	queryRegisterResponse = "RegisterResponse"

	registerClientWaitTimeout  = 30 * time.Second
	registerClientPollInterval = 150 * time.Millisecond
)

var nonRetryableErrorTypes = []string{onboardingInputErrorType}

type (
	onboardingWorkflow struct {
		workflow temporal_client.WorkflowExecution
		temporal temporal_client.Temporal

		onboardingRepository repositories.OnboardingRepository
		userRepository       repositories.UserRepository
	}

	OnboardingWorkflow interface {
		OnboardingRegisterWorkflow(ctx context.Context, request *requests.RegisterRequest) (*responses.RegisterResponse, *apperror.CustomError)
		ActivateUserRegisterWorkflow(ctx context.Context, request *requests.ActivateRequest) (*responses.ActivateResponse, *apperror.CustomError)
	}
)

func (o *onboardingWorkflow) registerWorkflowAndActivities() {
	activityOption := &workflow.ActivityOptions{
		StartToCloseTimeout: 25 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:        5,
			InitialInterval:        2 * time.Second,
			BackoffCoefficient:     2,
			MaximumInterval:        30 * time.Second,
			NonRetryableErrorTypes: nonRetryableErrorTypes,
		},
	}

	correctionOnly := temporal_client.ResumableOptions{
		CorrectionSignal: sgOnboardingCorrection,
	}
	registerThenActivate := temporal_client.ResumableOptions{
		CorrectionSignal: sgOnboardingCorrection,
		ApprovalSignal:   sgActivatedUser,
	}

	o.workflow.AddResumableTransitionActivityWithOptions(
		persistOnboardingRequestActivity,
		o.persistRequest,
		activityOption,
		correctionOnly,
		registerUserActivity,
	)

	o.workflow.AddResumableTransitionActivityWithOptions(
		registerUserActivity,
		o.registerUser,
		activityOption,
		registerThenActivate,
		updateOnboardingStatusActivity,
	)

	o.workflow.AddResumableTransitionActivityWithOptions(
		updateOnboardingStatusActivity,
		o.updateOnboardingStatus,
		activityOption,
		correctionOnly,
	)

	o.workflow.RegisterWorkflow(onboardingWorkflowExecution, o.register)
}

func (o *onboardingWorkflow) OnboardingRegisterWorkflow(ctx context.Context, request *requests.RegisterRequest) (*responses.RegisterResponse, *apperror.CustomError) {
	workflowID := fmt.Sprintf("%s-%s", onboardingWorkflowExecution, ctx.Value(enums.RequestID))
	startWorkflowOption := &temporal_client.StartWorkflowOptions{
		WorkflowID: workflowID,
	}
	weRun, err := o.workflow.StartWorkflow(ctx, *startWorkflowOption, o.register, workflowID, request)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, "there is wrong with workflow", err)
	}

	response, err := o.waitForRegisterResponse(ctx, weRun.GetID(), weRun.GetRunID())
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, "there is wrong with get result", err)
	}
	return response, nil
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
		status, statusErr := o.queryWorkflowStatus(waitCtx, workflowID, runID)
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

func (o *onboardingWorkflow) queryWorkflowStatus(ctx context.Context, workflowID, runID string) (string, error) {
	raw, err := o.temporal.QueryWorkflow(ctx, workflowID, runID, temporal_client.QueryGetStatus)
	if err != nil {
		return "", err
	}
	status, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("unexpected status query type %T", raw)
	}
	return status, nil
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

func (o *onboardingWorkflow) ActivateUserRegisterWorkflow(ctx context.Context, request *requests.ActivateRequest) (*responses.ActivateResponse, *apperror.CustomError) {
	obj, err := decodeActivationToken(request.Token)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrBadRequest, "invalid token", err)
	}

	if err := o.workflow.SignalWorkflow(ctx, obj.WorkflowID, obj.RunID, sgActivatedUser, true); err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, "activation failed", err)
	}

	var registerResult *responses.RegisterResponse
	if err := o.workflow.GetWorkflowResult(ctx, obj.WorkflowID, obj.RunID, &registerResult); err != nil {
		return nil, apperror.NewCustomError(apperror.ErrInternalServer, "activation did not complete", err)
	}

	return &responses.ActivateResponse{Status: "success", Activate: true}, nil
}

func NewOnboardingWorkflow(
	workflow temporal_client.WorkflowExecution,
	temporal temporal_client.Temporal,
	onboardingRepository repositories.OnboardingRepository,
	userRepository repositories.UserRepository,
) *onboardingWorkflow {
	ow := &onboardingWorkflow{
		workflow:             workflow,
		temporal:             temporal,
		onboardingRepository: onboardingRepository,
		userRepository:       userRepository,
	}
	ow.registerWorkflowAndActivities()
	return ow
}
