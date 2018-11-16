// .volumecontroller
package main

import (
	"fmt"
	"io/ioutil"

	"github.com/getlantern/systray"
	"github.com/labstack/echo"
)

func main() {
	systray.Run(onReady, onExit)

	//Add server for settings-manager in webapp
	e := echo.New()
	e.Static("/", "settingsTemplate.html")
	e.Logger.Fatal(e.Start(":1337"))
}

func onReady() {
	systray.SetIcon(getIcon("icon/volumecontroller.ico"))
	systray.SetTitle("VolumController")
	systray.SetTooltip("Arduino volume controller for windows")
	mOpen := systray.AddMenuItem("Config", "Configure settings")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Exit", "Exit application")

	for {
		select {
		case <-mOpen.ClickedCh:
			fmt.Println("Config clicked")
		case <-mQuit.ClickedCh:
			fmt.Println("Quitting")
			systray.Quit()
		}
	}
}

func onExit() {}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}
