package vms

import "errors"

var drivers map[string]Driver

func Register(name string, driver Driver) {
	if drivers == nil {
		drivers = make(map[string]Driver, 0)
	}
	drivers[name] = driver
}

func Open(driverName string, options map[string]string) (*VM, error) {
	driver := drivers[driverName]
	if driver != nil {
		lala, _ := driver.Open(options)
		return &lala, nil
	}
	return nil, errors.New("Invalid driver name")
}
