package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/getlantern/systray"
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

	//Add server for settings-manager in webapp
	indexHandler := http.NewServeMux()
	indexHandler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("SettingsTemplate.html")
		t.Execute(w, nil)
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
			srv.Close()
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
