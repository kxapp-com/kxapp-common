package utilz

import "sync"

/*
*
任务工厂，对数据datas，每个数据需要进行task处理。分配workcount个工人（go程）进行并行工作
*/
func TaskWork[T any](datas []T, task func(workId int, data T), workCount int) {
	size := len(datas)
	if size == 0 {
		return
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
				task(id, v)
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
}
