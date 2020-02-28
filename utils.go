package hzgorm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
	"strings"
	"unicode"
)

type hzGormUtils struct {
}

func (hzutils *hzGormUtils) stringBetween(value string, start string, end string) string {
	posFirst := strings.Index(value, start)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, end)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(start)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

func (hzutils *hzGormUtils) stringBefore(value string, c string) string {
	pos := strings.Index(value, c)
	if pos == -1 {
		return ""
	}
	return value[0:pos]
}

func (hzutils *hzGormUtils) stringAfter(value string, c string) string {
	pos := strings.LastIndex(value, c)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(c)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}

func (hzutils *hzGormUtils) stringCaseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}

func (hzutils *hzGormUtils) structGetFieldNamesDeep(t reflect.Type, fieldNames *[]string) []string {
	structMap := map[reflect.Type]string{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name
		fieldType := field.Type
		firstRune := []rune(fieldName)[0]
		if unicode.IsUpper(firstRune) && fieldType != reflect.TypeOf(gorm.Model{}) {
			*fieldNames = append(*fieldNames, fieldName)
		}
		if fieldType.Kind() == reflect.Struct {
			structMap[fieldType] = fieldName
		}
	}
	for k, _ := range structMap {
		hzutils.structGetFieldNamesDeep(k, fieldNames)
	}
	return *fieldNames
}

func (hzutils *hzGormUtils) determinePrimaryKeyValue(val reflect.Value, primaryKey string) string {
	var temp []reflect.Value
	var primaryKeyValue interface{}
	for i := 0; i < val.Type().NumField(); i++ {
		if reflect.Struct == val.Field(i).Kind() {
			temp = append(temp, val.Field(i))
		}
		f := val.Type().Field(i)
		if strings.EqualFold(f.Name, primaryKey) {
			primaryKeyValue = val.Field(i)
			break
		}
	}
	if primaryKeyValue == nil {
		for _, temp := range temp {
			if primaryKeyValue = hzutils.determinePrimaryKeyValue(temp, primaryKey); primaryKeyValue != nil {
				break
			}
		}
	}
	return fmt.Sprintf("%v", primaryKeyValue)
}

func (hzutils *hzGormUtils) createNewStructType(value reflect.Value) reflect.Type {

	switch value.Kind() {
	case reflect.Slice:
		return reflect.New(value.Type().Elem()).Elem().Type()
	case reflect.Struct:
		return reflect.New(value.Type()).Elem().Type()
	}

	return value.Type()
}

func (hzutils *hzGormUtils) createNewStructInterface(value reflect.Value) interface{} {
	switch value.Kind() {
	case reflect.Slice:
		return reflect.New(value.Type().Elem()).Interface()
	case reflect.Struct:
		return reflect.New(value.Type()).Interface()
	}
	return value.Interface()
}
