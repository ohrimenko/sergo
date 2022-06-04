package controllers

type BroadcastWebsocket struct {
	WsKey     uint64
	StackKey  uint64
	ClientKey uint64
	Message   string
}

func NewBroadcastWebsocket(wsKey uint64, stackKey uint64, clientKey uint64, message string) *BroadcastWebsocket {
	wsBroadcast := &BroadcastWebsocket{
		WsKey:     wsKey,
		StackKey:  stackKey,
		ClientKey: clientKey,
		Message:   message,
	}

	return wsBroadcast
}
