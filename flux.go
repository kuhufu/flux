package flux

import "sync"

type Flux struct {
	c             chan interface{}
	concurrentNum int
}

type Filter func(e interface{}) bool
type ForEach func(e interface{})
type Map func(e interface{}) interface{}

func newFlux(c chan interface{}, concurrentNum int) Flux {
	return Flux{
		concurrentNum: concurrentNum,
		c:             c,
	}
}

func Range(begin, end int) Flux {
	next := make(chan interface{}, 1)
	go func() {
		for i := begin; i < end; i++ {
			next <- i
		}
		close(next)
	}()
	return newFlux(next, 1)
}

func Chan(c chan interface{}) Flux {
	next := make(chan interface{}, 1)
	go func() {
		for e := range c {
			next <- e
		}
		close(next)
	}()

	return newFlux(next, 1)
}

func Slice(c []interface{}) Flux {
	next := make(chan interface{}, 1)
	go func() {
		for _, e := range c {
			next <- e
		}
		close(next)
	}()

	return newFlux(next, 1)
}

func Of(c ...interface{}) Flux {
	next := make(chan interface{}, 1)
	go func() {
		for _, e := range c {
			next <- e
		}
		close(next)
	}()

	return newFlux(next, 1)
}

func (f Flux) Filter(filter Filter) Flux {
	num := f.concurrentNum
	if num == 0 {
		num = 1
	}
	next := make(chan interface{}, num)

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

func (f Flux) Map(m Map) Flux {
	num := f.concurrentNum
	if num == 0 {
		num = 1
	}
	next := make(chan interface{}, num)

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

//调整并发度，下一个方法生效，非终止方法的并发度至少为1，
func (f Flux) Parallel(concurrentNum int) Flux {
	if concurrentNum < 0 {
		concurrentNum = 0
	}
	f.concurrentNum = concurrentNum
	return f
}

//------------------------------------------------
//以下为终止方法

func (f Flux) ForEach(each ForEach) {
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
	for _ = range f.c {
		count += 1
	}
	return count
}
