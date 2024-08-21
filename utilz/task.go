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

func MultiThreadTask1[T any](slice []T, numThreads int, taskFunc func(data T)) {
	if len(slice) <= 0 {
		return
	}
	if numThreads < 0 {
		numThreads = runtime.NumCPU()
	}
	if numThreads == 0 || numThreads > len(slice) {
		numThreads = len(slice)
	}
	// Calculate the number of elements per thread
	numElementsPerThread := len(slice) / numThreads
	// Create a wait group to wait for all threads to finish
	var wg sync.WaitGroup
	// Iterate over the number of threads
	for i := 0; i < numThreads; i++ {
		// Increment the wait group
		wg.Add(1)
		// Calculate the start and end indices for the current thread
		start := i * numElementsPerThread
		end := (i + 1) * numElementsPerThread
		// For the last thread, include any remaining elements
		if i == numThreads-1 {
			end = len(slice)
		}
		// Launch a goroutine to process the elements for the current thread
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				taskFunc(slice[j])
			}
		}(start, end)
	}
	// Wait for all threads to finish
	wg.Wait()
}

// MultiThreadTask2 执行多任务处理
func MultiThreadTask[T any](slice []T, numThreads int, taskFunc func(sliceIndex int, data T) error) []error {
	if len(slice) == 0 {
		return nil
	}

	if numThreads <= 0 {
		numThreads = runtime.NumCPU()
	}
	if numThreads > len(slice) {
		numThreads = len(slice)
	}

	var errors []error
	var mu sync.Mutex
	var wg sync.WaitGroup

	chunkSize := (len(slice) + numThreads - 1) / numThreads
	errorChan := make(chan error, len(slice))

	for i := 0; i < numThreads; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(slice) {
			end = len(slice)
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				if err := taskFunc(j, slice[j]); err != nil {
					errorChan <- err
				}
			}
		}(start, end)
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	for err := range errorChan {
		mu.Lock()
		errors = append(errors, err)
		mu.Unlock()
	}

	return errors
}

/*
任务工厂，对数据datas，每个数据需要进行task处理。分配workcount个工人（go程）进行并行工作
workCount 设置0表示使用len(datas)个工人，设置-1表示有多少个cpu内核就多少个工人，如果工人数超过datas的size，则工人数会设置为datas的size.
返回的错误数组，如果在执行任务过程中产生了奔溃错误，例如空指针错误，且在任务函数内没处理，则会在taskwork中recovery，避免整个程序奔溃，并且错误会放在返回结果中
*/
//func TaskWork[T any](datas []T, task func(workId int, data T), workCount int) []*TaskWorkCrash[T] {
//	var errList []*TaskWorkCrash[T]
//	size := len(datas)
//	if size == 0 {
//		return errList
//	}
//	if workCount < 0 {
//		workCount = runtime.NumCPU()
//	}
//	if workCount == 0 || workCount > size {
//		workCount = size
//	}
//	var waiter sync.WaitGroup
//	waiter.Add(size)
//	queue := make(chan T)
//
//	writer := func() {
//		for i := 0; i < size; i++ {
//			queue <- datas[i]
//		}
//		close(queue)
//	}
//	reader := func(id int) {
//		for {
//			v, ok := <-queue
//			if ok {
//				func() {
//					defer func() {
//						if err := recover(); err != nil {
//							switch err.(type) {
//							case nil:
//							default:
//								crash := &TaskWorkCrash[T]{data: v, error: err}
//								errList = append(errList, crash)
//							}
//						}
//					}()
//					task(id, v)
//				}()
//				waiter.Done()
//			} else {
//				break
//			}
//		}
//	}
//	go writer()
//
//	for i := 0; i < workCount; i++ {
//		go reader(i)
//	}
//	waiter.Wait()
//	return errList
//}
