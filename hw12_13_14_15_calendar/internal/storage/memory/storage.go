package memorystorage

import (
	"context"
	"sync"
	"time"

	appError "github.com/clawfinger/hw12_13_14_15_calendar/internal/errors"
	data "github.com/clawfinger/hw12_13_14_15_calendar/internal/event"
)

type MemoryStorage struct {
	m       sync.RWMutex
	storage map[string]*data.Event
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		storage: make(map[string]*data.Event),
	}
}

//nolint
func (s *MemoryStorage) Create(e *data.Event) error {
	s.m.Lock()
	defer s.m.Unlock()
	for _, event := range s.storage {
		if (event.Date.After(e.Date) && event.Date.Before(e.Date.Add(e.Duration))) ||
			(e.Date.After(event.Date) && e.Date.Before(event.Date.Add(event.Duration))) {
			return appError.ErrDateBusy
		}
	}
	s.storage[e.ID] = e
	return nil
}

//nolint
func (s *MemoryStorage) Update(e *data.Event) error {
	s.m.Lock()
	defer s.m.Unlock()
	_, ok := s.storage[e.ID]
	if ok {
		s.storage[e.ID] = e
	} else {
		return appError.ErrNoSuchEvent
	}
	return nil
}

//nolint
func (s *MemoryStorage) Delete(e *data.Event) error {
	s.m.Lock()
	defer s.m.Unlock()
	_, ok := s.storage[e.ID]
	if ok {
		delete(s.storage, e.ID)
	} else {
		return appError.ErrNoSuchEvent
	}
	return nil
}

func (s *MemoryStorage) getEventsBetweenDates(from time.Time, to time.Time) []*data.Event {
	result := make([]*data.Event, 0)
	s.m.RLock()
	defer s.m.RUnlock()
	for _, event := range s.storage {
		if event.Date.UTC().After(from) && event.Date.UTC().Before(to) {
			result = append(result, event)
		}
	}
	return result
}

func (s *MemoryStorage) GetEventsForDay(day time.Time) []*data.Event {
	from := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 0, 1)
	return s.getEventsBetweenDates(from, to)
}

func (s *MemoryStorage) GetEventsForWeek(weekStart time.Time) []*data.Event {
	from := time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 0, 7)
	return s.getEventsBetweenDates(from, to)
}

func (s *MemoryStorage) GetEventsForMonth(monthStart time.Time) []*data.Event {
	from := time.Date(monthStart.Year(), monthStart.Month(), monthStart.Day(), 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, 0)
	return s.getEventsBetweenDates(from, to)
}

func (s *MemoryStorage) Close(ctx context.Context) error {
	return nil
}
