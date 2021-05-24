package main

import (
	"context"
	"log"
	"net"

	"github.com/go-redis/redis"
	pb "github.com/jklaw90/shinfo/internal/pb/room"
	"github.com/jklaw90/shinfo/internal/room"
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

	logger.Info("starting room service")

	redisClient := redis.NewClient(&redis.Options{})

	session, err := database.NewCassandra(ctx, cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer session.Close()

	service := room.NewService(session, redisClient)

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		logger.Fatal(err)
	}

	s := grpc.NewServer()
	server, err := room.NewRoomServer(service)
	if err != nil {
		logger.Fatal(err)
	}

	defer logger.Info("closing room service")

	pb.RegisterRoomServer(s, server)
	if err := s.Serve(lis); err != nil {
		logger.Fatal(err)
	}
}
