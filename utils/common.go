package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

func MapJSONToInterface(jsonRequest json.RawMessage, result interface{}) error {
	if len(jsonRequest) == 0 {
		return json.Unmarshal([]byte("{}"), result)
	}
	return json.Unmarshal(jsonRequest, result)
}

func MergeName(firstName, lastName string) string {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	if lastName == "" {
		return firstName
	}
	if firstName == "" {
		return lastName
	}
	return fmt.Sprintf("%s %s", firstName, lastName)
}

func AdditionalInfoOrEmpty(raw json.RawMessage) json.RawMessage {
	if len(raw) == 0 {
		return json.RawMessage(`{}`)
	}
	return raw
}
