// bratenthermometerserver project main.go
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-yaml/yaml"
	"github.com/unrolled/render"
	"github.com/zigapeda/porkmeter/server/cm"
	"github.com/zigapeda/porkmeter/server/db"
)

var configfile *string = flag.String("conf", "config.yaml", "sets the config file used")

type Config struct {
	Meters     []*Meter    `yaml:"meters"`
	Limits     []Limit     `yaml:"limits"`
	PushAPIKey string      `yaml:"pushapikey"`
	DB         db.Dbconfig `yaml:"dbconfig"`
}

type Meter struct {
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	PhysID string `yaml:"physid"`
}

type Limit struct {
	Min  float64 `yaml:"min"`
	Max  float64 `yaml:"max"`
	Type string  `yaml:"type"`
}

type Temps struct {
	Time  time.Time
	Temps []Temp
}

type Temp struct {
	Meter *Meter
	Temp  float64
}

var (
	temps  []Temps
	config Config
	rend   = render.New()
	pcms   *cm.CM
	pkeys  []string
)

func setTemps(vals url.Values) error {
	t := make([]Temp, 0, 5)
	var err error
	for _, element := range config.Meters {
		var tv float64
		if tv, err = strconv.ParseFloat(vals.Get(element.PhysID), 64); err != nil {
			return err
		}
		t = append(t, Temp{Meter: element, Temp: tv})
	}
	ts := Temps{Time: time.Now(), Temps: t}
	fmt.Println("set temps")
	temps = append(temps, ts)
	if len(pkeys) > 0 && pcms != nil {
		fmt.Println("send notification to", pkeys)
		err = pcms.Send(map[string]string{}, pkeys)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func registerKey(vals url.Values) error {
	key := vals.Get("key")
	if key == "" {
		return errors.New("no key parameter")
	}
	fmt.Println("register Key", key)
	for _, pkey := range pkeys {
		if pkey == key {
			fmt.Println("Key already registered")
			return errors.New("key already registered")
		}
    }
	pkeys = append(pkeys, key)
	return nil
}

func removeKey(vals url.Values) error {
	key := vals.Get("key")
	if key == "" {
		return errors.New("no key parameter")
	}
	fmt.Println("remove Key", key)
	for i, pkey := range pkeys {
		if pkey == key {
			pkeys = append(pkeys[:i], pkeys[i+1:]...)
			return nil
		}
    }
	return errors.New("key not found")
}

func checkKey(vals url.Values) (string, error) {
	key := vals.Get("key")
	if key == "" {
		return "", errors.New("no key parameter")
	}
	fmt.Println("check Key", key)
	for _, pkey := range pkeys {
		if pkey == key {
			return "on", nil
		}
    }
	return "off", nil
}

func getTemps() (Temps, error) {
	if len(temps) > 0 {
		return temps[len(temps)-1], nil
	}
	return Temps{}, errors.New("no temperatures")
}

func apiGetTemps(w http.ResponseWriter, r *http.Request) {
	temps, err := getTemps()
	if err != nil {
		rend.JSON(w, 200, map[string]interface{}{"success": nil, "error": err.Error()})
	} else {
		rend.JSON(w, 200, map[string]interface{}{"success": temps, "error": nil})
	}
}

func apiSetTemps(w http.ResponseWriter, r *http.Request) {
	err := setTemps(r.URL.Query())
	if err != nil {
		rend.JSON(w, 200, map[string]interface{}{"success": nil, "error": err.Error()})
	} else {
		rend.JSON(w, 200, map[string]interface{}{"success": "ok", "error": nil})
	}
}

func apiRegisterKey(w http.ResponseWriter, r *http.Request) {
	err := registerKey(r.URL.Query())
	if err != nil {
		rend.JSON(w, 200, map[string]interface{}{"success": nil, "error": err.Error()})
	} else {
		rend.JSON(w, 200, map[string]interface{}{"success": "ok", "error": nil})
	}
}

func apiRemoveKey(w http.ResponseWriter, r *http.Request) {
	err := removeKey(r.URL.Query())
	if err != nil {
		rend.JSON(w, 200, map[string]interface{}{"success": nil, "error": err.Error()})
	} else {
		rend.JSON(w, 200, map[string]interface{}{"success": "ok", "error": nil})
	}
}

func apiCheckKey(w http.ResponseWriter, r *http.Request) {
	str, err := checkKey(r.URL.Query())
	if err != nil {
		rend.JSON(w, 200, map[string]interface{}{"success": nil, "error": err.Error()})
	} else {
		rend.JSON(w, 200, map[string]interface{}{"success": str, "error": nil})
	}
}

func main() {

	flag.Parse()

	configbytes, err := ioutil.ReadFile(*configfile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configbytes, &config)
	if err != nil {
		panic(err)
	}
	spew.Dump(config)

	if len(config.PushAPIKey) > 0 {
		pushnotificationserver := cm.NewCM(config.PushAPIKey)
		pcms = &pushnotificationserver
	} else {
		pcms = nil
	}

	temps = make([]Temps, 0, 100)
	//	temps = append(temps, Temps{Date: time.Now(), C1: 11.1, C2: 11.2, C3: 11.3, C4: 11.4, C5: 11.5, C6: 11.6, C7: 11.7, C8: 11.8})
	//	temps = append(temps, Temps{Date: time.Now(), C1: 12.1, C2: 12.2, C3: 12.3, C4: 12.4, C5: 12.5, C6: 12.6, C7: 12.7, C8: 12.8})

	http.Handle("/", http.FileServer(http.Dir("html")))

	http.HandleFunc("/api/GetTemps", apiGetTemps)
	http.HandleFunc("/api/SetTemps", apiSetTemps)
	http.HandleFunc("/api/RegisterKey", apiRegisterKey)
	http.HandleFunc("/api/RemoveKey", apiRemoveKey)
	http.HandleFunc("/api/CheckKey", apiCheckKey)

	err = http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
