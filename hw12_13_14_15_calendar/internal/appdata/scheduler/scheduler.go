package schedulerapp

import (
	"context"
	"encoding/json"
	"time"

	pb "github.com/clawfinger/hw12_13_14_15_calendar/api/generated"
	schedulerconfig "github.com/clawfinger/hw12_13_14_15_calendar/internal/config/scheduler"
	data "github.com/clawfinger/hw12_13_14_15_calendar/internal/event"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/logger"
	rabbit "github.com/clawfinger/hw12_13_14_15_calendar/internal/rabbitmq"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Scheduler struct { // TODO
	cfg            *schedulerconfig.Config
	logger         logger.Logger
	grpcClient     pb.CalendarClient
	radditProducer *rabbit.Producer
	cycleTime      time.Duration
	events         []*pb.Event
}

func New(cfg *schedulerconfig.Config, logger logger.Logger, grpcClient pb.CalendarClient,
	radditProducer *rabbit.Producer, cycleTime time.Duration) *Scheduler {
	return &Scheduler{
		cfg:            cfg,
		logger:         logger,
		grpcClient:     grpcClient,
		radditProducer: radditProducer,
		cycleTime:      cycleTime,
	}
}

func (s *Scheduler) requestEvents(ctx context.Context) {
	requestCtx, cancFunc := context.WithTimeout(ctx, time.Second)
	defer cancFunc()

	result, err := s.grpcClient.GetEventsForDay(requestCtx, timestamppb.Now(), grpc.EmptyCallOption{})
	if err != nil {
		s.logger.Error("Error on requesting events for a day: ", err.Error())
		return
	}
	s.events = result.GetEvents()
	if s.events != nil {
		s.logger.Error("Error on getting events from request result")
	}
}

func (s *Scheduler) processEvents() {
	for _, rawEvent := range s.events {
		event := data.EventFromPBData(rawEvent)
		// if time.Until(event.Date) <= event.NotifyTime {
		// 	res, err := json.Marshal(event)
		// 	if err != nil {
		// 		s.logger.Error("Failed to marshall event")
		// 		return
		// 	}
		// 	s.radditProducer.Send(string(res))
		// 	deleteCtx, cancFunc := context.WithTimeout(context.Background(), time.Second)
		// 	defer cancFunc()
		// 	// тут нужно дописать нормальную обработку ивентов для которых уже отправлены нотификации но времени не хватает
		// 	s.grpcClient.Delete(deleteCtx, rawEvent)
		// }
		res, err := json.Marshal(event)
		if err != nil {
			s.logger.Error("Failed to marshall event")
			return
		}
		s.logger.Info("Sending notification for event:\n %s", string(res))

	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	go s.radditProducer.Handle(ctx)
	done := ctx.Done()

	requestTicker := time.NewTicker(s.cycleTime)
	processTicker := time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-requestTicker.C:
				s.requestEvents(ctx)
			case <-processTicker.C:
				s.processEvents()
			}
		}
	}()
	<-ctx.Done()
	return nil
}
