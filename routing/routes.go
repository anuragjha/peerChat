package routing

import "net/http"

type Route struct {
	Method      string
	Path        string
	Name        string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		Method:      "GET",
		Path:        "/hello",
		Name:        "hello",
		HandlerFunc: Hello,
	},
	Route{
		Method:      "GET",
		Path:        "/wit",
		Name:        "wit/{textInput}",
		HandlerFunc: Wit,
	},
	Route{
		Method:      "GET",
		Path:        "/chat",
		Name:        "chat",
		HandlerFunc: Chat,
	},
	Route{
		Method:      "POST",
		Path:        "/chat",
		Name:        "chat",
		HandlerFunc: Chat,
	},
	Route{
		Method:      "GET",
		Path:        "/chat/show",
		Name:        "chatShow",
		HandlerFunc: ChatShow,
	},
	Route{
		Method:      "POST",
		Path:        "/peeralive",
		Name:        "peeralive",
		HandlerFunc: PeerAlive,
	},
}
