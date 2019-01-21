package main

import (
    "time"
    "log"
    "math/rand"
    "fmt"
)

func Job(i int) {
    <-time.After(time.Duration(1+rand.Intn(2)) * time.Second)
    log.Println("Job ", i, " done")
    return
}

func JobMayFail(i int) int {
    <-time.After(time.Duration(10+rand.Intn(10)) * time.Microsecond)
    flip := rand.Int31n(100)
    if flip > 80 {
        log.Println("Job ", i, " failed")
        return -1
    }
    log.Println("Job ", i, " done")
    return 0
}

func DoJobMayFailSequencial() {
    jobCount := 500
    undoneChan := make(chan int, jobCount)
    doneChan := make(chan int, jobCount)
    giveUpChan := make(chan int, jobCount)
    failChan := make(chan int, jobCount)
    failedCount := make(map[int]int)
    retries := 2
    for i := 0; i < jobCount; i++ {
        undoneChan <- i
    }
    givenUpCount := 0

    for len(doneChan)+len(giveUpChan) != jobCount {
        select {
        case i := <-undoneChan:
            ret := JobMayFail(i)
            if ret == 0 {
                doneChan <- i
            } else {
                failChan <- i
            }
        case i := <-failChan:
            failedCount[i]++
            if failedCount[i] == retries {
                log.Println("Job ", i, " givenup")
                giveUpChan <- i
            } else {
                undoneChan <- i
            }
        }
    }
    fmt.Println("len(undoneChan)", len(undoneChan))
    fmt.Println("len(doneChan)", len(doneChan))
    fmt.Println("len(giveUpChan)", len(giveUpChan), givenUpCount)
    fmt.Println("len(failChan)", len(failChan))
}

func DoJobMayFailParallel() {
    jobCount := 500000
    undoneChan := make(chan int, jobCount)
    doneChan := make(chan int, jobCount)
    giveUpChan := make(chan int, jobCount)
    failChan := make(chan int, jobCount)
    failedCount := make(map[int]int)
    retries := 2
    for i := 0; i < jobCount; i++ {
        undoneChan <- i
    }
    givenUpCount := 0

    workerCount := 500
    idle := make(chan int, workerCount)
    for i := 0; i < workerCount; i++ {
        idle <- i
    }

    for len(doneChan)+len(giveUpChan) != jobCount {
        select {
        case i := <-undoneChan:
            worker := <-idle
            go func() {
                worker := worker
                ret := JobMayFail(i)
                if ret == 0 {
                    doneChan <- i
                } else {
                    failChan <- i
                }
                idle <- worker
            }()
        case i := <-failChan:
            failedCount[i]++
            if failedCount[i] == retries {
                log.Println("Job ", i, " givenup")
                giveUpChan <- i
            } else {
                undoneChan <- i
            }

        default:
            //wait for return
            // in sequencial,we don't need this case,cause at any time,one of
            //<-undoneChan,<-failChan or len(doneChan)+len(giveUpChan) != jobCount (at last) will hold.
            // here,because of go func,doneChan <- i or failChan <- i may happen after select.
            time.Sleep(200 * time.Millisecond)
            continue
        }
    }
    fmt.Println("len(undoneChan)", len(undoneChan))
    fmt.Println("len(doneChan)", len(doneChan))
    fmt.Println("len(giveUpChan)", len(giveUpChan), givenUpCount)
    fmt.Println("len(failChan)", len(failChan))
}

func DoJobSequencial() {
    jobCount := 10
    //workerCount := 1
    for i := 0; i < jobCount; i++ {
        Job(i)
    }
}

func DoJobParallel() {
    jobCount := 500
    workerCount := 1500
    if workerCount > jobCount {
        workerCount = jobCount
    }
    idle := make(chan struct{}, workerCount)
    allDone := make(chan struct{})
    for i := 0; i < workerCount; i++ {
        idle <- struct{}{}
    }
    for i := 0; i < jobCount; i++ {
        <-idle
        i := i
        go func() {
            Job(i)
            idle <- struct{}{}
            if i == jobCount-1 {
                for len(idle) != workerCount {
                }
                allDone <- struct{}{}
            }
        }()
    }
    <-allDone
}

func randomDispatch(from <-chan struct{}, outs []chan<- struct{}) {
    for {
        select {
        case <-from:
            outs[rand.Intn(len(outs))] <- struct{}{}
        }
    }
}

func PingPongGoroutine() {
    c := make(chan struct{})
    counter := 0
    go func() {
        for {
            c <- struct{}{}
            <-c
            counter++
        }
    }()
    go func() {
        for {
            <-c
            c <- struct{}{}
        }
    }()
    <-time.After(1 * time.Second)
    fmt.Println(counter)
}
