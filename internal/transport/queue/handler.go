package queue

type Handler interface {
	Register(r *Router)
}

type Router struct {
	routeMap map[string]func(body []byte)
}

func NewRouter() *Router {
	return &Router{routeMap: map[string]func(body []byte){}}
}

func (r *Router) Accept(messageType string, hf func(body []byte)) {
	r.routeMap[messageType] = hf
}

func (r *Router) find(messageType string) (hf func(body []byte), ok bool) {
	hf, ok = r.routeMap[messageType]

	return hf, ok
}
