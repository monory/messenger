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
	pattern string
	clients map[string]*Client
	addCh   chan *Client
	delCh   chan *Client
	msgCh   chan *ClientMessage
	cmdCh   chan *ClientCommand
	doneCh  chan bool
	errCh   chan error
	db      *sql.DB
}

// Create new chat server.
func NewServer(pattern string) *Server {
	clients := make(map[string]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	msgCh := make(chan *ClientMessage)
	cmdCh := make(chan *ClientCommand)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		pattern,
		clients,
		addCh,
		delCh,
		msgCh,
		cmdCh,
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

func (s *Server) HandleCommand(cmd *ClientCommand) {
	s.cmdCh <- cmd
}

func (s *Server) HandleMessage(msg *ClientMessage) {
	s.msgCh <- msg
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errCh <- err
}

func (s *Server) sendPastMessages(c *Client, contact string) {
	var msg database.Message
	msg.Cmd = &database.Command{Name: "send-messages", Args: database.GetMessages(s.db, c.id, contact)}

	c.Write(&msg)
}

func (s *Server) sendContacts(c *Client) {
	var msg database.Message
	msg.Cmd = &database.Command{Name: "send-contacts", Args: database.GetContacts(s.db, c.id)}

	c.Write(&msg)
}

func (s *Server) handleMessage(msg *ClientMessage) {
	log.Println("MESSAGE RECEIVED", msg.Message.Message, msg.Message.Author, msg.Message.Receiver)
	database.AddMessage(s.db, msg.Message)

	if msg.Client.activeChat == msg.Message.Author ||
		(msg.Client.activeChat == msg.Message.Receiver && msg.Client.name == msg.Message.Author) {
		msg.Client.Write(&database.Message{Txt: msg.Message, Cmd: nil})
	}

	if _, ok := s.clients[msg.Message.Receiver]; ok {
		if s.clients[msg.Message.Receiver].activeChat == msg.Message.Author {
			s.clients[msg.Message.Receiver].Write(&database.Message{Txt: msg.Message, Cmd: nil})
		}
	}
}

func (s *Server) handleCommand(cmd *ClientCommand) {
	switch cmd.Command.Name {
	case "chat-select":
		cmd.Client.activeChat, _ = cmd.Command.Args.(string)
		s.sendPastMessages(cmd.Client, cmd.Client.activeChat)
	case "new-chat":
		err := database.AddContact(s.db, cmd.Client.id, cmd.Command.Args.(string))
		if err != nil {
			log.Println("ERROR:", err)
			return
		}
		s.sendContacts(cmd.Client)
	}
}

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
			s.clients[c.name] = c
			log.Println("Now", len(s.clients), "clients connected.")
			s.sendContacts(c)

		// del a client
		case c := <-s.delCh:
			log.Println("Delete client")
			delete(s.clients, c.name)

		case msg := <-s.msgCh:
			log.Println("Message:", *msg)
			s.handleMessage(msg)

		case cmd := <-s.cmdCh:
			log.Println("Command:", *cmd)
			s.handleCommand(cmd)

		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}
}
