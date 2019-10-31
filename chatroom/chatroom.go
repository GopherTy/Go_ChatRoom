package chatroom

import (
	"chatroom/logger"
	grpc_c "chatroom/protocol/chatroom"
	"context"
	"fmt"

	"go.uber.org/zap"
)

// RoomService ...
type RoomService struct {
}

// 定义一个 map 来存放连接已经成功登录的用户
var user = make(map[string]grpc_c.ChatRoomServe_ChatServer)
var name string

// Login ...
func (RoomService) Login(ctx context.Context, req *grpc_c.LoginRequest) (resp *grpc_c.LoginResponse, err error) {
	name = req.GetUname()
	fmt.Println(name)
	return &grpc_c.LoginResponse{
		Result: req.GetUname() + " weclome",
	}, nil

}

// Chat ...
func (RoomService) Chat(stream grpc_c.ChatRoomServe_ChatServer) (err error) {
	user[name] = stream
	fmt.Println(user)
	for {
		req, err := stream.Recv()
		if err != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, "recv fail "); ce != nil {
				ce.Write(zap.Error(err))
			}
			break
		}
		if req.GetCmsg() != "Q" {
			go func() {
				for k, v := range user {
					v.Send(&grpc_c.ChatResponse{
						Smsg: k + " say: " + req.GetCmsg(),
					})
				}
			}()

			fmt.Println(req.GetCmsg())
		} else {
			for k, v := range user {
				if k != name {
					v.Send(&grpc_c.ChatResponse{
						Smsg: req.GetCmsg() + "leaving chatroom",
					})
				}
			}
			break
		}
	}

	return
}
