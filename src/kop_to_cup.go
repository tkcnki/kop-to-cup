package kop2cup

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"
)

/*
構造体中の同一名称または同一タグの項目の内容をコピーします。
  - &dest コピー先ポインタ
  - &src コピー元ポインタ
  - tfmt タグで指定されていない場合のデフォルト日付フォーマット（省略時"2006-01-02T15:04:05+09:00")
*/
func CopyFrom(dest interface{}, src interface{}, tfmt ...TimeFormat) error {
	destValue := reflect.ValueOf(dest).Elem()
	srcValue := reflect.ValueOf(src).Elem()
	defer func() error {
		if r := recover(); r != nil {
			return errors.New(fmt.Sprint("Recovered from panic: ", r))
		}
		return nil
	}()

	var wg sync.WaitGroup
	for i := 0; i < srcValue.NumField(); i++ {
		wg.Add(1)
		tf := []TimeFormat{RFC3339B}
		if len(tfmt) != 0 {
			tf[0] = tfmt[0]
		}
		srcField := srcValue.Field(i)
		if formatter := srcValue.Type().Field(i).Tag.Get("kopcup-dateformat"); formatter != "" {
			tf[0] = TimeFormat(formatter)
		}

		go func(i int) {
			defer wg.Done()
			if destField := destValue.FieldByName(srcValue.Type().Field(i).Tag.Get("kopcup-alias")); destField.IsValid() {
				destField.Set(convertDestToSrcType(destField, srcField, tf[0]))
			} else if destField := destValue.FieldByName(srcValue.Type().Field(i).Name); destField.IsValid() {
				destField.Set(convertDestToSrcType(destField, srcField, tf[0]))
			}
		}(i)
	}
	wg.Wait()
	return nil
}

func convertDestToSrcType(destField reflect.Value, srcField reflect.Value, tfmt ...TimeFormat) reflect.Value {
	if destField.Type() != srcField.Type() {
		switch destField.Type() {
		case reflect.TypeOf(""):
			return reflect.ValueOf(convertToString(srcField))
		case reflect.TypeOf(1):
			return reflect.ValueOf(convertToInt(srcField))
		case reflect.TypeOf(time.Time{}):
			return reflect.ValueOf(convertToTime(srcField, tfmt...))
		}
	}
	return srcField

}

func convertToString(srcField reflect.Value, tfmt ...TimeFormat) string {
	switch srcField.Type() {
	case reflect.TypeOf(int(1)):
		return strconv.Itoa(srcField.Interface().(int))
	case reflect.TypeOf(time.Time{}):
		return srcField.Interface().(time.Time).Format(tfmt[0].String())
	case reflect.TypeOf(true):
		return strconv.FormatBool(srcField.Interface().(bool))
	default:
		panic(errors.New("convert error"))
	}
}

func convertToInt(srcField reflect.Value) int {
	switch srcField.Type() {
	case reflect.TypeOf(""):
		if val, err := strconv.Atoi(srcField.Interface().(string)); err != nil {
			panic(err)
		} else {
			return val
		}
	case reflect.TypeOf(true):
		src := 0
		if srcField.Interface().(bool) {
			src = 1
		}
		return src
	default:
		panic(errors.New("convert error"))
	}

}

func convertToTime(srcField reflect.Value, tfmt ...TimeFormat) time.Time {
	switch srcField.Type() {
	case reflect.TypeOf(int(1)):
		return time.Unix(int64(srcField.Interface().(int)), 0)
	case reflect.TypeOf(""):
		jst, _ := time.LoadLocation("Asia/Tokyo")
		if t, err := time.ParseInLocation(tfmt[0].String(), srcField.Interface().(string), jst); err != nil {
			panic(errors.New("convert error"))
		} else {
			return t
		}
	default:
		panic(errors.New("convert error"))
	}
}
