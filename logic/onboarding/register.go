package onboarding_workflow

import (
	"context"
	"fmt"

	temporal_client "github.com/RandySteven/go-cook/temporal"
	"github.com/RandySteven/onboard-be/entities/payloads/requests"
	"github.com/RandySteven/onboard-be/entities/payloads/responses"
	"go.temporal.io/sdk/workflow"
)

func (w *onboardingWorkflow) register(ctx workflow.Context, workflowID string, request *requests.RegisterRequest) (*responses.RegisterResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("register request is nil")
	}

	info := workflow.GetInfo(ctx)
	executionData := &ExecutionData{
		Request:    request,
		WorkflowID: workflowID,
		RunID:      info.WorkflowExecution.RunID,
	}
	if executionData.WorkflowID == "" {
		executionData.WorkflowID = info.WorkflowExecution.ID
	}

	if err := workflow.SetQueryHandler(ctx, queryRegisterResponse, func() (*responses.RegisterResponse, error) {
		return executionData.Response, nil
	}); err != nil {
		return nil, err
	}

	if err := w.workflow.Execute(ctx, executionData); err != nil {
		return nil, err
	}

	return executionData.Response, nil
}

func (w *onboardingWorkflow) registerUser(ctx context.Context, executionData *ExecutionData) (*ExecutionData, error) {
	user := mappingUserInfo(executionData.Request.AdditionalInfo)
	if user == nil {
		return nil, temporal_client.NewNonRetryableError(onboardingInputErrorType, "invalid user additional_info")
	}
	if executionData.Onboarding == nil {
		return nil, temporal_client.NewNonRetryableError(onboardingInputErrorType, "onboarding record is missing")
	}

	user.OnboardingID = executionData.Onboarding.ID
	user.UserName = executionData.Request.Username
	user.Status = "PROCESSING"

	user, err := w.userRepository.Save(ctx, user)
	if err != nil {
		return nil, err
	}

	onboardingData := executionData.Onboarding
	onboardingData.Status = "ACTIVE"
	onboardingData, err = w.onboardingRepository.Update(ctx, onboardingData)
	if err != nil {
		return nil, err
	}
	executionData.Onboarding = onboardingData
	executionData.User = user

	executionData.ActivationToken = activationTokenGenerate(user.ID, executionData.Onboarding.ID, executionData.WorkflowID, executionData.RunID)
	executionData.setPendingActivationResponse()
	executionData.SetActivity(updateOnboardingStatusActivity)
	return executionData, nil
}

func (e *ExecutionData) setPendingActivationResponse() {
	resp := &responses.RegisterResponse{
		ActivationToken: e.ActivationToken,
	}
	if e.Onboarding != nil {
		resp.OnboardingID = e.Onboarding.ID
		resp.UserName = e.Onboarding.Name
		resp.UserEmail = e.Onboarding.Email
	}
	if e.User != nil {
		resp.UserID = e.User.ID
		if e.User.UserName != "" {
			resp.UserName = e.User.UserName
		}
	}
	e.Response = resp
}
