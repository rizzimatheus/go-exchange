package api

import (
	"go-exchange/util"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}

var validPair validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if pair, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedPair(pair)
	}
	return false
}
