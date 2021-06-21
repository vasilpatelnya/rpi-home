package errors

import (
	"errors"
	"fmt"
	"github.com/vasilpatelnya/rpi-home/tool/translate"
)

func ErrorMsg(code int, err error) string {
	return fmt.Sprintf("%s: %s", translate.Txt(code), err.Error())
}

func ErrWrap(code int, err error) error {
	return errors.New(fmt.Sprintf("%s: %s", translate.Txt(code), err.Error()))
}
