package common

type StatusUpdatedReport interface {
	GetId() string
	GetData() any
}

type ModuleStatusUpdatedReport struct {
	Id     string
	Report StatusUpdatedReport
}
