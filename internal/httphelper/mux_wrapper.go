package httphelper

import "net/http"

// WrapMux satisfies the generated api.ServeMux interface. every route a service
// registers through it lands on the shared mux, wrapped in that service's middleware.
type WrapMux struct {
	mux *http.ServeMux
	mw  func(http.Handler) http.Handler
}

func (m WrapMux) HandleFunc(pattern string, h func(http.ResponseWriter, *http.Request)) {
	m.mux.Handle(pattern, m.mw(http.HandlerFunc(h)))
}

func (m WrapMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

func WithMiddleware(mux *http.ServeMux, mw func(http.Handler) http.Handler) WrapMux {
	return WrapMux{mux: mux, mw: mw}
}
