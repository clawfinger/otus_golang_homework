package data

import (
	"time"

	pb "github.com/clawfinger/hw12_13_14_15_calendar/api/generated"
	"github.com/gofrs/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Событие - основная сущность, содержит в себе поля:
// * ID - уникальный идентификатор события (можно воспользоваться UUID);
// * Заголовок - короткий текст;
// * Дата и время события;
// * Длительность события (или дата и время окончания);
// * Описание события - длинный текст, опционально;
// * ID пользователя, владельца события;
// * За сколько времени высылать уведомление, опционально.

type Event struct {
	ID          string        `db:"ID"`
	Title       string        `db:"Title"`
	Date        time.Time     `db:"Date"`
	Duration    time.Duration `db:"Duration"`
	Description string        `db:"Description"`
	OwnerID     string        `db:"OwnerID"`
	NotifyTime  time.Duration `db:"NotifyTime"`
}

func NewEvent(title string, date time.Time, duration time.Duration, owner string) (*Event, error) {
	id, err := uuid.NewV1()
	if err != nil {
		return nil, err
	}
	event := &Event{
		ID:       id.String(),
		Title:    title,
		Date:     date,
		Duration: duration,
		OwnerID:  owner,
	}
	return event, nil
}

func EventFromPBData(event *pb.Event) *Event {
	return &Event{
		ID:          event.ID,
		Title:       event.Title,
		Date:        event.Date.AsTime(),
		Duration:    time.Duration(event.Duration),
		Description: event.Description,
		OwnerID:     event.OwnerID,
		NotifyTime:  time.Duration(event.NotifyTime),
	}
}

func PBDataFromEvent(event *Event) *pb.Event {
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
