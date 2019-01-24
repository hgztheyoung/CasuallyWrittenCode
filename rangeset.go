package collect

import (
	"errors"
	"log"
	"reflect"
)

// a range set of [f,l) ranges
type Range struct {
	F, L int64
}

//we actually store [f,l] range,but supply [f,l+1) interface
type RangeSet struct {
	Ranges []Range
}

func LowerBound(array []int64, first, last int, value int64) int {
	for first < last {
		mid := first + int((last-first)/2)
		if array[mid] < value {
			first = mid + 1
		} else {
			last = mid
		}
	}
	return first
}

func (rs *RangeSet) GetIntersect(input Range) (ranges []Range, rangei []int) {
	input.L -= 1
	ranges = make([]Range, 0)
	rangei = make([]int, 0)
	Fs := make([]int64, 0, len(rs.Ranges))
	Ls := make([]int64, 0, len(rs.Ranges))
	for _, r := range rs.Ranges {
		Fs = append(Fs, r.F)
		Ls = append(Ls, r.L)
	}
	ilinFs := LowerBound(Fs, 0, len(Fs), input.L)
	ifinLs := LowerBound(Ls, 0, len(Ls), input.F)
	//fmt.Println(Fs)
	//fmt.Println(Ls)
	//fmt.Println("F,L", input.F, input.L)
	//fmt.Println("F,L", ilinFs, ifinLs)
	if ilinFs == ifinLs {
		if ilinFs >= 1 && (input.F <= Ls[ilinFs-1]) {
			ranges = append(ranges, rs.Ranges[ilinFs-1])
			rangei = append(rangei, ilinFs-1)
		}
		if ilinFs < len(Fs) && (input.L >= Fs[ilinFs]) {
			ranges = append(ranges, rs.Ranges[ilinFs])
			rangei = append(rangei, ilinFs)
		}
	}
	for i := ifinLs; i < ilinFs; i++ {
		ranges = append(ranges, rs.Ranges[i])
		rangei = append(rangei, i)
	}
	return
}


func (rs *RangeSet) AddRange(input Range) (err error) {
	inter, _ := rs.GetIntersect(input)
	if len(inter) > 0 {
		log.Println(input, inter)
		return errors.New("intersect range exists,Add Failed")
	}
	Fs := make([]int64, 0, len(rs.Ranges))
	for _, r := range rs.Ranges {
		Fs = append(Fs, r.F)
	}

	lb := LowerBound(Fs, 0, len(Fs)-1, input.L)
	//Insert(rs.Ranges, lb, input)
	return
}
