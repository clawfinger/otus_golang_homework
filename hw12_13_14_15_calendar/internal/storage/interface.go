package storage

import (
	"context"
	"time"
)

// * Создать (событие);
// * Обновить (ID события, событие);
// * Удалить (ID события);
// * СписокСобытийНаДень (дата);
// * СписокСобытийНаНеделю (дата начала недели);
// * СписокСобытийНaМесяц (дата начала месяца).
type Storage interface {
	Create(e *Event) error
	Update(e *Event) error
	Delete(e *Event) error
	GetEventsForDay(day time.Time) []*Event
	GetEventsForWeek(weekStart time.Time) []*Event
	GetEventsForMonth(monthStart time.Time) []*Event
	Close(ctx context.Context) error
}
