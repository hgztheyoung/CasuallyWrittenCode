package main
import (
	"time"
	"fmt"
	// "fmt"
	// "time"
	"github.com/nsf/termbox-go"
	keyboard "github.com/julienroland/keyboard-termbox"
)

var num int
var inc int
var nextIncFee int
var nextInc int
func main(){
	//Init termbox
	err := termbox.Init()
	if err != nil {
			panic(err)
	}
	defer termbox.Close()

	setUpKey() //using type assertion
	setUpGlobalGameVar()
	fmt.Printf("\r\t%10v\t%10v\t\t%10v\t\t%10v\t\t\n","num","inc","nextInc","nextIncFee")
    for {
		num += inc
		time.Sleep(400*time.Millisecond)		
		fmt.Printf("\r\t%10v\t%10v\t\t%10v\t\t%10v\t\t",num,inc,nextInc,nextIncFee)
    }
}

func setUpGlobalGameVar() {
	num = 0
	inc = 1
	nextIncFee = 10
	nextInc = 1
}

func setUpKey() {
	kb := keyboard.New()
	kb.Bind(func() { panic("Exiting!!!!!!!!!!!!!!!!!!!!!!") }, "escape", "q")
    kb.Bind(func() {
		inccount := 0
		if(num > nextIncFee){	
			num -= inc*10		
			inc += nextInc
			inccount++
			nextIncFee = inc*10+inccount
			nextInc = inc/(8+inccount) + 1
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
