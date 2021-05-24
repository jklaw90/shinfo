package main

import (
	"context"
	"log"
	"net"

	"github.com/go-redis/redis"
	"github.com/jklaw90/shinfo/internal/message"
	pb "github.com/jklaw90/shinfo/internal/pb/message"
	"github.com/jklaw90/shinfo/pkg/config"
	"github.com/jklaw90/shinfo/pkg/database"
	"github.com/jklaw90/shinfo/pkg/logging"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.NewLogger(cfg)
	ctx := logging.WithLogger(context.Background(), logger)

	logger.Info("starting message service")

	redisClient := redis.NewClient(&redis.Options{})

	session, err := database.NewCassandra(ctx, cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer session.Close()

	service := message.NewService(session, redisClient)

	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		logger.Fatal(err)
	}

	s := grpc.NewServer()
	server, err := message.NewMessageServer(service)
	if err != nil {
		logger.Fatal(err)
	}

	defer logger.Info("closing message service")

	pb.RegisterMessageServer(s, server)
	if err := s.Serve(lis); err != nil {
		logger.Fatal(err)
	}
}
