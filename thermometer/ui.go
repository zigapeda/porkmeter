package main

import (
	"github.com/zigapeda/raspui"
	"fmt"
	"os"
	"log"
	"strconv"
	"time"
)

func initUI() {
	if _, err := os.Stat("calibration.conf"); err != nil {
		fmt.Println("start calibration")
		raspui.Calibrate()
	}
	err := raspui.CreateRaspUI()
	if err != nil {
		log.Fatal(err)
	}
}

func closeUI() {
	raspui.CloseRaspUI()
}

func timeLoop() {
	for {
		tb := raspui.CreateTextbox(256, 0, 64, 16, time.Now().Format("15:04:05"))
		raspui.AddElement(tb)
		time.Sleep(1 * time.Second)
		raspui.RemoveElement(tb)
	}
}

func createNoConfigUI() {
	raspui.AddElement(raspui.CreateTextbox(10, 20, 300, 20, "Konnte keine Konfiguration finden!"))
	raspui.AddElement(raspui.CreateButton(20, 65, 280, 40, "   Jetzt Konfiguration anlegen   ", func() {go createLimitUI()}))
	interfaces, err := getInterfaces()
	if err != nil {
		log.Println(err)
	}
	for i, ife := range interfaces {
		raspui.AddElement(raspui.CreateTextbox(5, 215 - (i * 25), 310, 20, ife))
	}
}

func createLimitUI() {
	time.Sleep(500 * time.Millisecond)
	go raspui.Clear()
	time.Sleep(500 * time.Millisecond)
	
	txtgmin := raspui.CreateTextbox(5, 5, 310, 20, "Garraum Mintemp.: 100 °C")
	slidgmin := raspui.CreateSlider(5, 25, 310, 20, 10, 30)
	txtgmax := raspui.CreateTextbox(5, 50, 310, 20, "Garraum Maxtemp.: 130 °C")
	slidgmax := raspui.CreateSlider(5, 70, 310, 20, 10, 30)
	txtf1 := raspui.CreateTextbox(5, 95, 310, 20, "Fleisch Zieltemp. 1: 95 °C")
	slidf1 := raspui.CreateSlider(5, 115, 310, 20, 10, 30)
	txtf2 := raspui.CreateTextbox(5, 140, 310, 20, "Fleisch Zieltemp. 2: 95 °C")
	slidf2 := raspui.CreateSlider(5, 160, 310, 20, 10, 30)
	raspui.AddElement(txtgmin)
	raspui.AddElement(slidgmin)
	raspui.AddElement(txtgmax)
	raspui.AddElement(slidgmax)
	raspui.AddElement(txtf1)
	raspui.AddElement(slidf1)
	raspui.AddElement(txtf2)
	raspui.AddElement(slidf2)
	slidgmin.SetValue(20)
	slidgmax.SetValue(26)
	slidf1.SetValue(19)
	slidf2.SetValue(19)
	slidgmin.SetChangeFunc(func(v int) {txtgmin.SetText("Garraum Mintemp.: " + strconv.Itoa(v * 5) + " °C")})
	slidgmax.SetChangeFunc(func(v int) {txtgmax.SetText("Garraum Maxtemp.: " + strconv.Itoa(v * 5) + " °C")})
	slidf1.SetChangeFunc(func(v int) {txtf1.SetText("Fleisch Zieltemp. 1: " + strconv.Itoa(v * 5) + " °C")})
	slidf2.SetChangeFunc(func(v int) {txtf2.SetText("Fleisch Zieltemp. 2: " + strconv.Itoa(v * 5) + " °C")})
	raspui.AddElement(raspui.CreateButton(20, 190, 280, 35, "             Weiter             ", func() {
		config.Gmin = slidgmin.GetValue() * 5
		config.Gmax = slidgmax.GetValue() * 5
		config.F1 = slidf1.GetValue() * 5
		config.F2 = slidf2.GetValue() * 5
		config.Meters = make([]Meter, 0, 8)
		go createAssignmentUI(1)
	}))
}

