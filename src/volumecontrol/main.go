package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/getlantern/systray"
	"github.com/zserge/webview"
)

//globals for functions for knob and buttons
var knobModes = []string{"System volume", "Application volume"}
var buttonModes = []string{"Music Control", "Keystroke", "Sound Device", "Command"}

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
			url := startServer()
			w := webview.New(webview.Settings{
				Width:     700,
				Height:    600,
				Title:     "Volume controller settings",
				Resizable: false,
				URL:       url,
			})
			defer w.Exit()
			w.Run()
		case <-mQuit.ClickedCh:
			fmt.Println("Quitting")
			systray.Quit()
			fmt.Println("Finished quitting")
			return
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

func startServer() string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer ln.Close()
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(indexHTML))
		})
		log.Fatal(http.Serve(ln, nil))
	}()
	return "http://" + ln.Addr().String()
}

func populateHTML(label string, modes []string, selected int) string {
	var output = `
	<label for="` + label + `"> Settings ` + label + `</label>
	<select class="form-control" id="` + label + `">`
	for i, item := range modes {
		if selected != i {
			output += "<option>" + item + "</option>"
		} else {
			output += "<option selected>" + item + "</option>"
		}
	}
	output += "</select>"
	return output
}

var indexHTML = `
<!doctype html>
<html>
	<head>
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<!-- Latest compiled and minified CSS -->
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css">
		
		<!-- jQuery library -->
		<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js"></script>
		<!-- Popper JS -->
		<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.3/umd/popper.min.js"></script>
		<!-- Latest compiled JavaScript -->
		<script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.1.3/js/bootstrap.min.js"></script>

	</head>
	<body>
	<div class="container">
		<!-- Nav tabs -->
		<ul class="nav nav-tabs" role="tablist">
			<li class="nav-item">
				<a class="nav-link active" data-toggle="tab" href="#knob1">Knob 1</a>
			</li>
			<li class="nav-item">
				<a class="nav-link" data-toggle="tab" href="#knob2">Knob 2</a>
			</li>
			<li class="nav-item">
				<a class="nav-link" data-toggle="tab" href="#knob3">Knob 3</a>
			</li>
			<li class="nav-item">
				<a class="nav-link" data-toggle="tab" href="#knob4">Knob 4</a>
			</li>
		</ul>

		<!-- Tab panes -->
		<div class="tab-content">
			<div id="knob1" class ="container tab-pane active">
				<div class="form-group col-sm-6">
					` + populateHTML("Knob 1", knobModes, 0) + `
					` + populateHTML("Button 1", buttonModes, 3) + `
				</div>
			</div>
			<div id="knob2" class ="container tab-pane fade">
				<div class="form-group col-sm-6">
					` + populateHTML("Knob 2", knobModes, 1) + `
					` + populateHTML("Button 2", buttonModes, 1) + `
				</div>
			</div>
			<div id="knob3" class ="container tab-pane fade">		
				<div class="form-group col-sm-6">
					` + populateHTML("Knob 3", knobModes, 1) + `
					` + populateHTML("Button 3", buttonModes, 2) + `
				</div>
			</div>
			<div id="knob4" class ="container tab-pane fade">
				<div class="form-group col-sm-6">
					` + populateHTML("Knob 4", knobModes, 1) + `
					` + populateHTML("Button 1", buttonModes, 0) + `
				</div>
			</div>
		</div>
	</div>
	</body>
</html>`
