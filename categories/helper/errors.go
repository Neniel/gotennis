package helper

import "errors"

var ErrCategoryNameIsEmpty = errors.New("field 'name' of category is empty")

const CacheDBType = "redis"
const DBType = "mongodb"
