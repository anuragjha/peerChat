# peerChat
Peer to peer chat application

### Usage
- go run main.go <'port'>
- /chat to open chatbox page


### Features
1. ChatBox
2. Decentralized

### Features To be added
3. Security


#### 1. Chatbox
a. GET /chat    
> Opens a chat interface to send messages and files across to other peers.
#### 2. Decentralized
a. Each bot is a server and holds own memory for processing.
b. Mechanism to know alive peers and identify dead peers by sending and receieving alive beats.
#### 3. Identity & Peers
a. Information of a bot on Peer network, like addr, etc
b. Peers detection mechanism share alieve beat to other peers and also detects if a peer is down. 
#### 4. Security
a. Public-Private Key for data authentication and integrity (signature).
b. Symetric Key for confidentiality (encryption-decryption).
#### 5. Communication
> Communication will be done in Beats. 
Beats will contain among other things ->      
  * a. Message (M) will be encrypted with a new secret symetric key (Sk)
  * b. Sk will be encrypted with Public Key of Peer (Pk)
