package util

import "errors"

var ErrCategoryNameIsEmpty = errors.New("field 'name' of category is empty")
var ErrPlayerFirstNameIsEmpty = errors.New("field 'first_name' of player is empty")
var ErrPlayerLastNameIsEmpty = errors.New("field 'last_name' of player is empty")
