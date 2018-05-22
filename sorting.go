package main

import (
	"fmt"
	"math/rand"
	"time"
	"sync"
)

func main() {
	big := 100000000
	A := rand.Perm(big)	
	s := time.Now()
	sort(A)
	l := time.Now()
	fmt.Println(l.Sub(s))
	// fmt.Println(sort(A))
}

// top-down approach
func sort(A []int) []int {
    if len(A) <= 20000 {
        return sort_s(A)
    }
	mid := len(A) / 2
	left, right := A[0:mid], A[mid:]
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(){
		defer wg.Done()
		left = sort(left)
	}()
	right = sort(right)
	wg.Wait()
    return merge_s(left, right)
}

//_s stands for single thread
func sort_s(A []int) []int {
    if len(A) <= 1 {
        return A
    }
	mid := len(A) / 2
	left, right := A[0:mid], A[mid:]
	left = sort_s(left)
	right = sort_s(right)
	return merge_s(left, right)
}

// assumes that A and B are sorted
func merge(A, B []int) []int {
	lA,lB := len(A),len(B)
	if lA+lB<100000 {
		return merge_s(A,B)
	}else if lA<lB {
		return merge(B,A)
	}else if lA == 0{
		return B
	}else {
		
		ma := lA / 2	
		mb := binarySearch(B,A[ma])
		wg := sync.WaitGroup{}
		wg.Add(1)
		var l []int
		go func(){
			defer wg.Done()
			l = merge(A[:ma],B[:mb])
		}()	
		r := merge(A[ma+1:],B[mb:])
		wg.Wait()
		l = append(l,A[ma])
		l = append(l,r...)
		return l
	}
}
func merge_s(A, B []int) []int {
	lA,lB := len(A),len(B)
	ai,bi,ri :=0,0,-1
	ret := make([]int, lA+lB)
	for ai!=lA && bi !=lB {
		ri++
		if(A[ai]<B[bi]){
			ret[ri] = A[ai]			
			ai++
		}else{
			ret[ri] = B[bi]
			bi++
		}
	}
	for ai!=lA {
		ri++
		ret[ri] = A[ai]			
		ai++
	}
	for bi!=lB {
		ri++
		ret[ri] = B[bi]			
		bi++
	}
	return ret
}

// 作者：Simth
// 链接：https://www.jianshu.com/p/fe4728688adb
// 來源：简书
// 著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。
func binarySearch(sortedArray []int, lookingFor int) int {
    var low int = 0
    var high int = len(sortedArray) - 1
    for low <= high {
        var mid int =low + (high - low)/2
        var midValue int = sortedArray[mid]
        if midValue == lookingFor {
            return mid
        } else if midValue > lookingFor {
            high = mid -1
        } else {
            low = mid + 1
        }
    }
    return low
}
