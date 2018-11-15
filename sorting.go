package main
import (
  "fmt"
  "math/rand"
  "time"
)

func main() {
  // arr := []int{3,2,1,5,4,6,8,7,9})
  arr := rand.Perm(1000000)
  t1 := time.Now()
  MergeSort(arr)
  fmt.Println(time.Since(t1))
  t2 := time.Now()
  MergeSort_P(arr)
  fmt.Println(time.Since(t2))
}




func MergeSort(arr []int) []int{
  if len(arr)<=1{
    return arr
  }
  mid := len(arr)/2
  l := MergeSort(arr[:mid])
  r := MergeSort(arr[mid:])
  return Merge(l,r)
}

func MergeSort_P(arr []int) []int{
  if len(arr)<=1{
    return arr
  }
  mid := int(len(arr)/2)
  ch := make(chan struct{})
  l := arr[:mid]
  go func(){
    l = MergeSort_P(arr[:mid])
    ch <- struct{}{}
  }()
  r := MergeSort_P(arr[mid:])
  <-ch
  return Merge_P(l,r)
}

func Merge(l,r []int) []int{
  ll,lr :=len(l),len(r)
  length :=ll+lr 
  ret := make([]int,length,length)
  k,i,j:=0,0,0
  for i!=ll &&j!=lr{
    if l[i] < r[j]{
      ret[k] = l[i]
      i++
    }else{
      ret[k] = r[j]
      j++
    }
    k++
  }
  for i!=ll{
    ret[k] = l[i]
    i++
    k++
  }
  for j!=lr{
    ret[k] = r[j]
    j++
    k++
  }
  return ret
}

// binary search ,return index i where
// forall x in arr[:i] n >= x and 
// forall y in arr[i:] n < y
// notice that arr[:i] or arr[i:] can be empty 
func BinarySearch(arr []int, n int) int{
  f,l := 0,len(arr)
  for f<l{
    mid := f + (l-f)/2
    if arr[mid] <= n{
      f = mid+1
    }else{
      l = mid
    }
  }
  return f
}


func Merge_P(a1,a2 []int) []int{
  // base case
  if len(a1)+len(a2) < 100000{
    return Merge(a1,a2)
  }
  // parallel case
  if len(a2)<len(a1){
    a1,a2 = a2,a1
  }
  mida2 := a2[len(a2)/2]
  mida2_in_a1_i := BinarySearch(a1,mida2)
  done := make(chan struct{})
  var large []int
  go func(){
    large = Merge_P(a1[mida2_in_a1_i:],a2[len(a2)/2:])    
    done <- struct{}{}
  }()
  <-done
  small := Merge_P(a1[:mida2_in_a1_i],a2[:len(a2)/2])
  
  return append(small,large...)
}
