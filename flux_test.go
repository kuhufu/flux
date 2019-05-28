package flux

import (
	"fmt"
	"testing"
)

func TestRangeFlux_Range(t *testing.T) {
	count := Range(1, 10).Parallel().
		Filter(func(e interface{}) bool {
			return e.(int)%2 == 1
		}).
		Map(func(e interface{}) interface{} {
			return e.(int) * 2
		}).Parallel(0).
		Count()
	fmt.Println(count)
}

func TestFlux_FilterMap(t *testing.T) {
	count := Range(1, 10).Parallel().
		FilterMap(func(e interface{}) (interface{}, bool) {
			v := e.(int)
			return v * 2, v%2 == 1
		}).
		Parallel(0).
		Count()
	fmt.Println(count)
}
