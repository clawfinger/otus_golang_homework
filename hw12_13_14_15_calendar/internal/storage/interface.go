package storage

import (
	"context"
	"time"

	data "github.com/clawfinger/hw12_13_14_15_calendar/internal/event"
)

// * Создать (событие);
// * Обновить (ID события, событие);
// * Удалить (ID события);
// * СписокСобытийНаДень (дата);
// * СписокСобытийНаНеделю (дата начала недели);
// * СписокСобытийНaМесяц (дата начала месяца).
type Storage interface {
	Create(e *data.Event) error
	Update(e *data.Event) error
	Delete(e *data.Event) error
	GetEventsForDay(day time.Time) []*data.Event
	GetEventsForWeek(weekStart time.Time) []*data.Event
	GetEventsForMonth(monthStart time.Time) []*data.Event
	Close(ctx context.Context) error
}
