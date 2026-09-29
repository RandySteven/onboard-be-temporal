package onboarding_workflow

import (
	"context"

	temporal_client "github.com/RandySteven/go-cook/temporal"
	"github.com/RandySteven/onboard-be/enums"
)

const onboardingInputErrorType = "OnboardingInputError"

func (o *onboardingWorkflow) persistRequest(ctx context.Context, execData *ExecutionData) (*ExecutionData, error) {
	if execData.Request == nil {
		return nil, temporal_client.NewNonRetryableError(onboardingInputErrorType, "register request is missing")
	}

	onboarding, err := onboardingRegisterRequestMapping(execData.Request)
	if err != nil {
		return nil, temporal_client.NewNonRetryableError(onboardingInputErrorType, "failed to hash password")
	}
	execData.Onboarding = onboarding

	onboardingResult, err := o.onboardingRepository.Save(ctx, execData.Onboarding)
	if err != nil {
		return nil, err
	}
	execData.Onboarding = onboardingResult

	if execData.Request.RegisterAs != enums.RegisterAsUser && execData.Request.RegisterAs != "" {
		return nil, temporal_client.NewNonRetryableError(onboardingInputErrorType, "unsupported register_as")
	}

	execData.SetActivity(registerUserActivity)
	return execData, nil
}
