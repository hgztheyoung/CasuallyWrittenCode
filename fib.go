package main

import(
	"fmt"
	"time"
	"sync"
)

func fib(n int64) int64{
	if(n == 0 || n == 1){
		return 1
	}
	f1 := fib(n-1)
	f2 := fib(n-2)
	return f1+f2
}

func fib_wg_go(n int64,res *int64) {
	// base case
	// 20 is like a magic number,can be tuned to get better performance.
	// the trade off between base case time and the allocation of all this waitgroups and go threads
	if(n<20){
		*res = fib(n)
		return 
	}
	var f1,f2 int64
	var wg sync.WaitGroup = sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		fib_wg_go(n-1,&f1)
	}()
	fib_wg_go(n-2,&f2)
	// cilk_sync
	wg.Wait()
	*res = f1+f2
}

func fib_ch(n int64) int64{
	if(n<20){
		return fib(n)
	}
	f1 := make(chan int64,1)
	defer close(f1)
	go func(){
		f1 <- fib(n-1)
	}()
	f2 := fib(n-2)
	return <-f1+f2
}

func fib_wg_go_wrap(n int64) int64{
	var res int64
	fib_wg_go(n,&res)
	return res
}

func perf(f func(int64)int64,n int64){
	fst := time.Now()
	fmt.Println(f(n))
	l := time.Now()
	fmt.Println(l.Sub(fst))
	fmt.Println()
}

func main() {
	var n int64 = 43
	perf(fib,n)
	perf(fib_wg_go_wrap,n)
	perf(fib_ch,n)
}

// go run fib.go
// 701408733
// 3.9187619s

// 701408733
// 1.0933758s

// 701408733
// 2.3916358s
