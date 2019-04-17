package chat

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"

	//"bytes"
	"html/template"
	//"io/ioutil"

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

	"../filehelp"
	"../httphelp"
	"../identity"
	//"./data"
	"../peers"
)

const CHATFILESSENTDIR = "/chat/files/sent/"
const CHATFILESRECVDIR = "/chat/files/recv/"
const CHATFILEPREFIX = "/chat/files/chat-"
const CHATBOXHTML = "/resources/view/html/chatBox.html"

var chatIdentity identity.Identity

var isStarted bool

type Chat struct {
	From       string     `json:"from"`
	To         []string   `json:"to"`
	Message    string     `json:"message"`
	LoadedFile LoadedFile `json:"loadedFile"`
	Timestamp  time.Time  `json:"timestamp"`
}

type ChatShow struct {
	From           string    `json:"from"`
	To             []string  `json:"to"`
	Message        string    `json:"message"`
	LoadedFileName string    `json:"loadedFileName"`
	Timestamp      time.Time `json:"timestamp"`
}

type LoadedFile struct {
	FileName string `json:"filename"`
	FileData []byte `json:"filedata"`
}

type Chats struct {
	ChatList []Chat `json:"chats"`
}

type ChatsShow struct {
	ChatsShowList []ChatShow `json:"chatsShow"`
}

type ChatPage struct {
	IdDS      identity.Identity
	PeersDS   peers.Peers
	ChatsShow ChatsShow
	//Chats Chats//adding

}

func NewChat(identity identity.Identity, to []string, message string, loadedFile LoadedFile) Chat {
	c := Chat{}
	c.Timestamp = time.Now()
	c.From = identity.Id
	c.To = to
	c.Message = message
	c.LoadedFile = loadedFile
	return c
}

func NewChatShow(from string, to []string, message string, loadedFileName string) ChatShow {
	c := ChatShow{}
	c.Timestamp = time.Now()
	c.From = from
	c.To = to
	c.Message = message
	c.LoadedFileName = loadedFileName
	return c
}

func NewLoadedFile(name string, b []byte) LoadedFile {
	return LoadedFile{
		FileName: name,
		FileData: b,
	}
}

func NewChats() Chats {
	c := Chats{}

	return c
}

func NewChatsShow() ChatsShow {
	c := ChatsShow{}

	return c
}

func NewChatPage(idDS identity.Identity, peersDS peers.Peers, chatsShow ChatsShow) ChatPage {
	return ChatPage{
		IdDS:      idDS,
		PeersDS:   peersDS,
		ChatsShow: chatsShow,
	}
}

//internal funcs
func setIdentity(identity identity.Identity) {
	chatIdentity = identity
}

func getIdentity() identity.Identity {
	return chatIdentity
}

func getChatFIlePath() string {
	cwd, _ := os.Getwd()
	chatId := getIdentity()
	filename := cwd + CHATFILEPREFIX + chatId.Id + ".txt"
	return filename
}

//get /chat
func Begin(w http.ResponseWriter, r *http.Request, identity identity.Identity, peerDS peers.Peers) {

	if isStarted == false {
		setIdentity(identity)
		isStarted = true
	}

	showChatsShow(w, identity, peerDS)

}

// showChatsShow connects to HTML and executes the template
func showChatsShow(w http.ResponseWriter, identity identity.Identity, peerDS peers.Peers) {

	chatsToShow := chatFileToChatsShow()

	page := NewChatPage(identity, peerDS, chatsToShow) // to be used in template execute
	cwd, _ := os.Getwd()
	chatBoxHtml := path.Join(cwd, CHATBOXHTML) //"/resources/view/html/chatBox.html")
	t, _ := template.ParseFiles(chatBoxHtml)
	_ = t.Execute(w, page) // todo - form - encType

}

// chatFileToChatsShow func reads the File - chat-id.txt and return a ChatsShow struct
func chatFileToChatsShow() ChatsShow {

	filename := getChatFIlePath()

	inputFile, err := os.Open(filename)
	if err != nil {
		log.Println("Error in readChat() : err - ", err)
	}

	chatsShow := NewChatsShow()
	chatsShow.ChatsShowList = make([]ChatShow, 0)
	chatsShow.initialize(inputFile)

	//fmt.Println("chatFileToChatsShow -  ", chatsShow)

	return chatsShow
}

