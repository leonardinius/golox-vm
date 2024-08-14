package vmstd

import (
	"errors"
	"fmt"
	"time"

	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
)

var errArgumentNotNumber = errors.New("formatNumber: argument is not a number")

func StdClockNative() (vmvalue.Value, error) {
	v := float64(time.Now().UnixMilli()) / 1000.0
	return vmvalue.NumberAsValue(v), nil
}

func StdFormatNumber(value vmvalue.Value) (vmvalue.Value, error) {
	if !vmvalue.IsNumber(value) {
		return vmvalue.NilValue, errArgumentNotNumber
	}

	number := vmvalue.ValueAsNumber(value)
	str := fmt.Sprintf("%#v", number)
	obj := vmvalue.StringInternCopy([]byte(str))
	return vmvalue.ObjAsValue(obj), nil
}
