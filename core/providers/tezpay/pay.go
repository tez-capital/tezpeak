package tezpay

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gocarina/gocsv"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hjson/hjson-go/v4"
	"github.com/tez-capital/tezbake/apps"
	"github.com/tez-capital/tezbake/apps/pay"
	"github.com/tez-capital/tezpay/common"
	"github.com/tez-capital/tezpeak/configuration"
	peakCommon "github.com/tez-capital/tezpeak/core/common"
)

type TezpayProvider struct {
	configuration *configuration.TezpayModuleConfiguration

	tezpay *pay.Tezpay
}

type TezpayVersion struct {
	AmiTezpay string `json:"ami-tezpay"`
	Tezpay    string `json:"tezpay"`
}

func (tezpayProvider *TezpayProvider) RegisterApi(app *fiber.Group) error {
	app.Get("/tezpay/can-pay", func(c *fiber.Ctx) error {
		return c.JSON(tezpayProvider.CanPay())
	})

	app.Get("/tezpay/info", func(c *fiber.Ctx) error {
		version, err := tezpayProvider.Version()
		if err != nil {
			slog.Error("failed to get version", "error", err.Error())
			version = "unknown"
		}

		var versions TezpayVersion
		err = hjson.Unmarshal([]byte(version), &versions)
		if err != nil {
			slog.Error("failed to parse version", "error", err.Error())
		}

		configurationString, err := tezpayProvider.GetTezpayConfiguration()
		if err != nil {
			slog.Error("failed to get configuration", "error", err.Error())
			configurationString = "{}"
		}

		var configuration map[string]interface{}
		err = hjson.Unmarshal([]byte(configurationString), &configuration)
		if err != nil {
			slog.Error("failed to parse configuration", "error", err.Error())
		}

		return c.JSON(map[string]interface{}{
			"version":       versions,
			"configuration": configuration,
		})
	})

	app.Get("/tezpay/generate-payouts", func(c *fiber.Ctx) error {
		if !tezpayProvider.CanPay() {
			return c.Status(fiber.StatusForbidden).SendString("not allowed")
		}

		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		cycle := int64(-1)
		cycleQuery := c.Query("cycle")
		if cycleQuery != "" {
			cycle, err := strconv.ParseInt(cycleQuery, 10, 64)
			if err != nil || cycle < 0 {
				return c.Status(fiber.StatusBadRequest).SendString("Invalid 'cycle' parameter")
			}
		}

		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			outputChannel := make(chan string)
			go tezpayProvider.GeneratePayouts(cycle, outputChannel)

			for output := range outputChannel {
				fmt.Fprintf(w, "%v\n", output)
				w.Flush()
			}
		})

		return nil
	})

	app.Post("/tezpay/pay", func(c *fiber.Ctx) error {
		if !tezpayProvider.CanPay() {
			return c.Status(fiber.StatusForbidden).SendString("not allowed")
		}
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		var blueprint CyclePayoutBlueprint
		if err := c.BodyParser(&blueprint); err != nil {
			slog.Error("failed to parse generate-payouts params", "error", err.Error())
			return c.Status(400).SendString("invalid request")
		}

		dry := c.Query("dry") == "true"

		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			outputChannel := make(chan string)
			go tezpayProvider.Pay(&blueprint, outputChannel, dry)

			for output := range outputChannel {
				fmt.Fprintf(w, "%v\n", output)
				w.Flush()
			}
		})

		return nil
	})

	app.Get("/tezpay/statistics", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		var numberOfCycles int64 = 0
		var lastCycle int64 = 0

		if numCycles := c.Query("numberOfCycles"); numCycles != "" {
			var err error
			numberOfCycles, err = strconv.ParseInt(numCycles, 10, 64)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).SendString("Invalid 'numberOfCycles' parameter")
			}
		}

		// Parse lastCycle from query parameter
		if lastCyc := c.Query("lastCycle"); lastCyc != "" {
			var err error
			lastCycle, err = strconv.ParseInt(lastCyc, 10, 64)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).SendString("Invalid 'lastCycle' parameter")
			}
		}

		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			outputChannel := make(chan string)

			go tezpayProvider.Statistics(numberOfCycles, lastCycle, outputChannel)

			for output := range outputChannel {
				fmt.Fprintf(w, "%v\n", output)
				w.Flush()
			}
		})

		return nil
	})

	app.Post("/tezpay/test-notify", func(c *fiber.Ctx) error {
		if !tezpayProvider.CanPay() {
			return c.Status(fiber.StatusForbidden).SendString("not allowed")
		}
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		notificator := c.Query("notificator")

		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			outputChannel := make(chan string)

			go tezpayProvider.TestNotify(notificator, outputChannel)

			for output := range outputChannel {
				fmt.Fprintf(w, "%v\n", output)
				w.Flush()
			}
		})
		return nil
	})

	app.Post("/tezpay/test-extensions", func(c *fiber.Ctx) error {
		if !tezpayProvider.CanPay() {
			return c.Status(fiber.StatusForbidden).SendString("not allowed")
		}
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			outputChannel := make(chan string)

			go tezpayProvider.TestExtensions(outputChannel)

			for output := range outputChannel {
				fmt.Fprintf(w, "%v\n", output)
				w.Flush()
			}
		})
		return nil
	})

	app.Get("/tezpay/list-reports", func(c *fiber.Ctx) error {
		reports, err := tezpayProvider.ListReports(c.Query("dry") == "true")
		if err != nil {
			return c.Status(500).SendString("failed to list reports")
		}

		return c.JSON(reports)
	})

	app.Get("/tezpay/report", func(c *fiber.Ctx) error {
		report, err := tezpayProvider.GetReport(c.Query("id"), c.Query("dry") == "true")
		if err != nil {
			return c.Status(500).SendString("failed to get report")
		}

		return c.JSON(report)
	})

	app.Get("/tezpay/stop-continual", func(c *fiber.Ctx) error {
		if !tezpayProvider.CanPay() {
			return c.Status(fiber.StatusForbidden).SendString("not allowed")
		}
		err := tezpayProvider.StopContinualPayouts()
		if err != nil {
			return c.Status(500).SendString("failed to stop service: " + err.Error())
		}
		return c.Status(200).SendString("service stopped")
	})

	app.Get("/tezpay/start-continual", func(c *fiber.Ctx) error {
		if !tezpayProvider.CanPay() {
			return c.Status(fiber.StatusForbidden).SendString("not allowed")
		}
		err := tezpayProvider.StartContinualPayouts()
		if err != nil {
			return c.Status(500).SendString("failed to start service: " + err.Error())
		}
		return c.Status(200).SendString("service started")
	})

	return nil
}

