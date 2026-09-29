package usecases

import (
	"context"
	"fmt"
	"strings"

	temporal_client "github.com/RandySteven/go-cook/temporal"
	"github.com/RandySteven/onboard-be/apperror"
	"github.com/RandySteven/onboard-be/entities/payloads/requests"
	"github.com/RandySteven/onboard-be/entities/payloads/responses"
	"github.com/RandySteven/onboard-be/enums"
	onboarding_workflow "github.com/RandySteven/onboard-be/logic/onboarding"
	"github.com/RandySteven/onboard-be/repositories"
)

type (
	OnboardingUsecase interface {
		Register(ctx context.Context, request *requests.RegisterRequest) (*responses.RegisterResponse, *apperror.CustomError)
		ActivateUser(ctx context.Context, request *requests.ActivateRequest) (*responses.ActivateResponse, *apperror.CustomError)
	}

	onboardingUsecase struct {
		onboardingWorkflow onboarding_workflow.OnboardingWorkflow
	}
)

func (o *onboardingUsecase) Register(ctx context.Context, request *requests.RegisterRequest) (*responses.RegisterResponse, *apperror.CustomError) {
	if err := validateRegisterRequest(request); err != nil {
		return nil, apperror.NewCustomError(apperror.ErrBadRequest, err.Error(), err)
	}
	return o.onboardingWorkflow.OnboardingRegisterWorkflow(ctx, request)
}

func (o *onboardingUsecase) ActivateUser(ctx context.Context, request *requests.ActivateRequest) (*responses.ActivateResponse, *apperror.CustomError) {
	if request == nil || strings.TrimSpace(request.Token) == "" {
		err := fmt.Errorf("token is required")
		return nil, apperror.NewCustomError(apperror.ErrBadRequest, err.Error(), err)
	}
	return o.onboardingWorkflow.ActivateUserRegisterWorkflow(ctx, request)
}

func validateRegisterRequest(request *requests.RegisterRequest) error {
	if request == nil {
		return fmt.Errorf("register request is required")
	}
	if strings.TrimSpace(request.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if strings.TrimSpace(request.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if strings.TrimSpace(request.Password) == "" {
		return fmt.Errorf("password is required")
	}
	if strings.TrimSpace(request.PhoneNumber) == "" {
		return fmt.Errorf("phone_number is required")
	}
	if request.RegisterAs == "" {
		request.RegisterAs = enums.RegisterAsUser
	}
	if request.RegisterAs != enums.RegisterAsUser {
		return fmt.Errorf("unsupported register_as %s", request.RegisterAs)
	}
	return nil
}

func newOnboardingUsecase(
	onboardingRepository repositories.OnboardingRepository,
	userRepository repositories.UserRepository,
	temporal temporal_client.Temporal,
) *onboardingUsecase {
	onboardingWorkflow := onboarding_workflow.NewOnboardingWorkflow(
		temporal_client.NewWorkflowExecution(temporal),
		temporal,
		onboardingRepository,
		userRepository,
	)
	return &onboardingUsecase{
		onboardingWorkflow: onboardingWorkflow,
	}
}

var _ OnboardingUsecase = &onboardingUsecase{}
