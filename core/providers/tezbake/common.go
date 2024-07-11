package tezbake

import (
	"github.com/tez-capital/tezpeak/core/common"
)

var (
	rightsEventSource = common.NewEventSource[*RightsStatus](func(rights *RightsStatus) bool {
		return rights.Level > lastProcessedRightsHeight
	})
	lastProcessedRightsHeight = int64(0)
)

func init() {
	go rightsEventSource.Run()
}
