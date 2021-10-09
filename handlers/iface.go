package handlers

type Handler interface {
	AsyncHandleMessage(peerId int, message string)
}
