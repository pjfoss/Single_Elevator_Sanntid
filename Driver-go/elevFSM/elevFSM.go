package elevFSM

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"Driver-go/orders"
	"fmt"
)

func RunElevFSM(numFloors int, myElevator elevator.Elevator, orderPanel [orders.ConstNumFloors][3]int) {

	for f := 0; f < numFloors; f++ {
		for b := 0; b < 3; b++ {
			elevio.SetButtonLamp(elevio.ButtonType(b), f, false)
		}
	}

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	//	go orders.UpdateOrders(&orderPanel, drv_buttons)
	//	go orders.PollPriorityOrder(priOrderChan, orderPanel, myElevator.GetCurrentFloor(), myElevator.GetDirection())

	for {
		select {

		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			orders.SetOrder(&orderPanel, a.Floor, int(a.Button), orders.OT_Order)
			myElevator.DriveTo(a)

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			myElevator.SetFloor(a)
			elevio.SetFloorIndicator(a)
			if a == numFloors-1 {
				myElevator.SetDirection(elevio.MD_Stop)
			} else if a == 0 {
				myElevator.SetDirection(elevio.MD_Stop)
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

		}
	}
}
