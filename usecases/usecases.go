package usecases

import (
	temporal_client "github.com/RandySteven/go-cook/temporal"
	"github.com/RandySteven/onboard-be/repositories"
)

type Usecases struct {
	OnboardingUsecase OnboardingUsecase
}

func NewUsecases(repositories *repositories.Repositories, temporal temporal_client.Temporal) *Usecases {
	return &Usecases{
		OnboardingUsecase: newOnboardingUsecase(
			repositories.OnboardingRepository,
			repositories.UserRepository,
			temporal,
		),
	}
}
