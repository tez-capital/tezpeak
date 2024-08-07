package tezbake

import (
	"context"
	"maps"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tez-capital/tezpeak/configuration"
	"github.com/tez-capital/tezpeak/core/common"
)

type Status struct {
	Rights   RightsStatus                    `json:"rights,omitempty"`
	Services common.AplicationServicesStatus `json:"services,omitempty"`
	Bakers   BakersStatus                    `json:"bakers,omitempty"`
	Ledgers  LedgerStatus                    `json:"ledgers,omitempty"`
}

func (status *Status) Clone() *Status {
	return &Status{
		// no need to clone RightsStatus
		status.Rights,
		common.AplicationServicesStatus{
			Applications: maps.Clone(status.Services.Applications),
			Timestamp:    status.Services.Timestamp,
		},
		status.Bakers,  // no need to clone BakersStatus
		status.Ledgers, // no need to clone LedgerStatus
	}
}

func GetEmptyStatus() *Status {
	return &Status{
		Rights: RightsStatus{
			Level:  0,
			Rights: []BlockRights{},
		},
		Services: common.AplicationServicesStatus{
			Applications: make(map[string]common.ApplicationServices),
			Timestamp:    time.Now().Unix(),
		},
		Bakers: BakersStatus{
			Level:  0,
			Bakers: map[string]*BakerStakingStatus{},
		},
		Ledgers: LedgerStatus{
			Level: 0,
		},
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

func SetupModule(ctx context.Context, configuration *configuration.TezbakeModuleConfiguration, app *fiber.Group, statusChannel chan<- common.StatusUpdate) error {
	err := setupGovernanceProvider(configuration, app)
	if err != nil {
		return err
	}

	tezbakeStatus := GetEmptyStatus()
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
					tezbakeStatus.Services.Applications[application] = statusUpdate.Status
					tezbakeStatus.Services.Timestamp = time.Now().Unix()
				case *RightsStatusUpdate:
					tezbakeStatus.Rights = statusUpdate.RightsStatus
				case *BakersStatusUpdate:
					tezbakeStatus.Bakers = statusUpdate.BakersStatus
					// case *LedgerStatusUpdate:
					// TODO: LedgerStatusUpdate
				}

				statusChannel <- &StatusUpdate{
					Status: tezbakeStatus.Clone(),
				}
			}
		}
	}()

	startRightsStatusProviders(ctx, configuration.Bakers, configuration.RightsBlockWindow, tezbakeStatusChannel)
	setupBakerStatusProviders(ctx, configuration.Bakers, tezbakeStatusChannel)
	common.StartServiceStatusProviders(ctx, configuration.Applications, tezbakeStatusChannel)

	return nil
}
