package service

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"

	logger_util "github.com/free5gc/util/logger"
	"github.com/sirupsen/logrus"
	"github.com/softmurata/nef/internal/logger"
	"github.com/softmurata/nef/internal/util"
	"github.com/softmurata/nef/pkg/factory"
	"github.com/urfave/cli"

	"github.com/softmurata/nef/internal/context"

	"github.com/softmurata/nef/internal/sbi/assessionwithqos"
	"github.com/softmurata/nef/internal/sbi/serviceparameter"

	"github.com/free5gc/util/httpwrapper"
)

type NEF struct {
	KeyLogPath string
}

type (
	// Commands information.
	Commands struct {
		config string
	}
)

var commands Commands

var cliCmd = []cli.Flag{
	cli.StringFlag{
		Name:  "config, c",
		Usage: "Load configuration from `FILE`",
	},
	cli.StringFlag{
		Name:  "log, l",
		Usage: "Output NF log to `FILE`",
	},
	cli.StringFlag{
		Name:  "log5gc, lc",
		Usage: "Output free5gc log to `FILE`",
	},
}

func (*NEF) GetCliCmd() (flags []cli.Flag) {
	return cliCmd
}

func (nef *NEF) Initialize(c *cli.Context) error {
	commands = Commands{
		config: c.String("config"),
	}

	if commands.config != "" {
		if err := factory.InitConfigFactory(commands.config); err != nil {
			return err
		}
	} else {
		if err := factory.InitConfigFactory(util.NefDefaultConfigPath); err != nil {
			return err
		}
	}

	nef.SetLogLevel()

	if err := factory.CheckConfigVersion(); err != nil {
		return err
	}

	return nil
}

func (nef *NEF) SetLogLevel() {
	if factory.NefConfig.Logger == nil {
		logger.InitLog.Warnln("NEF config without log level setting!!!")
		return
	}

	if factory.NefConfig.Logger.NEF != nil {
		if factory.NefConfig.Logger.NEF.DebugLevel != "" {
			if level, err := logrus.ParseLevel(factory.NefConfig.Logger.NEF.DebugLevel); err != nil {
				logger.InitLog.Warnf("NEF Log level [%s] is invalid, set to [info] level",
					factory.NefConfig.Logger.NEF.DebugLevel)
				logger.SetLogLevel(logrus.InfoLevel)
			} else {
				logger.InitLog.Infof("NEF Log level is set to [%s] level", level)
				logger.SetLogLevel(level)
			}
		} else {
			logger.InitLog.Infoln("NEF Log level not set. Default set to [info] level")
			logger.SetLogLevel(logrus.InfoLevel)
		}
		logger.SetReportCaller(factory.NefConfig.Logger.NEF.ReportCaller)
	}
}

func (nef *NEF) FilterCli(c *cli.Context) (args []string) {
	for _, flag := range nef.GetCliCmd() {
		name := flag.GetName()
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (nef *NEF) Start() {

	logger.InitLog.Infoln("Server started")

	router := logger_util.NewGinWithLogrus(logger.GinLog)

	// ToDo: add service name
	assessionwithqos.AddService(router)
	serviceparameter.AddService(router)

	pemPath := util.NefDefaultPemPath
	keyPath := util.NefDefaultKeyPath
	sbi := factory.NefConfig.Configuration.Sbi
	if sbi.Tls != nil {
		pemPath = sbi.Tls.Pem
		keyPath = sbi.Tls.Key
	}

	self := context.NEF_Self()
	context.InitNefContext()
	addr := fmt.Sprintf("%s:%d", self.BindingIPv4, self.SBIPort)

	/*
		// ToDo: maybe not need unit test
		addr := fmt.Sprintf("%s:%d", self.BindingIPv4, self.SBIPort)
		profile, err := consumer.BuildNFProfile(self)

		if err != nil {
			logger.InitLog.Error("Failed to build NSSF profile")
		}

		_, self.NfId, err = consumer.SendRegisterNFInstance(self.NrfUri, self.NfId, profile)
		if err != nil {
			logger.InitLog.Errorf("Failed to register NSSF to NRF: %s", err.Error())
		}
	*/

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				// Print stack for panic to log. Fatalf() will let program exit.
				logger.InitLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
			}
		}()

		<-signalChannel
		nef.Terminate()
		os.Exit(0)
	}()

	server, err := httpwrapper.NewHttp2Server(addr, nef.KeyLogPath, router)
	if server == nil {
		logger.InitLog.Errorf("Initialize HTTP server failed: %+v", err)
		return
	}

	if err != nil {
		logger.InitLog.Warnf("Initialize HTTP server: %+v", err)
	}

	serverScheme := factory.NefConfig.Configuration.Sbi.Scheme
	if serverScheme == "http" {
		err = server.ListenAndServe()
	} else if serverScheme == "https" {
		err = server.ListenAndServeTLS(pemPath, keyPath)
	}

	if err != nil {
		logger.InitLog.Fatalf("HTTP server setup failed: %+v", err)
	}
}

func (nef *NEF) Exec(c *cli.Context) error {
	logger.InitLog.Traceln("args:", c.String("nefcfg"))
	args := nef.FilterCli(c)
	logger.InitLog.Traceln("filter: ", args)
	command := exec.Command("./nef", args...)

	stdout, err := command.StdoutPipe()
	if err != nil {
		logger.InitLog.Fatalln(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				// Print stack for panic to log. Fatalf() will let program exit.
				logger.InitLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
			}
		}()

		in := bufio.NewScanner(stdout)
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	stderr, err := command.StderrPipe()
	if err != nil {
		logger.InitLog.Fatalln(err)
	}
	go func() {
		defer func() {
			if p := recover(); p != nil {
				// Print stack for panic to log. Fatalf() will let program exit.
				logger.InitLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
			}
		}()

		in := bufio.NewScanner(stderr)
		for in.Scan() {
			fmt.Println(in.Text())
		}
		wg.Done()
	}()

	go func() {
		defer func() {
			if p := recover(); p != nil {
				// Print stack for panic to log. Fatalf() will let program exit.
				logger.InitLog.Fatalf("panic: %v\n%s", p, string(debug.Stack()))
			}
		}()

		if err = command.Start(); err != nil {
			fmt.Printf("NEF Start error: %v", err)
		}
		wg.Done()
	}()

	wg.Wait()

	return err
}

func (nef *NEF) Terminate() {
	logger.InitLog.Infof("Terminating NEF...")
	/*
		// deregister with NRF
		problemDetails, err := consumer.SendDeregisterNFInstance()
		if problemDetails != nil {
			logger.InitLog.Errorf("Deregister NF instance Failed Problem[%+v]", problemDetails)
		} else if err != nil {
			logger.InitLog.Errorf("Deregister NF instance Error[%+v]", err)
		} else {
			logger.InitLog.Infof("Deregister from NRF successfully")
		}
	*/

	logger.InitLog.Infof("NEF terminated")
}
