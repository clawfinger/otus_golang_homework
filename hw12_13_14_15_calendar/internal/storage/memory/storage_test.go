package memorystorage

import (
	"errors"
	"testing"
	"time"

	appError "github.com/clawfinger/hw12_13_14_15_calendar/internal/errors"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {

	now := time.Now()

	t.Run("Add the same", func(t *testing.T) {
		testStorage := New()
		event, err := storage.NewEvent("title", time.Now(), 5*time.Minute, "owner")
		require.NoError(t, err)
		err = testStorage.Create(event)
		require.NoError(t, err)
		err = testStorage.Create(event)
		require.True(t, errors.Is(err, appError.ErrDateBusy))
	})
	t.Run("Get for day", func(t *testing.T) {
		testStorage := New()
		event, err := storage.NewEvent("title", now.Add(time.Hour*6), 5*time.Minute, "owner")
		require.NoError(t, err)
		err = testStorage.Create(event)
		require.NoError(t, err)
		dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		events := testStorage.GetEventsForDay(dayStart)
		require.True(t, len(events) == 1)
	})
	t.Run("Update", func(t *testing.T) {
		testStorage := New()
		event, err := storage.NewEvent("title", time.Now(), 5*time.Minute, "owner")
		require.NoError(t, err)
		err = testStorage.Create(event)
		require.NoError(t, err)

		event.OwnerID = "odd"
		err = testStorage.Update(event)
		require.NoError(t, err)
		dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		events := testStorage.GetEventsForDay(dayStart)
		require.True(t, len(events) == 1)
		require.Equal(t, events[0].OwnerID, "odd")
	})
	t.Run("Get for week", func(t *testing.T) {
		testStorage := New()
		event, err := storage.NewEvent("title", now, 5*time.Minute, "owner")
		require.NoError(t, err)
		err = testStorage.Create(event)
		require.NoError(t, err)
		weekStart := now
		for weekStart.Weekday() != time.Monday { // iterate back to Monday
			weekStart = weekStart.AddDate(0, 0, -1)
		}
		events := testStorage.GetEventsForWeek(weekStart)
		require.True(t, len(events) == 1)
	})
	t.Run("Get for month", func(t *testing.T) {
		testStorage := New()
		event, err := storage.NewEvent("title", now, 5*time.Minute, "owner")
		require.NoError(t, err)
		err = testStorage.Create(event)
		require.NoError(t, err)
		currentYear, currentMonth, _ := now.Date()
		currentLocation := now.Location()

		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		events := testStorage.GetEventsForMonth(firstOfMonth)
		require.True(t, len(events) == 1)
	})
}
