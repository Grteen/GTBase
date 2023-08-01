package utils

type Message struct {
	Type int
}

func CreateMessage(tp int) *Message {
	return &Message{Type: tp}
}
