/*
// test code for hello.go
package cmd

func Command() string {

	cmd := "git commit -m first commit"

	return cmd

}
*/

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"runtime/debug"

	"github.com/asaskevich/govalidator"
	"github.com/free5gc/util/version"
	"github.com/softmurata/nef/internal/logger"
	"github.com/softmurata/nef/internal/util"
	nef_service "github.com/softmurata/nef/pkg/service"
	"github.com/urfave/cli"
)

var NEF = &nef_service.NEF{}

func main() {

	defer func() {
		if p := recover(); p != nil {
			// Print stack for panic to log. Fatalf() will let program exit.
			logger.AppLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
		}
	}()

	logger.AppLog.Infof("Main Function")
	app := cli.NewApp()
	app.Name = "nef"
	app.Usage = "5G Network Exposure Function (NEF)"
	app.Action = action
	app.Flags = NEF.GetCliCmd()

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("NEF Run Error: %v\n", err)
	}

}

func action(c *cli.Context) error {
	if err := initLogFile(c.String("log"), c.String("log5gc")); err != nil {
		logger.AppLog.Errorf("%+v", err)
		return err
	}

	if err := NEF.Initialize(c); err != nil {
		switch errType := err.(type) {
		case govalidator.Errors:
			validErrs := err.(govalidator.Errors).Errors()
			for _, validErr := range validErrs {
				logger.CfgLog.Errorf("%+v", validErr)
			}
		default:
			logger.CfgLog.Errorf("%+v", errType)
		}
		logger.CfgLog.Errorf("[-- PLEASE REFER TO SAMPLE CONFIG FILE COMMENTS --]")
		return fmt.Errorf("Failed to initialize !!")
	}

	logger.AppLog.Infoln(c.App.Name)
	logger.AppLog.Infoln("NEF version: ", version.GetVersion())

	NEF.Start()

	return nil

}

func initLogFile(logNfPath, log5gcPath string) error {
	NEF.KeyLogPath = util.NefDefaultKeyLogPath

	if err := logger.LogFileHook(logNfPath, log5gcPath); err != nil {
		return err
	}

	if logNfPath != "" {
		nfDir, _ := filepath.Split(logNfPath)
		tmpDir := filepath.Join(nfDir, "key")
		if err := os.MkdirAll(tmpDir, 0775); err != nil {
			logger.InitLog.Errorf("Make directory %s failed: %+v", tmpDir, err)
			return err
		}
		_, name := filepath.Split(util.NefDefaultKeyLogPath)
		NEF.KeyLogPath = filepath.Join(tmpDir, name)
	}

	return nil
}
