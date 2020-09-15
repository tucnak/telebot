package telebot

// MiddlewareFunc represents a middleware processing function,
// which get called before the endpoint group or specific handler.
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// Group is a separated group of handlers, united by the general middleware.
type Group struct {
	middleware []MiddlewareFunc
	handlers   map[string]HandlerFunc
}

// Use adds middleware to the chain.
func (g *Group) Use(middleware ...MiddlewareFunc) {
	g.middleware = append(g.middleware, middleware...)
}

// Handle adds endpoint handler to the bot, combining group's middleware
// with the optional given middleware.
func (g *Group) Handle(endpoint interface{}, h HandlerFunc, m ...MiddlewareFunc) {
	if len(g.middleware) > 0 {
		m = append(g.middleware, m...)
	}

	handler := func(c Context) error {
		return applyMiddleware(h, m...)(c)
	}

	switch end := endpoint.(type) {
	case string:
		g.handlers[end] = handler
	case CallbackEndpoint:
		g.handlers[end.CallbackUnique()] = handler
	default:
		panic("telebot: unsupported endpoint")
	}
}
