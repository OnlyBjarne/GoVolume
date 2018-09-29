package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/getlantern/systray"
)

//globals for functions for knob and buttons
var knobModes = []functionStruct{{"System volume", nil, true}, {"Application volume", []string{"Spotify.exe"}, false}}
var buttonModes = []functionStruct{{"Music Control", nil, true}, {"Keystroke", []string{"alt+F4", "F"}, false}, {"Toggle sound device", []string{"Headset", "Speakers", "Earbuds"}, false}}

//SettingsTemplate exports to the template for settings
type SettingsTemplate struct {
	CountKnobs      []struct{}
	KnobModesList   []functionStruct
	ButtonModesList []functionStruct
}

type functionStruct struct {
	Text        string
	Application []string
	Active      bool
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
		p := SettingsTemplate{CountKnobs: make([]struct{}, 4), KnobModesList: knobModes, ButtonModesList: buttonModes}
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
			log.Fatal(srv.ListenAndServe())
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
