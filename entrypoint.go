package hzgorm

import (
	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/config"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type hzGorm struct {
	db      *gorm.DB
	Client  hazelcast.Client
	Options *Options
	utils   hzGormUtils
}

type Options struct {
	CacheAfterPersist     bool
	Ttl                   time.Duration
	HazelcastClientConfig *config.Config
	queryTtl              time.Duration
}

const (
	All             = "ALL"
	ReadWriteUpdate = "read-write-update"
	Read            = "read"
	Write           = "write"
	Update          = "update"
)

func Register(db *gorm.DB, options *Options) (*hzGorm, error) {
	var hz hzGorm
	hz.db = db
	if options == nil {
		// Default options
		hz.Options = &Options{
			CacheAfterPersist: true,
			Ttl:               -1,
		}
	} else {
		hz.Options = options
	}
	hz.utils = hzGormUtils{}
	if db == nil {
		log.Println("DB is nil")
		return &hz, nil
	}

	var client hazelcast.Client
	var err error

	if options != nil {
		if options.HazelcastClientConfig != nil {
			client, err = hazelcast.NewClientWithConfig(options.HazelcastClientConfig)
			if err != nil {
				log.Println(err)
				return &hz, err
			}
		} else {
			client, err = hazelcast.NewClient()
			if err != nil {
				log.Println(err)
				return &hz, err
			}
		}
	} else {
		client, err = hazelcast.NewClient()
		if err != nil {
			log.Println(err)
			return &hz, err
		}
	}

	cb := db.Callback()
	if cb.Create().Get("hzgorm:before_create") == nil {
		cb.Create().Before("gorm:before_create").Register("hzgorm:before_create", hz.hazelcastBeforeCreateCallback)
	}
	if cb.Create().Get("hzgorm:after_create") == nil {
		cb.Create().After("gorm:after_create").Register("hzgorm:after_create", hz.hazelcastAfterCreateCallback)
	}

	if cb.Update().Get("hzgorm:before_update") == nil {
		cb.Update().Before("gorm:before_update").Register("hzgorm:before_update", hz.hazelcastBeforeUpdateCallback)
	}
	if cb.Update().Get("hzgorm:after_update") == nil {
		cb.Update().After("gorm:after_update").Register("hzgorm:after_update", hz.hazelcastAfterUpdateCallback)
	}

	if cb.Query().Get("hzgorm:before_query") == nil {
		cb.Query().Before("gorm:query").Register("hzgorm:before_query", hz.hazelcastBeforeQueryCallback)
	}
	if cb.Query().Get("hzgorm:after_query") == nil {
		cb.Query().After("gorm:query").Register("hzgorm:after_query", hz.hazelcastAfterQueryCallback)
	}
	hz.Client = client
	return &hz, nil
}

func (hz *hzGorm) DB() *gorm.DB {
	return hz.db
}

func (hz *hzGorm) DisableCache(disabledType string) *hzGorm {
	if disabledType == ReadWriteUpdate || disabledType == All {
		hz.db.Callback().Create().Replace("hzgorm:before_create", voidCallback)
		hz.db.Callback().Create().Replace("hzgorm:after_create", voidCallback)
		hz.db.Callback().Update().Replace("hzgorm:before_update", voidCallback)
		hz.db.Callback().Update().Replace("hzgorm:after_update", voidCallback)
		hz.db.Callback().Query().Replace("hzgorm:before_query", voidCallback)
		hz.db.Callback().Query().Replace("hzgorm:after_query", voidCallback)
		return hz
	}
	if disabledType == Read {
		hz.db.Callback().Query().Replace("hzgorm:before_query", voidCallback)
		hz.db.Callback().Query().Replace("hzgorm:after_query", voidCallback)
		return hz
	}
	if disabledType == Write {
		hz.db.Callback().Create().Replace("hzgorm:before_create", voidCallback)
		hz.db.Callback().Create().Replace("hzgorm:after_create", voidCallback)
		return hz
	}
	if disabledType == Update {
		hz.db.Callback().Update().Replace("hzgorm:before_update", voidCallback)
		hz.db.Callback().Update().Replace("hzgorm:after_update", voidCallback)
	}
	return hz
}

