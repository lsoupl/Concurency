package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	// defer block
	//defer fmt.Println(1)
	//defer fmt.Println(2)
	//fmt.Println(sum(2, 3))
	//deferValues()

	// goroutines block
	//runtime.GOMAXPROCS(1)
	//fmt.Println(runtime.NumCPU())

	//go showNumbers(100)

	//runtime.Gosched()
	//time.Sleep(time.Second)
	//fmt.Println("exit")
	//
	//makePanic()

	//withoutWait()
	//withWait()
	//wrongAdd()
	//writeWithoutConcurrent()
	//writeWithoutMutex()
	//writeWithMutex()
	//readWithMutex()
	//readWithRWMutex()

	// channel:
	//nilChannel()
	//unbufferedChannel()
	//bufferedChannel()
	//forRange()
	//baseSelect()
	//gracefulShutdown()

	// context
	//baseKnowledge()
	//workerPool()

	//chanAsPromise()
	//chanAsMutex()
	//withoutErrGroup()
	//errGroup()

	//AddMutex()
	//AddAtomic()

	StoreLoadSwap()
	//compareAndSwap()

	//atomicVal()
}

func showNumbers(num int) {
	for i := 0; i < num; i++ {
		fmt.Println(i)
	}
}

func sum(x, y int) (sum int) {
	defer func() {
		sum *= 2
	}()

	sum = x + y
	return
}

func deferValues() {
	for i := 0; i < 10; i++ {
		defer fmt.Println("first", i)
	}
	for i := 0; i < 10; i++ {
		defer func() {
			fmt.Println("second", i)
		}()
	}

	for i := 0; i < 10; i++ {
		k := i
		defer func() {
			fmt.Println("third", k)
		}()
	}
	for i := 0; i < 10; i++ {
		defer func(k int) {
			fmt.Println("fourth", k)
		}(i)
	}
}

func makePanic() {
	defer func() {
		panicValue := recover()
		fmt.Println(panicValue)
	}()

	panic("some panic")
	fmt.Println("Unreachable code")
}

func withoutWait() {
	for i := 0; i <= 10; i++ {
		go fmt.Println(i + 1)
	}

	fmt.Println("exit")
}

func withWait() {
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		//wg.Add(1)
		go func(i int) {
			//defer wg.Done()
			fmt.Println(i + 1)
			wg.Done()
		}(i)
	}

	wg.Wait()
	fmt.Println("exit")
}

func wrongAdd() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		go func(i int) {
			wg.Add(1)

			defer wg.Done()

			fmt.Println(i + 1)
		}(i)
	}

	wg.Wait()
	fmt.Println("exit")
}

func writeWithoutConcurrent() {
	start := time.Now()
	var counter int

	for i := 0; i < 1000; i++ {
		time.Sleep(time.Nanosecond)
		counter++
	}

	fmt.Println(counter)
	fmt.Println(time.Now().Sub(start).Seconds())
}

func writeWithoutMutex() {
	start := time.Now()
	var counter int
	var wg sync.WaitGroup

	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(time.Nanosecond)
			counter++
		}()
	}
	wg.Wait()

	fmt.Println(counter)
	fmt.Println(time.Now().Sub(start).Seconds())
}

func writeWithMutex() {
	start := time.Now()
	var counter int
	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(time.Nanosecond)

			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()

	fmt.Println(counter)
	fmt.Println(time.Now().Sub(start).Seconds())
}

func readWithMutex() {
	start := time.Now()
	var (
		counter int
		wg      sync.WaitGroup
		mu      sync.Mutex
	)

	wg.Add(100)

	for i := 0; i < 50; i++ {
		go func() {
			defer wg.Done()

			mu.Lock()

			time.Sleep(time.Nanosecond)
			_ = counter

			mu.Unlock()
		}()

		go func() {
			defer wg.Done()

			mu.Lock()

			time.Sleep(time.Nanosecond)
			counter++

			mu.Unlock()
		}()
	}

	wg.Wait()

	fmt.Println(counter)
	fmt.Println(time.Now().Sub(start).Seconds())
}

func readWithRWMutex() {
	start := time.Now()
	var (
		counter int
		wg      sync.WaitGroup
		mu      sync.RWMutex
	)

	wg.Add(100)

	for i := 0; i < 50; i++ {
		go func() {
			defer wg.Done()

			mu.RLock()

			time.Sleep(time.Nanosecond)
			_ = counter

			mu.RUnlock()
		}()

		go func() {
			defer wg.Done()

			mu.Lock()

			time.Sleep(time.Nanosecond)
			counter++

			mu.Unlock()
		}()
	}

	wg.Wait()

	fmt.Println(counter)
	fmt.Println(time.Now().Sub(start).Seconds())
}