func setupTezpayProvider(configuration *configuration.TezpayModuleConfiguration, app *fiber.Group) error {
	tezpayPath, ok := configuration.Applications["tezpay"]
	if !ok {
		return errors.New("tezpay path not found in configuration")
	}

	tezpayProvider := &TezpayProvider{
		configuration: configuration,

		tezpay: apps.TezpayFromPath(tezpayPath),
	}

	return tezpayProvider.RegisterApi(app)
}

type ExecutionFinishedMessage struct {
	Succeess bool   `json:"success"`
	ExitCode int    `json:"exit_code"`
	Error    string `json:"error"`
	Phase    string `json:"phase"`
}

func buildFinishMessage(exitCode int, err error) string {
	message := ExecutionFinishedMessage{
		Succeess: exitCode == 0 && err == nil,
		ExitCode: exitCode,
		Error:    "",
		Phase:    "execution_finished",
	}
	if err != nil {
		message.Error = err.Error()
	}

	messageBytes, _ := json.Marshal(message)
	return string(messageBytes)
}

func (t *TezpayProvider) GeneratePayouts(cycle int64, outputChannel chan<- string) {
	switch {
	case cycle < 0:
		exitcode, err := t.tezpay.ExecuteWithOutputChannel(outputChannel, "generate-payouts", "--output-format", "json")
		outputChannel <- buildFinishMessage(exitcode, err)
		close(outputChannel)
	default:
		exitCode, err := t.tezpay.ExecuteWithOutputChannel(outputChannel, "generate-payouts", "--cycle", fmt.Sprintf("%d", cycle))
		outputChannel <- buildFinishMessage(exitCode, err)
		close(outputChannel)
	}
}

