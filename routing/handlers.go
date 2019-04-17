package routing

import (
	//"github.com/gorilla/mux"
	//"strconv"

	//"net"
	"../chat"
	//"bytes"
	"fmt"
	"os"
	"strings"

	//"strings"
	"time"

	"../identity"
	"../peers"
	"../witai"
	"net/http"
)

var isStarted bool

var PeersDS peers.Peers
var DeletedPeersDS peers.Peers

var idDS identity.Identity

func init() {

	var w http.ResponseWriter
	var r http.Request

	Start(w, &r)

}

func Start(w http.ResponseWriter, r *http.Request) { //we take responsewriter and *request as parameter

	idDS = identity.Identity{
		Id:   os.Args[1],
		Addr: "localhost:" + os.Args[1],
	}

	PeersDS = peers.NewPeers()
	DeletedPeersDS = peers.NewPeers()

	if strings.Compare(idDS.Id, "8080") != 0 {
		PeersDS.PeerMap["8080"] = peers.Peer{
			Id:   "8080",
			Addr: "localhost:8080",
		}
	}

	fmt.Println("starting alive share")
	go SendAliveBeat()

}

func SendAliveBeat() {
	for true {
		peers.SendAliveBeat(idDS, &PeersDS, &DeletedPeersDS)
		time.Sleep(5 * time.Second)
	}
}

func PeersAlive(w http.ResponseWriter, r *http.Request) { //we take responsewriter and *request as parameter

	switch r.Method {
	case "POST":
		peers.RecvPeerAlive(w, r, idDS, &PeersDS, &DeletedPeersDS)
	case "GET":
		peers.ShowPeerAlive(w, r, idDS, PeersDS)
	}

}

func Wit(w http.ResponseWriter, r *http.Request) { //we take responsewriter and *request as parameter

	witai.Conn(w, r)
}

//// chat start
func Chat(w http.ResponseWriter, r *http.Request) { //we take responsewriter and *request as parameter
	//public api point

	switch r.Method {
	case "GET":
		chat.Begin(w, r, idDS, PeersDS)

	case "POST":
		chat.Continue(w, r, idDS, PeersDS) //process chat form submit
	}

}

func ChatBeatRecv(w http.ResponseWriter, r *http.Request) { //ChatBeatRecv takes in message beat from other peers
	//// private api point - chatBeatSend() - defined in chat.go

	switch r.Method {
	case "POST":
		chat.BeatRecv(w, r, idDS, PeersDS)

	}

}

////// chat over

func Hello(w http.ResponseWriter, r *http.Request) { //we take responsewriter and *request as parameter
	w.WriteHeader(200)
	_, _ = w.Write([]byte("hello"))
}