func nilChannel() {
	var nilChannel chan int
	fmt.Printf("Len: %d\tCap: %d\n", nilChannel, nilChannel)

	/*nilChannel <- 1

	<-nilChannel*/
	close(nilChannel)
}

func unbufferedChannel() {
	unbufferedChannel := make(chan int)
	fmt.Printf("Len: %d\tCap: %d\n", unbufferedChannel, unbufferedChannel)

	//unbufferedChannel <- 1

	//<-unbufferedChannel

	go func(chanForWriting chan<- int) {
		time.Sleep(time.Second)
		//<-chanForWriting
		unbufferedChannel <- 1
	}(unbufferedChannel)

	value := <-unbufferedChannel
	fmt.Println(value)

	go func(chanForReading <-chan int) {
		time.Sleep(time.Second)
		value := <-chanForReading
		fmt.Println(value)
	}(unbufferedChannel)

	unbufferedChannel <- 2

	go func() {
		time.Sleep(time.Millisecond * 500)
		close(unbufferedChannel)
	}()
	go func() {
		time.Sleep(time.Second)
		fmt.Println(<-unbufferedChannel)
	}()

	//unbufferedChannel <- 3
	//close(unbufferedChannel)
	//close(unbufferedChannel)

}

func bufferedChannel() {
	bufferedChannel := make(chan int, 2)
	fmt.Printf("Length: %d\tCap: %d\n", len(bufferedChannel), cap(bufferedChannel))

	bufferedChannel <- 1
	bufferedChannel <- 2

	fmt.Printf("Length: %d\tCap: %d\n", len(bufferedChannel), cap(bufferedChannel))

	//bufferedChannel <- 3

	fmt.Println(<-bufferedChannel)
	fmt.Println(<-bufferedChannel)

	fmt.Printf("Length: %d\tCap: %d\n", len(bufferedChannel), cap(bufferedChannel))

	//fmt.Println(<-bufferedChannel)

}

func forRange() {
	bufferedChannel := make(chan int, 3)

	numbers := []int{5, 6, 7, 8}

	// show all elements with for

	go func() {
		for _, num := range numbers {
			bufferedChannel <- num
		}
		close(bufferedChannel)
	}()

	for {
		v, ok := <-bufferedChannel
		fmt.Println(v, ok)

		if !ok {
			break
		}
	}

	// show with for range buffered
	bufferedChannel = make(chan int, 3)

	go func() {
		for _, num := range numbers {
			bufferedChannel <- num
		}
		close(bufferedChannel)
	}()

	for v := range bufferedChannel {
		fmt.Println("buffered", v)
	}

	unbufferedChannel := make(chan int)
	go func() {
		for _, num := range numbers {
			unbufferedChannel <- num
		}
		close(unbufferedChannel)
	}()

	for value := range unbufferedChannel {
		fmt.Println("unbuffered", value)
	}
}

func baseSelect() {
	bufferedChan := make(chan string, 3)
	bufferedChan <- "first"

	select {
	case str := <-bufferedChan:
		fmt.Println("read", str)
	case bufferedChan <- "second":
		fmt.Println("write", <-bufferedChan, <-bufferedChan)
	}

	unbufChan := make(chan int)

	go func() {
		time.Sleep(time.Second)
		unbufChan <- 1
	}()

	select {
	case bufferedChan <- "third":
		fmt.Println("unblocking writing")
	case val := <-unbufChan:
		fmt.Println("Blocking reading:", val)
	case <-time.After(time.Millisecond * 1500):
		fmt.Println("time is up")
	default:
		fmt.Println("default case")
	}

	resultChan := make(chan int)
	timer := time.After(time.Second)

	go func() {
		defer close(resultChan)

		for i := 1; i <= 1000; i++ {
			select {
			case <-timer:
				fmt.Println("time is up")
				return
			default:
				time.Sleep(time.Nanosecond)
				resultChan <- i
			}
		}
	}()

	for v := range resultChan {
		fmt.Println(v)
	}

}

func gracefulShutdown() {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	timer := time.After(10 * time.Second)

	select {
	case <-timer:
		fmt.Println("time is up")
		return
	case sig := <-sigChan:
		fmt.Println("Stopped by signal:", sig)
		return
	}
}

func baseKnowledge() {
	ctx := context.Background()
	fmt.Println(ctx)

	toDo := context.TODO()
	fmt.Println(toDo)

	withValue := context.WithValue(ctx, "name", "vasya")
	fmt.Println(withValue.Value("name"))

	withCancel, cancel := context.WithCancel(ctx)
	fmt.Println(withCancel.Err())
	cancel()
	fmt.Println(withCancel.Err())

	withDeadline, cancel := context.WithDeadline(ctx, time.Now().Add(time.Second*3))
	defer cancel()
	fmt.Println(withDeadline.Deadline())
	fmt.Println(withDeadline.Err())
	fmt.Println(<-withDeadline.Done())

	withTimeout, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	fmt.Println(<-withTimeout.Done())
}

