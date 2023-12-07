package kop2cup

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	tFmt "github.com/tkcnki/kop-to-cup/time_format"
)

func TestConvertToTime(t *testing.T) {
	type TestData struct {
		SrcField   reflect.Value
		Expected   time.Time
		ShouldFail bool
	}

	// テストデータ
	intField := reflect.ValueOf(1637277000)                          // Unix timestamp: 2021-11-19 12:30:00
	stringField := reflect.ValueOf("2021-11-19T12:30:00+09:00")      // RFC3339 formatted string
	invalidTypeField := reflect.ValueOf(3.14)                        // Invalid type (float64)
	invalidValueField := reflect.ValueOf("invalidDateString")        // Invalid date string
	missingTimeFormatField := reflect.ValueOf("2021-11-19T12:30:00") // Missing time format

	testCases := []TestData{
		{SrcField: intField, Expected: time.Unix(1637277000, 0)},
		{SrcField: stringField, Expected: time.Date(2021, 11, 19, 12, 30, 0, 0, time.FixedZone("JST", 9*60*60))},
		{SrcField: invalidTypeField, ShouldFail: true},
		{SrcField: invalidValueField, ShouldFail: true},
		{SrcField: missingTimeFormatField, ShouldFail: true},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tc.ShouldFail {
						t.Errorf("Unexpected panic: %v", r)
					}
				}
			}()

			result := convertToTime(tc.SrcField, tFmt.RFC3339)

			if tc.ShouldFail {
				t.Error("Expected an error, but got none.")
				return
			}

			if !result.Equal(tc.Expected) {
				t.Errorf("Unexpected result. Got: %v, Expected: %v", result, tc.Expected)
			}
		})
	}
}

func TestConvertToInt(t *testing.T) {
	testCases := []struct {
		Input      reflect.Value
		Expected   int
		ShouldFail bool
	}{
		{Input: reflect.ValueOf("42"), Expected: 42},
		{Input: reflect.ValueOf(true), Expected: 1},
		{Input: reflect.ValueOf(false), Expected: 0},
		{Input: reflect.ValueOf(42), ShouldFail: true},        // Invalid type
		{Input: reflect.ValueOf(3.14), ShouldFail: true},      // Invalid type
		{Input: reflect.ValueOf("invalid"), ShouldFail: true}, // Invalid value
	}

	for _, tc := range testCases {
		t.Run(tc.Input.Type().String(), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tc.ShouldFail {
						t.Errorf("Unexpected panic: %v", r)
					}
				}
			}()

			result := convertToInt(tc.Input)

			if tc.ShouldFail {
				t.Error("Expected an error, but got none.")
				return
			}

			if result != tc.Expected {
				t.Errorf("Unexpected result. Got: %v, Expected: %v", result, tc.Expected)
			}
		})
	}
}

func TestConvertToString(t *testing.T) {
	testCases := []struct {
		Input       reflect.Value
		Expected    string
		TimeFormat  tFmt.TimeFormat
		ShouldFail  bool
		ShouldPanic bool
	}{
		{Input: reflect.ValueOf(42), Expected: "42"},
		{Input: reflect.ValueOf(true), Expected: "true"},
		{Input: reflect.ValueOf(false), Expected: "false"},
		{Input: reflect.ValueOf(time.Date(2021, 11, 19, 12, 30, 0, 0, time.UTC)), Expected: "2021-11-19T12:30:00Z", TimeFormat: tFmt.RFC3339},
		{Input: reflect.ValueOf("test"), ShouldPanic: true},                         // Invalid type
		{Input: reflect.ValueOf(3.14), ShouldPanic: true},                           // Invalid type
		{Input: reflect.ValueOf(time.Now()), ShouldPanic: true},                     // Missing time format
		{Input: reflect.ValueOf(struct{ Field int }{Field: 42}), ShouldPanic: true}, // Invalid type
	}

	for _, tc := range testCases {
		t.Run(tc.Input.Type().String(), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tc.ShouldPanic {
						t.Errorf("Unexpected panic: %v", r)
					}
				}
			}()

			result := convertToString(tc.Input, tc.TimeFormat)

			if tc.ShouldFail {
				t.Error("Expected an error, but got none.")
				return
			}

			if result != tc.Expected {
				t.Errorf("Unexpected result. Got: %v, Expected: %v", result, tc.Expected)
			}
		})
	}
}

