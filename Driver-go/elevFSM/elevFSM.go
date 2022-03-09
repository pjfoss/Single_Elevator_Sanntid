package elevFSM

import (
	"Driver-go/elevator"
	"Driver-go/elevio"
	"Driver-go/orders"
	"fmt"
	"time"
)

func RunElevFSM(numFloors int, myElevator elevator.Elevator, orderPanel [orders.ConstNumFloors][3]int) {

	init := 0
	myElevator.SetFloor(numFloors + 1)

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

		case p := <-priOrderChan:
			myElevator.DriveTo(p)

		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			orders.SetOrder(&orderPanel, a.Floor, int(a.Button), orders.OT_Order)

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			myElevator.SetFloor(a)
			elevio.SetFloorIndicator(a)

			if init == 0 {
				elevio.SetMotorDirection(elevio.MD_Stop)
				init = 1
			}

			bt := elevio.BT_HallDown
			if myElevator.GetDirection() == elevio.MD_Up {
				bt = int(elevio.BT_HallUp)
			}
			if orders.GetOrder(&orderPanel, a, bt) == orders.OT_Order {
				elevio.SetMotorDirection(elevio.MD_Stop)
				time.Sleep(3 * time.Second)
				orders.SetOrder(&orderPanel, a, bt, orders.OT_NoOrder)
				elevio.SetButtonLamp(elevio.ButtonType(bt), a, false)
				a {
					elevio.SetMotorDirection(elevio.MD_Stop)
				} else {
					elevio.SetMotorDirection(myElevator.GetDirection())
				}
	
			case a := <-drv_stop:
				fmt.Printf("%+v\n", a)
	
				// default:
				// 	if myElevator.GetCurrentFloor() > numFloors && init == 0 {
				// else if orders.GetOrder(&orderPanel, a, int(elevio.BT_Cab)) == orders.OT_Order {
				elevio.SetMotorDirection(elevio.MD_Stop)
				time.Sleep(3 * time.Second)
				orders.SetOrder(&orderPanel, a, int(elevio.BT_Cab), orders.OT_NoOrder)
				elevio.SetButtonLamp(elevio.BT_Cab, a, false)
			}

			if a == numFloors-1 {
				myElevator.SetDirection(elevio.MD_Down)
				elevio.SetMotorDirection(elevio.MD_Down)
			} else if a == 0 {
				myElevator.SetDirection(elevio.MD_Up)
				elevio.SetMotorDirection(elevio.MD_Up)
			}

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(myElevator.GetDirection())
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)

			// default:
			// 	if myElevator.GetCurrentFloor() > numFloors && init == 0 {
			// 		myElevator.SetDirection(elevio.MD_Down)
			// 		elevio.SetMotorDirection(elevio.MD_Down)
			// 	}
		}
	}
}
