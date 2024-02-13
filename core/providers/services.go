package providers

import (
	"context"
	"fmt"
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

func StartServiceStatusProvider(ctx context.Context, tezbakeHome string, statusChannel chan<- common.ProviderStatusUpdatedReport) {
	if tezbakeHome == "" {
		slog.Warn("tezbake home not set, not starting service status provider")
		return
	}
	node := apps.NodeFromPath(path.Join(tezbakeHome, "node"))
	signer := apps.SignerFromPath(path.Join(tezbakeHome, "signer"))

	go func() {
		status := ServicesStatus{
			Timestamp:      time.Now().Unix(),
			NodeServices:   map[string]base.AmiServiceInfo{},
			SignerServices: map[string]base.AmiServiceInfo{},
		}

		for {
			nodeServiceInfoChannel := make(chan map[string]base.AmiServiceInfo)
			signerServiceInfoChannel := make(chan map[string]base.AmiServiceInfo)

			go func() {
				// recover
				defer func() {
					if r := recover(); r != nil {
						slog.Warn("recovered from panic", "error", r)
						nodeServiceInfoChannel <- map[string]base.AmiServiceInfo{}
					}
				}()

				nodeServiceInfo, err := node.GetServiceInfo()
				if err != nil {
					slog.Warn("failed to get node service info", err)
					nodeServiceInfoChannel <- map[string]base.AmiServiceInfo{}
					return
				}
				nodeServiceInfoChannel <- nodeServiceInfo
			}()

			go func() {
				defer func() {
					if r := recover(); r != nil {
						slog.Warn("recovered from panic", "error", r)
						nodeServiceInfoChannel <- map[string]base.AmiServiceInfo{}
					}
				}()

				signerServiceInfo, err := signer.GetServiceInfo()
				if err != nil {
					slog.Warn("failed to get node service info", err)
					signerServiceInfoChannel <- map[string]base.AmiServiceInfo{}
					return
				}
				signerServiceInfoChannel <- signerServiceInfo
			}()

			status.SignerServices = <-signerServiceInfoChannel
			status.NodeServices = <-nodeServiceInfoChannel
			fmt.Println("status", time.Now().Unix())
			status.Timestamp = time.Now().Unix()

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
