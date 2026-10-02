package requests

import "encoding/json"

type (
	AuthPaymentRequest struct {
		UserID         int             `json:"user_id"`
		Amount         int             `json:"amount"`
		Currency       string          `json:"currency"`
		AdditionalInfo json.RawMessage `json:"additional_info"`
	}

	CapturePaymentRequest struct {
		UserID         int             `json:"user_id"`
		AuthPaymentID  int             `json:"auth_payment_id"`
		Amount         int             `json:"amount"`
		Currency       string          `json:"currency"`
		AdditionalInfo json.RawMessage `json:"additional_info"`
	}
)
