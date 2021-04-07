package room

import (
	"context"
	"time"

	"github.com/gocql/gocql"
	corepb "github.com/jklaw90/shinfo/internal/pb"
	pb "github.com/jklaw90/shinfo/internal/pb/room"
	"github.com/jklaw90/shinfo/pkg/model"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	service Service
	pb.UnimplementedRoomServer
}

func NewRoomServer(service Service) (*GrpcServer, error) {
	return &GrpcServer{
		service: service,
	}, nil
}

func (s *GrpcServer) Get(ctx context.Context, req *pb.RoomGetRequest) (*pb.RoomGetResponse, error) {
	roomID, err := gocql.ParseUUID(req.RoomID)
	if err != nil {
		return nil, err
	}

	resp, err := s.service.Get(ctx, roomID)
	if err != nil {
		return nil, err
	}

	return &pb.RoomGetResponse{
		Room: mapToPBRoom(resp),
	}, nil
}

func (s *GrpcServer) GetByUserID(ctx context.Context, req *pb.RoomGetByUserIDRequest) (*pb.RoomGetByUserIDResponse, error) {
	userID, err := gocql.ParseUUID(req.UserID)
	if err != nil {
		return nil, err
	}

	var nextID *gocql.UUID
	if req.NextID != "" {
		tmp, err := gocql.ParseUUID(req.NextID)
		if err != nil {
			return nil, err
		}
		nextID = &tmp
	}
	userRooms, err := s.service.GetByUserID(ctx, userID, req.Limit, nextID)
	if err != nil {
		return nil, err
	}
	resp := &pb.RoomGetByUserIDResponse{}
	for _, room := range userRooms.Rooms {
		resp.Rooms = append(resp.Rooms, mapUserRoomToPBRoom(room))
	}
	if userRooms.NextID != nil {
		resp.NextID = userRooms.NextID.String()
	}

	return resp, nil
}

func (s *GrpcServer) GetUsers(ctx context.Context, req *pb.RoomGetUsersRequest) (*pb.RoomGetUsersResponse, error) {
	roomID, err := gocql.ParseUUID(req.RoomID)
	if err != nil {
		return nil, err
	}

	var lastSeenID *gocql.UUID
	if req.NextID != "" {
		tmp, err := gocql.ParseUUID(req.NextID)
		if err != nil {
			return nil, err
		}
		lastSeenID = &tmp
	}

	roomUsers, err := s.service.GetUsers(ctx, roomID, req.Limit, lastSeenID)
	if err != nil {
		return nil, err
	}

	resp := &pb.RoomGetUsersResponse{}
	for _, ru := range roomUsers.Users {
		resp.Users = append(resp.Users, &corepb.RoomUser{
			User:     mapToPBUser(ru.User),
			JoinedAt: ru.JoinedAt.String(),
		})
	}

	if roomUsers.NextID != nil {
		resp.NextID = roomUsers.NextID.String()
	}

	return resp, nil
}

func (s *GrpcServer) Create(ctx context.Context, req *pb.RoomCreateRequest) (*pb.RoomCreateResponse, error) {
	r, err := s.service.Create(ctx, model.RoomCreate{
		Name:     req.Name,
		Type:     (*model.RoomType)(&req.Type),
		Public:   &req.Public,
		Archived: &req.Public,
	})
	if err != nil {
		return nil, err
	}

	return &pb.RoomCreateResponse{
		Room: mapToPBRoom(r),
	}, nil
}

func (s *GrpcServer) AddUser(ctx context.Context, req *pb.RoomAddUserRequest) (*pb.RoomAddUserResponse, error) {
	roomID, err := gocql.ParseUUID(req.RoomID)
	if err != nil {
		return nil, err
	}

	user, err := mapToUser(req.User)
	if err != nil {
		return nil, err
	}

	if err = s.service.AddUser(ctx, roomID, user); err != nil {
		return nil, err
	}
	return &pb.RoomAddUserResponse{}, err
}

type GrpcClient struct {
	client pb.RoomClient
}

var _ Service = (*GrpcClient)(nil)

func NewRoomClient(address string) (*GrpcClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	service := &GrpcClient{
		client: pb.NewRoomClient(conn),
	}

	return service, nil
}

func (s *GrpcClient) Get(ctx context.Context, roomID gocql.UUID) (model.Room, error) {
	var r model.Room
	resp, err := s.client.Get(ctx, &pb.RoomGetRequest{
		RoomID: roomID.String(),
	})
	if err != nil {
		return r, err
	}
	return mapToRoom(*resp.Room)
}

