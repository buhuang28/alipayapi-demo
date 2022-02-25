package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func GetRequest(u string, cookie *map[string]string, header, params map[string]string) (bool, HttpResp) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()

	if u == "" {
		return false, HttpResp{}
	}
	request, _ := http.NewRequest("GET", u, nil)

	if cookie != nil && len(*cookie) != 0 {
		if header == nil {
			header = make(map[string]string)
		}
		header["Cookie"] = CookieMap2Str(*cookie)
	}

	if header != nil {
		for k, v := range header {
			request.Header.Set(k, v)
		}
	}

	//加入get参数
	q := request.URL.Query()
	if params != nil {
		for k, v := range params {
			q.Add(k, v)
		}
	}

	request.URL.RawQuery = q.Encode()
	timeout := time.Duration(6 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	//禁止重定向，防止cookie丢失
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("Redirect")
	}

	resp, err := client.Do(request)
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	if err != nil && resp.StatusCode != 301 && resp.StatusCode != 302 {
		return false, HttpResp{Error: err}
	}

	resCookie := resp.Cookies()
	if cookie == nil {
		tempCk := make(map[string]string)
		cookie = &tempCk
	}
	for _, v := range resCookie {
		if _, ok := (*cookie)[v.Name]; ok {
			(*cookie)[v.Name] += ";" + v.Value
		} else {
			(*cookie)[v.Name] = v.Value
		}
	}

	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return false, HttpResp{Error: err2}
	}

	var httpResp HttpResp
	success := false
	//判断是否需要重定向
	localtion, _ := resp.Location()
	if localtion != nil && localtion.String() != "" {
		success, httpResp = GetRequest(localtion.String(), cookie, nil, nil)
	} else {
		success = true
		httpResp.Error = nil
		httpResp.Data = data
		httpResp.Localtion = resp.Request.URL.String()
		httpResp.Cookie = cookie
	}
	return success, httpResp
}

func PostJson(api string, cookie *map[string]string, header map[string]string, data interface{}) (bool, HttpResp) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()

	bytesData := []byte(`{}`)

	if data != nil {
		bytesData, _ = json.Marshal(data)
	}

	request, err := http.NewRequest("POST", api, bytes.NewReader(bytesData))
	timeout := time.Duration(6 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	//禁止重定向，防止cookie丢失
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("Redirect")
	}

	if cookie != nil && len(*cookie) != 0 {
		if header == nil {
			header = make(map[string]string)
		}
		header["Cookie"] = CookieMap2Str(*cookie)
	}

	if header != nil {
		for k, v := range header {
			request.Header.Set(k, v)
		}
	}
	request.Header.Set("content-type", "application/json")

	res, err := client.Do(request)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()

	if err != nil && res.StatusCode != 302 && res.StatusCode != 301 {
		return false, HttpResp{}
	}
	respBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, HttpResp{}
	}
	resCookie := res.Cookies()
	if cookie == nil {
		tempCk := make(map[string]string)
		cookie = &tempCk
	}
	for _, v := range resCookie {
		if _, ok := (*cookie)[v.Name]; ok {
			(*cookie)[v.Name] += ";" + v.Value
		} else {
			(*cookie)[v.Name] = v.Value
		}
	}

	var resp HttpResp

	location, err := res.Location()
	if location != nil && location.String() != "" {
		resp.Localtion = location.String()
	}
	resp.Cookie = cookie
	resp.Data = respBytes
	resp.Error = nil
	return true, resp
}

func PostForm(api string, cookie *map[string]string, header, params map[string]string) (bool, HttpResp) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()

	data := make(url.Values)
	if params != nil {
		for k, v := range params {
			data[k] = []string{v}
		}
	}
	request, err := http.NewRequest("POST", api, bytes.NewReader([]byte(data.Encode())))
	timeout := time.Duration(6 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	//禁止重定向，防止cookie丢失
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("Redirect")
	}
	if cookie != nil && len(*cookie) != 0 {
		if header == nil {
			header = make(map[string]string)
		}
		header["Cookie"] = CookieMap2Str(*cookie)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if header != nil {
		for k, v := range header {
			request.Header.Set(k, v)
		}
	}

	res, err := client.Do(request)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()

	if err != nil && res.StatusCode != 302 && res.StatusCode != 301 {
		return false, HttpResp{}
	}
	respBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("读取内容失败:", err.Error())
		//这里不返回是因为可能是一个可以继续重定向的链接
	}
	resCookie := res.Cookies()
	if cookie == nil {
		tempCk := make(map[string]string)
		cookie = &tempCk
	}
	for _, v := range resCookie {
		if _, ok := (*cookie)[v.Name]; ok {
			(*cookie)[v.Name] += ";" + v.Value
		} else {
			(*cookie)[v.Name] = v.Value
		}
	}

	var resp HttpResp

	location, err := res.Location()
	if location != nil && location.String() != "" {
		resp.Localtion = location.String()
	}
	resp.Cookie = cookie
	resp.Data = respBytes
	resp.Error = nil
	return true, resp
}

func CookieMap2Str(ck map[string]string) string {
	ckStr := ""
	if ck == nil || len(ck) == 0 {
		return ckStr
	}
	for k, v := range ck {
		ckStr += k + "=" + v + ";"
	}
	return strings.Trim(ckStr, ";")
}

func CookieStr2Map(ck string) map[string]string {
	ckMap := make(map[string]string)
	if ck == "" {
		return nil
	}

	split := strings.Split(ck, ";")
	for _, v := range split {
		split2 := strings.Split(v, "=")
		ckMap[split2[0]] = split2[1]
	}
	return ckMap
}
