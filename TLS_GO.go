package main


import (
"fmt"
)

//in the world of lisp,everything is either atom or list

type LinkList struct {
	Car interface{}
	Cdr *LinkList
}

func List(arr ...interface{}) *LinkList {
	ret := EmptyList()
	for i := len(arr) - 1; i != -1; i-- {
		ret = Cons(arr[i], ret)
	}
	return ret
}

func (l *LinkList) PrettyPrint() {
	l.RPrint()
	fmt.Println()
}
func (l *LinkList) RPrint() {
	fmt.Print("'(")
	defer fmt.Print(")")
	for p := l; !p.IsEmptyList(); p = p.Cdr {
		if pl, ok := p.Car.(*LinkList); ok {
			pl.RPrint()
		} else {
			fmt.Print(p.Car)
		}
		if !p.Cdr.IsEmptyList() {
			fmt.Print(" ")
		}
	}
}

func Cons(val interface{}, l *LinkList) *LinkList {
	ret := &LinkList{Car: val}
	ret.Cdr = l
	return ret
}

func EmptyList() *LinkList {
	str := "EmptyList_Unique_ID"
	return &LinkList{Car: str, Cdr: nil}
}

func (l *LinkList) IsEmptyList() bool {
	str := "EmptyList_Unique_ID"
	return l.Cdr == nil && l.Car == str
}

func IsNull(a interface{}) bool {
	if IsAtom(a) {
		return false
	}
	l := a.(*LinkList)
	return l.IsEmptyList()
}

func IsAtom(a interface{}) bool {
	switch a.(type) {
	case *LinkList:
		return false
	default:
		return true
	}
}

func IsListOfAtom(a interface{}) bool {
	if IsNull(a) {
		return true
	}
	if IsAtom(a) {
		return false
	}
	l := a.(*LinkList)
	return IsAtom(l.Car) && IsListOfAtom(l.Cdr)
}

func IsMember(a interface{}, lat *LinkList) bool {
	if IsNull(lat) {
		return false
	}
	if a == lat.Car {
		return true
	}
	return IsMember(a, lat.Cdr)
}

func Rember(a interface{}, lat *LinkList) *LinkList {
	if IsNull(lat) {
		return EmptyList()
	}
	if a == lat.Car {
		return lat.Cdr
	}
	return Cons(lat.Car, Rember(a, lat.Cdr))
}

func MultiRember(a interface{}, lat *LinkList) *LinkList {
	if IsNull(lat) {
		return EmptyList()
	}
	if a == lat.Car {
		return MultiRember(a, lat.Cdr)
	}
	return Cons(lat.Car, MultiRember(a, lat.Cdr))
}

// continuation passing style
// collect whatever you want in the collector
func MultiRemberAndCollect(a interface{}, lat *LinkList, col func(l1, l2 *LinkList) interface{}) interface{} {
	if IsNull(lat) {
		return col(EmptyList(), EmptyList())
	}
	if a == lat.Car {
		return MultiRemberAndCollect(a, lat.Cdr,
			func(newlat, seen *LinkList) interface{} {
				return col(newlat, Cons(lat.Car, seen))
			})
	}
	return MultiRemberAndCollect(a, lat.Cdr,
		func(newlat, seen *LinkList) interface{} {
			return col(Cons(lat.Car, newlat), seen)
		})
}

func InsertR(new interface{}, old interface{}, lat *LinkList) *LinkList {
	if IsNull(lat) {
		return EmptyList()
	}
	if old == lat.Car {
		return Cons(lat.Car, Cons(new, lat.Cdr))
	}
	return Cons(lat.Car, InsertR(new, old, lat.Cdr))
}

func RemberStar(a interface{}, l *LinkList) *LinkList {
	if IsNull(l) {
		return EmptyList()
	}
	if IsAtom(l.Car) {
		if a == l.Car {
			return RemberStar(a, l.Cdr)
		} else {
			return Cons(l.Car, RemberStar(a, l.Cdr))
		}
	}
	lcar := l.Car.(*LinkList) //in the world of lisp,everything is either atom or list
	return Cons(RemberStar(a, lcar), RemberStar(a, l.Cdr))
}

func InsertRStar(new interface{}, old interface{}, l *LinkList) *LinkList {
	if IsNull(l) {
		return l
	}
	if IsAtom(l.Car) {
		if old == l.Car {
			return Cons(old, Cons(new, InsertRStar(new, old, l.Cdr)))
		}
		return Cons(l.Car, InsertRStar(new, old, l.Cdr))
	}
	lCar := l.Car.(*LinkList)
	return Cons(InsertRStar(new, old, lCar), InsertRStar(new, old, l.Cdr))
}

func OccurStar(a interface{}, l *LinkList) int {
	if IsNull(l) {
		return 0
	}
	if IsAtom(l.Car) {
		if a == l.Car {
			return 1 + OccurStar(a, l.Cdr)
		}
		return OccurStar(a, l.Cdr)
	}
	lCar := l.Car.(*LinkList)
	return OccurStar(a, lCar) + OccurStar(a, l.Cdr)
}