func (s *GrpcClient) GetByUserID(ctx context.Context, userID gocql.UUID, limit int64, lastSeen *gocql.UUID) (model.UserRooms, error) {
	var userRooms model.UserRooms

	req := &pb.RoomGetByUserIDRequest{
		UserID: userID.String(),
		Limit:  limit,
	}
	if lastSeen != nil {
		req.NextID = lastSeen.String()
	}

	resp, err := s.client.GetByUserID(ctx, req)
	if err != nil {
		return userRooms, err
	}

	for _, ru := range resp.Rooms {
		roomUser, err := mapPBRoomUserToRoomUser(*ru)
		if err != nil {
			return userRooms, err
		}
		userRooms.Rooms = append(userRooms.Rooms, roomUser)
	}
	if resp.NextID != "" {
		tmp, _ := gocql.ParseUUID(resp.NextID)
		userRooms.NextID = &tmp
	}

	return userRooms, nil
}

func (s *GrpcClient) GetUsers(ctx context.Context, roomID gocql.UUID, limit int64, lastSeen *gocql.UUID) (model.RoomUsers, error) {
	var roomUsers model.RoomUsers
	req := &pb.RoomGetUsersRequest{
		RoomID: roomID.String(),
		Limit:  limit,
	}
	if lastSeen != nil {
		req.NextID = lastSeen.String()
	}

	resp, err := s.client.GetUsers(ctx, req)
	if err != nil {
		return roomUsers, err
	}

	for _, ru := range resp.Users {
		u, err := mapToUser(ru.User)
		if err != nil {
			return roomUsers, err
		}
		joinedAt, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", ru.JoinedAt)
		if err != nil {
			return roomUsers, err
		}

		roomUsers.Users = append(roomUsers.Users, model.RoomUser{
			User:     u,
			JoinedAt: joinedAt,
		})
	}

	if resp.NextID != "" {
		tmp, _ := gocql.ParseUUID(resp.NextID)
		roomUsers.NextID = &tmp
	}

	return roomUsers, nil
}

func (s *GrpcClient) Create(ctx context.Context, room model.RoomCreate) (model.Room, error) {
	resp, err := s.client.Create(ctx, &pb.RoomCreateRequest{
		Name:     room.Name,
		Type:     int64(*room.Type),
		Public:   *room.Public,
		Archived: *room.Archived,
	})
	if err != nil {
		return model.Room{}, err
	}
	return mapToRoom(*resp.Room)
}

func (s *GrpcClient) AddUser(ctx context.Context, roomID gocql.UUID, user model.User) error {
	if _, err := s.client.AddUser(ctx, &pb.RoomAddUserRequest{
		RoomID: roomID.String(),
		User:   mapToPBUser(user),
	}); err != nil {
		return err
	}
	return nil
}

func (s *GrpcClient) OnMessage(ctx context.Context, msg model.Message) error {
	return nil
}

func mapToRoom(r corepb.Room) (model.Room, error) {
	var room model.Room
	var err error

	room.ID, err = gocql.ParseUUID(r.Id)
	if err != nil {
		return room, err
	}
	room.Name = r.Name
	t := model.RoomType(r.Type)
	room.Type = &t
	room.Public = &r.Public
	room.Archived = &r.Archived
	room.CreatedAt, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", r.CreatedAt)
	if err != nil {
		return room, err
	}

	return room, nil
}

func mapToPBRoom(r model.Room) *corepb.Room {
	return &corepb.Room{
		Id:        r.ID.String(),
		Name:      r.Name,
		Type:      int64(*r.Type),
		Public:    *r.Public,
		Archived:  *r.Archived,
		CreatedAt: r.CreatedAt.String(),
	}
}

func mapToUser(u *corepb.User) (model.User, error) {
	var user model.User
	id, err := gocql.ParseUUID(u.Id)
	if err != nil {
		return user, err
	}

	user = model.User{
		ID:     id,
		Name:   u.Name,
		Avatar: u.Avatar,
	}

	return user, nil
}

func mapToPBUser(u model.User) *corepb.User {
	user := &corepb.User{
		Id:     u.ID.String(),
		Name:   u.Name,
		Avatar: u.Avatar,
	}
	return user
}

func mapUserRoomToPBRoom(r model.RoomByUser) *corepb.UserRoom {
	return &corepb.UserRoom{
		Room: &corepb.Room{
			Id:        r.RoomID.String(),
			Name:      r.RoomName,
			Type:      int64(*r.Type),
			Public:    *r.Public,
			Archived:  *r.Archived,
			CreatedAt: r.CreatedAt.String(),
		},
		JoinedAt: r.JoinedAt.String(),
	}
}

func mapPBRoomUserToRoomUser(r corepb.UserRoom) (model.RoomByUser, error) {
	ru := model.RoomByUser{
		RoomName: r.Room.Name,
	}
	tmp, err := gocql.ParseUUID(r.Room.Id)
	if err != nil {
		return ru, err
	}
	ru.RoomID = tmp
	ru.RoomName = r.Room.Name
	t := model.RoomType(r.Room.Type)
	ru.Type = &t
	ru.Public = &r.Room.Public
	ru.Archived = &r.Room.Archived
	ru.CreatedAt, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", r.Room.CreatedAt)
	if err != nil {
		return ru, err
	}

	ru.JoinedAt, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", r.JoinedAt)
	if err != nil {
		return ru, err
	}

	return ru, nil
}
