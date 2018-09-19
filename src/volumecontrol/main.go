package main

import (
	"fmt"
	"io/ioutil"

	"github.com/getlantern/systray"
	"github.com/zserge/webview"
)

func main() {
	systray.Run(onReady, onExit)

}
func onReady() {
	systray.SetIcon(getIcon("../icon/volumecontroller.ico"))
	systray.SetTitle("VolumController")
	systray.SetTooltip("Arduino volume controller for windows")

	mOpen := systray.AddMenuItem("Config", "Configure settings")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Exit", "Exit application")

	for {
		select {
		case <-mOpen.ClickedCh:
			fmt.Println("Config clicked")
			webview.Open("Volumecontroller", "https://google.com", 400, 400, false)

		case <-mQuit.ClickedCh:
			fmt.Print("Quitting")
			systray.Quit()
			fmt.Print("Finished quitting")
			return
		}
	}

}

func onExit() {

}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}
