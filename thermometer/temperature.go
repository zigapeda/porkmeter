package main

import (
	"github.com/zigapeda/raspui"
	"fmt"
	"math"
	"time"
	"sort"
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
	fmt.Println(adc, un, rn)
	return rn
}

func readTemps() {
	t25 := 25.0
	r25 := 1000000.0
	t100 := 100.0
	r100 := 47000.0
	b := getBValue(t25, r25, t100, r100)
	for {
		adc := readADCValue(0)
		un := getNTCVoltage(adc)
		rn := getNTCResistance(un)
		t := getTemperature(rn, r25, t25, b)
		fmt.Println(t)
		//time.Sleep(1 * time.Second)
		// readChannel(1)
		// time.Sleep(1 * time.Second)
		// readChannel(2)
		// time.Sleep(1 * time.Second)
		// readChannel(3)
		// time.Sleep(1 * time.Second)
		// readChannel(4)
		// time.Sleep(1 * time.Second)
		// readChannel(5)
		// time.Sleep(1 * time.Second)
		// readChannel(6)
		// time.Sleep(1 * time.Second)
		// readChannel(7)
		// time.Sleep(1 * time.Second)
	}
}