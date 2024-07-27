package common

type StatusUpdate interface {
	GetId() string
	GetData() any
}

type ModuleStatusUpdate interface {
	StatusUpdate
	GetModule() string
	GetStatusUpdate() StatusUpdate
}

type moduleStatusUpdate struct {
	StatusUpdate
	module string
}

func (m *moduleStatusUpdate) GetModule() string {
	return m.module
}

func (m *moduleStatusUpdate) GetStatusUpdate() StatusUpdate {
	return m.StatusUpdate
}

func NewModuleStatusUpdate(module string, statusUpdate StatusUpdate) ModuleStatusUpdate {
	return &moduleStatusUpdate{
		StatusUpdate: statusUpdate,
		module:       module,
	}
}
