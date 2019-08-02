package sysdes

import (
	"reflect"
	"strings"

	"github.com/geekymedic/neon-cli/types"
	"gopkg.in/go-playground/validator.v9"
)

var bffValidate = validator.New()

func init() {
	err := bffValidate.RegisterValidation("nx_contains", func(fl validator.FieldLevel) bool {
		//fmt.Println("=======>", fl.Param(), fl.Field().String())
		switch fl.Field().Kind() {
		case reflect.String:
			for _, s := range strings.Split(fl.Param(), "-") {
				if s == fl.Field().String() {
					return true
				}
			}
			return false
		default:
			return false
		}
	})
	types.AssertNil(err)
}
