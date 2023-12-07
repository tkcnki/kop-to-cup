package kop2cup

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	tFmt "github.com/tkcnki/kop-to-cup/time_format"
)

/*
構造体中の同一名称または同一タグの項目の内容をコピーします。
  - &dest コピー先ポインタ
  - &src コピー元ポインタ
  - tfmt タグで指定されていない場合のデフォルト日付フォーマット（省略時"2006-01-02T15:04:05+09:00")
*/
func CopyFrom(dest interface{}, src interface{}, tfmt ...tFmt.TimeFormat) error {
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
		tf := []tFmt.TimeFormat{tFmt.RFC3339B}
		if len(tfmt) != 0 {
			tf[0] = tfmt[0]
		}
		srcField := srcValue.Field(i)
		if formatter := srcValue.Type().Field(i).Tag.Get("kopcup-dateformat"); formatter != "" {
			if t, err := tFmt.StrToTimeFormat(formatter); err != nil {
				return err
			} else {
				tf[0] = t
			}
		}

		go func(i int) {
			defer wg.Done()
			if destField := destValue.FieldByName(srcValue.Type().Field(i).Tag.Get("kopcup-alias")); destField.IsValid() {
				destField.Set(convertDestToSrcType(destField, srcField, tf[0]).Convert(destField.Type()))
			} else if destField := destValue.FieldByName(srcValue.Type().Field(i).Name); destField.IsValid() {
				destField.Set(convertDestToSrcType(destField, srcField, tf[0]).Convert(destField.Type()))
			}
		}(i)
	}
	wg.Wait()
	return nil
}

func convertDestToSrcType(destField reflect.Value, srcField reflect.Value, tfmt ...tFmt.TimeFormat) reflect.Value {
	if destField.Type().Kind() != srcField.Type().Kind() {
		switch destField.Type().Kind() {
		case reflect.TypeOf("").Kind():
			return reflect.ValueOf(convertToString(srcField))
		case reflect.TypeOf(1).Kind():
			return reflect.ValueOf(convertToInt(srcField))
		case reflect.TypeOf(3.14).Kind():
			return reflect.ValueOf(convertToFloat(srcField))
		case reflect.TypeOf(time.Time{}).Kind():
			return reflect.ValueOf(convertToTime(srcField, tfmt...))
		}
	}
	return srcField

}

func convertToString(srcField reflect.Value, tfmt ...tFmt.TimeFormat) string {
	switch srcField.Type().Kind() {
	case reflect.TypeOf(int(1)).Kind():
		return strconv.Itoa(srcField.Interface().(int))
	case reflect.TypeOf(time.Time{}).Kind():
		return srcField.Interface().(time.Time).Format(tfmt[0].String())
	case reflect.TypeOf(true).Kind():
		return strconv.FormatBool(srcField.Interface().(bool))
	default:
		panic(errors.New("convert error: cannot convert to string type"))
	}
}

func convertToInt(srcField reflect.Value) int {
	switch srcField.Type().Kind() {
	case reflect.TypeOf("").Kind():
		if val, err := strconv.Atoi(srcField.Interface().(string)); err != nil {
			panic(err)
		} else {
			return val
		}
	case reflect.TypeOf(true).Kind():
		if srcField.Interface().(bool) {
			return 1
		}
		return 0
	default:
		panic(errors.New("convert error: cannot convert to int type"))
	}

}

func convertToTime(srcField reflect.Value, tfmt ...tFmt.TimeFormat) time.Time {
	switch srcField.Type().Kind() {
	case reflect.TypeOf(int(1)).Kind():
		return time.Unix(int64(srcField.Interface().(int)), 0)
	case reflect.TypeOf("").Kind():
		jst, _ := time.LoadLocation("Asia/Tokyo")
		if t, err := time.ParseInLocation(tfmt[0].String(), srcField.Interface().(string), jst); err != nil {
			panic(errors.New("convert error: time perse err"))
		} else {
			return t
		}
	default:
		panic(errors.New("convert error: cannot convert to time type"))
	}
}

func convertToFloat(srcField reflect.Value) float64 {
	switch srcField.Type().Kind() {
	case reflect.TypeOf(int(1)).Kind():
		return float64(srcField.Interface().(int))
	case reflect.TypeOf("").Kind():
		if f, err := strconv.ParseFloat(srcField.Interface().(string), 64); err != nil {
			panic(errors.New("convert error"))
		} else {
			return f
		}
	case reflect.TypeOf(true).Kind():
		if srcField.Interface().(bool) {
			return 1.0
		}
		return 0.0
	default:
		panic(errors.New("convert error: cannot convert to float64 type"))
	}
}
