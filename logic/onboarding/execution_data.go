package onboarding_workflow

import (
	"bytes"
	"encoding/json"

	temporal_client "github.com/RandySteven/go-cook/temporal"
	"github.com/RandySteven/onboard-be/entities/models"
	"github.com/RandySteven/onboard-be/entities/payloads/requests"
	"github.com/RandySteven/onboard-be/entities/payloads/responses"
)

type ExecutionData struct {
	Request    *requests.RegisterRequest   `json:"request"`
	Response   *responses.RegisterResponse `json:"response"`
	WorkflowID string                      `json:"workflow_id"`
	RunID      string                      `json:"run_id"`
	Status     string                      `json:"status"`

	Onboarding      *models.Onboarding `json:"onboarding"`
	User            *models.User       `json:"user"`
	ActivationToken string             `json:"activation_token"`
	CurrActivity    string             `json:"curr_activity"`
	NextActivity    string             `json:"next_activity"`
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

	// Activation signals `true` / `false` only to resume the parked step.
	var flag bool
	if err := json.Unmarshal(data, &flag); err == nil {
		return nil
	}

	var asRequest requests.RegisterRequest
	if err := json.Unmarshal(data, &asRequest); err == nil && isRegisterRequestPatch(&asRequest) {
		e.Request = &asRequest
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

func isRegisterRequestPatch(req *requests.RegisterRequest) bool {
	return req.Email != "" || req.Username != "" || req.FirstName != "" || req.PhoneNumber != ""
}

var _ temporal_client.ExecutionData = (*ExecutionData)(nil)
