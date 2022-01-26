package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"git.manfredschreiber.de/8021x-caffeinate/icon"
	"github.com/getlantern/systray"
)

const searchProzess = "eapolclient"

var c caffeinate

func findProcess() bool {
	cmd := exec.Command("ps", "-ax")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(stdout)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		a := strings.Join(strings.Fields(scanner.Text())[3:], " ")
		if strings.Contains(a, searchProzess) {
			return true
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return false
}

func main() {
	c = newCaffeinate()
	defer c.stop()
	onExit := func() {

	}
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		fmt.Println("Requesting quit")
		systray.Quit()
		fmt.Println("Finished quitting")
	}()

	go func() {
		ticker := time.NewTicker(60 * time.Second)
		for {
			if findProcess() {
				if c.start() {
					systray.SetIcon(icon.Data)
				}
			} else {
				if c.stop() {
					systray.SetTemplateIcon(icon.Data, icon.Data)
				}
			}
			<-ticker.C
		}
	}()
}
