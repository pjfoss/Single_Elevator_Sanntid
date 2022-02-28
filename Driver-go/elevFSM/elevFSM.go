package elevFSM

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"Driver-go/orders"
	"fmt"
)

func RunElevFSM(numFloors int, myElevator elevator.Elevator, orderPanel [orders.ConstNumFloors][3]int) {

	priOrderChan := make(chan elevio.ButtonEvent)
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	go orders.UpdateOrders(&orderPanel, drv_buttons)
	go orders.PollPriorityOrder(priOrderChan, orderPanel, myElevator.GetCurrentFloor(), myElevator.GetDirection())

	for {
		select {

		case p := <-priOrderChan:
			if p.Floor != -1 {
				myElevator.DriveTo(p)
			}

		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)

		case a := <-drv_floors:
			myElevator.SetFloor(a)
			fmt.Printf("%+v\n", a)
			if myElevator.GetDirection() == elevio.MD_Up {
				elevio.SetButtonLamp(elevio.BT_HallUp, a, false)
			} else if myElevator.GetDirection() == elevio.MD_Down {
				elevio.SetButtonLamp(elevio.BT_HallDown, a, false)
			}
			if a == numFloors-1 {
				myElevator.SetDirection(elevio.MD_Down)
			} else if a == 0 {
				myElevator.SetDirection(elevio.MD_Up)
			}
			elevio.SetMotorDirection(myElevator.GetDirection())

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(myElevator.GetDirection())
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
		}
	}
}
