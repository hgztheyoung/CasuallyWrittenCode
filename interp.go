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

func ParseYaml2(input map[interface{}]interface{}) (expr Expr, err error) {
	_, nameok := input["name"]
	if nameok {
		expr, err = ParseSymbol(input)
		return
	}

	_, numberok := input["number"]
	if numberok {
		expr, err = ParseNumber(input)
		return
	}

	_, lok := input["lambda"]
	_, actualargok := input["actualarg"]
	if lok && actualargok {
		expr, err = ParseApp(input)
		return
	}
	_, argok := input["arg"]
	_, bodyok := input["body"]
	if argok && bodyok {
		expr, err = ParseLambda(input)
		return
	}
	_, opok := input["op"]
	_, e1ok := input["e1"]
	_, e2ok := input["e2"]
	if opok && e1ok && e2ok {
		expr, err = ParseBinOp(input)
		return
	}
	return nil, errors.New("ParseYaml2 failed")
}
func ParseLambda(input map[interface{}]interface{}) (expr Expr, err error) {
	arg := input["arg"]
	body := input["body"]
	var se Symbol
	if ma, ok := arg.(map[interface{}]interface{}); !ok {
		err = errors.New(fmt.Sprint("ParseSymbol failed", arg))
		return
	} else {
		if pse, perr := ParseSymbol(ma); perr != nil {
			err = perr
			return
		} else {
			se = pse
		}
	}
	bodye, err := ParseYaml2(body.(map[interface{}]interface{}))
	if err != nil {
		return
	}
	expr = Lambda{
		Arg:  se,
		Body: bodye,
	}
	return
}

func ParseSymbol(arg map[interface{}]interface{}) (expr Symbol, err error) {
	se := Symbol{}
	sn, ok := arg["name"].(string)
	if !ok {
		err = errors.New(fmt.Sprint("ParseSymbol failed", arg))
	}
	se.Name = sn
	expr = se
	return
}

func ParseNumber(arg map[interface{}]interface{}) (expr Number, err error) {
	ne := Number{}
	fsn, ok := arg["number"].(float64)
	if ok {
		ne.Number = fsn
		expr = ne
		return
	}
	sn, ok := arg["number"].(int)
	if ok {
		ne.Number = float64(sn)
		expr = ne
		return

	}
	err = errors.New(fmt.Sprint("ParseNumber failed", arg))
	return
}

func ParseBinOp(arg map[interface{}]interface{}) (expr BinExpr, err error) {
	sn, ok := arg["op"].(string)
	if !ok {
		err = errors.New(fmt.Sprint("ParseBinOp", arg["op"]))
	}
	e1e, err := ParseYaml2(arg["e1"].(map[interface{}]interface{}))
	if err != nil {
		return
	}
	e2e, err := ParseYaml2(arg["e2"].(map[interface{}]interface{}))
	if err != nil {
		return
	}
	expr = BinExpr{
		E1: e1e,
		E2: e2e,
		Op: sn,
	}
	return
}

func ParseApp(input map[interface{}]interface{}) (Expr, error) {
	l := input["lambda"]
	actualarg := input["actualarg"]
	le, e1 := ParseYaml2(l.(map[interface{}]interface{}))
	actualarge, e2 := ParseYaml2(actualarg.(map[interface{}]interface{}))
	log.Println(actualarge, e2)
	if e1 != nil {
		return nil, e1
	}
	if e2 != nil {
		return nil, e2
	}

	return App{
		Lambda:    le,
		ActualArg: actualarge,
	}, nil
}

func main2() {
	//testInterpSymbol()
	y := `add1expr: &f
    arg:
        name: x
    body:
      op: +
      e1:
        name: x
      e2:
        number: 1
lambda:
   <<: *f
actualarg:
  lambda:
     <<: *f
  actualarg:
    number: 5.5`
	var res map[interface{}]interface{}
	yaml.Unmarshal([]byte(y), &res)
	expr, err := ParseYaml2(res)
	fmt.Println(expr, err)
	fmt.Println(expr.interp(MapEnv{}.empty_env()))
}
