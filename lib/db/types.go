package db

import "time"

type Date time.Time

func (d Date) MarshalJSON() ([]byte, error) {
	if time.Time(d).IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + time.Time(d).Format(time.DateOnly) + `"`), nil
}
