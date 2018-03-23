package main

import (
	"Infrastructure/RedisContext"
	"Service/Submail"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	Nocheck = iota
	Checking
	NotPass
	NotClear
	SendSuccess
	Error
)

var Rconn RedisContext.RedisContext

type RespResult struct {
	Status      string `json:"status"`
	Send_id     string `json:"send_id"`
	Fee         int    `json:"fee"`
	Sms_credits string `json:"sms_credits"`
}

type Status int

type Result struct {
	Status string `json:"status"`
}

type SmsCode struct {
	Number        string `json:"number"`
	Sms_signature string `json:"sms_signature"`
}

type SmsCodeValidation struct {
	Number string
	Code   string
}

func (self *Result) SetSendStatus(s Status) {
	self.Status = s.String()
}

func (self *Result) SetValidationStatus(s int) {
	self.Status = strconv.Itoa(s)
}

func (s Status) String() string {
	strings := []string{"NoCheck", "Checking", "NotPass", "NotClear", "SendSuccess", "Error"}
	return strings[s]
}

func SendSmsCode(number, code string) bool {

	messageconfig := make(map[string]string)
	messageconfig["appid"] = "20532"
	messageconfig["appkey"] = "7acca58a1aa29a8edde2bb844fdbc00d"
	messageconfig["signtype"] = "md5"

	//messagexsend
	messagexsend := submail.CreateMessageXSend()
	submail.MessageXSendAddTo(messagexsend, number)
	submail.MessageXSendSetProject(messagexsend, "T1UxU3")
	submail.MessageXSendAddVar(messagexsend, "code", code)
	state := submail.MessageXSendRun(submail.MessageXSendBuildRequest(messagexsend), messageconfig)
	v := RespResult{}
	json.Unmarshal([]byte(state), &v)
	if v.Status == "success" {
		fmt.Println("Success")
		return true
	} else {
		fmt.Println("Fail")
	}
	return false
}

func Log(v ...interface{}) {
	log.Println(v...)
}

func BuildRandom() string {
	rand.Seed(time.Now().UnixNano())
	temp := rand.Int63n(899999) + 100000
	return strconv.FormatInt(temp, 10)
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
		code := BuildRandom()
		status := SendSmsCode(sms.Number, code)

		if status {
			saveStatus := Rconn.Write(sms.Number, code, "600")
			if saveStatus {
				Reqs := Result{}
				Reqs.SetSendStatus(SendSuccess)
				b, _ := json.Marshal(Reqs)
				io.WriteString(w, string(b))
			} else {
				Rconn.Write(sms.Number, code, "1")
			}

		} else {
			Reqs := Result{}
			Reqs.SetSendStatus(Error)
			b, _ := json.Marshal(Reqs)
			io.WriteString(w, string(b))
		}

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

		isExist := Rconn.IsExist(smsValidation.Number)
		Reqs := Result{}
		//存在该号码
		if isExist {
			//通过 0
			status, DbCode := Rconn.Read(smsValidation.Number)

			if status {
				if DbCode == smsValidation.Code {

					Reqs.SetValidationStatus(0)
					b, _ := json.Marshal(Reqs)
					io.WriteString(w, string(b))
				} else {
					Reqs.SetValidationStatus(2)
					b, _ := json.Marshal(Reqs)
					io.WriteString(w, string(b))
				}
			} else {
				Reqs.SetValidationStatus(2)
				b, _ := json.Marshal(Reqs)
				io.WriteString(w, string(b))
			}
		} else {
			//超时或未验证 1
			Reqs.SetValidationStatus(3)
			b, _ := json.Marshal(Reqs)
			io.WriteString(w, string(b))
		}
	} else {
		io.WriteString(w, "Other")
	}
}

func main() {

	Rconn.Open("127.0.0.1:6379")

	defer Rconn.Close()

	http.HandleFunc("/api/smscode", SendService)
	http.HandleFunc("/api/smscodeconfirm", ValidationService)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
