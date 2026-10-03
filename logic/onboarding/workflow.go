package onboarding_workflow

import (
	"context"
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

	sgDefaultEvent       = ""
	sgActivatedUserEvent = "activated_user_signal"

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

	// Empty signalEvent = run immediately. Non-empty signalEvent parks until
	// that Signal arrives, then runs the activity (go-cook @c6981a4).
	o.workflow.AddTransitionActivityWithOptions(persistOnboardingRequestActivity, sgDefaultEvent, o.persistRequest, activityOption, registerUserActivity)
	o.workflow.AddTransitionActivityWithOptions(registerUserActivity, sgDefaultEvent, o.registerUser, activityOption, updateOnboardingStatusActivity)
	o.workflow.AddTransitionActivityWithOptions(updateOnboardingStatusActivity, sgActivatedUserEvent, o.updateOnboardingStatus, activityOption)

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

func (o *onboardingWorkflow) ActivateUserRegisterWorkflow(ctx context.Context, request *requests.ActivateRequest) (*responses.ActivateResponse, *apperror.CustomError) {
	obj, err := decodeActivationToken(request.Token)
	if err != nil {
		return nil, apperror.NewCustomError(apperror.ErrBadRequest, "invalid token", err)
	}

	if err := o.workflow.SignalWorkflow(ctx, obj.WorkflowID, obj.RunID, sgActivatedUserEvent, true); err != nil {
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