// initialize func takes file as param and inits the chatsShow struct
func (chatsShow *ChatsShow) initialize(inputFile *os.File) {
	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines) // Read until a newline for each Scan() (Default)
	for inputScanner.Scan() {
		//fmt.Println(inputScanner.Text())   // get the buffered content as string
		//fmt.Println(inputScanner.Bytes())  // same content as above but as []byte
		chatShowJSON := inputScanner.Text()
		chatToShow := ChatShow{}
		jerr := json.Unmarshal([]byte(chatShowJSON), &chatToShow)
		if jerr != nil {
			log.Println("Error while converting json to chatShow - err : ", jerr)
			continue
		}
		chatsShow.ChatsShowList = append(chatsShow.ChatsShowList, chatToShow)
	}
}

//func showChats(w http.ResponseWriter, identity identity.Identity, peerDS peers.Peers) {
//
//	//chats := returnChatJSONstoChats()
//	fmt.Println("BEFORE starting SHOWCHAT")
//
//	showChats := returnChatJSONstoChatsWithoutLoadedFIle()
//	fmt.Println("chats from returnChatJSONstoChatsWithoutLoadedFIle :", showChats)
//
//	page := NewChatPage(identity, peerDS, showChats)//
//	fmt.Println("After generating CHAT PAGE")
//
//	cwd, _ := os.Getwd()
//	chatBoxHtml := path.Join(cwd, CHATBOXHTML) //"/resources/view/html/chatBox.html")
//	t, _ := template.ParseFiles(chatBoxHtml)
//
//	//_ = t.Execute(w, chats.Chats)
//	_ = t.Execute(w, page)
//	//t.ExecuteTemplate(w,"page", p)
//
//	fmt.Println("BEFORE Finishing SHOWCHAT")
//}

//func returnChatJSONstoChats() Chats {
//	//cwd, _ := os.Getwd()
//
//	//chatId := getIdentity()
//	//filename := cwd + CHATFILEPREFIX + chatId.Id +".txt"//"/chat/files/chat.txt"
//	filename := getChatFIlePath()
//
//	inputFile, err := os.Open(filename)
//	if err != nil {
//		log.Println("Error in readChat() : err - ", err)
//	}
//
//	chats := NewChats()
//	chats.ChatList = make([]Chat, 0)
//	chats.buildChatsShow(inputFile)
//
//	return chats
//}

//
//func returnChatJSONstoChatsWithoutLoadedFIle() Chats {
//	fmt.Println("BEFORE starting returnChatJSONstoChatsWithoutLoadedFIle")
//	chats := returnChatJSONstoChats()
//	for _, chat := range chats.ChatList {
//		chat.LoadedFile = LoadedFile{}
//	}
//
//	fmt.Println("Just Finishing returnChatJSONstoChatsWithoutLoadedFIle")
//	return chats
//}

//func (chats *Chats) buildChatsShow(inputFile *os.File) {
//	inputScanner := bufio.NewScanner(inputFile)
//	inputScanner.Split(bufio.ScanLines) // Read until a newline for each Scan() (Default)
//	for inputScanner.Scan() {
//		//fmt.Println(inputScanner.Text())   // get the buffered content as string
//		//fmt.Println(inputScanner.Bytes())  // same content as above but as []byte
//		chatJSON := inputScanner.Text()
//		chat := Chat{}
//		jerr := json.Unmarshal([]byte(chatJSON), &chat)
//		if jerr != nil {
//			log.Println("Error while converting json to chat - err : ", jerr)
//			continue
//		}
//		chats.ChatList = append(chats.ChatList, chat)
//	}
//
//}

// Called when a node sends a Chat object  - In Req body there is chat json
// Chat struct
// Continue func gets param - w, r, identity, peersDS
// Continue func should do -
// 0. save chat to local chat log
// 1. parse recv Chat
// 2.
func Continue(w http.ResponseWriter, r *http.Request, identity identity.Identity, peerDS peers.Peers) {

	//save the chat in req body - generated from submit of chatform
	processChatFormSubmit(r, identity, peerDS)

	showChatsShow(w, identity, peerDS)
}

