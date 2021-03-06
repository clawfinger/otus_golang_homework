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
