package storage

import (
	"time"

	"github.com/gofrs/uuid"
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
	ID          string
	Title       string
	Date        time.Time
	Duration    time.Duration
	Description string
	OwnerID     string
	NotifyTime  time.Duration
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
