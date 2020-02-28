gorm-hazelcast
=====================

The primary goal of the `hzgorm` project is to make it easier to cache [gorm](https://github.com/jinzhu/gorm) data results with a single line of code on Hazelcast. This module provides integration with [Hazelcast](http://github.com/hazelcast/hazelcast).

[Download Hazelcast](https://hazelcast.org/download/)

[Hazelcast Reference Manual](https://docs.hazelcast.org/docs/latest/manual/html-single/index.html)

Run hazelcast with docker:
``` 
docker run hazelcast/hazelcast:3.12.6
```

# Installation

`go get github.com/ilkerkorkut/gorm-hazelcast`

# Usage

```go
package main

import (
	"log"
	"github.com/ilkerkorkut/hazelcast-gorm"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
    "time"
)

type User struct {
	gorm.Model
	Username string
	Orders   []Order
}
type Order struct {
	gorm.Model
	UserID uint
	Price  float64
	Type   string
}

func main() {
    db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=password search_path=schema_name sslmode=disable")
    if err != nil {
        log.Println("Error while postgres connection !!!")
        return
    }

    db.AutoMigrate(&User{}, &Order{})

    hz := hzgorm.Register(db, &hzgorm.Options{
    		CacheAfterPersist: true,
    		Ttl: 120 * time.Second,
    	})
    log.Printf("Hz Instance %v", hz)
    
    orders := []Order{{
		Type: "Software",
	}}
	db.Save(&User{
		Username: "ilker",
		Orders:   orders,
	})

    var users []User
    if err := db.Table("users").Preload("Orders").Where("username = ?", "ilker").Or("username = ?", "ilker").Find(&users).Error; err != nil {
        log.Printf("Err: %v", err)
    }

    log.Printf("Result : %v", users)
}
```

#### Options
If `CacheAfterPersist` option is `true` , caches data after its persistence to db otherwise persists cache before persisting data on db. By default `CacheAfterPersist` is `true`

You can set `Ttl (Time to Live)` parameter to your options for your whole queries. By default this option is infinite.

```
hzgorm.Register(db, &hzgorm.Options{
    CacheAfterPersist: true,
    Ttl: 120 * time.Second,
})
```

#### Api
After registering `hzgorm` instance, you will be able to use following api methods.

```go
hz.EvictAll(tableName) // Evict all values in cache with its tablename
hz.EvictWithPrimaryKey() // Evict single cache entry with tablename and primarykey
hz.DisableCache(hzgorm.ReadWriteUpdate) // Disables cache with type, ReadWriteUpdate, Read, Write ,Update
hz.EnableCache(hzgorm.ReadWriteUpdate) // Enables cache with type, ReadWriteUpdate, Read, Write ,Update
hz.SetQueryTtl(120 * time.Second) // Single query based TTL
hz.Client // Reach native hazelcast-go-client
hz.Options // Change or get options dynamically
```

### Supported SQL Syntax  
  
**AND/OR:** `<expression> AND <expression> AND <expression>...`
  
- `active AND age > 30  `
- `active = false OR age = 45 OR name = 'Joe'`
- `active AND ( age > 20 OR salary < 60000 )`
- `Equality: =, !=, <, ⇐, >, >=`
  
`<expression> = value`

`age <= 30`

`name = 'Joe'`

`salary != 50000`

**BETWEEN:** `<attribute> [NOT] BETWEEN <value1> AND <value2>`
  
- `age BETWEEN 20 AND 33 (same as age >= 20 AND age ⇐ 33)`
- `age NOT BETWEEN 30 AND 40 (same as age < 30 OR age > 40)`

**IN:** `<attribute> [NOT] IN (val1, val2,...)`
  
- `age IN ( 20, 30, 40 )`
- `age NOT IN ( 60, 70 )`
- `active AND ( salary >= 50000 OR ( age NOT BETWEEN 20 AND 30 ) )`
- `age IN ( 20, 30, 40 ) AND salary BETWEEN ( 50000, 80000 )`

**LIKE:** `<attribute> [NOT] LIKE 'expression'`
  
The % (percentage sign) is the placeholder for multiple characters, an _ (underscore) is the placeholder for only one character.  
  
- `name LIKE 'Jo%' (true for 'Joe', 'Josh', 'Joseph' etc.)`
- `name LIKE 'Jo_' (true for 'Joe'; false for 'Josh')`
- `name NOT LIKE 'Jo_' (true for 'Josh'; false for 'Joe')`
- `name LIKE 'J_s%' (true for 'Josh', 'Joseph'; false 'John', 'Joe')`

**ILIKE:** `<attribute> [NOT] ILIKE 'expression'`
  
ILIKE is similar to the LIKE predicate but in a case-insensitive manner.  
  
- `name ILIKE 'Jo%' (true for 'Joe', 'joe', 'jOe','Josh','joSH', etc.)`
- `name ILIKE 'Jo_' (true for 'Joe' or 'jOE'; false for 'Josh')`

**REGEX:** `<attribute> [NOT] REGEX 'expression'`

- `name REGEX 'abc-.*' (true for 'abc-123'; false for 'abx-123')`


**Supports Preloading entities and cache them as another map.** 
