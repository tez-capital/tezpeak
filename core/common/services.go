package common

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/tez-capital/tezbake/ami"
	"github.com/tez-capital/tezbake/apps/base"
	"github.com/tez-capital/tezpeak/constants"
)

type ApplicationServices *map[string]base.AmiServiceInfo

type ServicesStatusUpdate struct {
	Application string
	Status      ApplicationServices
}

func (s *ServicesStatusUpdate) GetId() string {
	return s.Application
}

func (s *ServicesStatusUpdate) GetData() any {
	return s.Status
}

type registeredApplicationContext struct {
	refreshTrigger chan struct{}
	refreshPending sync.RWMutex
}

type applicationRegistry struct {
	apps map[string]*registeredApplicationContext
	mtx  sync.RWMutex
}

func (r *applicationRegistry) RegisterApplication(application string) <-chan struct{} {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if app, ok := r.apps[application]; ok {
		return app.refreshTrigger
	}

	refreshTrigger := make(chan struct{})
	r.apps[application] = &registeredApplicationContext{
		refreshTrigger: refreshTrigger,
		refreshPending: sync.RWMutex{},
	}

	return refreshTrigger
}

func (r *applicationRegistry) RefreshApplication(application string) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if app, ok := r.apps[application]; ok {
		if !app.refreshPending.TryLock() {
			return
		}
		defer app.refreshPending.Unlock()
		app.refreshTrigger <- struct{}{}
	}
}

func NewApplicationRegistry() *applicationRegistry {
	return &applicationRegistry{
		apps: make(map[string]*registeredApplicationContext),
	}
}

var (
	serviceRegisteredApplications = NewApplicationRegistry()
)

type GenericAmiApp struct{}

func GetServiceInfo(appPath string) (map[string]base.AmiServiceInfo, error) {
	infoBytes, _, err := ami.ExecuteInfo(appPath, "--services")
	if err != nil {
		return nil, fmt.Errorf("failed to collect app info (%s)", err.Error())
	}

	info, err := base.ParseInfoOutput(infoBytes)
	infoString, _ := json.Marshal(info["services"])
	result := map[string]base.AmiServiceInfo{}
	json.Unmarshal(infoString, &result)

	return result, err
}

func getApplicationServiceStatus(_ context.Context, application string) (result map[string]base.AmiServiceInfo) {
	defer func() {
		if r := recover(); r != nil {
			slog.Warn("recovered from panic", "error", r)
			result = map[string]base.AmiServiceInfo{}
		}
	}()

	result, err := GetServiceInfo(application)
	if err != nil {
		slog.Warn("failed to get service info", "error", err.Error(), "application", application)
		return map[string]base.AmiServiceInfo{}
	}

	return result
}

func refreshApplicationStatus(ctx context.Context, application string, serviceStatusChannel chan<- map[string]base.AmiServiceInfo) (toSleep time.Duration) {
	serviceStatus := getApplicationServiceStatus(ctx, application)
	for _, serviceStatus := range serviceStatus {
		if serviceStatus.Status != "running" {
			toSleep = time.Second * constants.MIN_SERVICES_REFRESH_INTERVAL
		} else if started, err := time.Parse("Mon 2006-01-02 15:04:05 UTC", serviceStatus.Started); err == nil {
			diff := time.Since(started)
			toSleep = max(time.Second*constants.MIN_SERVICES_REFRESH_INTERVAL, min(time.Second*diff, constants.MAX_SERVICES_REFRESH_INTERVAL))
		} else {
			toSleep = time.Second * constants.MIN_SERVICES_REFRESH_INTERVAL
		}
	}
	go func() {
		serviceStatusChannel <- serviceStatus
	}()
	return max(toSleep, time.Second*constants.MIN_SERVICES_REFRESH_INTERVAL)
}

func startApplicationServiceStatusProvider(ctx context.Context, application string, serviceStatusChannel chan<- map[string]base.AmiServiceInfo) {
	toSleep := time.Duration(0)
	refreshTrigger := serviceRegisteredApplications.RegisterApplication(application)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(toSleep):
				toSleep = refreshApplicationStatus(ctx, application, serviceStatusChannel)
			case <-refreshTrigger:
				toSleep = refreshApplicationStatus(ctx, application, serviceStatusChannel)
			}
		}
	}()
}

func StartServiceStatusProviders(ctx context.Context, applications map[string]string, statusChannel chan<- StatusUpdate) {
	for name, path := range applications {
		serviceStatusChannel := make(chan map[string]base.AmiServiceInfo)
		startApplicationServiceStatusProvider(ctx, path, serviceStatusChannel)

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case serviceStatus := <-serviceStatusChannel:
					statusChannel <- &ServicesStatusUpdate{
						Application: name,
						Status:      &serviceStatus,
					}
				}
			}
		}()
	}
}

func UpdateServiceStatus(application string) {
	serviceRegisteredApplications.RefreshApplication(application)
}
