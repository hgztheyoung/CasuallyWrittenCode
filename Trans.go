package handy

import (
	"testing"
	"fmt"
	"time"
	"math/rand"
)

type AskBid int
type MarketLimit int

const (
	Ask AskBid = iota
	Bid
)

const (
	Market MarketLimit = iota
	Limit
)

type Order struct {
	AskBid      AskBid
	MarketLimit MarketLimit
	Amount      int
	Price       float64 //used by Market Order
}

var Orders chan Order

func OrderBinarySearch(array []Order, first, last int, value Order) int {
	//deal with Market Order
	if value.MarketLimit == Market {
		if value.AskBid == Ask {
			return 0
		}
		if value.AskBid == Bid {
			return last
		}
	}

	for first < last {
		mid := first + (last-first)/2
		if array[mid].Price < value.Price {
			first = mid + 1
		} else {
			last = mid
		}
	}
	return first
}

func InsertOrder(orders []Order, order Order) []Order {
	i := OrderBinarySearch(orders, 0, len(orders), order)
	orders = append(orders, Order{})
	copy(orders[i+1:], orders[i:])
	orders[i] = order
	return orders
}

func BrokerMainLoop() {
	Asks := make([]Order, 0)
	Bids := make([]Order, 0)
	for {
		select {
		case order := <-Orders:
			switch order.AskBid {
			case Ask:
				Asks = InsertOrder(Asks, order)
			case Bid:
				Bids = InsertOrder(Bids, order)
			}
		}
		//trade off the Asks and Bids
		fmt.Println(len(Asks), len(Bids), len(Orders))

		for len(Asks) > 0 && len(Bids) > 0 &&
			(Asks[0].Price <= Bids[len(Bids)-1].Price ||
				Asks[0].MarketLimit == Market ||
				Bids[len(Bids)-1].MarketLimit == Market) {
			if Asks[0].Amount < Bids[len(Bids)-1].Amount {
				Bids[len(Bids)-1].Amount = Bids[len(Bids)-1].Amount - Asks[0].Amount
				fmt.Println("Deal!!!", Asks[0], " -> ", Bids[len(Bids)-1])
				Asks = Asks[1:]
				continue
			}
			if Asks[0].Amount > Bids[len(Bids)-1].Amount {
				Asks[0].Amount = Asks[0].Amount - Bids[len(Bids)-1].Amount
				fmt.Println("Deal!!!", Asks[0], " -> ", Bids[len(Bids)-1])
				Bids = Bids[:len(Bids)-1]
				continue
			}
			if Asks[0].Amount == Bids[len(Bids)-1].Amount {
				fmt.Println("Deal!!!", Asks[0], " -> ", Bids[len(Bids)-1])
				Asks = Asks[1:]
				Bids = Bids[:len(Bids)-1]
				continue
			}
		}
		if len(Asks) > 5 && len(Bids) > 5 {
			fmt.Println("Loweast 5 Asks", Asks[0:5])
			fmt.Println("Highest 5 Bids", Bids[len(Bids)-5:])
			fmt.Println("-----------------------------------")
		}
	}
}

func DealerMainLoop() {
	for {
		select {
		case <-time.After(time.Millisecond * 15):
			//create dummy orders and send to order chan
			order := Order{AskBid: AskBid(rand.Intn(2)), MarketLimit: MarketLimit(rand.Intn(2)), Amount: rand.Intn(1000), Price: float64(rand.Intn(1000))}
			if order.MarketLimit == Market {
				order.Price = 0
			}
			Orders <- order
		}
	}
}

func TestPlayWithChan(t *testing.T) {
	//Init Orders
	Orders = make(chan Order, 100000)
	go BrokerMainLoop()
	for i := 0; i < 1000; i++ {
		go DealerMainLoop()
	}
	select {}
}
