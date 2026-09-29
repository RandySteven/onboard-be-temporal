package handlers

import (
	"context"
	"net/http"

	"github.com/RandySteven/onboard-be/entities/payloads/requests"
	"github.com/RandySteven/onboard-be/enums"
	"github.com/RandySteven/onboard-be/usecases"
	"github.com/RandySteven/onboard-be/utils"
	"github.com/google/uuid"
)

type (
	IOnboardingHandler interface {
		RegisterUser(w http.ResponseWriter, r *http.Request)
		ActivateUser(w http.ResponseWriter, r *http.Request)
	}

	OnboardingHandler struct {
		onboardingUsecase usecases.OnboardingUsecase
	}
)

func (o *OnboardingHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), enums.RequestID, uuid.NewString())

	dataKey := "activate"
	request := &requests.ActivateRequest{}
	if err := utils.BindJSON(r, request); err != nil {
		utils.ResponseHandler(w, http.StatusBadRequest, "invalid request", nil, nil, err)
		return
	}

	response, customErr := o.onboardingUsecase.ActivateUser(ctx, request)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.Error(), nil, nil, customErr)
		return
	}
	utils.ResponseHandler(w, http.StatusOK, "success activate user", &dataKey, response, nil)
}

func (o *OnboardingHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), enums.RequestID, uuid.NewString())

	dataKey := "register"
	request := &requests.RegisterRequest{}
	if err := utils.BindJSON(r, request); err != nil {
		utils.ResponseHandler(w, http.StatusBadRequest, "invalid request", nil, nil, err)
		return
	}

	response, customErr := o.onboardingUsecase.Register(ctx, request)
	if customErr != nil {
		utils.ResponseHandler(w, customErr.ErrCode(), customErr.Error(), nil, nil, customErr)
		return
	}
	utils.ResponseHandler(w, http.StatusCreated, "success register user", &dataKey, response, nil)
}

var _ IOnboardingHandler = &OnboardingHandler{}

func NewOnboardingHandler(onboardingUsecase usecases.OnboardingUsecase) *OnboardingHandler {
	return &OnboardingHandler{onboardingUsecase: onboardingUsecase}
}
