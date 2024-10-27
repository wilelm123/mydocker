package subsystems

type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

type Subsystem interface {
	Name() string
	Set(patch string, res *ResourceConfig) error
	Apply(patch string, pid int) error
	Remove(path string) error
}

var (
	SubsystemsIns = []Subsystem{
		&CpuSubSystem{},
		&MemorySubSystem{},
		&CpusetSubSystem{},
	}
)
