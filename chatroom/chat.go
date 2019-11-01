package chatroom

import (
	"chatroom/chat"
	"chatroom/logger"
	"chatroom/model"
	grpc_c "chatroom/protocol/chatroom"
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc/metadata"

	"go.uber.org/zap"
)

// ChatServer ...
type ChatServer struct {
}

var _Rooms = make(map[int64]*chat.Room)

// var _Room = chat.NewRoom()
// var _Name string

// func init() {
// 	go _Room.Run()
// }

// Login ...
func (ChatServer) Login(ctx context.Context, req *grpc_c.LoginRequest) (resp *grpc_c.LoginResponse, err error) {
	uName := req.Uname
	user := model.User{
		Name: uName,
	}
	ok, err := user.Check()
	if ok {
		return &grpc_c.LoginResponse{
			Result: req.GetUname() + " weclome to chatroom system",
		}, nil
	}
	return &grpc_c.LoginResponse{
		Result: "fail",
	}, nil
}

// Create 创建聊天服务
func (ChatServer) Create(stream grpc_c.ChatRoomServe_CreateServer) (err error) {
	room := chat.NewRoom()
	rand.Seed(time.Now().Unix())
	rid := rand.Intn(10) + 1

	lock := sync.Mutex{}
	lock.Lock()
	_Rooms[int64(rid)] = room
	lock.Unlock()

	go room.Run()
	err = stream.Send(&grpc_c.CreateResponse{
		ID: int64(rid),
	})
	if err != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "create fail "); ce != nil {
			ce.Write(zap.Error(err))
		}
		return
	}

	for {
		_, err := stream.Recv()
		if err != nil {
			break
		}
	}
	delete(_Rooms, int64(rid))
	return
}

// Chat  聊天服务
func (ChatServer) Chat(stream grpc_c.ChatRoomServe_ChatServer) (err error) {
	// 从 context 取出数据
	var _Room *chat.Room
	var name string
	md, _ := metadata.FromIncomingContext(stream.Context())
	names := md["uname"]
	// 防止恶意攻击服务器
	if len(names) != 0 {
		name = names[0]
	} else {
		return
	}
	id, ok := md["roomid"]
	var req *grpc_c.ChatRequset
	if ok {
		if len(id) != 0 {
			x, err := strconv.ParseInt(id[0], 10, 64)
			fmt.Println("id", x)
			if err != nil {
				log.Fatalln(err)
			}
			if rm, ok := _Rooms[x]; ok {
				_Room = rm
			}
		} else {
			return
		}
	} else {
		req, err = stream.Recv()
		if err != nil {
			log.Println(err)
			return
		}
		if rm, ok := _Rooms[req.ID]; ok {
			_Room = rm
			err = _Room.SendMsg("system msg", req.Cmsg)
			if err != nil {
				return
			}
		} else {
			return
		}
	}

	ch := make(chan chat.Msg, 10)
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
						Who:  msg.Name,
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
