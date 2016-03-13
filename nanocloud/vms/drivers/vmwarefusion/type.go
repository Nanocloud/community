package vmwarefusion

type machineType struct {
	id    string
	label string
	size  string
	cpu   int
	ram   int
}

func (t *machineType) Id() string {
	return t.id
}

func (t *machineType) Label() string {
	return t.label
}

var defaultType *machineType

func init() {
	defaultType = &machineType{
		id:    "default",
		label: "Default",
		size:  "60GB",
		cpu:   2,
		ram:   4096,
	}
}
