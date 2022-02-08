package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"git.manfredschreiber.de/8021x-caffeinate/icon"
	"github.com/getlantern/systray"
)

const searchProzess = "eapolclient"

var c *process

func main() {
	f, err := os.OpenFile("/tmp/8021x-caffeinate.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	c = newProcess("caffeinate", "-s")
	defer c.stop()
	onExit := func() {

	}
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTemplateIcon(icon.Data, icon.Data)
	mChecked := systray.AddMenuItemCheckbox("Permanent", "Activate Caffeinate permanent", false)
	systray.AddSeparator()
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
			select {
			case <-mChecked.ClickedCh:
				if mChecked.Checked() {
					c.stop()
					mChecked.Uncheck()
					systray.SetTemplateIcon(icon.Data, icon.Data)
				} else {
					c.start()
					mChecked.Check()
					systray.SetIcon(icon.Data)
				}
			case <-ticker.C:
				if !mChecked.Checked() {
					if newProcess(searchProzess).findNameInProcesslist() {
						if c.started() || c.start() {
							systray.SetIcon(icon.Data)
						}
					} else {
						if !c.started() || c.stop() {
							systray.SetTemplateIcon(icon.Data, icon.Data)
						}
					}
				}
			}
		}
	}()
}
