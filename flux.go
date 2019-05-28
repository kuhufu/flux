package flux

import (
	"runtime"
	"sync"
)

type Flux struct {
	c             chan interface{}
	concurrentNum int
}

type FilterFunc func(e interface{}) bool
type ForEachFunc func(e interface{})
type MapFunc func(e interface{}) interface{}
type FilterMapFunc func(e interface{}) (interface{}, bool)

func newFlux(c chan interface{}, concurrentNum int) Flux {
	return Flux{
		concurrentNum: concurrentNum,
		c:             c,
	}
}

func Range(begin, end int) Flux {
	next := make(chan interface{}, runtime.NumCPU()*2)
	go func() {
		for i := begin; i < end; i++ {
			next <- i
		}
		close(next)
	}()
	return newFlux(next, 1)
}

func Chan(c chan interface{}) Flux {
	next := make(chan interface{}, runtime.NumCPU()*2)
	go func() {
		for e := range c {
			next <- e
		}
		close(next)
	}()

	return newFlux(next, 1)
}

func Slice(c []interface{}) Flux {
	next := make(chan interface{}, runtime.NumCPU()*2)
	go func() {
		for _, e := range c {
			next <- e
		}
		close(next)
	}()

	return newFlux(next, 1)
}

func Of(c ...interface{}) Flux {
	next := make(chan interface{}, runtime.NumCPU()*2)
	go func() {
		for _, e := range c {
			next <- e
		}
		close(next)
	}()

	return newFlux(next, 1)
}

func (f Flux) Filter(filter FilterFunc) Flux {
	num := f.concurrentNum
	if num == 0 {
		num = 1
	}
	next := make(chan interface{}, num*2)

	wg := &sync.WaitGroup{}
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			for e := range f.c {
				if filter(e) {
					next <- e
				}
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(next)
	}()
	return newFlux(next, f.concurrentNum)
}

func (f Flux) Map(m MapFunc) Flux {
	num := f.concurrentNum
	if num == 0 {
		num = 1
	}
	next := make(chan interface{}, num*2)

	wg := &sync.WaitGroup{}
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			for e := range f.c {
				next <- m(e)
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(next)
	}()
	return newFlux(next, f.concurrentNum)
}

//FilterMap 同时做filter 和 map 操作
func (f Flux) FilterMap(fm FilterMapFunc) Flux {
	num := f.concurrentNum
	if num == 0 {
		num = 1
	}
	next := make(chan interface{}, num*2)

	wg := &sync.WaitGroup{}
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			for e := range f.c {
				//当 ok == true 时，接受该元素
				if v, ok := fm(e); ok {
					next <- v
				}
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(next)
	}()
	return newFlux(next, f.concurrentNum)
}

//调整并发度，下一个方法生效，非终止方法的并发度至少为1，
func (f Flux) Parallel(args ...int) Flux {
	concurrentNum := runtime.NumCPU()
	if len(args) > 0 {
		concurrentNum = args[0]
	}
	if concurrentNum < 0 {
		concurrentNum = 0
	}
	f.concurrentNum = concurrentNum
	return f
}

//------------------------------------------------
//以下为终止方法

func (f Flux) ForEach(each ForEachFunc) {
	if f.concurrentNum == 0 {
		for e := range f.c {
			each(e)
		}
	}

	for i := 0; i < f.concurrentNum; i++ {
		go func() {
			for e := range f.c {
				each(e)
			}
		}()
	}
}

func (f Flux) ToSlice() []interface{} {
	arr := make([]interface{}, 0, 16)
	for e := range f.c {
		arr = append(arr, e)
	}

	return arr
}

func (f Flux) ToChan(args ...chan interface{}) <-chan interface{} {
	var res chan interface{}
	if len(args) > 0 {
		res = args[0]
	} else {
		res = make(chan interface{}, 0)
	}
	go func() {
		for e := range f.c {
			res <- e
		}
		close(res)
	}()
	return res
}

func (f Flux) Count() int {
	count := 0
	for = range f.c {
		count += 1
	}
	return count
}
