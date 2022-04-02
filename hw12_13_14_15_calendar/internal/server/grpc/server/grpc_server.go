package grpcserver

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/clawfinger/hw12_13_14_15_calendar/api/generated"
	servers "github.com/clawfinger/hw12_13_14_15_calendar/internal/server"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/storage"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcServer struct {
	context *servers.ServerContext
	server  *grpc.Server
	pb.UnimplementedCalendarServer
}

func NewGrpcServer(context *servers.ServerContext) *GrpcServer {
	return &GrpcServer{
		context: context,
	}
}

func (s *GrpcServer) Start() error {
	lsn, err := net.Listen("tcp", s.context.Cfg.Data.Grpc.Addr)
	if err != nil {
		log.Fatal(err)
	}

	s.server = grpc.NewServer(grpc.ChainUnaryInterceptor(LoggerInterceptor(s.context.Logger)))

	pb.RegisterCalendarServer(s.server, s)

	return s.server.Serve(lsn)
}

func (s *GrpcServer) Stop() error {
	s.server.Stop()
	return nil
}

func EventFromPBData(event *pb.Event) *storage.Event {
	return &storage.Event{
		ID:          event.ID,
		Title:       event.Title,
		Date:        event.Date.AsTime(),
		Duration:    time.Duration(event.Duration),
		Description: event.Description,
		OwnerID:     event.OwnerID,
		NotifyTime:  time.Duration(event.NotifyTime),
	}
}

func PBDataFromEvent(event *storage.Event) *pb.Event {
	return &pb.Event{
		ID:          event.ID,
		Title:       event.Title,
		Date:        timestamppb.New(event.Date),
		Duration:    uint64(event.Duration.Nanoseconds()),
		Description: event.Description,
		OwnerID:     event.OwnerID,
		NotifyTime:  uint64(event.NotifyTime.Nanoseconds()),
	}
}

func (s *GrpcServer) Create(ctx context.Context, e *pb.Event) (*pb.ModificationResult, error) {
	res := &pb.ModificationResult{}
	err := s.context.Storage.Create(EventFromPBData(e))
	if err != nil {
		res.Error = err.Error()
	} else {
		res.Error = ""
	}
	return res, err
}

func (s *GrpcServer) Update(ctx context.Context, e *pb.Event) (*pb.ModificationResult, error) {
	res := &pb.ModificationResult{}
	err := s.context.Storage.Update(EventFromPBData(e))
	if err != nil {
		res.Error = err.Error()
	} else {
		res.Error = ""
	}
	return res, err
}

func (s *GrpcServer) Delete(ctx context.Context, e *pb.Event) (*pb.ModificationResult, error) {
	res := &pb.ModificationResult{}
	err := s.context.Storage.Delete(EventFromPBData(e))
	if err != nil {
		res.Error = err.Error()
	} else {
		res.Error = ""
	}
	return res, err
}

func (s *GrpcServer) GetEventsForDay(ctx context.Context, time *timestamp.Timestamp) (*pb.RequestResult, error) {
	res := &pb.RequestResult{}
	events := s.context.Storage.GetEventsForDay(time.AsTime())
	for _, event := range events {
		res.Events = append(res.Events, PBDataFromEvent(event))
	}
	return res, nil
}

func (s *GrpcServer) GetEventsForWeek(ctx context.Context, time *timestamp.Timestamp) (*pb.RequestResult, error) {
	res := &pb.RequestResult{}
	events := s.context.Storage.GetEventsForWeek(time.AsTime())
	for _, event := range events {
		res.Events = append(res.Events, PBDataFromEvent(event))
	}
	return res, nil
}

func (s *GrpcServer) GetEventsForMonth(ctx context.Context, time *timestamp.Timestamp) (*pb.RequestResult, error) {
	res := &pb.RequestResult{}
	events := s.context.Storage.GetEventsForMonth(time.AsTime())
	for _, event := range events {
		res.Events = append(res.Events, PBDataFromEvent(event))
	}
	return res, nil
}
