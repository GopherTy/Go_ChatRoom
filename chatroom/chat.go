package chatroom

import (
	"chatroom/chat"
	"chatroom/logger"
	"chatroom/model"
	grpc_c "chatroom/protocol/chatroom"
	"context"
	"fmt"
	"log"
	"strconv"

	"google.golang.org/grpc/metadata"

	"go.uber.org/zap"
)

// ChatServer ...
type ChatServer struct {
}

var _Rooms = make(map[int64]*chat.Room)
var _Room = chat.NewRoom()
var _Name string

func init() {
	go _Room.Run()
}

// Login ...
func (ChatServer) Login(ctx context.Context, req *grpc_c.LoginRequest) (resp *grpc_c.LoginResponse, err error) {
	_Name = req.GetUname()
	user := model.User{
		Name: _Name,
	}
	ok, err := user.Check()
	if ok {
		err = _Room.SendMsg("system", _Name+"  come in")
		if err != nil {
			return
		}
		return &grpc_c.LoginResponse{
			Result: req.GetUname() + " weclome",
		}, nil
	}
	return &grpc_c.LoginResponse{
		Result: "fail",
	}, nil
}

// Create 创建聊天服务
func (ChatServer) Create(stream grpc_c.ChatRoomServe_CreateServer) (err error) {
	room := chat.NewRoom()
	_Rooms[1] = room
	go room.Run()
	err = stream.Send(&grpc_c.CreateResponse{
		ID: 1,
	})
	if err != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "create fail "); ce != nil {
			ce.Write(zap.Error(err))
		}
		return
	}
	return
}

// Chat  聊天服务
func (ChatServer) Chat(stream grpc_c.ChatRoomServe_ChatServer) (err error) {
	name := _Name
	fmt.Println("join", name)
	ch := make(chan chat.Msg, 10)
	md, ok := metadata.FromIncomingContext(stream.Context())
	var id int64
	var msg string
	if ok {
		for i, v := range md {
			x, err := strconv.Atoi(i)
			id = int64(x)
			if err != nil {
				log.Fatalln(err)
				break
			}
			if rm, ok := _Rooms[id]; ok {
				_Room = rm
			}
			msg = v[0]
		}
	}

	fmt.Println(msg)

	ctx, cancel := context.WithCancel(context.Background())
	f := chat.NewRecv(func(room *chat.Room, msg chat.Msg) {
		fmt.Println("callback", name, msg)
		select {
		case ch <- msg:
			fmt.Println("callback <-", name, msg)
		case <-ctx.Done():
		}
	})
	_Room.Registor(f)

	go func() {
		work := true
		defer cancel()
		for {
			select {
			case msg := <-ch:
				fmt.Println("msg", name, msg)
				if work {
					err := stream.Send(&grpc_c.ChatResponse{
						Smsg: msg.Content,
					})
					if err != nil {
						log.Println(err)
						cancel()
						work = false
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	var req *grpc_c.ChatRequset
	for {
		req, err = stream.Recv()
		if err != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, "recv fail "); ce != nil {
				ce.Write(zap.Error(err))
			}
			break
		}
		err = _Room.SendMsg(name, req.Cmsg)
		if err != nil {
			log.Println("send ..", err)
			break
		}
	}
	cancel()
	_Room.UnRegistor(f)
	log.Println("end ..")
	return
}
