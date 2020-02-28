package hzgorm

import (
	"errors"
	"github.com/jinzhu/gorm"
)

func (hz *hzGorm) hazelcastBeforeCreateCallback(scope *gorm.Scope) {
	if !hz.Options.CacheAfterPersist {
		defer hz.removeQueryTtl()
		hz.cachePut(scope)
	}
}

func (hz *hzGorm) hazelcastAfterCreateCallback(scope *gorm.Scope) {
	if hz.Options.CacheAfterPersist {
		defer hz.removeQueryTtl()
		hz.cachePut(scope)
	}
}

func (hz *hzGorm) hazelcastBeforeUpdateCallback(scope *gorm.Scope) {
	if !hz.Options.CacheAfterPersist {
		defer hz.removeQueryTtl()
		hz.cachePut(scope)
	}
}

func (hz *hzGorm) hazelcastAfterUpdateCallback(scope *gorm.Scope) {
	if hz.Options.CacheAfterPersist {
		defer hz.removeQueryTtl()
		hz.cachePut(scope)
	}
}

func (hz *hzGorm) hazelcastBeforeQueryCallback(scope *gorm.Scope) {
	scope.Err(errors.New("hazelcast cache"))
}

func (hz *hzGorm) hazelcastAfterQueryCallback(scope *gorm.Scope) {
	hz.cacheHit(scope)
	scope.Err(nil)
	scope.DB().Error = nil
}

func (hz *hzGorm) removeQueryTtl() {
	hz.Options.queryTtl = 0
}