// processChatFormSubmit saves the chat form message to self folders - recv and update self chat-id.txt copy
func processChatFormSubmit(r *http.Request, identity identity.Identity, peerDS peers.Peers) {
	//save chat in req body - a. save everything except file in chats-id.txt
	// 								- first create chatsShow
	//								- then save to chat file
	//						  if file attached -
	//						  	b. save the file - with name as loaded.name in sent folder
	//											 -  with data as loaded.data
	//
	//

	//err := r.ParseForm() //todo todo
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Fatal("Err in chat : Cannot process form, error - ", err)

	}
	fmt.Println("Chat Form ------ ", r.Form) //

	to := r.Form["peers"]
	if len(to) == 0 {
		to = []string{"All"}
	}
	fmt.Println("To in Chat Form ------ ", to)

	message := r.FormValue("message")
	fmt.Println("Message in Chat Form ------ ", message)

	var chatShow ChatShow // preparing chatshow
	//var fileBytes []byte // for storing file in byte array
	file, handler, ferr := r.FormFile("uploadfile")
	if ferr != nil {

		log.Println("Error in processLoadedOnSubmit func - can also mean no file was attached - Err : ", ferr)
		chatShow = NewChatShow(identity.Id, to, message, "")

	} else {

		//defer file.Close() //else if Loaded is present
		//fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		//fmt.Printf("File Size: %+v\n", handler.Size)
		//fmt.Printf("MIME Header: %+v\n", handler.Header)
		//fileBytes, err = ioutil.ReadAll(file)
		//if err != nil {
		//	log.Println("Error in reading bytes to file : Error - ", err)
		//}
		//log.Println("FILEBYTES pre: ", fileBytes)

		saveLoadedOnSubmit(file, handler)

		//NewChatShow(identity identity.Identity, to []string, message string, loadedFileName string)
		chatShow = NewChatShow(identity.Id, to, message, handler.Filename)
		file.Close()
	}

	filehelp.SaveToFile(getChatFIlePath(), string(ChatShowToJSON(&chatShow))) // saving chatShow to chat-id.txt
	fmt.Println("ChatShow JSON being saved - string(ChatShowToJSON(&chatShow)) : ", string(ChatShowToJSON(&chatShow)))

	////////////////////////////////////////
	//prepare and send ChatBeat //////////// todo
	prepareChatBeat(r, message, to, identity, peerDS)

}

// processLoadedOnSubmit func - save to new file if file is present in form request
func saveLoadedOnSubmit(file multipart.File, handler *multipart.FileHeader) {

	cwd, _ := os.Getwd()
	f, err := os.OpenFile(cwd+CHATFILESSENTDIR+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error in ceating file to save in sent folder : err - ", err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Println("Error in copying file to new created file - ", err)
		return
	}
}

func prepareChatBeat(r *http.Request, message string, to []string, identity identity.Identity, peerDS peers.Peers) {
	file, handler, ferr := r.FormFile("uploadfile")
	if ferr != nil {
		prepareAndSendChatBeat(LoadedFile{}, message, to, identity, peerDS)
	} else {
		fileName := handler.Filename
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println("Error in reading bytes to file : Error - ", err)
		}
		log.Println("FILEBYTES : ", fileBytes[:20])

		loadedFile := NewLoadedFile(fileName, fileBytes)
		//fmt.Println("loadedFile being sent : ", loadedFile)

		prepareAndSendChatBeat(loadedFile, message, to, identity, peerDS)
		file.Close()
	}

}

// prepares chatbeat and calls sendchatBeat
func prepareAndSendChatBeat(loadedFile LoadedFile, message string, to []string, identity identity.Identity, peerDS peers.Peers) {

	c := NewChat(identity, to, message, loadedFile)

	sendchatBeat(c, to, peerDS)

}

// sendchatBeat func sends the chat struct made from chat-form submit to specified peers
func sendchatBeat(c Chat, to []string, peersDS peers.Peers) {
	// sendchatBeat sends the chat generated to peers
	go chatBeat(c, to, peersDS)

}

