package tezpay

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tez-capital/tezpeak/configuration"
	"github.com/tez-capital/tezpeak/core/common"
	"golang.org/x/exp/maps"
)

type Status struct {
	Services common.AplicationServicesStatus `json:"services,omitempty"`
	Wallet   WalletStatus                    `json:"wallet,omitempty"`
}

func (status *Status) Clone() *Status {
	return &Status{
		Services: common.AplicationServicesStatus{
			Applications: maps.Clone(status.Services.Applications),
			Timestamp:    status.Services.Timestamp,
		},
		Wallet: status.Wallet,
	}
}

type StatusUpdate struct {
	Status *Status
}

func (statusUpdate *StatusUpdate) GetId() string {
	return "tezbake"
}

func (statusUpdate *StatusUpdate) GetData() any {
	return statusUpdate.Status
}

func GetEmptyStatus() *Status {
	return &Status{
		Services: common.AplicationServicesStatus{
			Applications: make(map[string]common.ApplicationServices),
			Timestamp:    time.Now().Unix(),
		},
	}
}

func SetupModule(ctx context.Context, configuration *configuration.TezpayModuleConfiguration, app *fiber.Group, statusChannel chan<- common.StatusUpdate) error {
	err := setupTezpayProvider(configuration, app)
	if err != nil {
		return err
	}

	tezpayStatus := GetEmptyStatus()
	tezbakeStatusChannel := make(chan common.StatusUpdate, 100)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case statusUpdate := <-tezbakeStatusChannel:
				switch statusUpdate := statusUpdate.(type) {
				case *common.ServicesStatusUpdate:
					application := statusUpdate.Application
					tezpayStatus.Services.Applications[application] = statusUpdate.Status
					tezpayStatus.Services.Timestamp = time.Now().Unix()
				case *WalletBalanceUpdate:
					tezpayStatus.Wallet = statusUpdate.Status
				}

				statusChannel <- &StatusUpdate{
					Status: tezpayStatus.Clone(),
				}
			}
		}
	}()

	common.StartServiceStatusProviders(ctx, configuration.Applications, tezbakeStatusChannel)
	startWalletStatusProviders(ctx, configuration.PayoutWallet, configuration.PayoutWalletPreferences, tezbakeStatusChannel)

	return nil
}
