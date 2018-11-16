package main

import (
    "fmt"
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
    ch := make(chan struct{})
    for count:=0;count<worker_count;count++{
        go func(){ //Gopher means go func
            for {
                i := <-c
                s := <-res
                res <-s+i
                ch<-struct{}{}
            }
        }()
    }
    for c:=0;c<data_size;c++{
        <-ch
    }
    fmt.Println(<-res)
}