func workerPool() {
	ctx, cancel := context.WithCancel(context.Background())
	ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*20)

	defer cancel()

	wg := &sync.WaitGroup{}
	numberToProcess, processedNumbers := make(chan int, 5), make(chan int, 5)

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(ctx, numberToProcess, processedNumbers)
		}()
	}

	go func() {
		for i := 0; i <= 1000; i++ {
			if i == 500 {
				cancel()
			}
			numberToProcess <- i
		}
		close(numberToProcess)
	}()

	go func() {
		wg.Wait()
		close(processedNumbers)
	}()

	var counter int
	for resultValue := range processedNumbers {
		counter++
		fmt.Println(resultValue)
	}

	fmt.Println(counter)

}

func worker(ctx context.Context, toProcess <-chan int, processed chan<- int) {
	for {
		select {
		case <-ctx.Done():
			return
		case value, ok := <-toProcess:
			if !ok {
				return
			}
			time.Sleep(time.Millisecond)
			processed <- value * value
		}
	}
}

func makeRequest(num int) <-chan string {
	responseChan := make(chan string)

	go func() {
		time.Sleep(time.Second)
		responseChan <- fmt.Sprintf("response number %d", num)
	}()
	return responseChan
}

func chanAsPromise() {
	firstResponseChan := makeRequest(1)
	secondResponseChan := makeRequest(2)
	// do something else
	fmt.Println("non blocking")

	fmt.Println(<-firstResponseChan, <-secondResponseChan)
}

func chanAsMutex() {
	var counter int
	mutexChan := make(chan struct{}, 1)
	wg := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			mutexChan <- struct{}{}

			counter++

			<-mutexChan
		}()
	}
	wg.Wait()
	fmt.Println(counter)
}

func withoutErrGroup() {
	var err error
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	wg.Add(3)

	go func() {
		time.Sleep(time.Second)
		defer wg.Done()

		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println("first started")
			time.Sleep(time.Second)
		}
	}()

	go func() {
		defer wg.Done()

		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println("second start")
			err = fmt.Errorf("any error")
			cancel()
		}
	}()

	go func() {
		defer wg.Done()

		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println("third start")
			time.Sleep(time.Second)
		}
	}()

	wg.Wait()
	fmt.Println(err)
}

func errGroup() {
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		time.Sleep(time.Second)

		select {
		case <-ctx.Done():
			return nil
		default:
			fmt.Println("first started")
			time.Sleep(time.Second)
			return nil
		}
	})
	g.Go(func() error {
		fmt.Println("second started")
		return fmt.Errorf("unexpected error in request 2")
	})
	g.Go(func() error {
		select {
		case <-ctx.Done():
		default:
			fmt.Println("third started")
			time.Sleep(time.Second)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}
}

func AddMutex() {
	start := time.Now()

	var (
		counter int64
		wg      sync.WaitGroup
		mu      sync.Mutex
	)

	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()

			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println(counter)
	fmt.Println("With mutex:", time.Now().Sub(start).Seconds())
}

func AddAtomic() {
	start := time.Now()
	var (
		counter int64
		wg      sync.WaitGroup
	)

	wg.Add(1000)

	for i := 0; i < 1000; i++ {
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}
	wg.Wait()
	fmt.Println(counter)
	fmt.Println("With atomic:", time.Now().Sub(start).Seconds())
}

func StoreLoadSwap() {
	var counter int64

	fmt.Println(atomic.LoadInt64(&counter))

	atomic.StoreInt64(&counter, 5)
	fmt.Println(atomic.LoadInt64(&counter))

	fmt.Println(atomic.SwapInt64(&counter, 10))
	fmt.Println(atomic.LoadInt64(&counter))
}

func compareAndSwap() {
	var (
		counter int64
		wg      sync.WaitGroup
	)
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()

			if !atomic.CompareAndSwapInt64(&counter, 0, 1) {
				return
			}

			fmt.Println("Swapped goroutine number is", i)
		}(i)
	}

	wg.Wait()
	fmt.Println(counter)
}

func atomicVal() {
	var (
		value atomic.Value
	)

	value.Store(1)
	fmt.Println(value.Load())

	fmt.Println(value.Swap(2))
	fmt.Println(value.Load())

	fmt.Println(value.CompareAndSwap(2, 3))
	fmt.Println(value.Load())
}
