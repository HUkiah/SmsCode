package HttpHandler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

//base
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
