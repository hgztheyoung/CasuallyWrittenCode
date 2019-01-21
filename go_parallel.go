package main

import (
	"time"
	"log"
	"math/rand"
)

func Job(i int) {
	<-time.After(time.Duration(1+rand.Intn(2)) * time.Second)
	log.Println("Job ", i, " done")
	return
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
