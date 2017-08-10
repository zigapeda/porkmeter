package cm

import (
	"fmt"
    "github.com/NaySoftware/go-fcm"
)

type CM struct {
	serverKey string
}

func NewCM(serverKey string) CM {
	return CM{serverKey: serverKey}
}

func (cm CM) Send(data interface{}, ids []string) error {

	c := fcm.NewFcmClient(cm.serverKey)
    c.NewFcmRegIdsMsg(ids, data)


	status, err := c.Send()

	fmt.Println(status)


	if err != nil {
		return err
	}

	return nil

}
