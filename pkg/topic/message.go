package topic

// wrapper for content to be kept in a channel
type Message struct {
	content []byte
}

func NewMessage(content []byte) *Message {
	return &Message{
		content: content,
	}
}

func (m *Message) String() string {
	return string(m.content)
}

func (m *Message) Bytes() []byte {
	return m.content
}
