package main

import (
	"net"
	"os"

	"github.com/go-redis/redis"
	"github.com/gocql/gocql"
	pb "github.com/jklaw90/shinfo/internal/pb/room"
	"github.com/jklaw90/shinfo/internal/room"
	"github.com/jklaw90/shinfo/pkg/logging"
	"github.com/jklaw90/shinfo/pkg/utils"
	"google.golang.org/grpc"
)

func main() {
	logger := logging.NewLogger(logging.Config{})
	addressList, ok := os.LookupEnv("CASSANDRA_ADDRESS")
	if !ok {
		logger.Fatal("missing cassandra config")
	}
	cluster := gocql.NewCluster(utils.StringToList(addressList)...)
	cluster.Keyspace = "shinfo"
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: os.Getenv("CASSANDRA_USERNAME"),
		Password: os.Getenv("CASSANDRA_PASSWORD"),
	}

	session, err := cluster.CreateSession()
	if err != nil {
		logger.Fatal(err)
	}
	defer session.Close()

	redisClient := redis.NewClient(&redis.Options{})

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

	logger.Info("starting room service")
	defer logger.Info("closing room service")

	pb.RegisterRoomServer(s, server)
	if err := s.Serve(lis); err != nil {
		logger.Fatal(err)
	}
}