func IsEqualList(l1, l2 *LinkList) bool {
	if IsNull(l1) && IsNull(l2) {
		return true
	}
	if IsNull(l1) || IsNull(l2) {
		return false
	}
	if IsAtom(l1.Car) && IsAtom(l2.Car) {
		return l1.Car == l2.Car && IsEqualList(l1.Cdr, l2.Cdr)
	}
	if IsAtom(l1.Car) || IsAtom(l2.Car) {
		return false
	}
	l1Car := l1.Car.(*LinkList)
	l2Car := l2.Car.(*LinkList)
	return IsEqualList(l1Car, l2Car) && IsEqualList(l1.Cdr, l2.Cdr)
}

func IsSet(lat *LinkList) bool {
	if IsNull(lat) {
		return true
	}
	if IsMember(lat.Cdr, lat.Cdr) {
		return false
	}
	return IsSet(lat.Cdr)
}

//
func MakeSet(lat *LinkList) *LinkList {
	if IsNull(lat) {
		return EmptyList()
	}
	if IsMember(lat.Car, lat.Cdr) {
		return MakeSet(lat.Cdr)
	}
	return Cons(lat.Car, MakeSet(lat.Cdr))
}

//is s1 subset of s2?
func IsSubSet(s1, s2 *LinkList) bool {
	if IsNull(s1) {
		return true
	}
	if IsMember(s1.Car, s2) {
		return IsSubSet(s1.Cdr, s2)
	}
	return false
}

// lisp has no loop and variable,only recursion(many tail recursions also)
func IsSubSetLoop(s1, s2 *LinkList) bool {
	rem := s1
	for !IsNull(rem) {
		if !IsMember(rem.Car, s2) {
			return false
		}
		rem = rem.Cdr
	}
	return true
}

func HasIntersect(s1, s2 *LinkList) bool {
	if IsNull(s1) {
		return false
	}
	if IsMember(s1.Car, s2) {
		return true
	}
	return HasIntersect(s1.Cdr, s2)
}

func Intersect(s1, s2 *LinkList) *LinkList {
	if IsNull(s1) {
		return EmptyList()
	}
	if IsMember(s1.Car, s2) {
		return Cons(s1.Car, Intersect(s1.Cdr, s2))
	}
	return Intersect(s1.Cdr, s2)
}

func Union(s1, s2 *LinkList) *LinkList {
	if IsNull(s1) {
		return s2
	}
	if IsMember(s1.Car, s2) {
		return Union(s1.Cdr, s2)
	}
	return Cons(s1.Car, Union(s1.Cdr, s2))
}

// income is a list of lists,Can't be represented in type
// type assertion needed
func IntersectAll(listofSet *LinkList) *LinkList {
	if IsNull(listofSet.Cdr) {
		return listofSet.Car.(*LinkList)
	}
	return Intersect(listofSet.Car.(*LinkList), IntersectAll(listofSet.Cdr))
}

func MultiRemberT(lat *LinkList, test func(interface{}) bool) *LinkList {
	if IsNull(lat) {
		return EmptyList()
	}
	if test(lat.Car) {
		return MultiRemberT(lat.Cdr, test)
	}
	return Cons(lat.Car, MultiRemberT(lat.Cdr, test))
}

func MultiInsertL(new, old interface{}, lat *LinkList) *LinkList {
	if IsNull(lat) {
		return EmptyList()
	}
	if lat.Car == old {
		return Cons(new, Cons(old, MultiInsertL(new, old, lat.Cdr)))
	}
	return Cons(lat.Car, MultiInsertL(new, old, lat.Cdr))
}

func MultiInsertR(new, old interface{}, lat *LinkList) *LinkList {
	if IsNull(lat) {
		return EmptyList()
	}
	if lat.Car == old {
		return Cons(old, Cons(new, MultiInsertR(new, old, lat.Cdr)))
	}
	return Cons(lat.Car, MultiInsertR(new, old, lat.Cdr))
}

//requires oldL != oldR
func MultiInsertLR(new, oldL, oldR interface{}, lat *LinkList) *LinkList {
	if IsNull(lat) {
		return EmptyList()
	}
	if lat.Car == oldL {
		return Cons(new, Cons(oldL, MultiInsertLR(new, oldL, oldR, lat.Cdr)))
	}
	if lat.Car == oldR {
		return Cons(oldR, Cons(new, MultiInsertLR(new, oldL, oldR, lat.Cdr)))
	}
	return Cons(lat.Car, MultiInsertLR(new, oldL, oldR, lat.Cdr))
}

