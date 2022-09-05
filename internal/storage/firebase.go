package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type FirebaseNotification struct {
	Notification FirebaseNotificationData `json:"notification"`
	To           string                   `json:"to"`
}
type FirebaseData struct {
	Data FirebaseDataUpdate `json:"data"`
	To   string             `json:"to"`
}
type FirebaseDataUpdate struct {
	Data   string `json:"data"`
	DevEui string `json:"deveui"`
}
type FirebaseNotificationData struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Time  int    `json:"time_to_live"`
	Delay bool   `json:"delay_while_idle"`
}

type OneSignalNotification struct {
	Ids               []string                  `json:"include_external_user_ids"`
	AppId             string                    `json:"app_id"`
	Headings          OneSignalNotificationEN   `json:"headings"`
	Contents          OneSignalNotificationEN   `json:"contents"`
	Data              OneSignalNotificationData `json:"data"`
	AndroidVisibility int                       `json:"android_visibility"`
	Priority          int                       `json:"priority"`
}
type OneSignalNotificationEN struct {
	En string `json:"en"`
}
type OneSignalNotificationData struct {
	Priority int `json:"priority"`
}

var firebaseAythKey = "AAAA5h0bGnM:APA91bHFEwqNn8auXh64E_z_cltvqmrPa6OygwVUQfmGctyuINkThNmNpBRT2X43yAByAn04MFI03oVYhYpMzU5gXYh2QZOI3oQh4NsQiGGTxdwIv20aoISOQQiOkaCVK8mTx-Eq8A5E"
var OneSignalAythKey = "Basic YzU2Yzg4NGMtZjQ2Yy00Nzg4LWFkNjYtNDNjNGI2YTM1MDgy"

// SendFirebaseNotification SendFirebaseNotification
func SendFirebaseNotification(u User, f FirebaseNotificationData) error {
	client := &http.Client{}

	if u.WebKey != "" {
		notification := FirebaseNotification{
			Notification: f,
			To:           u.WebKey,
		}
		var jsonBody []byte
		jsonBody, err := json.Marshal(notification)
		if err != nil {
			fmt.Println("json marshal error")
		}
		req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Println("NewRequest error")
		}
		req.Header.Set("Authorization", "key="+firebaseAythKey)
		req.Header.Set("Content-Type", "application/json")

		_, err = client.Do(req)
		if err != nil {
			fmt.Println("http error")
		}
	}
	if u.AndroidKey != "" {
		notification := FirebaseNotification{
			Notification: f,
			To:           u.AndroidKey,
		}
		var jsonBody []byte
		jsonBody, err := json.Marshal(notification)
		if err != nil {
			fmt.Println("json marshal error")
		}
		req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Println("NewRequest error")
		}
		req.Header.Set("Authorization", "key="+firebaseAythKey)
		req.Header.Set("Content-Type", "application/json")

		_, err = client.Do(req)
		if err != nil {
			fmt.Println("http error")
		}
		// One Signal Android

		// CreateOnseSignalNotification(strconv.FormatInt(u.ID, 10), f.Body)
		oneSignalNotification := OneSignalNotification{
			AppId:             "526c06f9-cac6-48f1-a092-96504bb779e1",
			Headings:          OneSignalNotificationEN{En: "Vaps"},
			Contents:          OneSignalNotificationEN{En: f.Body},
			Ids:               []string{strconv.FormatInt(u.ID, 10)},
			Data:              OneSignalNotificationData{Priority: 10},
			AndroidVisibility: 1,
			Priority:          10,
		}
		var json0 []byte
		json0, err = json.Marshal(oneSignalNotification)
		if err != nil {
			fmt.Println("json marshal error")
		}
		fmt.Println("JSOOON")
		fmt.Println(bytes.NewBuffer(json0))

		reqOne, err := http.NewRequest("POST", "https://onesignal.com/api/v1/notifications", bytes.NewBuffer(json0))
		reqOne.Header.Set("Authorization", OneSignalAythKey)
		reqOne.Header.Set("Content-Type", "application/json")
		reqOne.Header.Set("Accept", "application/json")

		_, err = client.Do(reqOne)
		if err != nil {
			fmt.Println("http error")

		}
	}
	return nil
}
