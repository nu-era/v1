package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

/*
// Send sends a twilio text message to the provided number from the
// number provided in 'numberFrom'. The message body is determined
// by 'msgBody'.
// NOTE: Since we are using a trial version of twilio, there will be
*/
func Send(numberTo string, numberFrom string, msgBody string) error {
	fmt.Println("Begginning to send twilio msg...")
	msgData := url.Values{}
	msgData.Set("To", numberTo)
	msgData.Set("From", numberFrom)
	msgData.Set("Body", msgBody)
	msgDataReader := *strings.NewReader(msgData.Encode())

	// Create HTTP request client
	client := &http.Client{}
	req, err := http.NewRequest("POST", twilURLString, &msgDataReader)

	if err != nil {
		fmt.Println("Error creating twilio request: ", err)
	}

	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make HTTP POST request and return message SID
	resp, err := client.Do(req)
	fmt.Println("Response: ", resp)
	fmt.Println("Error: ", err.Error)
	if err == nil {
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var data map[string]interface{}
			decoder := json.NewDecoder(resp.Body)
			err := decoder.Decode(&data)
			if err == nil {
				fmt.Println(data["sid"])
			}
		} else {
			fmt.Println(resp.Status)
		}
	} else {
		fmt.Println("Error getting twilio response: ", err)
	}

	return nil
}
