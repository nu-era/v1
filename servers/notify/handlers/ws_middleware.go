package handlers

import "net/http"

//WsHeader is a middleware handler that adds headers to satisfy
//Websocket upgrade constraints
type WsHeader struct {
	handler http.Handler
}

//ServeHTTP handles the request by passing it to the real
//handler and adding the necessary CORS headers
func (l *WsHeader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//add CORS headers to response
	r.Header.Set("Connection", r.Header.Get("X-Connection"))
	r.Header.Set("Upgrade", r.Header.Get("X-Upgrade"))
	r.Header.Set("Sec-Websocket-Key", r.Header.Get("X-Sec-Websocket-Key"))
	r.Header.Set("X-Sec-Websocket-Accept", r.Header.Get("Sec-Websocket-Accept"))
	l.handler.ServeHTTP(w, r)
}

//NewWsMiddleware constructs a new  Websocket Header middleware handler
func NewWsMiddleware(handlerToWrap http.Handler) *WsHeader {
	return &WsHeader{handlerToWrap}
}
