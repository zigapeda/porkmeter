// bratenthermometerserver project main.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Temps struct {
	Date time.Time
	C1   float64
	C2   float64
	C3   float64
	C4   float64
	C5   float64
	C6   float64
	C7   float64
	C8   float64
}

var (
	temps []Temps
)

func setTemps(vals url.Values) error {
	t := Temps{Date: time.Now()}
	var err error
	if t.C1, err = strconv.ParseFloat(vals.Get("C1"), 64); err != nil {
		return err
	}
	if t.C2, err = strconv.ParseFloat(vals.Get("C2"), 64); err != nil {
		return err
	}
	if t.C3, err = strconv.ParseFloat(vals.Get("C3"), 64); err != nil {
		return err
	}
	if t.C4, err = strconv.ParseFloat(vals.Get("C4"), 64); err != nil {
		return err
	}
	if t.C5, err = strconv.ParseFloat(vals.Get("C5"), 64); err != nil {
		return err
	}
	if t.C6, err = strconv.ParseFloat(vals.Get("C6"), 64); err != nil {
		return err
	}
	if t.C7, err = strconv.ParseFloat(vals.Get("C7"), 64); err != nil {
		return err
	}
	if t.C8, err = strconv.ParseFloat(vals.Get("C8"), 64); err != nil {
		return err
	}
	temps = append(temps, t)
	return nil
}

func getTemps() (Temps, error) {
	if len(temps) > 0 {
		return temps[len(temps)-1], nil
	}
	return Temps{}, errors.New("no temperatures")
}

func executeUrl(reqUrl *url.URL) (interface{}, error) {
	switch reqUrl.Path[5:] {
	case "GetTemps":
		return getTemps()
	case "SetTemps":
		setTemps(reqUrl.Query())
		return "ok", nil
	}
	return nil, errors.New("api not found")
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	success, err := executeUrl(r.URL)
	if err != nil {
		jsonString, err := json.Marshal(map[string]interface{}{"success": nil, "error": err.Error()})
		if err != nil {
			fmt.Println(1, err)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write(jsonString)
		}
	} else {
		jsonString, err := json.Marshal(map[string]interface{}{"success": success, "error": nil})
		if err != nil {
			fmt.Println(2, err)
		} else {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write(jsonString)
		}
	}
}

func main() {
	temps = make([]Temps, 0, 100)
	//	temps = append(temps, Temps{Date: time.Now(), C1: 11.1, C2: 11.2, C3: 11.3, C4: 11.4, C5: 11.5, C6: 11.6, C7: 11.7, C8: 11.8})
	//	temps = append(temps, Temps{Date: time.Now(), C1: 12.1, C2: 12.2, C3: 12.3, C4: 12.4, C5: 12.5, C6: 12.6, C7: 12.7, C8: 12.8})

	http.Handle("/", http.FileServer(http.Dir("html")))

	http.HandleFunc("/api/", apiHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
