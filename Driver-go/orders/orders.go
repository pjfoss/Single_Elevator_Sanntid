package orders

import (
	"Driver-go/elevio"
	"time"
)

const ConstNumFloors int = 4

const (
	//OT = OrderType
	OT_NoOrder = 0
	OT_Order   = 1
)
const (
	//CT = CostType
	CT_DistanceCost        = 10
	CT_DirSwitchCost       = 100
	CT_DoubleDirSwitchCost = 1000
)

func UpdateOrders(orderPanel *[ConstNumFloors][3]int, receiver <-chan elevio.ButtonEvent) {
	//Updates orderPanel matrix when receiver channel gets button calls
	for {
		orders := <-receiver
		orderPanel[orders.Floor][orders.Button] = OT_Order
		elevio.SetButtonLamp(orders.Button, orders.Floor, true)
	}
}

func calculateOrderCost(order elevio.ButtonEvent, elevFloor int, elevDirection elevio.MotorDirection) int {
	// Based on costed scenarios: on the order floor,above or below floor, type of requirede turns - calculate the cost of the given order
	var cost int = 0
	orderFloor := order.Floor
	if elevFloor == orderFloor {
		return cost
	}
	orderDirection := 0
	if elevFloor < orderFloor {
		orderDirection = int(elevio.MD_Up)
	} else if elevFloor > orderFloor {
		orderDirection = int(elevio.MD_Down)
	}
	newDirection := orderDirection
	if order.Button == elevio.BT_HallUp {
		newDirection = int(elevio.MD_Up)
	} else if order.Button == elevio.BT_HallDown {
		newDirection = int(elevio.MD_Down)
	}

	cost += CT_DistanceCost * intAbs(orderFloor-elevFloor)

	if orderDirection != int(elevDirection) {
		cost += CT_DirSwitchCost
		if newDirection != orderDirection {
			cost += CT_DoubleDirSwitchCost
		}
	} else if newDirection != orderDirection {
		cost += 0.8 * CT_DirSwitchCost
	}

	return cost
}

func priorityOrder(orderPanel *[ConstNumFloors][3]int, elevFloor int, elevDirection elevio.MotorDirection) elevio.ButtonEvent {
	//Calculate for given elevator which order it should take using calculateOrderCost for each current order.
	//fmt.Printf("Y00")
	var priorityOrder elevio.ButtonEvent = elevio.ButtonEvent{
		Floor:  -1,
		Button: -1,
	}
	var minCost int = 10000 //change to infinity <3
	for floor := 0; floor < len(orderPanel); floor++ {
		for btn := 0; btn < len(orderPanel[0]); btn++ {
			if orderPanel[floor][btn] != OT_NoOrder {
				order := elevio.ButtonEvent{
					Floor:  floor,
					Button: elevio.ButtonType(orderPanel[floor][btn]),
				}
				orderCost := calculateOrderCost(order, elevFloor, elevDirection)
				if orderCost < minCost {
					minCost = orderCost
					priorityOrder = order
				}
			}
		}
	}
	//fmt.Println(string(priorityOrder.Floor))
	return priorityOrder
}

func PollPriorityOrder(priOrderChan chan elevio.ButtonEvent, orderPanel [ConstNumFloors][3]int, elevFloor int, elevDirection elevio.MotorDirection) elevio.ButtonEvent {
	for {
		//fmt.Printf("Yoo")
		priOrderChan <- priorityOrder(&orderPanel, elevFloor, elevDirection)
		time.Sleep(time.Millisecond)
	}
}

func intAbs(x int) int {
	if x < 0 {
		x = -x
	}
	return x
}