//requires oldL != oldR
func MultiInsertLRAndCollect(new, oldL, oldR interface{}, lat *LinkList,
	col func(res *LinkList, oldLCount, oldRCount int) interface{}) interface{} {
	if IsNull(lat) {
		//return EmptyList()
		return col(EmptyList(), 0, 0)
	}
	if lat.Car == oldL {
		//return Cons(new, Cons(oldL, MultiInsertLRAndCollect(new, oldL, oldR, lat.Cdr)))
		return MultiInsertLRAndCollect(new, oldL, oldR, lat.Cdr,
			func(latCdrRes *LinkList, latCdrLCount, latCdrRCount int) interface{} {
				return col(Cons(new, Cons(oldL, latCdrRes)), latCdrLCount+1, latCdrRCount)
			})
	}
	if lat.Car == oldR {
		//return Cons(oldR, Cons(new, MultiInsertLRAndCollect(new, oldL, oldR, lat.Cdr)))
		return MultiInsertLRAndCollect(new, oldL, oldR, lat.Cdr,
			func(latCdrRes *LinkList, latCdrLCount, latCdrRCount int) interface{} {
				return col(Cons(oldR, Cons(new, latCdrRes)), latCdrLCount, latCdrRCount+1)
			})
	}
	//return Cons(lat.Car, MultiInsertLRAndCollect(new, oldL, oldR, lat.Cdr))
	return MultiInsertLRAndCollect(new, oldL, oldR, lat.Cdr,
		func(latCdrRes *LinkList, latCdrLCount, latCdrRCount int) interface{} {
			return col(Cons(lat.Car, Cons(new, latCdrRes)), latCdrLCount, latCdrRCount+1)
		})
}

func IsEven(n int) bool {
	return n%2 == 0
}

func EvensOnlyStar(l *LinkList) *LinkList {
	if IsNull(l) {
		return EmptyList()
	}
	if IsAtom(l.Car) {
		if IsEven(l.Car.(int)) {
			return Cons(l.Car, EvensOnlyStar(l.Cdr))
		}
		return EvensOnlyStar(l.Cdr)
	}
	return Cons(EvensOnlyStar(l.Car.(*LinkList)), EvensOnlyStar(l.Cdr))
}
func EvensOnlyStarAndCol(l *LinkList, col func(list *LinkList) interface{}) interface{} {
	if IsNull(l) {
		return col(EmptyList())
	}
	if IsAtom(l.Car) {
		if IsEven(l.Car.(int)) {
			//return Cons(l.Car, EvensOnlyStarAndCol(l.Cdr))
			return EvensOnlyStarAndCol(l.Cdr, func(lCdrres *LinkList) interface{} {
				return col(Cons(l.Car, lCdrres))
			})
		}
		//return EvensOnlyStarAndCol(l.Cdr)
		return EvensOnlyStarAndCol(l.Cdr, col)
	}
	//return Cons(EvensOnlyStarAndCol(l.Car.(*LinkList)), EvensOnlyStarAndCol(l.Cdr))
	return EvensOnlyStarAndCol(l.Car.(*LinkList), func(lCarres *LinkList) interface{} {
		return EvensOnlyStarAndCol(l.Cdr, func(lCdrres *LinkList) interface{} {
			return col(Cons(lCarres, lCdrres))
		})
	})
}


func main() {

	List().PrettyPrint()
	List(2, 3, 4, true, "dsaf").PrettyPrint()
	List(1, 2, List(3, 4, List(5, 6))).PrettyPrint()
	Union(List(1, 2, 3), List(2, 3, 4)).PrettyPrint()
	IntersectAll(List(
		List(1, 2, 3),
		List(2, 3, 4),
		List(3, 4, 5))).PrettyPrint()
	IntersectAll(List(
		List("Naive!"),
		List("Naive!", "Angry!"))).PrettyPrint()
	MultiRemberT(List(1, 2, 3), func(i interface{}) bool {
		iint := i.(int)
		return iint == 2
	}).PrettyPrint()
	MultiRemberT(List("Naive!", "Angry!", "Simple!"), func(s interface{}) bool {
		sstr := s.(string)
		return sstr == "Simple!"
	}).PrettyPrint()
	MultiRemberAndCollect(3, List(1, 3, 2, 3, 4), func(l1, l2 *LinkList) interface{} {
		l1.PrettyPrint()
		l2.PrettyPrint()
		return struct{}{} //we can return whatever we like
	})
	MultiInsertL(3, 2, List(1, 2, 2, 5, 2, 2, 4)).PrettyPrint()
	MultiInsertR(3, 2, List(1, 2, 2, 5, 2, 2, 4)).PrettyPrint()
	MultiInsertLR(3, 2, 5, List(1, 2, 2, 5, 2, 2, 4)).PrettyPrint()
	EvensOnlyStar(List(1, 2, List(3), 4, List(5, List(6, 7, 8)))).PrettyPrint()
	EvensOnlyStarAndCol(List(1, 2, List(3), 4, List(5, List(6, 7, 8))),
		func(list *LinkList) interface{} {
			return list
		}).(*LinkList).PrettyPrint()
}
