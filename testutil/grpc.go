package testutil

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServerRegistrationHook func(s *grpc.Server)
type Cleanup func()

func RunServicer(onServerRegistration ServerRegistrationHook) (*grpc.ClientConn, Cleanup, error) {
	s := grpc.NewServer()
	onServerRegistration(s)

	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, nil, err
	}

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
		}
	}()

	conn, err := grpc.Dial(lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		lis.Close()
		s.Stop()
		return nil, nil, err
	}

	cleanup := func() {
		conn.Close()
		s.Stop()
		lis.Close()
	}

	return conn, cleanup, nil
}
