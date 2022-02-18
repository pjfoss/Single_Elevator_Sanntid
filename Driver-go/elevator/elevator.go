package elevator

import (
	"Driver-go/elevio"
)

type Elevator struct {
	direction    elevio.MotorDirection
	currentFloor int
}

func (e Elevator) GetDirection() elevio.MotorDirection {
	return e.direction
}
func (e Elevator) GetCurrentFloor() int {
	return e.currentFloor
}
func (e Elevator) SetFloor(floor int) {
	e.currentFloor = floor
}

func (e Elevator) driveTo(floor int) {
	if e.GetCurrentFloor() == floor {
		return
	}
	if e.GetCurrentFloor() < floor {
		elevio.SetMotorDirection(elevio.MD_Up)
	} else {
		elevio.SetMotorDirection(elevio.MD_Up)
	}

}
