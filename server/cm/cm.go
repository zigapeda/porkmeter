package cm

import (
    "github.com/NaySoftware/go-fcm"
)

type CM struct {
	serverKey string
}

func NewCM(serverKey string) CM {
	return CM{serverKey: serverKey}
}

func (cm CM) Send(data map[string]string, ids []string) error {

	c := fcm.NewFcmClient(cm.serverKey)
    c.NewFcmRegIdsMsg(ids, data)


	_, err := c.Send()


	if err != nil {
		return err
	}

	return nil

}
