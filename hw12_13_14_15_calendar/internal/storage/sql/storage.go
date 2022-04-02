package sqlstorage

import (
	"context"
	"fmt"
	"time"

	calendarconfig "github.com/clawfinger/hw12_13_14_15_calendar/internal/config/calendar"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/logger"
	"github.com/clawfinger/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib" //nolint
	"github.com/jmoiron/sqlx"
)

const (
	createQuery = `insert into public."Events"("ID", "Title", "Date", "Duration", "Description", "OwnerID",
	"NotifyTime") values(:ID, :Title, :Date, :Duration, :Description, :OwnerID, :NotifyTime)`
	updateQuery = `update public."Events" set "ID"=:ID, "Title"=:Title, "Date"=:Date,
	"Duration"=:Duration, "Description"=:Description, "OwnerID"=:OwnerID, "NotifyTime"=:NotifyTime`
	deleteQuery = `delete from public."Events" where "ID"=$1`
	selectQuery = `select * from public."Events" where "Date"<$1 and "Date">$2`
)

type Storage struct { // TODO
	db     *sqlx.DB
	cfg    *calendarconfig.Config
	logger logger.Logger
}

func NewSQLStorage(cfg *calendarconfig.Config, logger logger.Logger) *Storage {
	return &Storage{
		db:     nil,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	var err error
	connectString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		s.cfg.Data.DBData.Username, s.cfg.Data.DBData.Password, "calendar")
	s.db, err = sqlx.Open("pgx", connectString)
	if err != nil {
		return err
	}
	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFunc()
	err = s.db.PingContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) Create(e *storage.Event) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	_, err := s.db.NamedExecContext(ctx, createQuery, *e)
	return err
}

func (s *Storage) Update(e *storage.Event) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	_, err := s.db.NamedExecContext(ctx, updateQuery, *e)
	return err
}

func (s *Storage) Delete(e *storage.Event) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	_, err := s.db.ExecContext(ctx, deleteQuery, e.ID)
	return err
}

func (s *Storage) GetEventsForInterval(from time.Time, to time.Time) []*storage.Event {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	rows, err := s.db.QueryxContext(ctx, selectQuery, to, from)
	if err != nil {
		s.logger.Error("Failed to query db. Reason %s", err.Error())
	}
	res := make([]*storage.Event, 0)
	defer rows.Close()
	for rows.Next() {
		event := storage.Event{}
		err := rows.StructScan(&event)
		if err != nil {
			s.logger.Error("Failed to deserialize row. Reason %s", err.Error())
		}
		res = append(res, &event)
	}
	return res
}

func (s *Storage) GetEventsForDay(day time.Time) []*storage.Event {
	from := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 0, 1)
	return s.GetEventsForInterval(from, to)
}

func (s *Storage) GetEventsForWeek(weekStart time.Time) []*storage.Event {
	from := time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 0, 7)
	return s.GetEventsForInterval(from, to)
}

func (s *Storage) GetEventsForMonth(monthStart time.Time) []*storage.Event {
	from := time.Date(monthStart.Year(), monthStart.Month(), monthStart.Day(), 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, 0)
	return s.GetEventsForInterval(from, to)
}
