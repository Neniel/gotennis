package util

import "errors"

var ErrCategoryNameIsEmpty = errors.New("field 'name' of category is empty")
var ErrPlayerFirstNameIsEmpty = errors.New("field 'first_name' of player is empty")
var ErrPlayerLastNameIsEmpty = errors.New("field 'last_name' of player is empty")
var ErrPlayerBirthdateIsEmpty = errors.New("field 'birthdate' of player has not been set")
var ErrPlayerBirthdateIsFutureDate = errors.New("field 'birthdate' of player has not occurred yet. Is the player comming from the future? :)")
