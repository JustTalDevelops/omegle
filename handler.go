package omegleapi

// Handler is an event handler for an Omegle slave.
type Handler interface {
	OnConnected(slave *Slave)
	OnCommonLikesReceived(slave *Slave, commonLikes []string)
	OnTyping(slave *Slave)
	OnTypingStopped(slave *Slave)
	OnMessage(slave *Slave, message string)
	OnStrangerDisconnected(slave *Slave)
	OnServerDisconnected(slave *Slave)
}

// NopHandler implements the Handler interface but does not execute any code when an event is called. The
// default handler of slaves is set to NopHandler. Users may embed NopHandler to avoid having to implement
// each method.
type NopHandler struct{}

func (h NopHandler) OnTypingStopped(_ *Slave) {}

func (h NopHandler) OnConnected(_ *Slave) {}

func (h NopHandler) OnCommonLikesReceived(_ *Slave, _ []string) {}

func (h NopHandler) OnTyping(_ *Slave) {}

func (h NopHandler) OnMessage(_ *Slave, _ string) {}

func (h NopHandler) OnStrangerDisconnected(_ *Slave) {}

func (h NopHandler) OnServerDisconnected(_ *Slave) {}

// Compile time check to make sure NopHandler implements Handler.
var _ Handler = (*NopHandler)(nil)
