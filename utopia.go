package utopiago

import (
	//"encoding/json"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"gopkg.in/grignaak/tribool.v1"
)

//Query is a filter for API requests
type Query struct {
	Method string            `json:"method"`
	Token  string            `json:"token"`
	Params map[string]string `json:"params"`
}

//UtopiaClient lets you connect to Utopia Client
type UtopiaClient struct {
	protocol, host, token string
	port                  int
}

//UtopiaClientInterface contains an enumeration of methods
type UtopiaClientInterface interface {
	apiQuery(methodName string) map[string]interface{}

	GetProfileStatus() map[string]interface{}
	GetSystemInfo() map[string]interface{}
	GetOwnContact() map[string]interface{}

	GetBalance() (float64, error)
	UseVoucher(voucherCode string) error
	GetFinanceHistory() map[string]interface{}
	CheckClientConnection() bool
	CreateVoucher(amount float64) error
}

func (c *UtopiaClient) apiQuery(methodName string, params map[string]string) map[string]interface{} {
	url := c.protocol + "://" + c.host + ":" + strconv.Itoa(c.port) + "/api/1.0/"
	var query = Query{
		Method: methodName,
		Token:  c.token,
	}
	if params != nil {
		query.Params = params
	}

	var jsonStr, err = json.Marshal(query)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var responseMap map[string]interface{}

	if !json.Valid(body) {
		//TODO: return error
		return responseMap
	}

	json.Unmarshal(body, &responseMap)
	return responseMap
}

//GetProfileStatus gets data about the status of the current account
func (c *UtopiaClient) GetProfileStatus() map[string]interface{} {
	return c.apiQuery("getProfileStatus", nil)
}

//GetSystemInfo retrieves client system information
func (c *UtopiaClient) GetSystemInfo() map[string]interface{} {
	return c.apiQuery("getSystemInfo", nil)
}

func (c *UtopiaClient) queryResultToString(methodName string, params map[string]string) string {
	var response map[string]interface{} = c.apiQuery(methodName, params)
	var resultstr string = fmt.Sprintf("%v", response["result"])
	return resultstr
}

func (c *UtopiaClient) queryResultToBool(methodName string, params map[string]string) bool {
	var resultstr string = c.queryResultToString(methodName, params)
	var result bool = tribool.FromString(resultstr).WithMaybeAsTrue()
	return result
}

//SetProfileStatus updates data about the status of the current account
func (c *UtopiaClient) SetProfileStatus(status string, mood string) bool {
	queryMap := make(map[string]string)
	queryMap["status"] = status
	queryMap["mood"] = mood

	return c.queryResultToBool("setProfileStatus", queryMap)
}

//GetOwnContact asks for full details of the current account
func (c *UtopiaClient) GetOwnContact() map[string]interface{} {
	return c.apiQuery("getOwnContact", nil)
}

//CheckClientConnection - checks if there are any errors when contacting the client
func (c *UtopiaClient) CheckClientConnection() bool {
	//TODO
	return false
}

//UseVoucher - uses the voucher and returns an error on failure
func (c *UtopiaClient) UseVoucher() error {
	//TODO
	return nil
}

//GetFinanceHistory request the necessary financial statistics
func (c *UtopiaClient) GetFinanceHistory() map[string]interface{} {
	//TODO
	//TODO: handle error
	return make(map[string]interface{})
}

//GetBalance request account Crypton balance
func (c *UtopiaClient) GetBalance() (float64, error) {
	//TODO
	return 0, nil
}

//CreateVoucher requests the creation of a new voucher
func (c *UtopiaClient) CreateVoucher(amount float64) error {
	//TODO
	return nil
}
