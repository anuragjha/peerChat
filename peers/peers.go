package peers

import (
	"encoding/json"
	"log"
	"sync"
)

//todo add and delete

type Peers struct {
	PeerMap map[string]Peer
	mux     sync.Mutex
}

type Peer struct {
	Id   string
	Addr string
}

func NewPeers() Peers {
	return Peers{
		PeerMap: make(map[string]Peer, 0),
	}
}

func (peers *Peers) CopyPeers() Peers {
	peers.mux.Lock()
	defer peers.mux.Unlock()

	copyOfPeers := *peers

	return copyOfPeers
}

func (peers *Peers) ConvertPeersToJSON() []byte {
	copyOfPeers := peers.CopyPeers()
	peersJSON, _ := json.Marshal(copyOfPeers)

	return peersJSON
}

func ConvertJSONToPeers(peersJSON string) Peers {

	newPeers := NewPeers()
	jerr := json.Unmarshal([]byte(peersJSON), &newPeers)
	if jerr != nil {
		log.Println("Err in unmarshalling ConvertPeersJSONToPeers - err :", jerr)
	}

	return newPeers

}

//func RecvPeerAlive(w http.ResponseWriter, r *http.Request) {
//
//
//
//}
