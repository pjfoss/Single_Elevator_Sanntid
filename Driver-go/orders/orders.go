package orders

import "Driver-go/elevio"

const ConstNumFloors int = 4

const (
	OT_NoOrder = 0
	OT_Order   = 1
)
const (
	CT_DistanceCost        = 10
	CT_DirSwitchCost       = 100
	CT_DoubleDirSwitchCost = 1000
)

func UpdateOrders(orderPanel [ConstNumFloors][3]int, receiver <-chan elevio.ButtonEvent) {
	for {
		orders := <-receiver
		orderPanel[orders.Floor][orders.Button] = OT_Order
	}
}

func calculateOrderCost(order elevio.ButtonEvent, elevFloor int, elevDirection elevio.MotorDirection) int {
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

func PriorityOrder(orderPanel [ConstNumFloors][3]int, elevFloor int, elevDirection elevio.MotorDirection) elevio.ButtonEvent {
	var priorityOrder elevio.ButtonEvent
	var minCost int = 10000
	for n := 0; n < len(orderPanel); n++ {
		for i := 0; i < len(orderPanel[0]); i++ {
			if orderPanel[n][i] != OT_NoOrder {
				order := elevio.ButtonEvent{
					Floor:  n,
					Button: elevio.ButtonType(orderPanel[n][i]),
				}
				orderCost := calculateOrderCost(order, elevFloor, elevDirection)
				if orderCost < minCost {
					minCost = orderCost
					priorityOrder = order
				}
			}
		}
	}
	return priorityOrder
}

func intAbs(x int) int {
	if x < 0 {
		x = -x
	}
	return x
}
