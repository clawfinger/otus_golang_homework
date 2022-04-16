package grpcserver

import (
	"context"
	"log"
	"net"

	pb "github.com/clawfinger/hw12_13_14_15_calendar/api/generated"
	data "github.com/clawfinger/hw12_13_14_15_calendar/internal/event"
	servers "github.com/clawfinger/hw12_13_14_15_calendar/internal/server"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
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

func (s *GrpcServer) Create(ctx context.Context, e *pb.Event) (*pb.ModificationResult, error) {
	res := &pb.ModificationResult{}
	err := s.context.Storage.Create(data.EventFromPBData(e))
	if err != nil {
		res.Error = err.Error()
	} else {
		res.Error = ""
	}
	return res, err
}

func (s *GrpcServer) Update(ctx context.Context, e *pb.Event) (*pb.ModificationResult, error) {
	res := &pb.ModificationResult{}
	err := s.context.Storage.Update(data.EventFromPBData(e))
	if err != nil {
		res.Error = err.Error()
	} else {
		res.Error = ""
	}
	return res, err
}

func (s *GrpcServer) Delete(ctx context.Context, e *pb.Event) (*pb.ModificationResult, error) {
	res := &pb.ModificationResult{}
	err := s.context.Storage.Delete(data.EventFromPBData(e))
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
		res.Events = append(res.Events, data.PBDataFromEvent(event))
	}
	return res, nil
}

func (s *GrpcServer) GetEventsForWeek(ctx context.Context, time *timestamp.Timestamp) (*pb.RequestResult, error) {
	res := &pb.RequestResult{}
	events := s.context.Storage.GetEventsForWeek(time.AsTime())
	for _, event := range events {
		res.Events = append(res.Events, data.PBDataFromEvent(event))
	}
	return res, nil
}

func (s *GrpcServer) GetEventsForMonth(ctx context.Context, time *timestamp.Timestamp) (*pb.RequestResult, error) {
	res := &pb.RequestResult{}
	events := s.context.Storage.GetEventsForMonth(time.AsTime())
	for _, event := range events {
		res.Events = append(res.Events, data.PBDataFromEvent(event))
	}
	return res, nil
}
