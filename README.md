gorm-hazelcast
=====================

The primary goal of the `hzgorm` project is to make it easier to cache [gorm](https://github.com/jinzhu/gorm) data results with a single line of code on Hazelcast. This module provides integration with [Hazelcast](http://github.com/hazelcast/hazelcast).

![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/ilkerkorkut/gorm-hazelcast) ![GitHub last commit](https://img.shields.io/github/last-commit/ilkerkorkut/gorm-hazelcast) ![GitHub](https://img.shields.io/github/license/ilkerkorkut/gorm-hazelcast) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ilkerkorkut/gorm-hazelcast)

[Download Hazelcast](https://hazelcast.org/imdg/download/archives/#hazelcast-imdg-3-12-6)

[Hazelcast Reference Manual](https://docs.hazelcast.org/docs/3.12.6/manual/html-single/index.html)


# Requirements
Hazelcast IMDG 3.6 or newer _(currently not supported 4.0)_

- Run hazelcast with docker:
``` 
docker run hazelcast/hazelcast:3.12.6
```

# Installation

`go get github.com/ilkerkorkut/gorm-hazelcast`

# Usage

Use gorm-hazelcast plugin with a single-line cache registration on gorm. To initialize with Options parameters, please look at [Options section]((https://github.com/ilkerkorkut/gorm-hazelcast#options)).

There are different types of usages in [examples](https://github.com/ilkerkorkut/gorm-hazelcast/tree/master/examples) directory. Look at [Api section](https://github.com/ilkerkorkut/gorm-hazelcast#api) for additional programmatic usage.

```go
import hzgorm "github.com/ilkerkorkut/gorm-hazelcast"

hzgorm.Register(db, nil)
```

# Options
gorm-hazelcast provides three type option parameters during initializing.

`CacheAfterPersist` : If `CacheAfterPersist` option is `true` , caches data after its persistence to db otherwise persists cache before persisting data on db. By default `CacheAfterPersist` is `true`

`Ttl`: You can set `Ttl (Time to Live)` parameter to your options for your whole queries. By default this option is infinite. Look at [Api section](https://github.com/ilkerkorkut/gorm-hazelcast#api) for query based ttl.

`HazelcastClientConfig`: In Options `HazelcastClientConfig` field. You are able to set your custom Hazelcast client configuration. 

```go
hzgorm.Register(db, &hzgorm.Options{
    CacheAfterPersist: true,
    Ttl: 120 * time.Second,
})
```

# Api
After registering `hzgorm` instance, you will be able to use following api methods.

```go
hz.EvictAll(tableName) // Evict all values in cache with its tablename
hz.EvictWithPrimaryKey(tableName, primaryKey) // Evict single cache entry with tablename and primarykey
hz.DisableCache(hzgorm.ReadWriteUpdate) // Disables cache with type, ReadWriteUpdate, Read, Write ,Update
hz.EnableCache(hzgorm.ReadWriteUpdate) // Enables cache with type, ReadWriteUpdate, Read, Write ,Update
hz.SetQueryTtl(120 * time.Second) // Single query based TTL
hz.Client // Reach native hazelcast-go-client
hz.Options // Change or get options dynamically (dynamic options changes are not recommended)
```

## Supported SQL Syntax for Hazelcast Cache 
  
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

**Data orders are not guaranteed.**


# Contributing

Feel free to contribute, or [create an issue](https://github.com/ilkerkorkut/gorm-hazelcast/issues) if you found a bug or need a feature request.
