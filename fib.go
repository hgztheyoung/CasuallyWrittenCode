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

var wg sync.WaitGroup

func fib_wg(n int64,res *int64) {
	if(n == 0 || n == 1){
		*res = 1
		return
	}
	var f1,f2 int64
	fib_wg(n-1,&f1)
	fib_wg(n-2,&f2)
	*res = f1+f2
}

func fib_wg_wrap(n int64) int64{
	var res int64
	fib_wg(n,&res)
	return res
}

func fib_wg_go(n int64,res *int64) {
	if(n == 0 || n == 1){
		*res = 1
		return
	}
	var f1,f2 int64
	go func() {
		wg.Add(1)
		defer wg.Done()
		fib_wg(n-1,&f1)
	}()
	fib_wg(n-2,&f2)
	// cilk_sync
	wg.Wait()
	*res = f1+f2
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
	perf(fib,43)
	wg = sync.WaitGroup{}
	perf(fib_wg_wrap,43)
	perf(fib_wg_go_wrap,43)
}