package main

import (
	"fmt"
	"math"
	"encoding/json"
)

var oldFileStr = []byte(`[
  {
    "key": 1,
    "value": "what"
  },
  {
    "key": 6,
    "value": "ever"
  },
  {
    "key": 10,
    "value": "value"
  },
  {
    "key": 4294967295,
    "value": "sentry"
  }
]`)

var transactionFileStr = []byte(`[
  {
    "key": 1,
    "value": "???",
    "operation": "update"
  },
  {
    "key": 1,
    "value": "lastUpdate",
    "operation": "update"
  },
  {
    "key": 5,
    "value": "insertV",
    "operation": "insert"
  },
  {
    "key": 7,
    "value": "insert7",
    "operation": "insert"
  },
  {
    "key": 10,
    "value": "lastUpdate",
    "operation": "delete"
  },
  {
    "key": 4294967295,
    "value": "sentry",
	"operation": "set_abnormal"
  }
]`)

type Record struct {
	Key   uint32
	Value string
}

func (r *Record) norm() bool {
	return r.Key < math.MaxUint32
}

func (t *Transaction) norm() bool {
	return t.Key < math.MaxUint32
}

func (t *Transaction) hasOp(s string) bool {
	return t.Operation == s
}

type Transaction struct {
	Record
	Operation string
}

func (x *Record) Update(y Transaction) {
	if x.norm() && x.Key == y.Key && y.hasOp("update") {
		//	do update
		x.Value = y.Value
	}
}

func (x *Record) Delete(y Transaction) {
	if x.norm() && x.Key == y.Key && y.hasOp("delete") {
		//	do update
		x.Key = math.MaxUint32
	}
}

func (x *Record) Insert(y Transaction) {
	if !x.norm() && y.hasOp("insert") {
		x.Key = y.Key
		x.Value = y.Value
	}
}

func (x *Record) SetAbnorm() {
	x.Key = math.MaxUint32
}

func execTrsactions(oldFile []Record, transactions []Transaction) (newFile []Record) {
	xf, yf := 0, 0
	x, y := oldFile[xf], transactions[yf]
	for x.norm() || y.norm() {
		var ckey uint32
		var xx Record
		if x.Key <= y.Key {
			ckey = x.Key
			xx.Key = x.Key
			xx.Value = x.Value
			xf++
			x = oldFile[xf]
		} else {
			ckey = y.Key
			xx.SetAbnorm()
		}
		for y.Key == ckey {
			if y.hasOp("update") && xx.norm() {
				xx.Update(y)
			}
			if y.hasOp("delete") && xx.norm() {
				xx.Delete(y)
			}
			if y.hasOp("insert") && !xx.norm() {
				xx.Insert(y)
			}
			if y.hasOp("insert") == xx.norm() {
				fmt.Errorf("error case")
			}
			yf++
			y = transactions[yf]
		}
		if xx.norm() {
			newFile = append(newFile, xx)
		}
	}
	newFile = append(newFile, x)
	return newFile
}

func main() {
	var oldFile []Record
	json.Unmarshal(oldFileStr, &oldFile)
	fmt.Println(oldFile)

	var transactions []Transaction
	json.Unmarshal(transactionFileStr, &transactions)
	fmt.Println(transactions)

	newfile := execTrsactions(oldFile, transactions)
	fmt.Println(newfile)
}
