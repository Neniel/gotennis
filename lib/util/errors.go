package util

import "errors"

var ErrCategoryNameIsEmpty = errors.New("field 'name' of category is empty")

var ErrPlayerGovernmentIDIsEmpty = errors.New("field 'governemnt_id' of player is empty")
var ErrPlayerEmailIsEmpty = errors.New("field 'email' of player is empty")
var ErrPlayerFirstNameIsEmpty = errors.New("field 'first_name' of player is empty")
var ErrPlayerLastNameIsEmpty = errors.New("field 'last_name' of player is empty")
var ErrPlayerAliasIsEmpty = errors.New("field 'alias' of player is empty")
var ErrPlayerBirthdateIsEmpty = errors.New("field 'birthdate' of player has not been set")
var ErrPlayerBirthdateIsFutureDate = errors.New("field 'birthdate' of player has not occurred yet. Is the player comming from the future? :)")

type AppError struct {
	Message string `json:"message"`
}

func (a *AppError) Error() string {
	return a.Message
}
