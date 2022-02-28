package main

import (
	"Driver-go/elevFSM"
	"Driver-go/elevator"
	"Driver-go/elevio"
	"Driver-go/orders"
)

func main() {
	numFloors := 4
	var myElevator elevator.Elevator
	var orderPanel [orders.ConstNumFloors][3]int


	elevio.Init("localhost:15657", numFloors)

	myElevator.SetDirection(elevio.MD_Up)
	//elevio.SetMotorDirection(d)
	//elevio.SetMotorDirection(d)
	//helt nytt, helt nytt igjen

	elevFSM.RunElevFSM(numFloors, myElevator, orderPanel)
}
