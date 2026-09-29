package enums

type ContextKey string

const (
	UserID    ContextKey = "user_id"
	RequestID ContextKey = "request_id"
	ClientIP  ContextKey = "client_ip"
)

func (c ContextKey) ToString() string {
	return string(c)
}
