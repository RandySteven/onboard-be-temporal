package requests

import (
	"encoding/json"

	"github.com/RandySteven/onboard-be/enums"
)

type (
	RegisterRequest struct {
		FirstName      string          `json:"first_name"`
		LastName       string          `json:"last_name"`
		Email          string          `json:"email"`
		Username       string          `json:"username"`
		Password       string          `json:"password"`
		PhoneNumber    string          `json:"phone_number"`
		Address        string          `json:"address"`
		RegisterAs     enums.RegisterAs `json:"register_as"`
		AdditionalInfo json.RawMessage `json:"additional_info"`
	}

	RegisterUserRequestData struct {
		Occupation     string `json:"occupation"`
		CompanyName    string `json:"company_name"`
		ReasonToUseApp string `json:"reason_to_use_app"`
	}

	ActivateRequest struct {
		Token string `json:"token"`
	}
)
