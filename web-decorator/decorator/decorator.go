package decorator

import "net/http"

type DecoratorFunc func(w http.ResponseWriter, r *http.Request, h http.Handler)

type DecoratorHnadler struct {
	fn DecoratorFunc
	h  http.Handler
}

func (self *DecoratorHnadler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.fn(w, r, self.h)
}

func NewDecoratorHandler(fn DecoratorFunc, h http.Handler) http.Handler {
	return &DecoratorHnadler{
		fn: fn,
		h:  h,
	}
}
