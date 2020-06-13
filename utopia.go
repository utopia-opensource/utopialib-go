package main

import (
	//"encoding/json"
    "fmt"
    "net/http"
    "bytes"
    "io/ioutil"
    "strconv"
    "encoding/json"
)

type Query struct {
	Method string `json:"method"`
	Token string  `json:"token"`
}

type UtopiaClient struct {
	protocol, host, token string
	port int
}

type UtopiaClientInterface interface {
	apiQuery(methodName string) map[string]interface{}
	GetProfileStatus() map[string]interface{}
	GetSystemInfo() map[string]interface{}
}

func (c UtopiaClient) apiQuery(methodName string) map[string]interface{} {
	url := c.protocol + "://" + c.host + ":" + strconv.Itoa(c.port) + "/api/1.0/"
	//fmt.Println(url) //debug

	var query = Query{
		Method: methodName,
		Token: c.token,
	}
	//var jsonStr = []byte(`{"token":"test"}`)
	var jsonStr, err = json.Marshal(query)
	if err != nil {
        panic(err)
    }
    //fmt.Println(jsonStr)
    //var jsonStr = []byte(json)

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    //req.Header.Set("X-Custom-Header", "myvalue")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    //fmt.Println("response Status:", resp.Status)
    //fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    var responseMap map[string]interface{}
    //fmt.Println("response Body:", string(body))
    //return string(body)

    //TODO: check json

    json.Unmarshal([]byte(body), &responseMap)
    return responseMap
}

func (c UtopiaClient) GetProfileStatus() map[string]interface{} {
	return c.apiQuery("getProfileStatus")
}

func (c UtopiaClient) GetSystemInfo() map[string]interface{} {
	return c.apiQuery("getSystemInfo")
}

func main() {
	client := UtopiaClient{
		protocol: "http",
		token:    "C17BF2E95821A6B545DC9A193CBB750B",
		host:     "127.0.0.1",
		port:     22791,
	}

	fmt.Println(client.apiQuery("getSystemInfo"))
}
