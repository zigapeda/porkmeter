package main

import (
	"net"
	"strings"
)

func getInterfaces() ([]string, error) {
	retval := make([]string, 0, 3)
	ifarr, err := net.Interfaces()
	if err != nil {
		return retval, err
	}
	for _, ife := range ifarr {
		if ife.Flags & net.FlagUp != 0  && ife.Flags & net.FlagLoopback == 0 {
			adresses, err := ife.Addrs()
			if err != nil {
				return retval, err
			}
			for _, a := range adresses {
				if strings.Contains(a.String(), ":") == false {
					adr := ife.Name + ": " + strings.Split(a.String(), "/")[0]
					retval = append(retval, adr)
				}
			}
		}
	}
	return retval, nil
}