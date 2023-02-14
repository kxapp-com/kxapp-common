package utilz

import (
	"runtime"
	"sync"
)

/*
*
taskWork使用的一个返回结果类，里面包含是哪条数据发送了奔溃，奔溃的error数据信息
*/
type TaskWorkCrash[T any] struct {
	data  T
	error any
}

/*
任务工厂，对数据datas，每个数据需要进行task处理。分配workcount个工人（go程）进行并行工作
workCount 设置0表示使用len(datas)个工人，设置-1表示有多少个cpu内核就多少个工人，如果工人数超过datas的size，则工人数会设置为datas的size.
返回的错误数组，如果在执行任务过程中产生了奔溃错误，例如空指针错误，且在任务函数内没处理，则会在taskwork中recovery，避免整个程序奔溃，并且错误会放在返回结果中
*/
func TaskWork[T any](datas []T, task func(workId int, data T), workCount int) []*TaskWorkCrash[T] {
	var errList []*TaskWorkCrash[T]
	size := len(datas)
	if size == 0 {
		return errList
	}
	if workCount < 0 {
		workCount = runtime.NumCPU()
	}
	if workCount == 0 || workCount > size {
		workCount = size
	}
	var waiter sync.WaitGroup
	waiter.Add(size)
	queue := make(chan T)

	writer := func() {
		for i := 0; i < size; i++ {
			queue <- datas[i]
		}
		close(queue)
	}
	reader := func(id int) {
		for {
			v, ok := <-queue
			if ok {
				func() {
					defer func() {
						if err := recover(); err != nil {
							switch err.(type) {
							case nil:
							default:
								crash := &TaskWorkCrash[T]{data: v, error: err}
								errList = append(errList, crash)
							}
						}
					}()
					task(id, v)
				}()
				waiter.Done()
			} else {
				break
			}
		}
	}
	go writer()

	for i := 0; i < workCount; i++ {
		go reader(i)
	}
	waiter.Wait()
	return errList
}
