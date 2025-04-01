package model

type Payload interface {
	GetPayload() string
	SetPayload(msg string)
}

type Message struct {
	Type    string  `json:"type"`
	Payload Payload `json:"payload"`
}

// Type: CHECK_TYPE
type CheckPayload struct {
	ID string `json:"id"`
}

func (c *CheckPayload) GetPayload() string {
	return c.ID
}

func (c *CheckPayload) SetPayload(msg string) {
	c.ID = msg
}

// Type: CHAT_MESSAGE_TPYE
type ChatMsgPayload struct {
	Message string `json:"message"`
}

func (c *ChatMsgPayload) GetPayload() string {
	return c.Message
}

func (c *ChatMsgPayload) SetPayload(msg string) {
	c.Message = msg
}
