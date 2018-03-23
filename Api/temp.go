package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	Nocheck = iota
	Checking
	NotPass
	NotClear
	SendSuccess
	Error
)

type Status byte

type Result struct {
	Status string `json:"status"`
}

type SmsCode struct {
	Number        string `json:"number,string"`
	Sms_signature string `json:"sms_signature,string"`
}

type SmsCodeValidation struct {
	Number string
	Code   string
}

func (self *Result) SetStatus(s Status) {
	self.Status = s.String()
}

func (s Status) String() string {
	strings := []string{"NoCheck", "Checking", "NotPass", "NotClear", "SendSuccess", "Error"}
	return strings[s]
}

func Log(v ...interface{}) {
	log.Println(v...)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		result, err := ioutil.ReadAll(r.Body)
		if err != nil {
			Log(err.Error())
		}
		r.Body.Close()

		var sms SmsCode

		json.Unmarshal(result, &sms)

		Log("smscode: ", sms)
		Log("number： ", sms.Number)
		Log("sms_signature:", sms.Sms_signature)
		// Log("Number:", r.Form["number"])
		// Log("sms_signature:", r.Form["sms_signature"])

		io.WriteString(w, "POST")
	} else {
		io.WriteString(w, "Other")
	}

}

func SendService(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		req, err := ioutil.ReadAll(r.Body)
		if err != nil {
			Log(err.Error())
		}
		r.Body.Close()

		sms := SmsCode{}

		json.Unmarshal(req, &sms)
		Log("number： ", sms.Number)
		Log("sms_signature:", sms.Sms_signature)

		Reqs := Result{}
		Reqs.SetStatus(Checking)
		b, _ := json.Marshal(Reqs)
		io.WriteString(w, string(b))
	} else {
		io.WriteString(w, "Other")
	}
}

func ValidationService(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		req, err := ioutil.ReadAll(r.Body)
		if err != nil {
			Log(err.Error())
		}
		r.Body.Close()

		smsValidation := SmsCodeValidation{}

		json.Unmarshal(req, &smsValidation)
		Log("number： ", smsValidation.Number)
		Log("code:", smsValidation.Code)

		Reqs := Result{}
		Reqs.SetStatus(NotClear)
		b, _ := json.Marshal(Reqs)
		Log("Reqs:", string(b[:]))

		//io.WriteString(w, b)
	} else {
		io.WriteString(w, "Other")
	}
}

// {
// 	"number": "13808761543",
// 	"sms_signature":"【XXXX】"
// }

func main() {
	http.HandleFunc("/api/smscode", SendService)
	http.HandleFunc("/api/smscodeconfirm", ValidationService)
	log.Fatal(http.ListenAndServe(":8080", nil))

	Log("continue..")
}
