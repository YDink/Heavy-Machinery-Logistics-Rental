package domain

type PhotoAngle string

const (
	PhotoAngleFront PhotoAngle = "FRONT"
	PhotoAngleRear  PhotoAngle = "REAR"
	PhotoAngleLeft  PhotoAngle = "LEFT"
	PhotoAngleRight PhotoAngle = "RIGHT"
	PhotoAngleCabin PhotoAngle = "CABIN"
	PhotoAngleMeter PhotoAngle = "METER"
)

var MandatoryPhotoAngles = []PhotoAngle{
	PhotoAngleFront,
	PhotoAngleRear,
	PhotoAngleLeft,
	PhotoAngleRight,
	PhotoAngleCabin,
	PhotoAngleMeter,
}
