package tezbake

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tez-capital/tezbake/ami"
	"github.com/tez-capital/tezbake/apps/base"
	"github.com/tez-capital/tezbake/apps/signer"
	"github.com/tez-capital/tezpeak/core/common"
)

type WalletsStatus map[string]base.AmiWalletInfo

type WalletsStatusUpdate struct {
	WalletsStatus
}

func (s *WalletsStatus) GetId() string {
	return "wallets"
}

func (s *WalletsStatus) GetData() any {
	return s
}

func (s *WalletsStatus) DisconnectedWallets() []string {
	result := make([]string, 0, len(*s))
	for id, wallet := range *s {
		if wallet.Kind != "ledger" {
			continue
		}
		if wallet.LedgerStatus == "disconnected" {
			result = append(result, id)
		}
	}
	return result
}

const (
	LEDGER_VENDOR_IDS = "2c97,2581"
)

type ArcEvent struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Fields    struct {
		Message   string `json:"message"`
		Path      string `json:"path"`
		VendorID  string `json:"vendor_id"`
		ProductID string `json:"product_id"`
	} `json:"fields"`
}

type CheckLedgerResult struct {
	Ledger     string `json:"ledger"`
	DevicePath string `json:"device_path"`
	Path       string `json:"path"`
	AppVersion string `json:"app_version"`
}

func RunArcMonitor(ctx context.Context, arcPath string) <-chan ArcEvent {
	outputChannel := make(chan ArcEvent)

	go func() {
		defer close(outputChannel)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			cmd := exec.Command(arcPath, "--vendors", LEDGER_VENDOR_IDS)
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				slog.Error("failed to get stdout pipe", "error", err.Error())
				time.Sleep(1 * time.Second)
				continue
			}

			// Start the command.
			if err := cmd.Start(); err != nil {
				slog.Error("failed to start command", "error", err.Error())
				time.Sleep(1 * time.Second)
				continue
			}

			done := make(chan struct{})
			// Watch for context cancellation. If the context is canceled,
			// send SIGINT to the process.
			go func() {
				select {
				case <-ctx.Done():
					// When context is cancelled, send SIGINT.
					if cmd.Process != nil {
						if err := cmd.Process.Signal(os.Interrupt); err != nil {
							slog.Warn("failed to send SIGINT", "error", err.Error())
						}
					}
				case <-done:
				}
			}()

			// Read from stdout and forward events.
			go func() {
				scanner := bufio.NewScanner(stdout)
				slog.Info("arc monitor started")
				for scanner.Scan() {
					select {
					case <-ctx.Done():
						return
					default:
					}

					var event ArcEvent
					if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
						slog.Error("failed to unmarshal event", "error", err.Error())
						continue
					}
					select {
					case outputChannel <- event:
					case <-ctx.Done():
						return
					}
				}
				if err := scanner.Err(); err != nil {
					slog.Error("error reading stdout", "error", err.Error())
				}
			}()

			cmd.Wait()
			close(done)

			select {
			case <-ctx.Done():
				// If the context was cancelled, exit without restarting.
				return
			default:
			}

			time.Sleep(1 * time.Second)
		}
	}()

	return outputChannel
}

var (
	activeWalletStatus    = WalletsStatus{}
	activeWalletStatusMtx = sync.RWMutex{}

	isCollectingWalletInfo = atomic.Bool{}
)

func collectWalletInfo(signerPath string, walletIds ...string) (map[string]base.AmiWalletInfo, error) {
	isCollectingWalletInfo.Store(true)
	defer isCollectingWalletInfo.Store(false)

	walletsIdsJoined := strings.Join(walletIds, ",")
	walletsArg := fmt.Sprintf("--wallets=%s", walletsIdsJoined)
	if len(walletIds) == 0 {
		walletsArg = "--wallets"
	}

	args := []string{walletsArg}
	infoBytes, _, err := ami.ExecuteInfo(signerPath, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to collect app info (%s)", err.Error())
	}
	info, err := base.ParseInfoOutput[signer.Info](infoBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse app info (%s)", err.Error())
	}
	return info.Wallets, nil
}

func RefresActiveWalletsStatus(signerPath string, wallets []string) error {
	if isCollectingWalletInfo.Load() {
		return nil
	}

	info, err := collectWalletInfo(signerPath)
	if err != nil {
		slog.Error("failed to collect wallet info", "error", err.Error())
		return err
	}
	newStatus := WalletsStatus{}
	for walletId, wallet := range info {
		if wallet.Kind != "ledger" {
			continue
		}
		if len(wallets) > 0 && !slices.Contains(wallets, walletId) {
			continue
		}
		newStatus[walletId] = wallet
	}

	activeWalletStatusMtx.Lock()
	defer activeWalletStatusMtx.Unlock()
	activeWalletStatus = newStatus
	return nil
}

func sendWalletStatusUpdate(statusChannel chan<- common.StatusUpdate) {
	statusChannel <- &WalletsStatusUpdate{maps.Clone(activeWalletStatus)}
}

func startWalletsStatusProvider(ctx context.Context, signerPath, arcPath string, wallets []string, statusChannel chan<- common.StatusUpdate) {
	arcEventChannel := RunArcMonitor(ctx, arcPath)

	for {
		if err := RefresActiveWalletsStatus(signerPath, wallets); err != nil {
			slog.Error("failed to collect wallet info", "error", err.Error())
			time.Sleep(1 * time.Second)
			continue
		}
		sendWalletStatusUpdate(statusChannel)
		break
	}

	go func() {

		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-arcEventChannel:
				if !ok {
					return
				}

				switch event.Fields.Message {
				case "connected":
					slog.Info("ledger connected", "path", event.Fields.Path)
					// TODO: check directly
					info, err := collectWalletInfo(signerPath, activeWalletStatus.DisconnectedWallets()...)
					if err != nil {
						slog.Error("failed to collect wallet info", "error", err.Error())
						continue
					}
					for walletId, wallet := range info {
						if wallet.Kind != "ledger" {
							continue
						}
						if len(wallets) > 0 && !slices.Contains(wallets, walletId) {
							continue
						}

						func() {
							activeWalletStatusMtx.Lock()
							defer activeWalletStatusMtx.Unlock()
							activeWalletStatus[walletId] = wallet
						}()
						sendWalletStatusUpdate(statusChannel)
					}
				case "disconnected":
					slog.Info("ledger disconnected", "path", event.Fields.Path)
					for walletId, wallet := range activeWalletStatus {
						if wallet.Kind != "ledger" {
							continue
						}

						if wallet.DevicePath != event.Fields.Path {
							continue
						}

						newWallet := wallet
						newWallet.LedgerStatus = "disconnected"
						newWallet.DevicePath = ""

						func() {
							activeWalletStatusMtx.Lock()
							defer activeWalletStatusMtx.Unlock()
							activeWalletStatus[walletId] = newWallet
						}()
						sendWalletStatusUpdate(statusChannel)
						continue
					}
				default:
					slog.Warn("unknown message", "message", event.Fields.Message)
					continue
				}

				sendWalletStatusUpdate(statusChannel)
			}
		}
	}()

}
