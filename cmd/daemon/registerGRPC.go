package daemon

import (
	"chatroom/chatroom"
	grpc_c "chatroom/protocol/chatroom"
	grpc_example "chatroom/protocol/example"
	"context"

	"google.golang.org/grpc"
)

func registerGRPC(s *grpc.Server) {
	grpc_example.RegisterServiceServer(s, _Example{})
	grpc_c.RegisterChatRoomServeServer(s, chatroom.ChatServer{})
}

type _Example struct {
}

func (_Example) Ping(ctx context.Context, request *grpc_example.PingRequest) (response *grpc_example.PingResponse, e error) {
	response = &grpc_example.PingResponse{}
	return
}
