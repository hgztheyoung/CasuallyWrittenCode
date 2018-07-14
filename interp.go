//inspired by www.yinwang.org/blog-cn/2012/08/01/interpreter

package main

import "fmt"

type Value struct {
	Value interface{}
}

type Env interface {
	empty_env() Env
	extend(k Symbol, v Value) Env
	lookup(k Symbol) Value
}

type Expr interface {
	interp(Env) Value
}

type Symbol struct {
	Name string
}

type Number struct {
	Number float64
}
type Lambda struct {
	Arg  Symbol
	Body Expr
}

type App struct {
	Rator Expr
	Rand  Expr
}

type BinExpr struct {
	op string
	e1 Expr
	e2 Expr
}

type MapEnv map[Symbol]Value

func deepcopyMap(originalMap MapEnv) MapEnv {
	newMap := make(MapEnv)
	for k, v := range originalMap {
		newMap[k] = v
	}
	return newMap
}

func (MapEnv) empty_env() Env {
	return make(MapEnv)
}
func (me MapEnv) extend(k Symbol, v Value) Env {
	ne := deepcopyMap(me)
	ne[k] = v
	return ne
}
func (me MapEnv) lookup(k Symbol) Value {
	return me[k]
}

func (s Symbol) interp(e Env) Value {
	return e.lookup(s)
}

func (n Number) interp(e Env) Value {
	return Value{n.Number}
}

type Closure struct {
	Lambda Lambda
	Env    Env
}

func (l Lambda) interp(e Env) Value {
	return Value{Closure{l, e}}
}

func (app App) interp(e Env) Value {
	c1 := app.Rator.interp(e).Value.(Closure)
	v2 := app.Rand.interp(e)
	return c1.Lambda.Body.interp(c1.Env.extend(c1.Lambda.Arg, v2))
}

func (b BinExpr) interp(e Env) Value {
	r1 := b.e1.interp(e).Value.(float64)
	r2 := b.e2.interp(e).Value.(float64)
	switch b.op {
	case "+":
		return Value{r1 + r2}
	case "-":
		return Value{r1 - r2}
	case "*":
		return Value{r1 * r2}
	case "/":
		return Value{r1 / r2}

	default:
		return Value{nil}
	}
}

type StringExpr struct {
	Str string
}

func (s StringExpr) interp(Env) Value {
	return Value{s.Str}
}

func testInterpSymbol() {
	var me MapEnv
	env0 := me.empty_env()
	fmt.Println(env0)
	s1 := Symbol{"X"}
	e1 := env0.extend(s1, Value{213})
	fmt.Println(Symbol{"X"}.interp(e1))
	expr1 := App{
		Lambda{Symbol{"x"}, Symbol{"x"}},
		Number{255},
	}
	fmt.Println(expr1.interp(env0))
	// lambda.x.x+1 5
	expr2 := App{
		Lambda{Symbol{"x"},
			BinExpr{"+", Symbol{"x"}, Number{1}},
		},
		Number{5},
	}
	fmt.Println(expr2.interp(env0))
}

func main() {
	testInterpSymbol()
}
