package utils

type Message struct {
	Type int
	Msg  interface{}
}

func CreateMessage(tp int, msg interface{}) *Message {
	return &Message{Type: tp, Msg: msg}
}
