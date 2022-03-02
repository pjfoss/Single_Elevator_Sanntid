package elevator

import (
	"Driver-go/elevio"
)

type Elevator struct {
	direction    elevio.MotorDirection
	currentFloor int
}

func (e *Elevator) GetDirection() elevio.MotorDirection {
	return e.direction
}
func (e *Elevator) GetCurrentFloor() int {
	return e.currentFloor
}
func (e *Elevator) SetFloor(floor int) {
	e.currentFloor = floor
}
func (e *Elevator) SetDirection(dir elevio.MotorDirection) {
	e.direction = dir
}

func (e *Elevator) DriveTo(order elevio.ButtonEvent) {
	floor := order.Floor
	dir := elevio.MotorDirection(elevio.MD_Stop)
	if e.GetCurrentFloor() < floor {
		dir = elevio.MD_Up
		e.SetDirection(dir)
	} else if e.GetCurrentFloor() > floor {
		dir = elevio.MD_Down
		e.SetDirection(dir)
	}
	elevio.SetMotorDirection(dir)
}
