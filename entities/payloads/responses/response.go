package responses

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func NewResponse(message string, data any, err error) *Response {
	resp := &Response{Message: message, Data: data}
	if err != nil {
		resp.Error = err.Error()
	}
	return resp
}
