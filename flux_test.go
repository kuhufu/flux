package flux

import (
	"fmt"
	"testing"
	"time"
)

func TestRangeFlux_Range(t *testing.T) {
	Range(1, 10).Parallel().
		Filter(func(e interface{}) bool {
			return e.(int)%2 == 1
		}).
		Map(func(e interface{}) interface{} {
			return e.(int) * 2
		}).Parallel(0).
		ForEach(func(e interface{}) {
			//fmt.Println("foreach")
			fmt.Printf("%v ", e)
		})

	time.Sleep(time.Second * 3)
}
