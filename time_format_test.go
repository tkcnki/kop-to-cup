package kop2cup

import (
	"testing"
	"time"
)

func TestTimeFormatString(t *testing.T) {
	testCases := []struct {
		Input    TimeFormat
		Expected string
	}{
		{Input: RFC3339, Expected: time.RFC3339},
		{Input: RFC3339Nano, Expected: time.RFC3339Nano},
		{Input: RFC3339, Expected: "2006-01-02T15:04:05Z07:00"},
		{Input: RFC3339Nano, Expected: "2006-01-02T15:04:05.999999999Z07:00"},
		{Input: RFC3339B, Expected: "2006/01/02T15:04:05Z07:00"},
		{Input: RFC3339BNano, Expected: "2006/01/02T15:04:05.999999999Z07:00"},
		{Input: RFC3339Block, Expected: "20060102150405Z07:00"},
		{Input: RFC3339BlockNano, Expected: "20060102150405.999999999Z07:00"},
		{Input: DateTime, Expected: time.DateTime},
		{Input: DateOnlyA, Expected: time.DateOnly},
		{Input: DateOnlyB, Expected: "2006/01/02"},
		{Input: DateOnlyBlock, Expected: "20060102"},
		{Input: TimeOnlyA, Expected: time.TimeOnly},
		{Input: TimeFormat("invalid"), Expected: "invalid"},
	}

	for _, tc := range testCases {
		t.Run(tc.Input.String(), func(t *testing.T) {
			result := tc.Input.String()

			if result != tc.Expected {
				t.Errorf("Unexpected result. Got: %v, Expected: %v", result, tc.Expected)
			}
		})
	}
}

func TestStrToTimeFormat(t *testing.T) {
	testCases := []struct {
		Input      string
		Expected   TimeFormat
		ShouldFail bool
	}{
		{Input: "2006-01-02T15:04:05Z07:00", Expected: RFC3339},
		{Input: "2006-01-02T15:04:05.999999999Z07:00", Expected: RFC3339Nano},
		{Input: "2006-01-02T15:04:05Z07:00", Expected: RFC3339},
		{Input: "2006-01-02T15:04:05.999999999Z07:00", Expected: RFC3339Nano},
		{Input: "2006/01/02T15:04:05Z07:00", Expected: RFC3339B},
		{Input: "2006/01/02T15:04:05.999999999Z07:00", Expected: RFC3339BNano},
		{Input: "20060102150405Z07:00", Expected: RFC3339Block},
		{Input: "20060102150405.999999999Z07:00", Expected: RFC3339BlockNano},
		{Input: "15:04:05", Expected: TimeOnlyA},
		{Input: "2006/01/02", Expected: DateOnlyB},
		{Input: "20060102", Expected: DateOnlyBlock},
		{Input: "2006-01-02 15:04:05", Expected: DateTime},
		{Input: "invalid", ShouldFail: true},
	}

	for _, tc := range testCases {
		t.Run(tc.Input, func(t *testing.T) {
			result, err := StrToTimeFormat(tc.Input)

			if tc.ShouldFail {
				if err == nil {
					t.Error("Expected an error, but got none.")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if result != tc.Expected {
					t.Errorf("Unexpected result. Got: %v, Expected: %v", result, tc.Expected)
				}
			}
		})
	}
}

func TestIsValueInTimeFormat(t *testing.T) {
	testCases := []struct {
		Input    TimeFormat
		Expected bool
	}{
		{Input: RFC3339, Expected: true},
		{Input: RFC3339Nano, Expected: true},
		{Input: RFC3339, Expected: true},
		{Input: RFC3339Nano, Expected: true},
		{Input: RFC3339B, Expected: true},
		{Input: RFC3339BNano, Expected: true},
		{Input: RFC3339Block, Expected: true},
		{Input: RFC3339BlockNano, Expected: true},
		{Input: DateTime, Expected: true},
		{Input: DateOnlyA, Expected: true},
		{Input: DateOnlyB, Expected: true},
		{Input: DateOnlyBlock, Expected: true},
		{Input: TimeOnlyA, Expected: true},
		{Input: TimeFormat("invalid"), Expected: false},
	}

	for _, tc := range testCases {
		t.Run(tc.Input.String(), func(t *testing.T) {
			result, ok := isValueInTimeFormat(tc.Input)

			if ok != tc.Expected {
				t.Errorf("Unexpected result. Got: %v, Expected: %v", ok, tc.Expected)
			}

			if tc.Expected && *result != tc.Input {
				t.Errorf("Unexpected result. Got: %v, Expected: %v", *result, tc.Input)
			}
		})
	}
}
