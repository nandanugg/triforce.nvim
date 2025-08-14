package db

import "time"

type Date time.Time

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(d).Format(time.DateOnly) + `"`), nil
}
