package main

import (
	"github.com/zigapeda/raspui"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"
	"sort"
	"strconv"
)

var (
	pinCS = 2 //6
	pinCLK = 3//7
	pinDO = 4 //9
	pinDI = 14//8
)

const bufferSize = 20
const ue = 5.0
const steps = 1024.0
const ups = ue / steps
const r1 = 47050

func initGpio() {
	raspui.SetGpioOutput(pinCS)
	raspui.SetGpioOutput(pinCLK)
	raspui.SetGpioOutput(pinDO)
	raspui.SetGpioInput(pinDI)

	raspui.SetGpio(pinCS, 1)
	raspui.SetGpio(pinCLK, 0)
	raspui.SetGpio(pinDO, 0)
}

func sendBit(bit int) {
	raspui.SetGpio(pinDO, bit)
	raspui.SetGpio(pinCLK, 1)
	time.Sleep(2 * time.Millisecond)
	raspui.SetGpio(pinCLK, 0)
	time.Sleep(2 * time.Millisecond)
}

func readBit() int {
	bit := raspui.GetGpio(pinDI)
	raspui.SetGpio(pinCLK, 1)
	time.Sleep(2 * time.Millisecond)
	raspui.SetGpio(pinCLK, 0)
	time.Sleep(2 * time.Millisecond)
	return bit
}

func readChannel(ch int) int {
	raspui.SetGpio(pinCLK, 0)
	raspui.SetGpio(pinDO, 0)
	raspui.SetGpio(pinCS, 0)
	time.Sleep(2 * time.Millisecond)

	//startbit
	sendBit(1)

	//mode
	sendBit(1)

	//data
	sendBit(ch & 0x4 >> 2) //d2
	sendBit(ch & 0x2 >> 1) //d1
	sendBit(ch & 0x1)      //d0

	//sample and hold
	sendBit(0)

	//null bit
	sendBit(0)

	//read
	val := 0
	val |= readBit()
	val = val << 1
	val |= readBit()
	val = val << 1
	val |= readBit()
	val = val << 1
	val |= readBit()
	val = val << 1
	val |= readBit()
	val = val << 1
	val |= readBit()
	val = val << 1
	val |= readBit()
	val = val << 1
	val |= readBit()
	val = val << 1
	val |= readBit()
	val = val << 1
	val |= readBit()
	// val = val << 1
	// val |= readBit()
	// val = val << 1
	// val |= readBit()

	raspui.SetGpio(pinCS, 1)
	time.Sleep(2 * time.Millisecond)

	return val
}

func getBValue(t1, r1, t2, r2 float64) float64 {
	t1 += 273.15
	t2 += 273.15
	r := math.Log(r1/r2)
	t := (t1 * t2) / (t2 - t1)
	b := r * t
	return b
}

func getNTCVoltage(adcvalue int) float64 {
	un := float64(adcvalue) * ups
	return un
}

func getNTCResistance(un float64) float64 {
	rn := (un * r1) / (ue - un)
	return rn
}

func getTemperature(rn, r25, t25, b float64) float64 {
	//T = 1 / ((lg(rn/r25) / lg(e)) / b + 1 / (t25 + 273.15)) - 273.15
	t25 += 273.15
	l := math.Log10(rn/r25) / math.Log10(math.E)
	//T = 1 / (l / b + 1 / t25) - 273.15
	z := l / b + 1 / t25
	t := 1 / z - 273.15
	return t
}

func readADCValue(channel int) int {
	list := make([]int, bufferSize * 4)
	for i := 0; i < bufferSize * 4; i++ {
		list[i] = readChannel(channel)
	}
	sort.Ints(list)
	sum := 0
	for i := bufferSize; i < bufferSize * 3; i++ {
		sum += list[i]
	}
	sum = sum / (bufferSize * 2)
	return sum
}

func readResistance(channel int) float64 {
	adc := readADCValue(channel)
	un := getNTCVoltage(adc)
	rn := getNTCResistance(un)
	fmt.Println(channel, adc, un, rn)
	return rn
}

func readTemps() {
	for {
		for i, v := range config.Meters {
			config.Meters[i].Reading = true
			r := readResistance(v.Channel)
			t := int(getTemperature(r, v.R1, config.T1, v.B))
			config.Meters[i].Temp = t
			config.Meters[i].Reading = false
		}
		go sendTemps()
	}
}

func sendTemps() {
	u, err := url.Parse("https://porkmeter.maplpapl.de/api/SetTemps")
	if err != nil {
		fmt.Println(err)
	}
	q := u.Query()
	for _, v := range config.Meters {
		q.Add(v.PhysID, strconv.Itoa(v.Temp))
	}
	u.RawQuery = q.Encode()
	_, err = http.Get(u.String())
	if err != nil {
		fmt.Println(err)
	}
}
