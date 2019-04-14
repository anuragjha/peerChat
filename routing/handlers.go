package routing

import (
	//"github.com/gorilla/mux"
	//"strconv"

	//"net"
	"../chat"
	"../peers"
	"../witai"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) { //we take responsewriter and *request as parameter
	w.WriteHeader(200)
	_, _ = w.Write([]byte("hello handler func"))
}

func Wit(w http.ResponseWriter, r *http.Request) { //we take responsewriter and *request as parameter

	witai.Conn(w, r)
}

//// chat start
func Chat(w http.ResponseWriter, r *http.Request) { //we take responsewriter and *request as parameter

	switch r.Method {
	case "GET":
		chat.Begin(w, r)
	case "POST":
		chat.Continue(w, r)
	}

}

func ChatShow(w http.ResponseWriter, r *http.Request) { //we take responsewriter and *request as parameter

	switch r.Method {
	case "GET":
		chat.ShowChatFile(w, r)
	}

}

////// chat over

func PeerAlive(w http.ResponseWriter, r *http.Request) { //we take responsewriter and *request as parameter

	switch r.Method {
	case "POST":
		chat.RecvPeerAlive(w, r)
	}

}
