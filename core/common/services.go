package common

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/tez-capital/tezbake/ami"
	"github.com/tez-capital/tezbake/apps/base"
	"github.com/tez-capital/tezpeak/constants"
)

type ApplicationServices map[string]base.AmiServiceInfo

type ServicesStatus struct {
	Timestamp    int64                          `json:"timestamp"`
	Applications map[string]ApplicationServices `json:"applications"`
}

type ServicesStatusUpdate struct {
	Status ServicesStatus
}

func (s *ServicesStatusUpdate) GetId() string {
	return "services"
}

func (s *ServicesStatusUpdate) GetData() any {
	return s.Status
}

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
		slog.Warn("failed to get node service info", "error", err.Error())
		return map[string]base.AmiServiceInfo{}
	}

	return result
}

func startApplicationServiceStatusProvider(ctx context.Context, application string, serviceStatusChannel chan<- map[string]base.AmiServiceInfo) {
	toSleep := time.Duration(0)

	go func() {
		for {
			serviceStatus := getApplicationServiceStatus(ctx, application)
			for _, serviceStatus := range serviceStatus {
				if toSleep > 0 {
					break
				}
				if serviceStatus.Status != "running" {
					toSleep = time.Second * constants.MIN_SERVICES_REFRESH_INTERVAL
				} else if started, err := time.Parse("Mon 2006-01-02 15:04:05 UTC", serviceStatus.Started); err == nil {
					diff := time.Since(started)
					toSleep = max(time.Second*constants.MIN_SERVICES_REFRESH_INTERVAL, min(diff, constants.MAX_SERVICES_REFRESH_INTERVAL))
				} else {
					toSleep = time.Second * constants.MIN_SERVICES_REFRESH_INTERVAL
				}
			}
			serviceStatusChannel <- serviceStatus
			time.Sleep(max(toSleep, time.Second*constants.MIN_SERVICES_REFRESH_INTERVAL))
		}
	}()

}

func StartServiceStatusProviders(ctx context.Context, applications map[string]string, statusChannel chan<- StatusUpdatedReport) {
	go func() {
		status := ServicesStatus{
			Timestamp:    time.Now().Unix(),
			Applications: map[string]ApplicationServices{},
		}

		for name, path := range applications {
			status.Applications[name] = ApplicationServices{}

			serviceStatusChannel := make(chan map[string]base.AmiServiceInfo)
			startApplicationServiceStatusProvider(ctx, path, serviceStatusChannel)

			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					case serviceStatus := <-serviceStatusChannel:
						status.Applications[name] = serviceStatus
						statusChannel <- &ServicesStatusUpdate{
							Status: status,
						}
					}
				}
			}()

		}
	}()

}
