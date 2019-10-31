package daemon

import (
	"chatroom/configure"
	"chatroom/logger"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func runGRPC() {
	cnf := configure.Single().GRPC

	l, e := net.Listen("tcp", cnf.Addr)
	if e != nil {
		if ce := logger.Logger.Check(zap.FatalLevel, "listen"); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		os.Exit(1)
	}

	var srv *grpc.Server
	opt := []grpc.ServerOption{}
	if cnf.H2() {
		creds, e := credentials.NewServerTLSFromFile(cnf.CertFile, cnf.KeyFile)
		if e != nil {
			if ce := logger.Logger.Check(zap.FatalLevel, "NewServerTLSFromFile"); ce != nil {
				ce.Write(
					zap.Error(e),
				)
			}
			os.Exit(1)
		}
		opt = append(opt, grpc.Creds(creds))
		srv = grpc.NewServer(
			opt...,
		)
		if logger.Logger.OutFile() {
			log.Println("h2 work at", cnf.Addr)
		}
		if ce := logger.Logger.Check(zap.InfoLevel, "h2 work"); ce != nil {
			ce.Write(
				zap.String("addr", cnf.Addr),
			)
		}
	} else {
		srv = grpc.NewServer(
			opt...,
		)
		if logger.Logger.OutFile() {
			log.Println("h2c work at", cnf.Addr)
		}
		if ce := logger.Logger.Check(zap.InfoLevel, "h2c work"); ce != nil {
			ce.Write(
				zap.String("addr", cnf.Addr),
			)
		}
	}

	registerGRPC(srv)
	reflection.Register(srv)

	go func() {
		ch := make(chan os.Signal, 2)
		signal.Notify(ch,
			os.Interrupt,
			os.Kill,
			syscall.SIGTERM)
		for {
			sig := <-ch
			switch sig {
			case os.Interrupt:
				srv.Stop()
				return
			case syscall.SIGTERM:
				srv.Stop()
				return
			}
		}
	}()
	if e := srv.Serve(l); e != nil {
		if ce := logger.Logger.Check(zap.FatalLevel, "grpc Serve"); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		os.Exit(1)
		return
	}
}
