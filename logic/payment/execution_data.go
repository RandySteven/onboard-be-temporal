package payment_workflow

import (
	"bytes"
	"encoding/json"

	temporal_client "github.com/RandySteven/go-cook/temporal"
	"github.com/RandySteven/onboard-be/entities/payloads/requests"
)

type ExecutionData struct {
	CurrActivity string `json:"curr_activity"`
	Status       string `json:"status"`

	AuthPayment    *requests.AuthPaymentRequest    `json:"auth_payment"`
	CapturePayment *requests.CapturePaymentRequest `json:"capture_payment"`

	AdditionalInfo json.RawMessage `json:"additional_info"`
}

func (e *ExecutionData) GetActivity() string {
	return e.CurrActivity
}

func (e *ExecutionData) SetActivity(activityName string) {
	e.CurrActivity = activityName
}

func (e *ExecutionData) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e *ExecutionData) Unmarshal(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	return json.Unmarshal(data, e)
}

func (e *ExecutionData) Clone() temporal_client.ExecutionData {
	if e == nil {
		return &ExecutionData{}
	}
	b, err := e.Marshal()
	if err != nil {
		c := *e
		return &c
	}
	out := &ExecutionData{}
	if err := json.Unmarshal(b, out); err != nil {
		c := *e
		return &c
	}
	return out
}

var _ temporal_client.ExecutionData = (*ExecutionData)(nil)
