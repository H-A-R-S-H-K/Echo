package models

import (
	"fmt"
	"net"
	"sync"
)

type ServerMap struct {
	mu sync.RWMutex
	servers map[string] *Server // map a room key -> server (models.Room)
}

func (sM *ServerMap) AddServer(room_key string, server *Server) {

	sM.mu.Lock()
	defer sM.mu.Unlock()

	if _, exists := sM.servers[room_key]; exists {
		fmt.Printf("room with key %v already exists\n", room_key)
		return
	}

	sM.servers[room_key] = server
	fmt.Printf("connected with room having room_key = %v\n", room_key)
}

func (sM *ServerMap) GetServer(room_key string) bool {

	sM.mu.RLock()
	defer sM.mu.RUnlock()

	_, exists := sM.servers[room_key]
	return exists;
}	

func (sM *ServerMap) RemoveServer(room_key string) {

	sM.mu.Lock()
	defer sM.mu.Unlock()

	if _, exists := sM.servers[room_key]; exists {
		delete(sM.servers, room_key)
		fmt.Printf("Removed the server with room_key = %v\n", room_key)
		return
	}
	fmt.Printf("The chat-room does not exist\n")
}

// A client can connect to multiple servers and it 
// will need a new room key to connect to the new server
// room key will have -> subnet + port encrypted into it.
func (client *Client) ConnectToNewServer(roomKey RoomKey) {

	room_key, subnet, port := roomKey.Room_key, roomKey.Subnet, roomKey.Port

	// check if the room key already exists, i.e the server 
	if client.ServerMap.GetServer(room_key) {
		fmt.Printf("the server with the room key = %v already exists", room_key)
		return
	}

	// create a new server
	server := Server{
		Port: port,
		Subnet: subnet,
	}

	// add the new server to the map.
	client.ServerMap.AddServer(room_key, &server)

	// establish connection with the new server
	client.connectToServer(&server);
}

// only responsibility is to connect to the given server.
func (client *Client) connectToServer(server *Server) (net.Conn, error) {

    ip := net.IP(server.Subnet[:]).String() 
	address := fmt.Sprintf("%s:%d", ip, server.Port)
	conn, err := net.Dial("tcp", address)

	if err != nil {
		fmt.Printf("error connecting to server %v", err)
	}

	return conn, err
	// defer conn.Close()
}

func (client *Client) SendMessage() {

}

func (client *Client) RecvMessage() {
	
}

// this is our client
type Client struct {
	Client_id   int
	Client_name string
	Rooms     []Room
	Status    ClientSTATUS
	ServerMap // store all the conencted servers.
}
