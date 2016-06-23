package chat

type Message struct {
	Author string `json:"author"`
	Body   string `json:"message"`
}

func (self *Message) String() string {
	return self.Author + " says " + self.Body
}
