package database

type TextMessage struct {
	Author   string `json:"author"`
	Message  string `json:"message"`
	Receiver string `json:"receiver"`
}

type Command struct {
	Name string      `json:"name"`
	Args interface{} `json:"args"`
}

type Message struct {
	Txt *TextMessage `json:"text,omitempty"`
	Cmd *Command     `json:"command,omitempty"`
}
