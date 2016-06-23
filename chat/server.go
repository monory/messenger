package chat

import (
	"database/sql"
	"io"
	"log"
	"net/http"

	"github.com/monory/messenger/auth"
	"github.com/monory/messenger/database"

	"golang.org/x/net/websocket"
)

// Chat server.
type Server struct {
	pattern   string
	messages  []*GeneralMessage
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *GeneralMessage
	doneCh    chan bool
	errCh     chan error
	db        *sql.DB
}

// Create new chat server.
func NewServer(pattern string) *Server {
	messages := []*GeneralMessage{}
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *GeneralMessage)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		pattern,
		messages,
		clients,
		addCh,
		delCh,
		sendAllCh,
		doneCh,
		errCh,
		nil,
	}
}

func (s *Server) Add(c *Client) {
	s.addCh <- c
}

func (s *Server) Del(c *Client) {
	s.delCh <- c
}

func (s *Server) SendAll(msg *GeneralMessage) {
	s.sendAllCh <- msg
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendPastMessages(c *Client) {
	for _, msg := range s.messages {
		c.Write(msg)
	}
}

func (s *Server) sendChats(c *Client) {
	var msg GeneralMessage
	msg.Chats = database.GetChats(s.db)

	c.Write(&msg)
}

func (s *Server) sendAll(msg *GeneralMessage) {
	for _, c := range s.clients {
		c.Write(msg)
	}
}

// Listen and serve.
// It serves client connection and broadcast request.
func (s *Server) Listen(db *sql.DB) {

	log.Println("Listening server...")
	s.db = db

	// websocket handler
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		var encodedToken string
		err := websocket.Message.Receive(ws, &encodedToken)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		token := auth.NewUserToken()
		err = token.FromString(encodedToken)
		if err != nil {
			log.Println(err)
			return
		}

		name, err := auth.CheckChatToken(db, token)
		if err != nil {
			log.Println(err)
			return
		}

		client := NewClient(ws, s, name)
		s.Add(client)
		client.Listen()
	}
	http.Handle(s.pattern, websocket.Handler(onConnected))
	log.Println("Created handler")

	for {
		select {

		// Add new a client
		case c := <-s.addCh:
			log.Println("Added new client")
			s.clients[c.id] = c
			log.Println("Now", len(s.clients), "clients connected.")
			s.sendPastMessages(c)
			s.sendChats(c)

		// del a client
		case c := <-s.delCh:
			log.Println("Delete client")
			delete(s.clients, c.id)

		// broadcast message for all clients
		case msg := <-s.sendAllCh:
			log.Println("Send all:", msg)
			s.messages = append(s.messages, msg)
			s.sendAll(msg)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
