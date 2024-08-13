package vmstd

import (
	"time"

	"github.com/leonardinius/goloxvm/internal/vm/vmvalue"
)

func StdClockNative(args ...vmvalue.Value) (vmvalue.Value, error) {
	v := float64(time.Now().UnixMilli()) / 1000.0
	return vmvalue.NumberAsValue(v), nil
}
