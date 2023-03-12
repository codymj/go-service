package web

// Middleware is a function designed to run before and/or after another Handler.
type Middleware func(Handler) Handler

// wrapMiddleware creates a new Handler by wrapping middleware around a final
// Handler.
func wrapMiddleware(mw []Middleware, handler Handler) Handler {
	// loop backwards through the middleware chain, invoking each one
	for i := len(mw) - 1; i >= 0; i-- {
		// get handler in chain
		h := mw[i]

		// wrap with new handler
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}
