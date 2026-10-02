package payment_workflow

import (
	"context"

	"github.com/RandySteven/onboard-be/apperror"
	"github.com/RandySteven/onboard-be/entities/payloads/requests"
	"github.com/RandySteven/onboard-be/entities/payloads/responses"
)

func (w *paymentWorkflow) auth(ctx context.Context, request *requests.AuthPaymentRequest) (*responses.AuthPaymentResponse, *apperror.CustomError) {
	return nil, nil
}
