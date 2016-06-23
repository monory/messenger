package chat

import (
	"fmt"
	"io"
	"log"

	"github.com/monory/messenger/database"

	"golang.org/x/net/websocket"
)

const channelBufSize = 100

type Client struct {
	id   int64
	name string

	ws     *websocket.Conn
	server *Server

	ch     chan *database.Message
	doneCh chan bool

	activeChat string
}

type ClientCommand struct {
	Client  *Client
	Command *database.Command
}

type ClientMessage struct {
	Client  *Client
	Message *database.TextMessage
}

func NewClient(ws *websocket.Conn, server *Server, name string) *Client {
	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}

	id, err := database.GetUserID(server.db, name)
	if err != nil {
		log.Println("Error:", err)
		return nil
	}

	ch := make(chan *database.Message, channelBufSize)
	doneCh := make(chan bool)

	return &Client{id, name, ws, server, ch, doneCh, ""}
}

func (c *Client) Conn() *websocket.Conn {
	return c.ws
}

func (c *Client) Write(msg *database.Message) {
	select {
	case c.ch <- msg:
	default:
		c.server.Del(c)
		err := fmt.Errorf("client %d is disconnected", c.id)
		c.server.Err(err)
	}
}

func (c *Client) Done() {
	c.doneCh <- true
}

func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

func (c *Client) listenWrite() {
	log.Println("Listening write to client")
	for {
		select {

		// send message to the client
		case msg := <-c.ch:
			log.Println("Send:", *msg)
			websocket.JSON.Send(c.ws, *msg)

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

func (c *Client) listenRead() {
	log.Println("Listening read from client")
	for {
		select {

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var msg database.Message
			err := websocket.JSON.Receive(c.ws, &msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				if msg.Cmd != nil {
					// log.Print("A command!", *msg.Cmd)
					log.Print("A command!")
					c.server.HandleCommand(&ClientCommand{c, msg.Cmd})
				} else {
					// log.Print("A message!", *msg.Txt)
					log.Print("A message!")
					msg.Txt.Author = c.name
					c.server.HandleMessage(&ClientMessage{c, msg.Txt})
				}
			}
		}
	}
}
