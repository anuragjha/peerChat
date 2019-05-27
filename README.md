# web-o-bot
Application to be installed on every drone

### Features Included
1. Decentralized
2. Identity
3. ChatBox

### Features To be added
3. Security
4. Communication - using security
5. Speech
6. Sentiment Analysis - API 
7. Intelligent Joker - If sentiment sad or angry -> then joke


#### 1. Decentralized
a. Each bot is a server and holds own memory for processing.
b. Mechanism to know alive peers and identify dead peers by sending and receieving alive beats.

#### 2. Identity
a. Information of a bot on Peer network, like addr, etc

#### 3. Chatbox
a. GET /chat    
> Opens a chat interface to send messages and files across to other peers.

#### 4. Security
a. Public-Private Key for data authentication and integrity (signature).
b. Symetric Key for confidentiality (encryption-decryption).

#### 5. Communication
> Communication will be done in Beats. 

Beats will contain among other things ->      
  a. Message (M) will be encrypted with a new secret symetric key (Sk)
  b. Sk will be encrypted with Public Key of Peer (Pk)
