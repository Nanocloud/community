package test

import (
	"log"
	"testing"
	"time"

	connector "github.com/Nanocloud/community/nanocloud/connectors/vms"
	"github.com/Nanocloud/community/nanocloud/vms"
)

var (
	machine_id string = ""
)

func getMachine() vms.Machine {
	machine, err := connector.Machine(machine_id)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if machine == nil {
		log.Fatalln("Machine is nil")
	}
	return machine
}

func TestOpen(t *testing.T) {
	v, err := vms.Open("test", nil)
	if err != nil {
		t.Fatal(err)
	}

	types, err := (*v).Types()
	if err != nil {
		t.Fatal(err)
	}

	if len(types) < 1 {
		t.Fatalf("empty type list retured")
	}

	connector.SetVM(v)
}

func TestType(t *testing.T) {
	machine_type, err := connector.Type("default-test-machine-type")
	if err != nil {
		t.Fatal(err)
	}
	if machine_type == nil {
		t.Errorf("Machine type is nil")
	}
}

func TestTypes(t *testing.T) {
	v, err := connector.Types()
	if err != nil {
		t.Fatal(err)
	}
	if len(v) < 1 {
		t.Fatalf("empty list returned")
	}
}

func TestCreate(t *testing.T) {
	machine_attributes := vms.MachineAttributes{
		Type:     nil,
		Name:     "machine-test",
		Username: "admin",
		Password: "secret",
		Ip:       "127.0.0.1",
	}

	new_machine, err := connector.Create(machine_attributes)
	if err != nil {
		log.Panicln(err)
	}
	if new_machine == nil {
		t.Fatalf("Machine was not created")
	}

	machine_id = new_machine.Id()
}

func TestStart(t *testing.T) {
	machine := getMachine()

	err := machine.Start()
	if err != nil {
		t.Errorf(err.Error())
	}
	status, err := machine.Status()
	if err != nil {
		t.Errorf(err.Error())
	}
	if status != vms.StatusUp {
		t.Errorf("VM status should be up, it is: %v\n", status)
	}
}

func TestStop(t *testing.T) {
	machine := getMachine()

	err := machine.Stop()
	if err != nil {
		t.Errorf(err.Error())
	}
	status, err := machine.Status()
	if err != nil {
		t.Errorf(err.Error())
	}
	if status != vms.StatusStopping && status != vms.StatusDown {
		t.Errorf("VM status should be stopping or down, it is: %v\n", status)
	}
}

func TestTerminate(t *testing.T) {
	machine := getMachine()

	err := machine.Terminate()
	if err != nil {
		t.Errorf(err.Error())
	}
	status, err := machine.Status()
	if err != nil {
		t.Errorf(err.Error())
	}
	if status != vms.StatusTerminated {
		t.Errorf("VM status should be terminated, it is: %v\n", status)
	}
}

func TestSetNil(t *testing.T) {
	machine_attributes := vms.MachineAttributes{
		Type:     nil,
		Name:     "machine-test",
		Username: "admin",
		Password: "secret",
		Ip:       "127.0.0.1",
	}

	SetNil()
	new_machine, _ := connector.Create(machine_attributes)
	if new_machine != nil {
		t.Errorf("Create(): Returned value should be nil")
	}

	SetNil()
	machine_type, _ := connector.Type("default-test-machine-type")
	if machine_type != nil {
		t.Errorf("Type(): Returned value should be nil")
	}

	SetNil()
	all_machines, _ := connector.Types()
	if all_machines != nil {
		t.Errorf("Types(): Returned value should be nil")
	}
}

func TestSetFail(t *testing.T) {
	machine := getMachine()
	machine_attributes := vms.MachineAttributes{
		Type:     nil,
		Name:     "machine-test",
		Username: "admin",
		Password: "secret",
		Ip:       "127.0.0.1",
	}

	SetFail()
	_, err := connector.Create(machine_attributes)
	if err == nil {
		t.Errorf("Create(): Error should not be nil")
	}

	SetFail()
	_, err = connector.Type("default-test-machine-type")
	if err == nil {
		t.Errorf("Type(): Error should not be nil")
	}

	SetFail()
	_, err = connector.Types()
	if err == nil {
		t.Errorf("Types(): Error should not be nil")
	}

	SetFail()
	err = machine.Start()
	if err == nil {
		t.Errorf("Start(): Error should not be nil")
	}

	SetFail()
	err = machine.Stop()
	if err == nil {
		t.Errorf("Stop(): Error should not be nil")
	}

	SetFail()
	err = machine.Terminate()
	if err == nil {
		t.Errorf("Terminate(): Error should not be nil")
	}
}

func TestSetDelay(t *testing.T) {
	machine := getMachine()
	machine_attributes := vms.MachineAttributes{
		Type:     nil,
		Name:     "machine-test",
		Username: "admin",
		Password: "secret",
		Ip:       "127.0.0.1",
	}

	SetDelay(1000)
	called := time.Now()
	connector.Create(machine_attributes)
	if time.Since(called) < time.Duration(1000) {
		t.Errorf("Create(): Delay should be observed")
	}

	SetDelay(1000)
	called = time.Now()
	connector.Type("default-test-machine-type")
	if time.Since(called) < time.Duration(1000) {
		t.Errorf("Type(): Delay should be observed")
	}

	SetDelay(1000)
	called = time.Now()
	connector.Types()
	if time.Since(called) < time.Duration(1000) {
		t.Errorf("Types(): Delay should be observed")
	}

	SetDelay(1000)
	called = time.Now()
	machine.Start()
	if time.Since(called) < time.Duration(1000) {
		t.Errorf("Start(): Delay should be observed")
	}

	SetDelay(1000)
	called = time.Now()
	machine.Stop()
	if time.Since(called) < time.Duration(1000) {
		t.Errorf("Stop(): Delay should be observed")
	}

	SetDelay(1000)
	called = time.Now()
	machine.Terminate()
	if time.Since(called) < time.Duration(1000) {
		t.Errorf("Terminate(): Delay should be observed")
	}
}