// chatBeat builds addresses to send the chat and then beats it
func chatBeat(c Chat, to []string, peersDS peers.Peers) {

	addrToSend := buildAddrToSend(to, peersDS)

	chatBeatingNow(addrToSend, c, peersDS)

}

// buildAddrToSend func buils a addrlist for a chat
func buildAddrToSend(to []string, peersDS peers.Peers) []string {

	var addrToSend []string
	for id, peer := range peersDS.PeerMap {
		for _, t := range to {
			if id == t {
				addrToSend = append(addrToSend, peer.Addr)
			}
		}
	}
	return addrToSend
}

// chatBeatingNow func beats a chat to all needed peers
func chatBeatingNow(addrToSend []string, c Chat, peersDS peers.Peers) {
	if len(addrToSend) == 0 { // to ALL
		c.To = []string{"All"}
		for _, peer := range peersDS.PeerMap {
			//fmt.Println("in SendMessage ChatToJSON(&c) --  all peers ------>  ", string(ChatToJSON(&c)))
			chatBeatingToAddress(peer.Addr, c)
		}
	} else { // to certain peers
		for _, addr := range addrToSend {
			//fmt.Println("in SendMessage ChatToJSON(&c) --  directed peers ------>  ", string(ChatToJSON(&c)))
			chatBeatingToAddress(addr, c)
		}
	}
}

// chatBeatingToAddress func beats a chat to an address
func chatBeatingToAddress(peerAddr string, c Chat) {
	_, err := http.Post("http://"+peerAddr+"/chat/recv", "json", bytes.NewBuffer([]byte(ChatToJSON(&c))))
	if err != nil {
		log.Println("Error in Send Message : ", err)
	}
}

//
//func Continue(w http.ResponseWriter, r *http.Request, identity identity.Identity, peerDS peers.Peers) { //todo refactor
//	// Continue - is called when user press submit method //
//	//read message - present in req body
//	err := r.ParseForm()
//	if err != nil {
//		log.Println("Err in chat : Cannot process form, error - ", err)
//	}
//
//	   /////////////////////                     //  /// taking data from HTML form field
//
//	to := r.Form["peers"]
//	if len(to) == 0 {
//		to = []string{"All"}
//	}
//
//	message := r.FormValue("message")
//	//upload := r.ParseMultipartForm() //todo
//
//					// we want to save the file before sending client - saving  now !!!
//	var c Chat // to be used in sed message !!
//	var f *os.File
//	file, handler, err := r.FormFile("uploadfile")
//	if err != nil {
//		log.Println("Error in Continue func - error : ", err)
//		c = NewChat(identity, to, message, LoadedFile{}/*NewLoadedFile(fileName, fileBytes)*/)
//
//
//	} else {
//		defer file.Close()
//		// !!! saving the file on chat file system
//		f, err = os.OpenFile(CHATFILESSENTDIR+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
//		if err != nil {
//			fmt.Println("Error in ceating file to save in sent folder : err - ", err)
//		} else {
//			defer f.Close()
//			_, _ = io.Copy(f, file)
//			// closing taking data from HTML form field ///
//
//			fmt.Println("::::; handler.Header ::: - handler.Header :", handler.Header)
//			fileName := handler.Filename
//			fileBytes, err := ioutil.ReadAll(f)
//			if err != nil {
//				log.Println("Error in reading bytes to file : Error - ", err)
//			}
//
//			c = NewChat(identity, to, message, NewLoadedFile(fileName, fileBytes))
//		}
//
//	}
//
//
//
//// todo up and down
//
//	filehelp.SaveToFile(string(ChatToJSON(&c)), getChatFIlePath()) //saving here
//
//	//send message to particular nodes
//	go SendMessage(c, to, peerDS) //
//	//////
//
//	//Begin(w, r, identity, peerDS)
//	showChatsShow(w, identity, peerDS)
//
//
//
//}

