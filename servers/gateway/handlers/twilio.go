package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}

var accountSid = os.Getenv("TWILIO_ACCOUNT_SID")
var authToken = os.Getenv("TWILIO_AUTH_TOKEN")
var twilURLString = "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
var twilCheckVerifyURLString = "https://verify.twilio.com/v2/Services/" + accountSid + "/VerificationCheck"

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

/*
Verify takes a number to send an sms message to, along with a number to send
it from, as well as a string
*/
func Verify(numberTo string, numberFrom string, msgBody string) (string, error) {
	fmt.Println("Begginning to send twilio verification msg...")
	msgData := url.Values{}
	msgData.Set("Channel", "sms")
	msgData.Set("To", numberTo)
	msgData.Set("serviceSid", serviceSID)
	msgData.Set("country_code", "1")
	msgData.Set("locale", "en")

	msgDataReader := *strings.NewReader(msgData.Encode())

	// Create HTTP request client
	client := &http.Client{}
	req, err := http.NewRequest("POST", twilAuthString, &msgDataReader)

	if err != nil {
		fmt.Println("Error creating twilio request: ", err)
		return "", err
	}

	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make HTTP POST request and return message SID
	resp, err := client.Do(req)

	// Save a copy of this request for debugging.
	// requestDump, err := httputil.DumpRequest(req, true)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// //fmt.Println(string(requestDump))
	if err == nil {
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var data map[string]string
			decoder := json.NewDecoder(resp.Body)
			err := decoder.Decode(&data)
			if err == nil {
				fmt.Println(data["sid"])
				return data["sid"], nil
			} else {
				return "", err
			}
		} else {
			responseDump, err := httputil.DumpResponse(resp, true)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(responseDump))
			return resp.Status, nil
		}
	} else {
		fmt.Println("Error getting twilio response: ", err)
		return "", err
	}
}

func CheckVerification(code string, phoneTo string) error {
	// Beginning to send twilio verification check message
	msgData := url.Values{}
	msgData.Set("Code", code)
	msgData.Set("To", phoneTo)
	msgData.Set("serviceSid", serviceSID)
	msgData.Set("country_code", "1")
	msgData.Set("locale", "en")

	msgDataReader := *strings.NewReader(msgData.Encode())

	// Create HTTP request client
	client := &http.Client{}
	req, err := http.NewRequest("POST", twilCheckVerifyURLString, &msgDataReader)

	if err != nil {
		fmt.Println("Error creating twilio request: ", err)
		return err
	}

	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make HTTP POST request and return message SID
	resp, err := client.Do(req)
	// Save a copy of this request for debugging.
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))
	if err == nil {
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var data map[string]interface{}
			decoder := json.NewDecoder(resp.Body)
			err := decoder.Decode(&data)
			if err == nil {
				fmt.Println(data["sid"])
			}
		} else {
			responseDump, err := httputil.DumpResponse(resp, true)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(responseDump))
		}
		return nil
	} else {
		fmt.Println("Error getting twilio response: ", err)
		return err
	}
}
