package db

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Date time.Time

func (d Date) MarshalJSON() ([]byte, error) {
	if time.Time(d).IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + time.Time(d).Format(time.DateOnly) + `"`), nil
}

func (d Date) ToPgtypeDate() pgtype.Date {
	t := time.Time(d)
	return pgtype.Date{Time: t, Valid: !t.IsZero()}
}

func (d *Date) UnmarshalText(data []byte) error {
	str := string(data)
	if str == "" {
		*d = Date(time.Time{})
		return nil
	}

	t, err := time.Parse(time.DateOnly, str)
	if err != nil {
		return err
	}

	*d = Date(t)
	return nil
}
