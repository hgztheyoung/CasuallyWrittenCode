package main

import "fmt"

//in the world of lisp,everything is either atom or list

type LinkList struct {
	Car interface{}
	Cdr *LinkList
}

func (l *LinkList) PrettyPrint() {
	l.RPrint()
	fmt.Println()
}
func (l *LinkList) RPrint() {
	fmt.Print("'(")
	defer fmt.Print(")")
	p := l
	for !p.IsEmptyList() {
		if pl, ok := p.Car.(*LinkList); ok {
			pl.RPrint()
		} else {
			fmt.Print(p.Car)
		}
		if !p.Cdr.IsEmptyList() {
			fmt.Print(" ")
		}
		p = p.Cdr
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

func main() {
	//l := Cons(Cons(1, Cons(3, EmptyList())), Cons(2, EmptyList()))
	//l.PrettyPrint()
	//e := EmptyList()
	//e.PrettyPrint()
	l := Cons(1, EmptyList())
	l2 := Cons(2, Cons(l, EmptyList()))
	//l3 := Cons(3, 4) //type system rule out this case
	fmt.Println(IsListOfAtom(l))
	fmt.Println(IsListOfAtom(l2))
	l3 := Cons(1, Cons(2, Cons(3, EmptyList())))
	l4 := Rember(2, l3)
	l4.PrettyPrint()
	l5 := InsertR(4, 2, l3)
	l5.PrettyPrint()
	l6 := Cons(1, Cons(2, Cons(2, Cons(3, EmptyList()))))
	MultiRember(2, l6).PrettyPrint()
	l7 := Cons(l3, l6)
	RemberStar(1, l7).PrettyPrint()
}
