// .volumecontroller
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/getlantern/systray"
)

type Settings struct {
	Settings []struct {
		Rotation struct {
			Func string        `json:"func"`
			Apps []interface{} `json:"apps"`
		} `json:"Rotation"`
		Button struct {
			Func string        `json:"func"`
			Apps []interface{} `json:"apps"`
		} `json:"Button"`
	} `json:"Settings"`
	Operations struct {
		Rotation []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"Rotation"`
		Button []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"Button"`
	} `json:"Operations"`
}

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(getIcon("icon/volumecontroller.ico"))
	systray.SetTitle("VolumeController")
	systray.SetTooltip("Arduino volume controller for windows")
	mOpen := systray.AddMenuItem("Config", "Configure settings")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Exit", "Exit application")

	//Add server for settings-manager in webapp
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/settings", settingsLoad)

	for {
		select {
		case <-mOpen.ClickedCh:
			fmt.Println("Config clicked")
			http.ListenAndServe(":1337", nil)
		case <-mQuit.ClickedCh:
			fmt.Println("Quitting")
			systray.Quit()
		}
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/static/index.html")
}

func settingsLoad(w http.ResponseWriter, r *http.Request) {
	config := readSettings()
	fmt.Println("Settings loaded")
	json.NewEncoder(w).Encode(config)
}

func readSettings() Settings {
	jsonFile, err := os.Open("../settings.cfg")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Config is opened successfully")
	//close file
	defer jsonFile.Close()
	var settings Settings
	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &settings)
	return settings
}

func onExit() {}

func getIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}