func (hz *hzGorm) EnableCache(enabledType string) *hzGorm {
	if enabledType == ReadWriteUpdate || enabledType == All {
		hz.db.Callback().Create().Replace("hzgorm:before_create", hz.hazelcastBeforeCreateCallback)
		hz.db.Callback().Create().Replace("hzgorm:after_create", hz.hazelcastAfterCreateCallback)
		hz.db.Callback().Update().Replace("hzgorm:before_update", hz.hazelcastBeforeUpdateCallback)
		hz.db.Callback().Update().Replace("hzgorm:after_update", hz.hazelcastAfterUpdateCallback)
		hz.db.Callback().Query().Replace("hzgorm:before_query", hz.hazelcastBeforeQueryCallback)
		hz.db.Callback().Query().Replace("hzgorm:after_query", hz.hazelcastAfterQueryCallback)
		return hz
	}
	if enabledType == Read {
		hz.db.Callback().Query().Replace("hzgorm:before_query", hz.hazelcastBeforeQueryCallback)
		hz.db.Callback().Query().Replace("hzgorm:after_query", hz.hazelcastAfterQueryCallback)
		return hz
	}
	if enabledType == Write {
		hz.db.Callback().Create().Replace("hzgorm:before_create", hz.hazelcastBeforeCreateCallback)
		hz.db.Callback().Create().Replace("hzgorm:after_create", hz.hazelcastAfterCreateCallback)
		return hz
	}
	if enabledType == Update {
		hz.db.Callback().Update().Replace("hzgorm:before_update", hz.hazelcastBeforeUpdateCallback)
		hz.db.Callback().Update().Replace("hzgorm:after_update", hz.hazelcastAfterUpdateCallback)
		return hz
	}
	return hz
}

func (hz *hzGorm) EvictAll(tableName string) *hzGorm {
	mp, err := hz.Client.GetMap(tableName)
	if err != nil {
		log.Printf("Couldn't reach %v map.", tableName)
		return hz
	}
	err = mp.EvictAll()
	if err != nil {
		log.Printf("Couldn't evict %v map.", tableName)
		return hz
	}
	return hz
}

func (hz *hzGorm) EvictWithPrimaryKey(tableName string, key interface{}) *hzGorm {
	mp, err := hz.Client.GetMap(tableName)
	if err != nil {
		log.Printf("Couldn't reach %v map.", tableName)
		return hz
	}
	isEvicted, err := mp.Evict(key)
	if err != nil {
		log.Printf("Couldn't evict %v entity with key %v.", tableName, key)
		return hz
	}
	if !isEvicted {
		log.Printf("Couldn't evict %v entity with key %v.", tableName, key)
		return hz
	}
	return hz
}

func (hz *hzGorm) SetQueryTtl(ttl time.Duration) {
	hz.Options.queryTtl = ttl
}

func (hz *hzGorm) getQueryTtl() time.Duration {
	if hz.Options.queryTtl == 0 {
		log.Println("default ttl")
		return hz.Options.Ttl
	} else {
		log.Println("using query ttl")
		return hz.Options.queryTtl
	}
}

func (hz *hzGorm) disableCallback(callbackName string) {

	switch callbackName {
	case "hzgorm:before_create":
		hz.db.Callback().Create().Replace("hzgorm:before_create", voidCallback)
	case "hzgorm:after_create":
		hz.db.Callback().Create().Replace("hzgorm:after_create", voidCallback)
	case "hzgorm:before_update":
		hz.db.Callback().Update().Replace("hzgorm:before_update", voidCallback)
	case "hzgorm:after_update":
		hz.db.Callback().Update().Replace("hzgorm:after_update", voidCallback)
	case "hzgorm:before_query":
		hz.db.Callback().Query().Replace("hzgorm:before_query", voidCallback)
	case "hzgorm:after_query":
		hz.db.Callback().Query().Replace("hzgorm:after_query", voidCallback)
	}
}

func (hz *hzGorm) enableCallback(callbackName string) {

	switch callbackName {
	case "hzgorm:before_create":
		hz.db.Callback().Create().Replace("hzgorm:before_create", hz.hazelcastBeforeCreateCallback)
	case "hzgorm:after_create":
		hz.db.Callback().Create().Replace("hzgorm:after_create", hz.hazelcastAfterCreateCallback)
	case "hzgorm:before_update":
		hz.db.Callback().Update().Replace("hzgorm:before_update", hz.hazelcastBeforeUpdateCallback)
	case "hzgorm:after_update":
		hz.db.Callback().Update().Replace("hzgorm:after_update", hz.hazelcastAfterUpdateCallback)
	case "hzgorm:before_query":
		hz.db.Callback().Query().Replace("hzgorm:before_query", hz.hazelcastBeforeQueryCallback)
	case "hzgorm:after_query":
		hz.db.Callback().Query().Replace("hzgorm:after_query", hz.hazelcastAfterQueryCallback)
	}
}

func voidCallback(scope *gorm.Scope) {
}
