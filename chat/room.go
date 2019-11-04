package chat

import (
	"errors"
	"sync"
)

// RecvF 事件
type RecvF func(*Room, Msg)

// Recv .
type Recv interface {
	Recv(*Room, Msg)
}

// _Recv ...
type _Recv struct {
	f RecvF
}

func (r *_Recv) Recv(room *Room, m Msg) {
	r.f(room, m)
}

// NewRecv ...
func NewRecv(f RecvF) Recv {
	return &_Recv{
		f: f,
	}
}

// Room ..
type Room struct {
	msg  chan Msg
	exit chan struct{}
	fun  map[Recv]bool
	lock sync.Mutex
}

// Msg ...
type Msg struct {
	Name    string
	Content string
}

// NewRoom 创建room对象
func NewRoom() *Room {
	room := &Room{
		msg:  make(chan Msg),
		exit: make(chan struct{}),
		fun:  make(map[Recv]bool),
	}
	return room
}

//SendMsg  ...
func (r *Room) SendMsg(name, content string) (err error) {
	select {
	case r.msg <- Msg{
		Name:    name,
		Content: content,
	}:
	case <-r.exit:
		err = errors.New("exit")
		return
	}
	return
}

// Run ...
func (r *Room) Run() {
	runing := true
	for runing {
		select {
		case msg := <-r.msg:
			r.broadcast(msg)
		case <-r.exit:
			runing = false
		}
	}
}

func (r *Room) broadcast(m Msg) {
	r.lock.Lock()
	for k := range r.fun {
		k.Recv(r, m)
	}
	r.lock.Unlock()
}

// Registor 注册
func (r *Room) Registor(recv Recv) {
	r.lock.Lock()
	r.fun[recv] = true
	r.lock.Unlock()
}

// UnRegistor ..
func (r *Room) UnRegistor(recv Recv) {
	r.lock.Lock()
	delete(r.fun, recv)
	r.lock.Unlock()
}

// Close ...
func (r *Room) Close() {
	r.broadcast(Msg{
		Name: "exit",
	})
	close(r.exit)
}
