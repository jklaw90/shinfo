package message

import (
	"context"

	"github.com/gocql/gocql"
	corepb "github.com/jklaw90/shinfo/internal/pb"
	pb "github.com/jklaw90/shinfo/internal/pb/message"
	"github.com/jklaw90/shinfo/pkg/model"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	service Service
	pb.UnimplementedMessageServer
}

func NewMessageServer(service Service) (*GrpcServer, error) {
	return &GrpcServer{
		service: service,
	}, nil
}

func (s *GrpcServer) List(ctx context.Context, req *pb.MessageListRequest) (*pb.MessageListResponse, error) {
	var resp pb.MessageListResponse

	roomID, err := gocql.ParseUUID(req.RoomID)
	if err != nil {
		return nil, err
	}

	params := ListParams{
		RoomID: roomID,
		Limit:  req.Limit,
	}

	if req.NextID != "" {
		tmp, err := gocql.ParseUUID(req.NextID)
		if err != nil {
			return nil, err
		}
		params.NextID = &tmp
	}

	messages, err := s.service.List(ctx, params)
	if err != nil {
		return nil, err
	}

	for _, m := range messages.Messages {
		resp.Messages = append(resp.Messages, messageToPBMessage(m))

	}
	if messages.NextID != nil {
		resp.NextID = messages.NextID.String()
	}

	return &resp, nil
}

func (s *GrpcServer) Add(ctx context.Context, req *pb.MessageAddRequest) (*pb.MessageAddResponse, error) {

	msg := model.MessageCreate{
		User: model.User{
			Name:   req.User.Name,
			Avatar: req.User.Avatar,
		},
		Entry:    req.Entry,
		Metadata: req.Metadata,
	}

	roomID, err := gocql.ParseUUID(req.RoomID)
	if err != nil {
		return nil, err
	}
	msg.RoomID = roomID

	userID, err := gocql.ParseUUID(req.User.Id)
	if err != nil {
		return nil, err
	}
	msg.User.ID = userID
	t := model.MessageType(req.Type)
	msg.Type = &t

	resp, err := s.service.AddMessage(ctx, msg)
	if err != nil {
		return nil, err
	}
	return &pb.MessageAddResponse{
		Message: messageToPBMessage(resp),
	}, nil
}

type GrpcClient struct {
	client pb.MessageClient
}

var _ Service = (*GrpcClient)(nil)

func NewMessageClient(address string) (*GrpcClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	service := &GrpcClient{
		client: pb.NewMessageClient(conn),
	}

	return service, nil
}

func (c *GrpcClient) List(ctx context.Context, params ListParams) (model.MessageList, error) {
	var msgList model.MessageList
	lp := &pb.MessageListRequest{
		RoomID: params.RoomID.String(),
		Limit:  params.Limit,
	}
	if params.NextID != nil {
		lp.NextID = params.NextID.String()
	}
	ml, err := c.client.List(ctx, lp)
	if err != nil {
		return msgList, err
	}

	for _, pbM := range ml.Messages {
		m, err := pbMessageToMessage(pbM)
		if err != nil {
			return msgList, err
		}
		msgList.Messages = append(msgList.Messages, m)
	}
	if ml.NextID != "" {
		tmp, err := gocql.ParseUUID(ml.NextID)
		if err != nil {
			return msgList, err
		}
		msgList.NextID = &tmp
	}

	return msgList, nil
}

func (c *GrpcClient) AddMessage(ctx context.Context, create model.MessageCreate) (model.Message, error) {
	var msg model.Message
	req := &pb.MessageAddRequest{
		RoomID: create.RoomID.String(),
		User: &corepb.User{
			Id:     create.User.ID.String(),
			Name:   create.User.Name,
			Avatar: create.User.Avatar,
		},
		Type:     int64(*create.Type),
		Entry:    create.Entry,
		Metadata: create.Metadata,
	}

	m, err := c.client.Add(ctx, req)
	if err != nil {
		return msg, err
	}

	msg, err = pbMessageToMessage(m.Message)
	if err != nil {
		return msg, err
	}

	return msg, nil
}

func messageToPBMessage(m model.Message) *corepb.Message {
	return &corepb.Message{
		Id:     m.ID.String(),
		RoomID: m.RoomID.String(),
		User: &corepb.User{
			Id:     m.User.ID.String(),
			Name:   m.User.Name,
			Avatar: m.User.Avatar,
		},
		Type:     int64(*m.Type),
		Entry:    m.Entry,
		Metadata: m.Metadata,
	}
}

func pbMessageToMessage(m *corepb.Message) (model.Message, error) {
	msg := model.Message{
		User: model.User{
			Name:   m.User.Name,
			Avatar: m.User.Avatar,
		},
		Entry:    m.Entry,
		Metadata: m.Metadata,
	}

	id, err := gocql.ParseUUID(m.Id)
	if err != nil {
		return msg, err
	}
	msg.ID = id

	roomID, err := gocql.ParseUUID(m.Id)
	if err != nil {
		return msg, err
	}
	msg.RoomID = roomID

	userID, err := gocql.ParseUUID(m.User.Id)
	if err != nil {
		return msg, err
	}
	msg.User.ID = userID
	t := model.MessageType(m.Type)
	msg.Type = &t

	return msg, nil
}
