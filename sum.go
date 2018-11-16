package main

import (
    "fmt"
    "sync"
    "math/rand"
)

func main() {
    size := 10000
    arr := rand.Perm(size)
    c := make(chan int,10)
    go func(){
        for _,i :=range arr{
            c <- i
        }
    }()
    res := make(chan int,1)
    res <-0
    var wg sync.WaitGroup
    wg.Add(size)
    for count:=0;count<10;count++{
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
