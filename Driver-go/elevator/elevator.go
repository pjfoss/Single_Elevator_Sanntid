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
func (e Elevator) SetDirection(dir elevio.MotorDirection) {
	e.direction = dir
}

func (e *Elevator) DriveTo(order elevio.ButtonEvent) {
	floor := order.Floor
	dir := int(elevio.MD_Stop)
	if e.GetCurrentFloor() < floor {
		dir = int(elevio.MD_Up)
	} else {
		dir = int(elevio.MD_Down)
	}
	e.SetDirection(elevio.MotorDirection(dir))
	elevio.SetMotorDirection(elevio.MotorDirection(dir))
}
