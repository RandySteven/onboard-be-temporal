package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/RandySteven/onboard-be/usecases"
)

type Handlers struct {
	OnboardingHandler IOnboardingHandler
}

func NewHandlers(usecases *usecases.Usecases) *Handlers {
	return &Handlers{
		OnboardingHandler: NewOnboardingHandler(usecases.OnboardingUsecase),
	}
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
