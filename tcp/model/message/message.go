package message

/*
Unmarshal essentially copies the byte array to a struct type.
However, in the case of the unmarshal of the type Payload, runtime could not find any destination struct.
The interface in Go just has behavior, not types.
-> 그래서 interface를 필드로 사용하고 싶다면, unmarshal하려는 구조체에 interface 정보를 추가해야함.
e.g. msg := Message{}; msg.Payload := &ChatMsgPayload{}
*/

type Payload interface {
	GetPayload() string
	SetPayload(msg string)
}

// 위 주석으로 인해, Payload interface를 사용하기보단 interface{} 타입을 사용
type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Type: CHECK_TYPE
type CheckPayload struct {
	Name string `json:"name"`
}

func (c *CheckPayload) GetPayload() string {
	return c.Name
}

func (c *CheckPayload) SetPayload(msg string) {
	c.Name = msg
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