//func SendMessage(c Chat, to []string, peersDS peers.Peers) { //send to peer or all
//
//	var addrToSend []string
//	for id, peer := range peersDS.PeerMap {
//		for _, t := range to {
//			if id == t {
//				addrToSend = append(addrToSend, peer.Addr)
//			}
//		}
//	}
//
//
//	if len(addrToSend) == 0 {  // to ALL
//		c.To = []string{"All"}
//		for _, peer := range peersDS.PeerMap {
//			fmt.Println("in SendMessage ChatToJSON(&c) --  all peers ------>  ", ChatToJSON(&c))
//			_, err := http.Post("http://"+peer.Addr+"/chat/recv", "json", bytes.NewBuffer([]byte(ChatToJSON(&c))))
//			if err != nil {
//				log.Println("Error in Send Message : ", err)
//			}
//		}
//	} else {  // to
//		for _, addr := range addrToSend {
//			fmt.Println("in SendMessage ChatToJSON(&c) --  directed peers ------>  ", ChatToJSON(&c))
//			_, err := http.Post("http://"+addr+"/chat/recv", "json", bytes.NewBuffer([]byte(ChatToJSON(&c))))
//			if err != nil {
//				log.Println("Error in Send Message : ", err)
//			}
//		}
//	}
//
//}

// Receives ChatBeat
func BeatRecv(w http.ResponseWriter, r *http.Request, identity identity.Identity, peerDS peers.Peers) { //receieve chatbeat from peers

	body := httphelp.ReadHttpRequestBody(r)

	chat := Chat{}
	jerr := json.Unmarshal(body, &chat)
	if jerr != nil {
		log.Println("Error in unmarshaling in BeatRecv !!!!!!! : ", jerr)

	}

	//if chat.LoadedFile.FileData != nil { // to save file in recv'ed chat Beat
	if chat.LoadedFile.FileName != "" { // to save file in recv'ed chat Beat
		loadedToSaveName := chat.LoadedFile.FileName
		loadedToSaveData := chat.LoadedFile.FileData
		saveLoadedFile(loadedToSaveName, loadedToSaveData)
	}

	//func NewChatShow(from string, to []string, message string, loadedFileName string) ChatShow
	cs := NewChatShow(chat.From, chat.To, chat.Message, chat.LoadedFile.FileName)
	filehelp.SaveToFile(getChatFIlePath(), string(ChatShowToJSON(&cs)))

}

func saveLoadedFile(name string, data []byte) {
	cwd, _ := os.Getwd()
	// !!! saving the file on chat file system
	f, err := os.OpenFile(cwd+CHATFILESRECVDIR+name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error in ceating file to save in Recv folder : err - ", err)
	}
	defer f.Close()

	f.Write(data)

}

//func ShowChatFile(w http.ResponseWriter, r *http.Request) {
//
//	chats := returnChatJSONstoChats()
//	chats.displayChats(w, r)
//
//}

//func (chats *Chats) displayChats(w http.ResponseWriter, r *http.Request) {
//
//	//log := ""
//	for _, chat := range chats.ChatList {
//		_, _ = fmt.Fprint(w, chat.Message+"\n")
//		//chat.Message+"\t\t\ttime:"+chat.Timestamp.String()+"\n")
//	}
//
//}

/////////////////////////////////

//func createChatsFile(filename string) {
//
//	var file, err = os.Create(filename)
//	if err != nil {
//		log.Println("Chat file not present also Cannot create chat file : err - ", err)
//		return
//	}
//	defer file.Close()
//}

//internal func

//func ChatToJSON(c *Chat) []byte {
//	j, _ := json.Marshal(c)
//	return j
//}
//
//func JSONToChat(b []byte) Chat {
//	c := Chat{}
//	jerr := json.Unmarshal(b, &c)
//	if jerr != nil {
//		log.Println("Error in unmarshalling json to Chat - err : ", jerr)
//	}
//
//	return c
//}
//
//
//func ChatShowToJSON(c *ChatShow) []byte {
//	j, _ := json.Marshal(c)
//	return j
//}
//
//func JSONToChatShow(b []byte) ChatShow {
//	c := ChatShow{}
//	jerr := json.Unmarshal(b, &c)
//	if jerr != nil {
//		log.Println("Error in unmarshalling json to Chat - err : ", jerr)
//	}
//
//	return c
//}
