package vmwarefusion

type machineType struct {
	id   string
	size string
	cpu  int
	ram  int
}

func (t *machineType) GetID() string {
	return t.id
}

var defaultType *machineType

func init() {
	defaultType = &machineType{
		id:   "default",
		size: "60GB",
		cpu:  2,
		ram:  4096,
	}
}
