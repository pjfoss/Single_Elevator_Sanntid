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

	var priorityOrder elevio.ButtonEvent
	var doorOpen bool = false

	for f := 0; f < numFloors; f++ {
		for b := 0; b < 3; b++ {
			elevio.SetButtonLamp(elevio.ButtonType(b), f, false)
		}
	}

	priOrderChan := make(chan elevio.ButtonEvent, 2)
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
			fmt.Println(a)
			if !doorOpen {
				priorityOrder = a
				myElevator.DriveTo(a)
			}

		case a := <-drv_buttons:
			fmt.Println(a)
			orders.SetOrder(&orderPanel, a.Floor, int(a.Button), orders.OT_Order)
			elevio.SetButtonLamp(a.Button, a.Floor, true)

		case a := <-drv_floors:
			fmt.Println(a)
			myElevator.SetFloor(a)
			elevio.SetFloorIndicator(a)
			event := elevio.ButtonEvent{
				Floor:  a,
				Button: elevio.BT_Cab,
			}
			if myElevator.GetDirection() == elevio.MD_Up {
				event.Button = elevio.BT_HallUp
			} else if myElevator.GetDirection() == elevio.MD_Down {
				event.Button = elevio.BT_HallDown
			}
			if priorityOrder.Floor == a && (priorityOrder.Button == event.Button || priorityOrder.Button == elevio.BT_Cab) {
				doorOpen = true
				myElevator.SetDirection(elevio.MD_Stop)
				time.Sleep(3 * time.Second)
				orders.SetOrder(&orderPanel, a, int(elevio.BT_HallUp), orders.OT_NoOrder)
				elevio.SetButtonLamp(event.Button, a, false)
				elevio.SetButtonLamp(elevio.BT_Cab, a, false)
				doorOpen = false
			}
			if myElevator.GetCurrentFloor() == 0 || myElevator.GetCurrentFloor() == numFloors-1 {
				myElevator.SetDirection(elevio.MD_Stop)
			}

		case a := <-drv_obstr:
			fmt.Println(a)

		case a := <-drv_stop:
			fmt.Println(a)
		}
	}
}
