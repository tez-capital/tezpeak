package tezpay

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/tez-capital/tezpeak/configuration"
	"github.com/tez-capital/tezpeak/core/common"
	"golang.org/x/exp/maps"
)

type Status struct {
	Services map[string]common.ApplicationServices `json:"services,omitempty"`
	Wallet   WalletStatus                          `json:"wallet,omitempty"`
}

func (status *Status) Clone() *Status {
	return &Status{
		Services: maps.Clone(status.Services),
		Wallet:   status.Wallet,
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
		Services: make(map[string]common.ApplicationServices),
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
					tezpayStatus.Services[application] = statusUpdate.Status
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
