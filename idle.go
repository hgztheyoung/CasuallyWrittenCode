package main
import (
	"time"
	"fmt"
	"github.com/nsf/termbox-go"
	keyboard "github.com/julienroland/keyboard-termbox"
)

var num int
var inc int

func main(){
	//Init termbox
	err := termbox.Init()
	if err != nil {
			panic(err)
	}
	defer termbox.Close()

	setUpKey()
	setUpGlobalGameVar()
	fmt.Printf("\t%v\t%v\t\t\n","num","inc")
	for {
		num += inc
		time.Sleep(200*time.Millisecond)		
		fmt.Printf("\r\t%v\t%v\t\t",num,inc)
	}
}

func setUpGlobalGameVar() {
	num = 0
	inc = 1
}

func setUpKey() {
	kb := keyboard.New()
	kb.Bind(func() { panic("Exiting!!!!!!!!!!!!!!!!!!!!!!") }, "escape", "q")
    kb.Bind(func() {
		inccount := 0
		if(num > (inc*10+inccount)){	
			num -= inc*10		
			inc += inc/(8+inccount) + 1
			inccount++
		}
	}, "up", "k")
    kb.Bind(func() { 
		fmt.Println("down")
	}, "down", "j")
    kb.Bind(func() {
		fmt.Println("left")
	}, "left", "h")
    kb.Bind(func() {
		fmt.Println("right")
	}, "right", "l")
	go func() {
		for{
			kb.Poll(termbox.PollEvent())
		}
	}()
}
