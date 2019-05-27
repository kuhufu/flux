# flux

源方法

```go
func Range(begin, end int) Flux
func Chan(c chan interface{}) Flux
func Slice(c []interface{}) Flux
func Of(c ...interface{}) Flux
```



流方法

```go
func (f Flux) Filter(filter Filter) Flux
func (f Flux) Map(m Map) Flux
func (f Flux) Parallel(concurrentNum int) Flux
```



终止方法

```go
func (f Flux) ForEach(each ForEach)
func (f Flux) ToSlice() []interface{}
func (f Flux) ToChan(args ...chan interface{}) <-chan interface{}
func (f Flux) Count() int
```



例

```go
Range(1, 10).
	Parallel(1).
	Filter(func(e interface{}) bool {
        return e.(int) % 2 == 1
	}).
	Map(func(e interface{}) interface{} {
		return e.(int) * 2
	}).
	ForEach(func(e interface{}) {
		fmt.Printf("%v ", e)
	})
```

输出结果：

```shell
2 6 10 14 18
```

因为并发度为1， 所以有序

当并发度为2时，无序

```shell
2 10 6 18 14
```

--------------------

### 注意

#### 1. 源方法

创建一个 flux

```go
Range(1, 4)	//区间[1, 4)
Slice([]interface{}{1, 2, 3})
Of(1, 2, 3)
```

#### 2. 流方法

返回值是 flux

流方法的并发度至少为1，即使你通过Parallel方法调整
##### Parallel

Parallel方法调整并发度，从下一个方法开始生效，非终止方法的并发度至少为1，即使你通过Parallel方法调整。

```go
Range(1, 10).
	Parallel(2).
	Filter(func(e interface{}) bool {
        return e.(int) % 2 == 1
	}).
	Map(func(e interface{}) interface{} {
		return e.(int) * 2
	}).
	Parallel(1).
	ForEach(func(e interface{}) {
		fmt.Printf("%v ", e)
	})
```

Filter 和 Map的并发度为2，ForEach 的并发度为1

如果想要调整某个方法的并发度，在该方法前调用 Parallel 即可

#### 3. 终止方法

返回值不是 flux，不能流式调用

> **注意 :**
>
> ToSlice 和 Count 方法是阻塞调用，将会阻塞调用方，直到所有元素转为slice或计数完成。
>
> ToChan 方法返回一个 chan，你可以从中不断读取，直到该chan关闭（表示所有元素已完成）

**ToSlice, Count, ToChan 不受并发度影响**

**ForEach 受并发度影响**

如果要使用自己的 chan 来就收结果，传入ToChan方法即可

```go
myChan := make(chan interface{}, 10)
xxx.ToChan(myChan)
```