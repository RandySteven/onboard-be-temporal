package payment_workflow

import (
	"context"

	"github.com/RandySteven/onboard-be/apperror"
	"github.com/RandySteven/onboard-be/entities/payloads/requests"
	"github.com/RandySteven/onboard-be/entities/payloads/responses"
)

func (w *paymentWorkflow) capture(ctx context.Context, request *requests.CapturePaymentRequest) (*responses.CapturePaymentResponse, *apperror.CustomError) {
	
	return nil, nil
}