func TestConvertDestToSrcType(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		DestType    reflect.Type
		SrcType     reflect.Type
		DestValue   interface{}
		SrcValue    interface{}
		TimeFormat  tFmt.TimeFormat
		Expected    interface{}
		ShouldPanic bool
	}{
		{DestType: reflect.TypeOf(""), SrcType: reflect.TypeOf("test"), DestValue: "", SrcValue: "test", Expected: "test"},
		{DestType: reflect.TypeOf(1), SrcType: reflect.TypeOf(42), DestValue: 0, SrcValue: 42, Expected: 42},
		{DestType: reflect.TypeOf(time.Time{}), SrcType: reflect.TypeOf(now), DestValue: time.Time{}, SrcValue: now, Expected: now},
		{DestType: reflect.TypeOf(1), SrcType: reflect.TypeOf("test"), DestValue: 0, SrcValue: "test", Expected: "test", ShouldPanic: true}, // Mismatched types, should return srcField
	}

	for _, tc := range testCases {
		t.Run(tc.DestType.String(), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tc.ShouldPanic {
						t.Errorf("Unexpected panic: %v", r)
					}
				}
			}()

			destField := reflect.New(tc.DestType).Elem()
			destField.Set(reflect.ValueOf(tc.DestValue))

			srcField := reflect.New(tc.SrcType).Elem()
			srcField.Set(reflect.ValueOf(tc.SrcValue))

			result := convertDestToSrcType(destField, srcField, tc.TimeFormat).Interface()

			if !reflect.DeepEqual(result, tc.Expected) {
				t.Errorf("Unexpected result. Got: %v, Expected: %v", result, tc.Expected)
			}
		})
	}
}

type SourceStruct struct {
	StringField  string    `kopcup-alias:"DestStringField"`
	IntField     int       `kopcup-alias:"DestIntField"`
	TimeField    time.Time `kopcup-dateformat:"2006-01-02T15:04:05Z"`
	StringField1 string    `kopcup-alias:"DestStringField1"`
	IntField1    int       `kopcup-alias:"DestIntField1"`
	TimeField1   time.Time `kopcup-dateformat:"2006-01-02T15:04:05Z"`
	StringField2 string    `kopcup-alias:"DestStringField2"`
	IntField2    int       `kopcup-alias:"DestIntField2"`
	TimeField2   time.Time `kopcup-dateformat:"2006-01-02T15:04:05Z"`
	StringField3 string    `kopcup-alias:"DestStringField3"`
	IntField3    int       `kopcup-alias:"DestIntField3"`
	TimeField3   time.Time `kopcup-dateformat:"2006-01-02T15:04:05Z"`
	StringField4 string    `kopcup-alias:"DestStringField4"`
	IntField4    int       `kopcup-alias:"DestIntField4"`
	TimeField4   time.Time `kopcup-dateformat:"2006-01-02T15:04:05Z"`
	StringField5 string    `kopcup-alias:"DestStringField5"`
	IntField5    int       `kopcup-alias:"DestIntField5"`
	TimeField5   time.Time `kopcup-alias:"TimeField6" kopcup-dateformat:"2006-01-02T15:04:05Z"`
	StringField6 string
	IntField6    int
	TimeField6   time.Time
}

type DestinationStruct struct {
	FloatField       float64
	DestStringField  string
	DestIntField     int
	TimeField        time.Time
	DestStringField1 string
	DestIntField1    int
	TimeField1       time.Time
	DestStringField2 string
	DestIntField2    int
	TimeField2       time.Time
	DestStringField3 string
	DestIntField3    int
	TimeField3       time.Time
	DestStringField4 string
	DestIntField4    int
	TimeField4       time.Time
	DestStringField5 string
	DestIntField5    int
	TimeField5       time.Time
	DestStringField6 string
	DestIntField6    int
	TimeField6       time.Time
}

