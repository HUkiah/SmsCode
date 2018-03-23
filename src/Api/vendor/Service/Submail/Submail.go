package submail

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

//messagexsend
type MessageXSend struct {
	to          []string
	addressbook []string
	project     string
	vars        map[string]string
}

func HttpGet(queryurl string) string {
	u, _ := url.Parse(queryurl)
	retstr, err := http.Get(u.String())
	if err != nil {
		return err.Error()
	}
	result, err := ioutil.ReadAll(retstr.Body)
	retstr.Body.Close()
	if err != nil {
		return err.Error()
	}
	return string(result)
}

func HttpPost(queryurl string, postdata map[string]string) string {
	data, err := json.Marshal(postdata)
	if err != nil {
		return err.Error()
	}

	body := bytes.NewBuffer([]byte(data))

	retstr, err := http.Post(queryurl, "application/json;charset=utf-8", body)

	if err != nil {
		return err.Error()
	}
	result, err := ioutil.ReadAll(retstr.Body)
	retstr.Body.Close()
	if err != nil {
		return err.Error()
	}
	return string(result)
}

func GetTimeStamp() string {
	resp := HttpGet("https://api.submail.cn/service/timestamp.json")
	var dict map[string]interface{}
	err := json.Unmarshal([]byte(resp), &dict)
	if err != nil {
		return err.Error()
	}
	return strconv.Itoa(int(dict["timestamp"].(float64)))
}

func CreateSignatrue(request map[string]string, config map[string]string) string {
	appkey := config["appkey"]
	appid := config["appid"]
	signtype := config["signtype"]
	request["sign_type"] = signtype
	keys := make([]string, 0, 32)
	for key, _ := range request {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	str_list := make([]string, 0, 32)
	for _, key := range keys {
		str_list = append(str_list, fmt.Sprintf("%s=%s", key, request[key]))
	}
	sigstr := strings.Join(str_list, "&")
	sigstr = fmt.Sprintf("%s%s%s%s%s", appid, appkey, sigstr, appid, appkey)
	if signtype == "normal" {
		return appkey
	} else if signtype == "md5" {
		mymd5 := md5.New()
		io.WriteString(mymd5, sigstr)
		return fmt.Sprintf("%x", mymd5.Sum(nil))
	} else {
		mysha1 := sha1.New()
		io.WriteString(mysha1, sigstr)
		return fmt.Sprintf("%x", mysha1.Sum(nil))
	}
}

func MessageXSendAddTo(messagexsend *MessageXSend, address string) {
	messagexsend.to = append(messagexsend.to, address)
}

func MessageXSendSetProject(messagexsend *MessageXSend, project string) {
	messagexsend.project = project
}

func MessageXSendAddVar(messagexsend *MessageXSend, key string, value string) {
	messagexsend.vars[key] = value
}

func MessageXSendBuildRequest(messagexsend *MessageXSend) map[string]string {
	request := make(map[string]string)
	if len(messagexsend.to) != 0 {
		request["to"] = strings.Join(messagexsend.to, ",")
	}
	if len(messagexsend.addressbook) != 0 {
		request["addressbook"] = strings.Join(messagexsend.addressbook, ",")
	}

	if messagexsend.project != "" {
		request["project"] = messagexsend.project
	}

	if len(messagexsend.vars) != 0 {
		data, err := json.Marshal(messagexsend.vars)
		if err == nil {
			request["vars"] = string(data)
		}
	}
	return request
}

func MessageXSendRun(request map[string]string, config map[string]string) string {
	url := "https://api.submail.cn/message/xsend.json"
	request["appid"] = config["appid"]
	request["timestamp"] = GetTimeStamp()
	request["signature"] = CreateSignatrue(request, config)
	return HttpPost(url, request)
}

func CreateMessageXSend() *MessageXSend {
	messagexsend := new(MessageXSend)
	messagexsend.vars = make(map[string]string)
	return messagexsend
}
