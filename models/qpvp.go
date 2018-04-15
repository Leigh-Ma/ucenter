package models

const (
	RobotPlayerLeave = int8(-1)
	RobotNoPlayer = int8(0)
	RobotIsRobot = int8(1)

)

type QPvpRoundInfo struct {

}

//Qualifying game pvp record
type QPvpLog struct {
	TCom

	IdA int64
	IdD int64

	RobotA int8
	RobotB int8

	StaminaA int
	StaminaB int

	HasRobot bool

	RoundLeave int
	RoundTotal  int

	StaminaBonus int

	RoundInfo string //json
}

