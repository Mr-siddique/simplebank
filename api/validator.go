package api

import (
	"simplebank/db/util"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if txRequest, ok := fl.Parent().Interface().(createTransferRequest); ok {
		return util.IsSupporteCurrency(txRequest.Currency)
	}
	return false
}
