package types

import "time"

type Date struct {
	time.Time
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.Format("2006-01-02") + `"`), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == `""` || s == "null" {
		return nil
	}
	parsedTime, err := time.Parse(`"2006-01-02"`, s)
	if err != nil {
		return err
	}
	d.Time = parsedTime
	return nil
}
