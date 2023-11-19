package kop2cup

import (
	"errors"
	"time"
)

type TimeFormat string

const (
	RFC3339A         TimeFormat = time.RFC3339
	RFC3339ANano     TimeFormat = time.RFC3339Nano
	RFC3339B         TimeFormat = "2006-01-02T15:04:05+07:00"
	RFC3339BNano     TimeFormat = "2006-01-02T15:04:05.999999999+07:00"
	RFC3339C         TimeFormat = "2006/01/02T15:04:05+07:00"
	RFC3339CNano     TimeFormat = "2006/01/02T15:04:05.999999999+07:00"
	RFC3339Block     TimeFormat = "20060102150405+07:00"
	RFC3339BlockNano TimeFormat = "20060102150405.999999999+07:00"
	DateTime         TimeFormat = time.DateTime
	DateOnlyA        TimeFormat = time.DateOnly
	DateOnlyB        TimeFormat = "2006/01/02"
	DateOnlyBlock    TimeFormat = "20060102"
	TimeOnlyA        TimeFormat = time.TimeOnly
)

func (t *TimeFormat) String() string {
	return string(*t)
}

func StrToTimeFormat(str string) (TimeFormat, error) {
	if tfmt, ok := isValueInTimeFormat(TimeFormat(str)); ok {
		return *tfmt, nil
	}
	return TimeFormat("invalid"), errors.New("convert error: not defined")
}

func isValueInTimeFormat(tfmt TimeFormat) (*TimeFormat, bool) {
	switch tfmt {
	case RFC3339A,
		RFC3339ANano,
		RFC3339B,
		RFC3339BNano,
		RFC3339C,
		RFC3339CNano,
		RFC3339Block,
		RFC3339BlockNano,
		DateTime,
		DateOnlyA,
		DateOnlyB,
		DateOnlyBlock,
		TimeOnlyA:
		return &tfmt, true
	}
	return nil, false
}
