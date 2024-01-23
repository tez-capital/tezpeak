package providers

import (
	"context"
	"log/slog"
	"path"
	"time"

	"github.com/tez-capital/tezbake/apps"
	"github.com/tez-capital/tezbake/apps/base"
	"github.com/tez-capital/tezpeak/constants"
	"github.com/tez-capital/tezpeak/core/common"
)

type ServicesStatus struct {
	Timestamp      int64                          `json:"timestamp"`
	NodeServices   map[string]base.AmiServiceInfo `json:"node_services"`
	SignerServices map[string]base.AmiServiceInfo `json:"signer_services"`
}

type ServicesStatusUpdate struct {
	Status ServicesStatus
}

func (s *ServicesStatusUpdate) GetId() string {
	return "services"
}

func (s *ServicesStatusUpdate) GetData() interface{} {
	return s.Status
}

func (s *ServicesStatusUpdate) GetKind() common.StatusUpdateKind {
	return common.ServicesStatusUpdateKind
}

func StartServiceStatusProvider(ctx context.Context, workingDirectory string, statusChannel chan<- common.ProviderStatusUpdatedReport) {
	node := apps.NodeFromPath(path.Join(workingDirectory, "node"))
	signer := apps.SignerFromPath(path.Join(workingDirectory, "signer"))

	go func() {

		status := ServicesStatus{
			Timestamp:      time.Now().Unix(),
			NodeServices:   map[string]base.AmiServiceInfo{},
			SignerServices: map[string]base.AmiServiceInfo{},
		}

		for {
			nodeServiceInfo, err := node.GetServiceInfo()
			if err != nil {
				slog.Warn("failed to get node service info", err)
			}
			status.NodeServices = nodeServiceInfo

			signerServiceInfo, err := signer.GetServiceInfo()
			if err != nil {
				slog.Warn("failed to get node service info", err)
			}
			status.SignerServices = signerServiceInfo

			statusChannel <- &ServicesStatusUpdate{
				Status: status,
			}

			toSleep := time.Duration(0)

			// parse Sun 2024-01-28 11:49:40 UTC
			for _, serviceStatus := range status.NodeServices {
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

			for _, serviceStatus := range status.SignerServices {
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

			time.Sleep(max(toSleep, time.Second*constants.MIN_SERVICES_REFRESH_INTERVAL))
		}

	}()

}
