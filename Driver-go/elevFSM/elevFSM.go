package elevFSM

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"Driver-go/orders"
	"fmt"
	"time"
)

func RunElevFSM(numFloors int, myElevator elevator.Elevator, orderPanel [orders.ConstNumFloors][3]int) {

	elevio.SetMotorDirection(elevio.MD_Down)

	var doorOpen bool = false
	var moving bool = true

	var priorityOrder elevio.ButtonEvent
	priorityOrder.Floor = -1

	for f := 0; f < numFloors; f++ {
		for b := 0; b < 3; b++ {
			elevio.SetButtonLamp(elevio.ButtonType(b), f, false)
		}
	}

	priOrderChan := make(chan elevio.ButtonEvent)
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	go orders.PollPriorityOrder(priOrderChan, &orderPanel, &myElevator)

	for {
		select {

		case a := <-priOrderChan:
			// fmt.Println("pri >> PRIORITY: " + fmt.Sprint(priorityOrder))
			// fmt.Println("pri >> ELEVATOR DIRECTION: " + fmt.Sprint(myElevator.GetDirection()))

			priorityOrder = a
			if !doorOpen && priorityOrder.Floor != -1 {
				//drive to the order
				myElevator.DriveTo(priorityOrder)
				if !moving && priorityOrder.Floor == myElevator.GetCurrentFloor() {
					if priorityOrder.Button == elevio.BT_HallUp {
						myElevator.SetDirection(elevio.MD_Up)
					} else if priorityOrder.Button == elevio.BT_HallDown {
						myElevator.SetDirection(elevio.MD_Down)
					}
					// create button event corresponding to current elev state
					event := elevio.ButtonEvent{
						Floor:  myElevator.GetCurrentFloor(),
						Button: elevio.BT_Cab,
					}
					if myElevator.GetDirection() == elevio.MD_Up {
						event.Button = elevio.BT_HallUp
					} else if myElevator.GetDirection() == elevio.MD_Down {
						event.Button = elevio.BT_HallDown
					}
					//open doors
					doorOpen = true
					//fmt.Println("pri >> door open")
					//timer
					time.Sleep(3 * time.Second)
					//clear the orders
					orders.SetOrder(&orderPanel, event.Floor, int(event.Button), orders.OT_NoOrder)
					orders.SetOrder(&orderPanel, event.Floor, int(elevio.BT_Cab), orders.OT_NoOrder)
					//set priority to an invalid order
					priorityOrder.Floor = -1
					//open door
					doorOpen = false
					//fmt.Println("pri >> door closed")
				}
			}
			if priorityOrder.Floor != myElevator.GetCurrentFloor() && priorityOrder.Floor != -1 {
				//stop moving
				moving = true
				//fmt.Println("floor >> moving")
			} else {
				moving = false
				//fmt.Println("floor >> not moving")
			}

		case a := <-drv_buttons:
			orders.SetOrder(&orderPanel, a.Floor, int(a.Button), orders.OT_Order)

		case a := <-drv_floors:

			//switch direction if at top or bottom floor
			if myElevator.GetCurrentFloor() == 0 {
				myElevator.SetDirection(elevio.MD_Up)
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else if myElevator.GetCurrentFloor() == numFloors-1 {
				myElevator.SetDirection(elevio.MD_Down)
				elevio.SetMotorDirection(elevio.MD_Stop)
			}

			//fmt.Println("floor >> " + fmt.Sprint(orderPanel))
			//updage the floor
			myElevator.SetFloor(a)
			//turn on the floor light
			elevio.SetFloorIndicator(a)

			//if this floor has an order
			if priorityOrder.Floor != a && priorityOrder.Floor != -1 {
				//stop moving
				moving = true
				//fmt.Println("floor >> moving")
			} else {
				moving = false
				//fmt.Println("floor >> not moving")
			}

		case a := <-drv_obstr:
			fmt.Println(a)

		case a := <-drv_stop:
			fmt.Println(a)
		}
	}
}
