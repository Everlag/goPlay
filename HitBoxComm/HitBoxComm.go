package HitBoxComm

import(
	"errors"
	"fmt"

	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/gorilla/websocket"	

)


type join struct{
	Method string `json:"method"`
	Params joinParams `json:"params"`
}
type joinParams struct{
	Channel string `json:"channel"`
	Name string `json:"name"`
	Token interface{} `json:"token"`
	IsAdmin bool `json:"isAdmin"`
}

type servers struct{
	IncludedServers []serverInfo
}

type serverInfo struct{
	Server_type string
	Server_ip string
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	//where we store our messages, feel free to block when desired to get input!
	Received chan []byte
}



//whenever a new message comes from the connection, pass it to received
func (c *connection) reader(readyChan chan bool) {
	readyChan <- true

	for{
		_, message, err:= c.ws.ReadMessage()
		if err!=nil{
			fmt.Println("Connection dropped")
			break
		}
		//fmt.Println(string(message))
		c.Received <- message
		//fmt.Println(c.Received)

	}
	c.ws.Close()
}

func (c *connection) WriteString(message string ) error {

	err:= c.ws.WriteMessage(websocket.TextMessage, []byte( message))
	if err!=nil{
		return errors.New("Failed to send websocket message")
	}

	return nil

}

//returns a string ready to be plugged into dialer
//to connect as a client to hitbox chat
func GetChatServer() string {
	//first we want to get the server to connect to.
	//that is given by the json api at http://api.hitbox.tv/chat/servers
	response, err:= http.Get("http://api.hitbox.tv/chat/servers")
	if err!=nil{
		panic("Failed to get server info!")
	}
	defer response.Body.Close()
	fmt.Println("Got Data")
	serverDataRaw, err:= ioutil.ReadAll(response.Body)

	serverDataRaw = []byte( "{\"IncludedServers\":"+ string(serverDataRaw)+ "}" )
	ioutil.WriteFile("ahhhh.json", serverDataRaw, 0664)

	fmt.Println("Read Data")

	var serverData servers
	fmt.Println("Unmarshalling data")
	//uhhh
	_ = json.Unmarshal(serverDataRaw, &serverData)

		//where we'll be connection to get the data streamed
	var serverDestination string

	// now go through the available servers till we hit a chat server... yes, I know this is redundant
	for _, i:= range serverData.IncludedServers{
		if i.Server_type == "chat"{
			serverDestination = i.Server_ip
			break
		}
	}

	fmt.Println(serverDestination)
	return serverDestination
}

//takes a channel name, a channel to act as a lock, and a chan where data is stored
func GetConnection(channel string, readyChan chan bool, commChan chan []byte) connection {
	serverDestination:= GetChatServer()
	
	//get the message to get info streamed ready
	params:= joinParams{Channel:channel, Name:"UnknownSoldier", Token:nil, IsAdmin:false}
	joinMessage:= join{Method:"joinChannel", Params:params}
	readyMessage, err:= json.Marshal(joinMessage)
	if err!=nil{
		fmt.Println("Failed to marhsal message")
		panic(err)
	}

	//readyMessage = []byte(`{"method":"joinChannel","params":{"channel":"kanjo","name":"UnknownSoldier","token":null,"isAdmin":false}}`)

	//get a connection to the server
		//setup the urls
	//origin:= "http://www.hitbox.tv"
	url:= "ws://" + serverDestination

	workingDialer:= websocket.Dialer{ReadBufferSize:4096, WriteBufferSize:4096}
	ws, _, err:= workingDialer.Dial(url, http.Header{})
	if err!=nil{
		fmt.Println(err)
		panic("Failed to connect properly")
	}

	clientConn:= connection{ws:ws, Received:commChan}

	//now get setup to read from the connection
	go clientConn.reader(readyChan)

	//finally, authenticate as an anonymous user
	clientConn.WriteString(string(readyMessage))

	return clientConn
}