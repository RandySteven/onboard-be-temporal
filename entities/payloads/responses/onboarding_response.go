package responses

type (
	RegisterResponse struct {
		UserID          uint64 `json:"user_id"`
		OnboardingID    uint64 `json:"onboarding_id"`
		UserName        string `json:"user_name"`
		UserEmail       string `json:"user_email"`
		ActivationToken string `json:"activation_token"`
	}

	ActivateResponse struct {
		Status   string `json:"status"`
		Activate bool   `json:"activate"`
	}
)
