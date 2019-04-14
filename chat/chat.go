package chat

import (
	"bufio"
	"bytes"
	"html/template"

	//"bytes"
	"encoding/json"
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	"log"
	"time"

	//"io/ioutil"
	//"io/ioutil"
	"net/http"
	"path"

	"os"

	"../peers"
)

const CHATFILE = "/chat/files/chat.txt"
const CHATBOXHTML = "/resources/view/html/chatBox.html"

var isStarted bool

var Peers peers.Peers

var identity = peers.Peer{
	Id:   os.Args[1],
	Addr: "localhost:" + os.Args[1],
}

type Alive struct {
	Identity peers.Peer
	PeerList peers.Peers
}

type Chat struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type Chats struct {
	Chats []Chat
}

func NewChat(message string) Chat {
	c := Chat{}
	c.Timestamp = time.Now()
	c.Message = message

	return c
}

func NewChats() Chats {
	c := Chats{}

	return c
}

func ShowJSON(c *Chat) []byte {
	j, _ := json.Marshal(c)
	return j
}

func Begin(w http.ResponseWriter, r *http.Request) {

	if isStarted == false {
		Peers = peers.NewPeers()

		master := peers.Peer{
			Id:   "8080",
			Addr: "localhost:" + os.Args[1] + "/",
		}

		Peers.PeerMap[master.Id] = master

		isStarted = true

		go shareAlive()
	}

	//http.ServeFile(w,r,chatBoxHtml)
	chats := returnChatJSONstoChats()

	cwd, _ := os.Getwd()
	chatBoxHtml := path.Join(cwd, CHATBOXHTML) //"/resources/view/html/chatBox.html")
	t, _ := template.ParseFiles(chatBoxHtml)

	_ = t.Execute(w, chats.Chats)

}

func returnChatJSONstoChats() Chats {
	cwd, _ := os.Getwd()

	filename := cwd + CHATFILE //"/chat/files/chat.txt"

	inputFile, err := os.Open(filename)
	if err != nil {
		log.Println("Error in readChat() : err - ", err)
	}

	chats := NewChats()
	chats.Chats = make([]Chat, 0)
	chats.build(inputFile)

	return chats
}

func (chats *Chats) build(inputFile *os.File) {
	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines) // Read until a newline for each Scan() (Default)
	for inputScanner.Scan() {
		//fmt.Println(inputScanner.Text())   // get the buffered content as string
		//fmt.Println(inputScanner.Bytes())  // same content as above but as []byte
		chatJSON := inputScanner.Text()
		chat := Chat{}
		jerr := json.Unmarshal([]byte(chatJSON), &chat)
		if jerr != nil {
			log.Println("Error while converting json to chat - err : ", jerr)
			continue
		}
		chats.Chats = append(chats.Chats, chat)
	}

}

func Continue(w http.ResponseWriter, r *http.Request) {

	//read message - present in req body
	err := r.ParseForm()
	if err != nil {
		log.Println("Err in chat : Cannot process form, error - ", err)
	}

	//fmt.Fprint(w, r.PostForm) //map that contains all fieldname as key and their value as value
	//message := r.PostForm["message"]
	message := r.FormValue("message")

	c := NewChat(message)
	save(string(ShowJSON(&c)))

	Begin(w, r)

}

func save(m string) {
	cwd, _ := os.Getwd()

	filename := cwd + "/chat/files/chat.txt"
	//return ioutil.WriteFile(filename, []byte(m), 0600)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(m + "\n"); err != nil {
		panic(err)
	}

}

func ShowChatFile(w http.ResponseWriter, r *http.Request) {

	chats := returnChatJSONstoChats()
	chats.displayChats(w, r)

}

func (chats *Chats) displayChats(w http.ResponseWriter, r *http.Request) {

	//log := ""
	for _, chat := range chats.Chats {
		_, _ = fmt.Fprint(w, chat.Message+"\n")
		//chat.Message+"\t\t\ttime:"+chat.Timestamp.String()+"\n")
	}

}

/////////////////////////////////

func shareAlive() {

	for true {

		client := &http.Client{
			/*CheckRedirect: redirectPolicyFunc,*/
		}

		peersCopy := Peers.CopyPeers()
		peersCopyJSON := peersCopy.ConvertPeersToJSON()

		for i, p := range peersCopy.PeerMap {
			req, err := http.NewRequest("POST", p.Addr+"/peeralive", bytes.NewBuffer(peersCopyJSON))
			req.Header.Add("Content-Type", `json`)

			_, err = client.Do(req)
			if err != nil {
				fmt.Println("Error in sharAlive - err : ", i, err)

				continue
			}
		}

		time.Sleep(10 * time.Second)
	}

}

//func ShowChatFile(w http.ResponseWriter, r *http.Request) {
//
//
//	chats := returnChatJSONstoChats()
//	chats.displayChats(w,r)
//
//}

func RecvPeerAlive(w http.ResponseWriter, r *http.Request) {

	for _, chat := range chats.Chats {
		_, _ = fmt.Fprint(w, chat.Message+"\n")
		//chat.Message+"\t\t\ttime:"+chat.Timestamp.String()+"\n")
	}

}
