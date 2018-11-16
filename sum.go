package main

import (
    "fmt"
    "sync"
    "math/rand"
)

func main() {
    data_size := 100000
    worker_count := 100
    queue_size := 50
    
    arr := rand.Perm(data_size)
    c := make(chan int,queue_size)
    go func(){
        for _,i :=range arr{
            c <- i
        }
    }()
    res := make(chan int,1)
    res <-0
    var wg sync.WaitGroup
    wg.Add(data_size)
    for count:=0;count<worker_count;count++{
        go func(){ //Gopher means go func
            for {
                i := <-c
                s := <-res
                res <-s+i
                wg.Done()
            }
        }()
    }
    wg.Wait()
    fmt.Println(<-res)
}
