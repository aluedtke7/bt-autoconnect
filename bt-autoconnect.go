package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/antigloss/go/logger"
)

var (
	btDevices []string
)

func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		return "~/"
	}
	return usr.HomeDir
}

func checkBluetooth() {
	result, err := exec.Command("bluetoothctl", "devices").Output()
	if err != nil {
		logger.Error(err.Error())
	} else {
		arr := strings.Split(string(result), "\n")
		logger.Info("Devices paired:")
		for _, s := range arr {
			parts := strings.Split(s, " ")
			if len(parts) > 1 {
				btDevices = append(btDevices, parts[1])
				logger.Info(parts[1])
			}
		}
	}

	var lastExitCode int = 999
	for {
		cmd := exec.Command("ls", "/dev/input/event0")
		err = cmd.Run()
		exitCode := cmd.ProcessState.ExitCode()
		if exitCode == 2 {
			// not connected
			if lastExitCode != 2 {
				logger.Info("Re-run mplayer... ")
			}
			for idx, btDevice := range btDevices {
				logger.Info(fmt.Sprintf("Trying to connect device #%d %s", idx, btDevice))
				cmd = exec.Command("bluetoothctl", "connect", btDevice)
				err = cmd.Run()
				connectExitCode := cmd.ProcessState.ExitCode()
				if connectExitCode == 0 {
					logger.Info("Success with device " + btDevice)
					break
				}
			}
		} else if exitCode == 0 {
			// connected
			if lastExitCode != 0 {
				logger.Info("Re-run mplayer... ")
			}
		}
		lastExitCode = exitCode
		time.Sleep(5 * time.Second)
	}
}

func main() {
	homePath := filepath.Join(getHomeDir(), ".bt-autoconnect")
	_ = os.MkdirAll(homePath, os.ModePerm)
	_ = logger.Init(filepath.Join(homePath, "log"), 10, 2, 10, true)
	logger.Info("Starting bt-autoconnect...")

	var ctrlChan = make(chan os.Signal)
	signal.Notify(ctrlChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	go checkBluetooth()

	// this is waiting for someone to stop piradio
	<-ctrlChan
	logger.Info("Ctrl+C received... Exiting")
	os.Exit(0)
}