type CyclePayoutBlueprint common.CyclePayoutBlueprint

func (t *TezpayProvider) Pay(blueprint *CyclePayoutBlueprint, outputChannel chan<- string, dry bool) {
	marshaledBlueprint, err := json.Marshal(blueprint)
	if err != nil {
		outputChannel <- buildFinishMessage(-1, err)
		close(outputChannel)
		return
	}

	// create temp file with blueprint
	fileName, err := uuid.NewRandom()
	if err != nil {
		outputChannel <- buildFinishMessage(-1, err)
		close(outputChannel)
		return
	}
	filePath := path.Join(os.TempDir(), fileName.String())
	err = os.WriteFile(filePath, marshaledBlueprint, 0644)
	if err != nil {
		outputChannel <- buildFinishMessage(-1, err)
		close(outputChannel)
		return
	}
	defer os.Remove(filePath)

	var exitcode int
	if dry || t.configuration.ForceDryRun {
		exitcode, err = t.tezpay.ExecuteWithOutputChannel(outputChannel, "pay", "--output-format", "json", "--from-file", filePath, "--confirm", "--disable-donation-prompt", "--dry-run")
	} else {
		exitcode, err = t.tezpay.ExecuteWithOutputChannel(outputChannel, "pay", "--output-format", "json", "--from-file", filePath, "--confirm", "--disable-donation-prompt")
	}
	outputChannel <- buildFinishMessage(exitcode, err)
	close(outputChannel)
}

func (t *TezpayProvider) Version() (string, error) {
	output, exitCode, err := t.tezpay.ExecuteGetOutput("version")
	if err != nil {
		return "", err
	}
	if exitCode != 0 {
		return "", errors.New("failed to get version")
	}

	return output, nil
}

func (t *TezpayProvider) Statistics(numberOfCycles, lastCycle int64, outputChannel chan<- string) {
	if numberOfCycles < 0 {
		numberOfCycles = 10
	}
	if lastCycle < 0 {
		lastCycle = 0
	}

	exitcode, err := t.tezpay.ExecuteWithOutputChannel(outputChannel, "tezpay", "--output-format", "json", "statistics", "--cycles", fmt.Sprintf("%d", numberOfCycles), "--last-cycle", fmt.Sprintf("%d", lastCycle))
	outputChannel <- buildFinishMessage(exitcode, err)
	close(outputChannel)
}

func (t *TezpayProvider) TestNotify(notificator string, outputChannel chan<- string) {
	var exitcode int
	var err error
	if notificator == "" {
		exitcode, err = t.tezpay.ExecuteWithOutputChannel(outputChannel, "tezpay", "--output-format", "json", "test-notify")
	} else {
		exitcode, err = t.tezpay.ExecuteWithOutputChannel(outputChannel, "tezpay", "--output-format", "json", "test-notify", "--notificator", notificator)
	}
	outputChannel <- buildFinishMessage(exitcode, err)
	close(outputChannel)
}

func (t *TezpayProvider) TestExtensions(outputChannel chan<- string) {
	exitcode, err := t.tezpay.ExecuteWithOutputChannel(outputChannel, "tezpay", "--output-format", "json", "test-extensions")
	outputChannel <- buildFinishMessage(exitcode, err)
	close(outputChannel)
}

