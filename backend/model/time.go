package model

import "time"

type DateTime time.Time

func newDateTime(s string) (DateTime, error) {
	var t DateTime
	if err := t.parse(s); err != nil {
		return DateTime{}, err
	}
	return t, nil
}

func (t *DateTime) parse(s string) error {
	parsed, err := time.Parse(time.DateTime, s)
	if err != nil {
		return err
	}
	*t = DateTime(parsed)
	return nil
}

func (t *DateTime) Format() string {
	return time.Time(*t).Format(time.DateTime)
}

// DateOnly

type DateOnly time.Time

func newDateOnly(s string) (DateOnly, error) {
	var t DateOnly
	if err := t.parse(s); err != nil {
		return DateOnly{}, err
	}
	return t, nil
}

func (t *DateOnly) parse(s string) error {
	parsed, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return err
	}
	*t = DateOnly(parsed)
	return nil
}

func (t *DateOnly) Format() string {
	return time.Time(*t).Format(time.DateOnly)
}
