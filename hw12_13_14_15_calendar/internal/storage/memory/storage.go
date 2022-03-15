package memorystorage

import (
	"context"
	"sync"
	"time"

	appError "github.com/clawfinger/hw12_13_14_15_calendar/internal/errors"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/storage"
)

type MemoryStorage struct {
	m       sync.RWMutex
	storage map[time.Time]*storage.Event
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		storage: make(map[time.Time]*storage.Event),
	}
}

//nolint
func (s *MemoryStorage) Create(e *storage.Event) error {
	s.m.Lock()
	defer s.m.Unlock()
	_, ok := s.storage[e.Date]
	if ok {
		return appError.ErrDateBusy
	}
	s.storage[e.Date] = e
	return nil
}

//nolint
func (s *MemoryStorage) Update(e *storage.Event) error {
	s.m.Lock()
	defer s.m.Unlock()
	_, ok := s.storage[e.Date]
	if ok {
		s.storage[e.Date] = e
	} else {
		return appError.ErrNoSuchEvent
	}
	return nil
}

//nolint
func (s *MemoryStorage) Delete(e *storage.Event) error {
	s.m.Lock()
	defer s.m.Unlock()
	_, ok := s.storage[e.Date]
	if ok {
		delete(s.storage, e.Date)
	} else {
		return appError.ErrNoSuchEvent
	}
	return nil
}

func (s *MemoryStorage) getEventsBetweenDates(from time.Time, to time.Time) []*storage.Event {
	result := make([]*storage.Event, 0)
	s.m.RLock()
	defer s.m.RUnlock()
	for _, event := range s.storage {
		if event.Date.UTC().After(from) && event.Date.UTC().Before(to) {
			result = append(result, event)
		}
	}
	return result
}

func (s *MemoryStorage) GetEventsForDay(day time.Time) []*storage.Event {
	from := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 0, 1)
	return s.getEventsBetweenDates(from, to)
}

func (s *MemoryStorage) GetEventsForWeek(weekStart time.Time) []*storage.Event {
	from := time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 0, 7)
	return s.getEventsBetweenDates(from, to)
}

func (s *MemoryStorage) GetEventsForMonth(monthStart time.Time) []*storage.Event {
	from := time.Date(monthStart.Year(), monthStart.Month(), monthStart.Day(), 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, 0)
	return s.getEventsBetweenDates(from, to)
}

func (s *MemoryStorage) Close(ctx context.Context) error {
	return nil
}
