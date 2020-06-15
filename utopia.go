package utopiago

import (
	//"encoding/json"
    "fmt"
    "net/http"
    "bytes"
    "io/ioutil"
    "strconv"
    "encoding/json"
    "gopkg.in/grignaak/tribool.v1"
)

type Query struct {
	Method string `json:"method"`
	Token string  `json:"token"`
	Params map[string]string `json:"params"`
}

type UtopiaClient struct {
	protocol, host, token string
	port int
}

type UtopiaClientInterface interface {
	apiQuery(methodName string) map[string]interface{}
	GetProfileStatus() map[string]interface{}
	GetSystemInfo() map[string]interface{}
	GetOwnContact() map[string]interface{}
}

func (c UtopiaClient) apiQuery(methodName string, params map[string]string) map[string]interface{} {
	url := c.protocol + "://" + c.host + ":" + strconv.Itoa(c.port) + "/api/1.0/"
	//fmt.Println(url) //debug

	var query = Query{
		Method: methodName,
		Token: c.token,
	}
	if params != nil {
		query.Params = params
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
	return c.apiQuery("getProfileStatus", nil)
}

func (c UtopiaClient) GetSystemInfo() map[string]interface{} {
	return c.apiQuery("getSystemInfo", nil)
}

func (c UtopiaClient) queryResultToString(methodName string, params map[string]string) string {
	var response map[string]interface{} = c.apiQuery(methodName, params)
	var resultstr string = fmt.Sprintf("%v", response["result"])
	return resultstr
}

func (c UtopiaClient) queryResultToBool(methodName string, params map[string]string) bool {
	var resultstr string = c.queryResultToString(methodName, params)
	var result bool = tribool.FromString(resultstr).WithMaybeAsTrue()
	return result
}

func (c UtopiaClient) SetProfileStatus(status string, mood string) bool {
	queryMap := make(map[string]string)
	queryMap["status"] = status
	queryMap["mood"] = mood

	return c.queryResultToBool("setProfileStatus", queryMap)
}

func (c UtopiaClient) GetOwnContact() map[string]interface{} {
	return response map[string]interface{} = c.apiQuery("getOwnContact", nil)
}

//func main() {
//	client := UtopiaClient{
//		protocol: "http",
//		token:    "C17BF2E95821A6B545DC9A193CBB750B",
//		host:     "127.0.0.1",
//		port:     22791,
//	}
//
//	fmt.Println(client.GetSystemInfo())
//	fmt.Println(client.SetProfileStatus("Available", "test mood"))
//}
