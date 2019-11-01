package cmd

import (
	"bufio"
	grpc_c "chatroom/protocol/chatroom"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func init() {
	var addr, protocol string
	cmd := &cobra.Command{
		Use:   "client",
		Short: "run a client",
		Run: func(cmd *cobra.Command, args []string) {
			var opts []grpc.DialOption
			switch protocol {
			case "h2c":
				opts = append(opts, grpc.WithInsecure())
			case "h2":
				opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
					InsecureSkipVerify: false,
				})))
			case "h2-skip":
				opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
					InsecureSkipVerify: true,
				})))
			default:
				log.Fatalln("not support", protocol)
			}
			c, e := grpc.Dial(addr, opts...)
			if e != nil {
				log.Fatalln(e)
			}
			// last := time.Now()
			// client := grpc_example.NewServiceClient(c)
			// _, e = client.Ping(context.Background(), &grpc_example.PingRequest{})
			// c.Close()
			// log.Println("ping", e, time.Now().Sub(last))
			chatRoomClient(c)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&addr, "addr",
		"a",
		"localhost:6000",
		"server addr",
	)
	flags.StringVarP(&protocol, "protocol",
		"p",
		"h2-skip",
		"net protocol [h2c h2 h2-skip]",
	)
	rootCmd.AddCommand(cmd)
}

// chatRoomClient 聊天室客户端
func chatRoomClient(c *grpc.ClientConn) (err error) {
	client := grpc_c.NewChatRoomServeClient(c)
	fmt.Println("please enter your name: ")
	input := bufio.NewReader(os.Stdin)
	msg, err := input.ReadString('\n')
	rs := strings.Split(msg, "\n")
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.Login(context.Background(), &grpc_c.LoginRequest{
		Uname: rs[0],
	})
	userName := rs[0]

	if resp.Result == "fail" {
		fmt.Println("please sign in")
	} else {
		fmt.Println(resp.Result)
		fmt.Println("do you want to create chatroom: y/n ?")
		input := bufio.NewReader(os.Stdin)
		msg, err := input.ReadString('\n')
		rs := strings.Split(msg, "\n")
		if err != nil {
			fmt.Println(err)
		}
		var ctx context.Context
		var md = make(map[string]string)
		var rid int
		var chat grpc_c.ChatRoomServe_ChatClient
		if rs[0] == "y" {
			roomID, _ := createRoom(client) // 创建聊天室 ...
			k := strconv.FormatInt(roomID, 10)
			md["roomID"] = k
			md["uname"] = userName
			mD := metadata.New(md)
			ctx = metadata.NewOutgoingContext(context.Background(), mD)
		} else {
			md["uname"] = userName
			mD := metadata.New(md)
			ctx = metadata.NewOutgoingContext(context.Background(), mD)
			fmt.Println("please enter room id: ")
			input = bufio.NewReader(os.Stdin)
			msg, err = input.ReadString('\n')
			rs = strings.Split(msg, "\n")
			if err != nil {
				fmt.Println(err)
			}
			rid, err = strconv.Atoi(rs[0])
			if err != nil {
				return err
			}
		}

		// 判断房间号是否存在，不存在提示用户输入正确的房间号
		if rid != 0 {
			chat, err = client.Chat(ctx)
			if err != nil {
				log.SetFlags(log.LstdFlags)
				log.Fatalln(err)
			}
			chat.Send(&grpc_c.ChatRequset{
				Cmsg: userName + " come in",
				ID:   int64(rid),
			})
			var id int64
			for {
				_, err := chat.Recv()
				if err != nil && err == io.EOF {
					fmt.Println("room id not exist, please enter correct room id: ")
					uReader := bufio.NewReader(os.Stdin)
					umsg, err := uReader.ReadString('\n')
					if err != nil {
						log.Fatalln(err)
						break
					}
					rs = strings.Split(umsg, "\n")
					id, err = strconv.ParseInt(rs[0], 10, 64)
					if err != nil {
						log.Fatalf("your enter error ! please enter correct room id format")
					}
					chat, err = client.Chat(ctx)
					if err != nil {
						log.SetFlags(log.LstdFlags)
						log.Fatalln(err)
					}
					err = chat.Send(&grpc_c.ChatRequset{
						ID:   id,
						Cmsg: userName + " come in",
					})

					if err != nil {
						fmt.Println("send fail")
						break
					}
				} else {
					break
				}
			}
		} else {
			chat, err = client.Chat(ctx)
			if err != nil {
				log.SetFlags(log.LstdFlags)
				log.Fatalln(err)
			}
		}

		go func() {
			for {
				resp1, err := chat.Recv()
				if err != nil {
					log.Fatal(err)
					break
				}
				fmt.Println(time.Now().Format("2006-1-2 15:04:05 ")+resp1.Who, ": "+resp1.Smsg)
			}
		}()

		for {
			uReader := bufio.NewReader(os.Stdin)
			umsg, err := uReader.ReadString('\n')
			if err != nil {
				log.Fatalln(err)
			}
			rs = strings.Split(umsg, "\n")
			if strings.ToUpper(rs[0]) != "Q" {
				err = chat.Send(&grpc_c.ChatRequset{
					Cmsg: rs[0],
				})
				if err != nil {
					log.Fatalln(err)
					break
				}
			} else {
				err = chat.Send(&grpc_c.ChatRequset{
					Cmsg: userName,
				})
				if err != nil {
					log.Fatalln(err)
				}
				break
			}
		}
	}
	return
}

// createRoom 创建房间
func createRoom(client grpc_c.ChatRoomServeClient) (roomID int64, err error) {
	room, err := client.Create(context.Background())
	if err != nil {
		log.Println(err)
		return roomID, err
	}
	rs, err := room.Recv()
	fmt.Println("your room id is:", rs.ID)
	return rs.ID, nil
}
