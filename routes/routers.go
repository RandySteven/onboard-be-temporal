package routes

import (
	"log"
	"net/http"

	"github.com/RandySteven/onboard-be/enums"
	"github.com/RandySteven/onboard-be/handlers"
	"github.com/RandySteven/onboard-be/middlewares"
	"github.com/gorilla/mux"
)

type (
	HandlerFunc func(w http.ResponseWriter, r *http.Request)

	Router struct {
		methodName  string
		path        string
		handler     HandlerFunc
		method      string
		middlewares []enums.Middleware
	}

	RouterPrefix map[enums.RouterPrefix][]*Router
)

func NewEndpointRouters(api *handlers.Handlers) RouterPrefix {
	endpointRouters := make(RouterPrefix)
	endpointRouters[enums.AuthPrefix] = []*Router{
		Post("RegisterUser", "/register", api.OnboardingHandler.RegisterUser),
		Post("ActivationUser", "/activated", api.OnboardingHandler.ActivateUser),
	}
	return endpointRouters
}

func InitRouter(routers RouterPrefix, r *mux.Router) {
	middleware := middlewares.NewMiddlewares()
	serverMiddleware := middlewares.RegisterServerMiddleware(middleware)

	r.Use(
		serverMiddleware.LoggingMiddleware,
		serverMiddleware.CorsMiddleware,
		serverMiddleware.TimeoutMiddleware,
	)

	r.HandleFunc("/health", handlers.Health).Methods(http.MethodGet)

	onboardingRouter := r.PathPrefix(enums.AuthPrefix.ToString()).Subrouter()
	for _, route := range routers[enums.AuthPrefix] {
		onboardingRouter.HandleFunc(route.path, route.handler).Methods(route.method)
	}

	routerLog(r)
}

func routerLog(r *mux.Router) {
	_ = r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		if len(methods) == 0 {
			log.Printf("         %-30s (subrouter)", path)
		} else {
			log.Printf("%-8s %-30s", methods[0], path)
		}
		return nil
	})
}
