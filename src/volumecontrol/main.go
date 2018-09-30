package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/getlantern/systray"
)

//SettingsStruct to populate html from json config file
type SettingsStruct struct {
	Settings []struct {
		Encoder []struct {
			Knob []struct {
				Text        string   `json:"text"`
				Application []string `json:"Application"`
				Active      bool     `json:"Active"`
			} `json:"Knob,omitempty"`
			Button []struct {
				Text        string   `json:"text"`
				Application []string `json:"Application"`
				Active      bool     `json:"Active"`
			} `json:"Button,omitempty"`
		} `json:"encoder"`
	} `json:"Settings"`
}

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

	//Add server for settings-manager in webapp
	indexHandler := http.NewServeMux()
	indexHandler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := getSettings()
		t, _ := template.ParseFiles("SettingsTemplate.html")
		t.Execute(w, p)
	})
	indexHandler.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Saved")
	})

	srv := &http.Server{
		Addr:    ":1337",
		Handler: indexHandler,
	}

	for {
		select {
		case <-mOpen.ClickedCh:
			fmt.Println("Config clicked")
			go log.Fatal(srv.ListenAndServe())
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

func getSettings() SettingsStruct {
	settingsFile, err := os.Open("settings.cfg.default")
	if err != nil {
		fmt.Println(err)
	}

	defer settingsFile.Close()

	byteValue, err := ioutil.ReadAll(settingsFile)
	var settings SettingsStruct

	json.Unmarshal(byteValue, &settings)
	fmt.Println(settings.Settings[0].Encoder[1].Knob[1].Text)
	return settings
}
