package responses

type (
	AuthPaymentResponse struct {
		PaymentID string `json:"payment_id"`
		Amount    int    `json:"amount"`
		Currency  string `json:"currency"`
	}

	CapturePaymentResponse struct {
		PaymentID string `json:"payment_id"`
		Amount    int    `json:"amount"`
		Currency  string `json:"currency"`
	}
)
