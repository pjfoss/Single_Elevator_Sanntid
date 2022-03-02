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

	myElevator.SetFloor(-1)
	elevio.Init("localhost:15657", numFloors)


	elevFSM.RunElevFSM(numFloors, myElevator, orderPanel)
}
