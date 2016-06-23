package chat

type TextMessage string

type ChatList []string

type Command struct {
	Command string `json:"command"`
	Arg     string `json:"argument"`
}

type GeneralMessage struct {
	Author  string      `json:"author,omitempty"`
	Message TextMessage `json:"message,omitempty"`
	Chats   ChatList    `json:"chats,omitempty"`
	Cmd     Command     `json:"command,omitempty"`
}

func (self *GeneralMessage) String() string {
	return self.Author + " says " + string(self.Message)
}
