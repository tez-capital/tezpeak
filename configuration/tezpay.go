package configuration

import (
	"github.com/tez-capital/tezpeak/constants"
	"github.com/trilitech/tzgo/tezos"
)

type PayoutWalletPreferences struct {
	BalanceWarningThreshold int64 `json:"balance_warning_threshold,omitempty"`
	BalanceErrorThreshold   int64 `json:"balance_error_threshold,omitempty"`
}

type TezpayModuleConfiguration struct {
	moduleConfigurationbase

	PayoutWallet            string                  `json:"payout_wallet,omitempty"`
	PayoutWalletPreferences PayoutWalletPreferences `json:"payout_wallet_preferences,omitempty"`
	ForceDryRun             bool                    `json:"force_dry_run,omitempty"`
}

func getDefaultTezpayModuleConfiguration() *TezpayModuleConfiguration {
	return &TezpayModuleConfiguration{
		moduleConfigurationbase: moduleConfigurationbase{
			Applications: map[string]string{
				"tezpay": constants.DEFAULT_TEZPAY_APP_PATH,
			},
		},
	}
}

func (c *TezpayModuleConfiguration) Hydrate() {

}

func (c *TezpayModuleConfiguration) Validate() error {
	if _, err := tezos.ParseAddress(c.PayoutWallet); err != nil {
		return constants.ErrInvalidPayoutWallet
	}

	if tezpayAppPath, ok := c.Applications["tezpay"]; !ok || tezpayAppPath == "" {
		return constants.ErrNoTezpayAppPath
	}

	return nil
}
