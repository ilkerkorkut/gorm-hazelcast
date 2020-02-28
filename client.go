package hzgorm

import (
	"fmt"
	"github.com/hazelcast/hazelcast-go-client/core"
	"github.com/hazelcast/hazelcast-go-client/core/predicate"
	"github.com/jinzhu/gorm"
	"log"
	"reflect"
)

func (hz *hzGorm) cachePut(scope *gorm.Scope) {
	hzMap, err := hz.Client.GetMap(scope.TableName())
	if err != nil {
		log.Printf("Couldn't get hazelcast map cache for put: %v", err.Error())
	} else {
		jsonValue, err := core.CreateHazelcastJSONValue(scope.Value)
		if err != nil {
			log.Printf("Couldn't serialize as json for hazelcast map cache: %v", err.Error())
		} else {
			primaryKey := fmt.Sprintf("%v", scope.PrimaryKeyValue())
			err := hzMap.PutTransient(primaryKey, jsonValue, hz.getQueryTtl())
			if err != nil {
				log.Printf("Couldn't put on hazelcast map cache: %v", err.Error())
			}
		}
	}
}

func (hz *hzGorm) cacheHit(scope *gorm.Scope) {
	if !scope.IndirectValue().IsValid() {
		return
	}
	obj := hz.utils.createNewStructType(scope.IndirectValue())

	var fieldNames []string
	hz.utils.structGetFieldNamesDeep(obj, &fieldNames)

	if scope.SQL == "" {
		log.Panic("THERE IS NO SQL QUERY !!!")
	}

	sqlPredicate := hz.predicateBuilder(scope.TableName(), scope.SQL, scope.SQLVars, fieldNames)

	hzMap, _ := hz.Client.GetMap(scope.TableName())
	if sqlPredicate == "" {
		values, err := hzMap.Values()
		if err != nil {
			// continues for db ops
			hz.continueForDbOperations(scope)
			return
		}
		if len(values) == 0 {
			// continues for db ops
			hz.continueForDbOperations(scope)
			return
		}
		hz.addJsonToScopeStruct(scope, values)
		return
	} else {
		values, err := hzMap.ValuesWithPredicate(predicate.SQL(sqlPredicate))
		if err != nil {
			// continues for db ops
			hz.continueForDbOperations(scope)
			return
		}
		if len(values) == 0 {
			// continues for db ops
			hz.disableCallback("hzgorm:before_query")
			hz.continueForDbOperations(scope)
			hz.enableCallback("hzgorm:before_query")
			return
		}
		hz.addJsonToScopeStruct(scope, values)
		return
	}
}

func (hz *hzGorm) continueForDbOperations(scope *gorm.Scope) {

	hz.db.Raw(scope.SQL, scope.SQLVars).Scan(scope.Value)

	switch reflect.TypeOf(scope.Value).Elem().Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(scope.Value).Elem()

		for i := 0; i < s.Len(); i++ {
			primaryKeyValue := hz.utils.determinePrimaryKeyValue(s.Index(i), scope.PrimaryKey())
			mp, _ := hz.Client.GetMap(scope.TableName())

			jsonValue, jerr := core.CreateHazelcastJSONValue(s.Index(i).Interface())
			if jerr != nil {
				log.Printf("Couldn't serialize as json for hazelcast map cache: %v", jerr.Error())
				continue
			}
			err := mp.PutTransient(primaryKeyValue, jsonValue, hz.getQueryTtl())
			if err != nil {
				log.Printf("Couldn't put value: %v", jsonValue)
			}
		}
	case reflect.Struct:
		s := reflect.ValueOf(scope.Value)
		primaryKeyValue := hz.utils.determinePrimaryKeyValue(s, scope.PrimaryKey())
		mp, _ := hz.Client.GetMap(scope.TableName())

		jsonValue, jerr := core.CreateHazelcastJSONValue(s.Interface())
		if jerr != nil {
			log.Printf("Create Hazelcast Json Value Error: %v", jerr)
			return
		}
		err := mp.PutTransient(primaryKeyValue, jsonValue, hz.getQueryTtl())
		if err != nil {
			log.Printf("Couldn't put value: %v", jsonValue)
		}
	}
}

func (hz *hzGorm) addJsonToScopeStruct(scope *gorm.Scope, values []interface{}) {
	for _, val := range values {
		cleanObject := hz.utils.createNewStructInterface(scope.IndirectValue())
		err := val.(*core.HazelcastJSONValue).Unmarshal(&cleanObject)
		if err != nil {

		} else {
			if scope.IndirectValue().Kind() == reflect.Struct {
				scope.IndirectValue().Set(reflect.ValueOf(cleanObject).Elem())
			} else if scope.IndirectValue().Kind() == reflect.Slice {
				result := reflect.Append(scope.IndirectValue(), reflect.ValueOf(cleanObject).Elem())
				scope.IndirectValue().Set(result)
			}
		}
	}
}
