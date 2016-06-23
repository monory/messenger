package chat

type TextMessage string

type ChatList []string

type GeneralMessage struct {
	Author  string      `json:"author,omitempty"`
	Message TextMessage `json:"message,omitempty"`
	Chats   ChatList    `json:"chats,omitempty"`
}

func (self *GeneralMessage) String() string {
	return self.Author + " says " + string(self.Message)
}
