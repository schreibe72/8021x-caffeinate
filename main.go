package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"git.manfredschreiber.de/8021x-caffeinate/icon"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

const searchProzess = "eapolclient"

var c *process
var version string

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
	mVersion := systray.AddMenuItem(version, "Version")
	mVersion.Disable()
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
		defer ticker.Stop()
		updateTicker := time.NewTicker(1 * time.Hour)
		defer updateTicker.Stop()

		updateFunc := func() string {
			u, nv := check4update(version)
			if u != "" {
				log.Printf("Enable %s %s", u, nv)
				mVersion.Enable()
				mVersion.SetTitle(fmt.Sprintf("%s -> %s", version, nv))
				mVersion.SetTooltip("Click to go to update page")
				return u
			}
			return ""
		}
		updateUrl := updateFunc()

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
			case <-mVersion.ClickedCh:
				log.Printf("Open %s", updateUrl)
				open.Run(updateUrl)
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
			case <-updateTicker.C:
				updateUrl = updateFunc()
			}
		}
	}()
}
