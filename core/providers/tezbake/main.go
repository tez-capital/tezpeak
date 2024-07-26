package tezbake

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/tez-capital/tezpeak/configuration"
	"github.com/tez-capital/tezpeak/core/common"
)

type Status struct {
	Governance GovernanceStatus      `json:"governance,omitempty"`
	Rights     RightsStatus          `json:"rights,omitempty"`
	Services   common.ServicesStatus `json:"services,omitempty"`
	Bakers     BakersStatus          `json:"bakers,omitempty"`
	Ledgers    LedgerStatus          `json:"ledgers,omitempty"`
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
		Rights: RightsStatus{
			Level:  0,
			Rights: []*BlockRights{},
		},
		Services: common.ServicesStatus{},
		Bakers: BakersStatus{
			Level:  0,
			Bakers: map[string]*BakerStakingStatus{},
		},
		Ledgers: LedgerStatus{
			Level: 0,
		},
	}
}

func SetupModule(ctx context.Context, configuration *configuration.TezbakeModuleConfiguration, app *fiber.Group, statusChannel chan<- common.StatusUpdatedReport) error {
	err := setupGovernanceProvider(configuration, app)
	if err != nil {
		return err
	}

	tezbakeStatus := GetEmptyStatus()
	tezbakeStatusChannel := make(chan common.StatusUpdatedReport, 100)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case statusUpdate := <-tezbakeStatusChannel:
				data := statusUpdate.GetData()
				switch data := data.(type) {
				case RightsStatus:
					rightsStatus := data
					tezbakeStatus.Rights = rightsStatus
				case common.ServicesStatus:
					servicesStatus := data
					tezbakeStatus.Services = servicesStatus
				case BakersStatus:
					bakersStatus := data
					tezbakeStatus.Bakers = bakersStatus
				case LedgerStatus:
					// not implemented
				}

				statusChannel <- &StatusUpdate{
					Status: tezbakeStatus,
				}
			}
		}
	}()

	startRightsStatusProviders(ctx, configuration.Bakers, configuration.RightsBlockWindow, tezbakeStatusChannel)
	setupBakerStatusProviders(ctx, configuration.Bakers, tezbakeStatusChannel)
	common.StartServiceStatusProviders(ctx, configuration.Applications, tezbakeStatusChannel)

	return nil
}