func createAssignmentUI(channel int) {
	time.Sleep(500 * time.Millisecond)
	go raspui.Clear()
	time.Sleep(500 * time.Millisecond)

	textbox := raspui.CreateTextbox(5, 5, 310, 20, "Zuweisung für Kanal " + strconv.Itoa(channel) + ":")
	assignChannel := func(as string) {
		if as != "o" {
			meter := Meter{}
			meter.Type = as
			meter.PhysID = "C" + strconv.Itoa(channel)
			meter.Channel = channel - 1
			config.Meters = append(config.Meters, meter)
		}
		channel += 1
		if channel < 9 {
			go createAssignmentUI(channel)
		} else {
			go createTempUI("t1")
		}
	}
	raspui.AddElement(textbox)
	raspui.AddElement(raspui.CreateButton(15, 40, 137, 85, "Offen", func() {assignChannel("o")}))
	raspui.AddElement(raspui.CreateButton(15, 140, 137, 85, "Garraum", func() {assignChannel("g")}))
	raspui.AddElement(raspui.CreateButton(167, 40, 138, 85, "Fleisch 1", func() {assignChannel("f1")}))
	raspui.AddElement(raspui.CreateButton(167, 140, 138, 85, "Fleisch 2", func() {assignChannel("f2")}))
}

func createTempUI(which string) {
	time.Sleep(500 * time.Millisecond)
	go raspui.Clear()
	time.Sleep(500 * time.Millisecond)

	temp := 25
	text := "Aktuelle Umgebungstemperatur: "
	if which == "t2" {
		temp = 100
		text = "Kalibrierungstemperatur: "
	}
	textbox := raspui.CreateTextbox(10, 20, 300, 20, text + strconv.Itoa(temp) + " °C")
	adjustTemp := func(how string) {
		if how == "u" {
			temp += 1
		} else if how == "d" {
			temp -= 1
		}
		textbox.SetText(text + strconv.Itoa(temp) + " °C")
	}
	raspui.AddElement(textbox)
	raspui.AddElement(raspui.CreateButton(15, 50, 137, 80, "-", func() {adjustTemp("d")}))
	raspui.AddElement(raspui.CreateButton(167, 50, 138, 80, "+", func() {adjustTemp("u")}))
	raspui.AddElement(raspui.CreateButton(15, 145, 290, 80, "Jetzt Kalibrieren...", func() {
		if which == "t1" {
			config.T1 = float64(temp)
		} else {
			config.T2 = float64(temp)
		}
		go createReadResistanceUI(which)
	}))
}

func createReadResistanceUI(which string) {
	time.Sleep(500 * time.Millisecond)
	go raspui.Clear()
	time.Sleep(500 * time.Millisecond)

	textbox := raspui.CreateTextbox(10, 20, 300, 20, "")
	raspui.AddElement(textbox)

	for i, v := range config.Meters {
		textbox.SetText("Lese Widerstand von Kanal " + strconv.Itoa(v.Channel + 1) + "...")
		r := readResistance(v.Channel)
		if which == "t1" {
			config.Meters[i].R1 = r
		} else {
			config.Meters[i].R2 = r
			config.Meters[i].B = getBValue(config.T1, config.Meters[i].R1, config.T2, config.Meters[i].R2)
		}
	}

	if which == "t1" {
		go createTempUI("t2")
	} else {
		writeConfiguration()
		fmt.Println("fertig")
		go createMenuUI()
	}
}

func createMenuUI() {
	time.Sleep(500 * time.Millisecond)
	go raspui.Clear()
	time.Sleep(500 * time.Millisecond)

	raspui.AddElement(raspui.CreateTextbox(10, 20, 300, 20, "Porken!"))
	interfaces, err := getInterfaces()
	if err != nil {
		log.Println(err)
	}
	for i, ife := range interfaces {
		raspui.AddElement(raspui.CreateTextbox(5, 215 - (i * 25), 310, 20, ife))
	}
}

//func createCalibrationUI() {
//	time.Sleep(500 * time.Millisecond)
//	go raspui.Clear()
//	time.Sleep(500 * time.Millisecond)
//	
//	raspui.AddElement(raspui.CreateButton(15, 40, 137, 85, "Konfigurieren", func() {}))
//	raspui.AddElement(raspui.CreateButton(15, 140, 137, 85, "Kalibrieren", func() {}))
//	raspui.AddElement(raspui.CreateButton(167, 40, 138, 85, "Starten", func() {}))
//	raspui.AddElement(raspui.CreateButton(167, 140, 138, 85, "Herunterfahren", func() {}))
//}