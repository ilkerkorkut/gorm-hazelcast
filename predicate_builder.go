package hzgorm

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"regexp"
	"strconv"
	"strings"
)

const (
	asc   = "asc"
	desc  = "desc"
	limit = "LIMIT"
)

// TODO: REMOVE - GROUP BY - HAVING - JOIN etc.
func (hz *hzGorm) predicateBuilder(tableName string, sql string, sqlVars []interface{}, fieldNames []string) string {
	sql = sql + "===end==="
	predicate := hz.utils.stringBetween(sql, "(", ")===end===")
	if predicate == "" {
		return predicate
	}

	predicate = hz.predicateNormalize(predicate, fieldNames)

	for i, sv := range sqlVars {
		i = i + 1
		comma := ","
		if len(sqlVars) == i || len(sqlVars) == 1 {
			comma = ""
		}
		iStr := fmt.Sprint(i)
		iStr = "\\$" + iStr + "(\\,|\\))"
		r, _ := regexp.Compile(iStr)
		sqlVar := fmt.Sprintf("%v", sv)
		predicate = r.ReplaceAllLiteralString(predicate, sqlVar+comma)
		predicate = strings.Replace(predicate, "(", "", -1)
	}

	predicate = strings.ReplaceAll(predicate, "IN ", "IN (")
	predicate = strings.ReplaceAll(predicate, "IN ((", "IN(")

	predicate = strings.ReplaceAll(predicate, "\""+tableName+"\".", "")
	predicate = strings.ReplaceAll(predicate, ", OR", " OR")
	return predicate
}

func (hz *hzGorm) predicateNormalize(predicate string, fieldNames []string) string {
	for _, fieldName := range fieldNames {
		columnName := gorm.ToColumnName(fieldName)
		predicate = strings.ReplaceAll(predicate, columnName, fieldName)
		predicate = strings.ReplaceAll(predicate, "\""+fieldName+"\"", fieldName)
	}
	return predicate
}

func (hz *hzGorm) parseLimitAndOrder(predicate string) (string, int) {
	if strings.Contains(predicate, limit) {
		limitValue := strings.TrimSpace(hz.utils.stringAfter(predicate, limit))
		if limitValue != "" {
			lv, err := strconv.Atoi(limitValue)
			if err != nil {
				return "", -1
			}
			if strings.Contains(predicate, asc) {
				return asc, lv
			} else if strings.Contains(predicate, desc) {
				return desc, lv
			} else {
				return "", lv
			}
		} else {
			if strings.Contains(predicate, asc) {
				return asc, -1
			} else if strings.Contains(predicate, desc) {
				return desc, -1
			} else {
				return "", -1
			}
		}
	} else {
		if strings.Contains(predicate, asc) {
			return asc, -1
		} else if strings.Contains(predicate, desc) {
			return desc, -1
		} else {
			return "", -1
		}
	}

}
