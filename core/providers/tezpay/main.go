package tezpay

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/tez-capital/tezpeak/configuration"
	"github.com/tez-capital/tezpeak/core/common"
)

type Status struct {
	Services common.ServicesStatus `json:"services,omitempty"`
	Wallet   WalletStatus          `json:"wallet,omitempty"`
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
		Services: common.ServicesStatus{},
	}
}

func SetupModule(ctx context.Context, configuration *configuration.TezpayModuleConfiguration, app *fiber.Group, statusChannel chan<- common.StatusUpdatedReport) error {
	err := setupTezpayProvider(configuration, app)
	if err != nil {
		return err
	}

	tezpayStatus := GetEmptyStatus()
	tezbakeStatusChannel := make(chan common.StatusUpdatedReport, 100)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case statusUpdate := <-tezbakeStatusChannel:
				data := statusUpdate.GetData()
				switch data := data.(type) {
				case common.ServicesStatus:
					servicesStatus := data
					tezpayStatus.Services = servicesStatus
				case WalletStatus:
					tezpayStatus.Wallet = data
				}

				statusChannel <- &StatusUpdate{
					Status: tezpayStatus,
				}
			}
		}
	}()

	common.StartServiceStatusProviders(ctx, configuration.Applications, tezbakeStatusChannel)
	startWalletStatusProviders(ctx, configuration.PayoutWallet, configuration.PayoutWalletPreferences, tezbakeStatusChannel)

	return nil
}