func listReportsInternal(reportsDirectoryPath string) ([]string, error) {
	dirInfo, err := os.Stat(reportsDirectoryPath)
	switch {
	case err != nil && os.IsNotExist(err):
		return []string{}, nil
	case err != nil:
		return nil, err
	case !dirInfo.IsDir():
		return []string{}, nil
	}

	// list folders in reports directory, which contains summary.json
	subDirectories, err := os.ReadDir(reportsDirectoryPath)
	if err != nil {
		return nil, err
	}

	acc := []string{}
	for _, dir := range subDirectories {
		path := filepath.Join(reportsDirectoryPath, dir.Name())
		fileInfo, err := os.Stat(filepath.Join(path, "summary.json"))
		if err != nil {
			continue
		}
		if fileInfo.IsDir() {
			continue
		}
		acc = append(acc, dir.Name())
	}
	return acc, nil
}

func (t *TezpayProvider) ListReports(dry bool) ([]string, error) {
	reportsDirectoryPath := path.Join(t.tezpay.GetPath(), "reports")
	if dry {
		reportsDirectoryPath = path.Join(reportsDirectoryPath, "dry")
	}
	return listReportsInternal(reportsDirectoryPath)
}

type PayoutReport struct {
	Name    string                    `json:"name"`
	Summary common.CyclePayoutSummary `json:"summary"`
	Payouts []common.PayoutReport     `json:"payouts"`
	Invalid []common.PayoutReport     `json:"invalid"`
}

func (t *TezpayProvider) GetReport(reportName string, dry bool) (*PayoutReport, error) {
	reportsDirectoryPath := path.Join(t.tezpay.GetPath(), "reports")
	if dry {
		reportsDirectoryPath = path.Join(reportsDirectoryPath, "dry")
	}

	reportPath := path.Join(reportsDirectoryPath, reportName)
	summaryPath := path.Join(reportPath, "summary.json")
	summary, err := os.ReadFile(summaryPath)
	if err != nil {
		return nil, err
	}
	var cycleSummary common.CyclePayoutSummary
	err = json.Unmarshal(summary, &cycleSummary)
	if err != nil {
		slog.Warn("failed to unmarshal summary", "error", err.Error())
	}

	// payouts.csv
	payoutsPath := path.Join(reportPath, "payouts.csv")
	payouts, _ := os.ReadFile(payoutsPath)
	if err != nil {
		slog.Warn("failed to read payouts", "error", err.Error())
	}
	var payoutReports []common.PayoutReport
	err = gocsv.UnmarshalBytes(payouts, &payoutReports)
	if err != nil {
		slog.Warn("failed to parse payouts", "error", err.Error())
	}

	// invalid.csv
	invalidPath := path.Join(reportPath, "invalid.csv")
	invalid, _ := os.ReadFile(invalidPath)
	if err != nil {
		slog.Warn("failed to read invalid", "error", err.Error())
	}
	var invalidReports []common.PayoutReport
	err = gocsv.UnmarshalBytes(invalid, &invalidReports)
	if err != nil {
		slog.Warn("failed to parse invalid", "error", err.Error())
	}

	return &PayoutReport{
		Name:    reportName,
		Summary: cycleSummary,
		Payouts: payoutReports,
		Invalid: invalidReports,
	}, nil
}

func (t *TezpayProvider) StopContinualPayouts() error {
	defer peakCommon.UpdateServiceStatus(t.tezpay.GetPath())
	exitCode, err := t.tezpay.Stop()
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return errors.New("failed to stop continual payouts")
	}
	return nil
}

func (t *TezpayProvider) StartContinualPayouts() error {
	defer peakCommon.UpdateServiceStatus(t.tezpay.GetPath())
	exitCode, err := t.tezpay.Start()
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return errors.New("failed to start continual payouts")
	}
	return nil
}

func (t *TezpayProvider) GetTezpayConfiguration() (string, error) {
	configurationFilePath := path.Join(t.tezpay.GetPath(), "config.hjson")
	configurationBytes, err := os.ReadFile(configurationFilePath)
	if err != nil {
		return "{}", err
	}
	return string(configurationBytes), nil
}

func (t *TezpayProvider) CanPay() bool {
	return t.configuration.Mode == configuration.PrivatePeakMode
}
