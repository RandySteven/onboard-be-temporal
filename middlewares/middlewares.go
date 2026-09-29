package middlewares

type Middlewares struct{}

func NewMiddlewares() *Middlewares {
	return &Middlewares{}
}

type ServerMiddleware struct {
	middlewares *Middlewares
}

func RegisterServerMiddleware(middlewares *Middlewares) *ServerMiddleware {
	return &ServerMiddleware{middlewares: middlewares}
}