func TestCopyFrom(t *testing.T) {
	src := SourceStruct{
		StringField:  "test",
		IntField:     42,
		TimeField:    time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
		StringField1: "test",
		IntField1:    42,
		TimeField1:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
		StringField2: "test",
		IntField2:    42,
		TimeField2:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
		StringField3: "test",
		IntField3:    42,
		TimeField3:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
		StringField4: "test",
		IntField4:    42,
		TimeField4:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
		StringField5: "test",
		IntField5:    42,
		TimeField5:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
		StringField6: "test",
		IntField6:    42,
		TimeField6:   time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC),
	}

	dest := DestinationStruct{}

	err := CopyFrom(&dest, &src)
	if err != nil {
		t.Fatalf("CopyFrom failed: %v", err)
	}

	// Validate copied values
	if dest.DestStringField != src.StringField {
		t.Errorf("Unexpected value for DestStringField. Got: %v, Expected: %v", dest.DestStringField, src.StringField)
	}

	if dest.DestIntField != src.IntField {
		t.Errorf("Unexpected value for DestIntField. Got: %v, Expected: %v", dest.DestIntField, src.IntField)
	}

	if !dest.TimeField.Equal(src.TimeField) {
		t.Errorf("Unexpected value for TimeField. Got: %v, Expected: %v", dest.TimeField, src.TimeField)
	}
}

type SourceStruct2 struct {
	FieldString string    `kopcup-alias:"AliasString" kopcup-dateformat:"2006-01-02T15:04:05"`
	FieldInt    int       `kopcup-alias:"AliasInt"`
	FieldBool   bool      `kopcup-alias:"AliasBool"`
	FieldTime   time.Time `kopcup-alias:"AliasTime" kopcup-dateformat:"2006-01-02T15:04:05"`
	Float       string
	// Add more fields as needed
}

type DestinationStruct2 struct {
	AliasString string
	AliasInt    int
	AliasBool   bool
	AliasTime   time.Time
	Float       float64
	// Add more fields as needed
}

func TestCopyFrom2(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		Name       string
		Src        interface{}
		Dest       interface{}
		Expected   interface{}
		TimeFormat tFmt.TimeFormat
		ShouldFail bool
	}{
		{
			Name:       "CopyFields",
			Src:        SourceStruct2{FieldString: "TestString", FieldInt: 42, FieldBool: true, FieldTime: now, Float: "3.14"},
			Dest:       DestinationStruct2{},
			Expected:   DestinationStruct2{AliasString: "TestString", AliasInt: 42, AliasBool: true, AliasTime: now, Float: 3.14},
			TimeFormat: tFmt.RFC3339,
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			src := tc.Src.(SourceStruct2)
			dest := tc.Dest.(DestinationStruct2)
			err := CopyFrom(&dest, &src, tc.TimeFormat)

			if tc.ShouldFail {
				if err == nil {
					t.Error("Expected an error, but got none.")
				} else {
					fmt.Println(err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(dest, tc.Expected) {
				t.Errorf("Unexpected result. Got: %+v, Expected: %+v", dest, tc.Expected)
			}
		})
	}
}

func TestConvertToFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected float64
		panics   bool // このケースでパニックが期待されているか
	}{
		{"ConvertInt", 42, 42.0, false},
		{"ConvertString", "3.14", 3.14, false},
		{"ConvertBoolTrue", true, 1.0, false},
		{"ConvertBoolFalse", false, 0.0, false},
		{"ConvertInvalidType", "invalid type", 0.0, true}, // このケースではパニックが期待されます
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				if test.panics {
					if r := recover(); r == nil {
						t.Errorf("Expected a panic, but didn't get one")
					}
				}
			}()

			result := convertToFloat(reflect.ValueOf(test.input))

			// 期待される結果と実際の結果を比較
			if result != test.expected && !test.panics {
				t.Errorf("Expected: %f, Got: %f", test.expected, result)
			}
		})
	}
}
