package utopiago

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"gopkg.in/grignaak/tribool.v1"
)

//Query is a filter for API requests
type Query struct {
	Method string                 `json:"method"`
	Token  string                 `json:"token"`
	Params map[string]interface{} `json:"params"`
}

//UtopiaClient lets you connect to Utopia Client
type UtopiaClient struct {
	Protocol, Host, Token string
	Port                  int
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

func (c *UtopiaClient) apiQuery(methodName string, params map[string]interface{}) (map[string]interface{}, error) {
	var responseMap map[string]interface{}
	url := c.Protocol + "://" + c.Host + ":" + strconv.Itoa(c.Port) + "/api/1.0/"
	var query = Query{
		Method: methodName,
		Token:  c.Token,
	}
	if params != nil {
		query.Params = params
	}

	var jsonStr, err = json.Marshal(query)
	if err != nil {
		return responseMap, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return responseMap, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if !json.Valid(body) {
		//TODO: return error
		return responseMap, errors.New("failed to validate json from client")
	}

	json.Unmarshal(body, &responseMap)
	return responseMap, nil
}

//GetProfileStatus gets data about the status of the current account
func (c *UtopiaClient) GetProfileStatus() (map[string]interface{}, error) {
	return c.apiQuery("getProfileStatus", nil)
}

//GetSystemInfo retrieves client system information
func (c *UtopiaClient) GetSystemInfo() (map[string]interface{}, error) {
	return c.apiQuery("getSystemInfo", nil)
}

func (c *UtopiaClient) queryResultToInterface(methodName string, params map[string]interface{}) (interface{}, error) {
	if !c.CheckClientConnection() {
		return nil, errors.New("client disconected")
	}
	response, err := c.apiQuery(methodName, params)
	if result, ok := response["result"]; ok {
		return result, err
	}
	return nil, errors.New("result field doesn't exists in client response")
}

func (c *UtopiaClient) queryResultToInterfaceArray(methodName string, params map[string]interface{}) ([]interface{}, error) {
	if !c.CheckClientConnection() {
		return nil, errors.New("client disconected")
	}
	response, err := c.apiQuery(methodName, params)
	if result, ok := response["result"]; ok {
		//TODO: check type assertion
		return result.([]interface{}), err
	}
	return nil, errors.New("accaptable result doesn't exists in client response")
}

func (c *UtopiaClient) queryResultToString(methodName string, params map[string]interface{}) (string, error) {
	if !c.CheckClientConnection() {
		return "", errors.New("client disconected")
	}
	response, err := c.apiQuery(methodName, params)
	if result, ok := response["result"]; ok {
		resultstr := fmt.Sprintf("%v", result)
		return resultstr, err
	}
	return "", errors.New("result field doesn't exists in client response")
}

func (c *UtopiaClient) queryResultToBool(methodName string, params map[string]interface{}) (bool, error) {
	resultstr, err := c.queryResultToString(methodName, params)
	resultBool := tribool.FromString(resultstr).WithMaybeAsTrue()
	return resultBool, err
}

func (c *UtopiaClient) queryResultToFloat64(methodName string, params map[string]interface{}) (float64, error) {
	resultstr, err := c.queryResultToString(methodName, params)
	if err != nil {
		return 0, err
	}
	resultFloat, err := strconv.ParseFloat(resultstr, 64)
	return resultFloat, err
}

func (c *UtopiaClient) queryResultToInt(methodName string, params map[string]interface{}) (int64, error) {
	resultstr, err := c.queryResultToString(methodName, params)
	if err != nil {
		return 0, err
	}
	result, err := strconv.ParseInt(resultstr, 10, 64)
	return result, err
}

//SetProfileStatus updates data about the status of the current account
func (c *UtopiaClient) SetProfileStatus(status string, mood string) error {
	queryMap := make(map[string]interface{})
	queryMap["status"] = status
	queryMap["mood"] = mood

	result, err := c.queryResultToBool("setProfileStatus", queryMap)
	if err != nil {
		return err
	}
	if result == false {
		return errors.New("failed to set profile status")
	}
	return nil
}

//GetOwnContact asks for full details of the current account
func (c *UtopiaClient) GetOwnContact() (map[string]interface{}, error) {
	return c.apiQuery("getOwnContact", nil)
}

//CheckClientConnection - checks if there are any errors when contacting the client
func (c *UtopiaClient) CheckClientConnection() bool {
	_, err := c.GetSystemInfo()
	if err != nil {
		return false
	}
	return true
}

//UseVoucher - uses the voucher and returns an error on failure
func (c *UtopiaClient) UseVoucher(voucherID string) (string, error) {
	params := map[string]interface{}{
		"voucherid": voucherID,
	}
	return c.queryResultToString("useVoucher", params)
}

//GetFinanceHistory request the necessary financial statistics
func (c *UtopiaClient) GetFinanceHistory(filters string, referenceNumber string) ([]interface{}, error) {
	params := map[string]interface{}{
		"filters":         filters,
		"referenceNumber": referenceNumber,
	}
	return c.queryResultToInterfaceArray("getFinanceSystemInformation", params)
}

//GetBalance request account Crypton balance
func (c *UtopiaClient) GetBalance() (float64, error) {
	result, err := c.queryResultToFloat64("getBalance", nil)
	if err != nil {
		return 0, err
	}
	return result, nil
}

//CreateVoucher requests the creation of a new voucher. it returns referenceNumber
func (c *UtopiaClient) CreateVoucher(amount float64) (string, error) {
	params := map[string]interface{}{
		"amount": amount,
	}
	result, err := c.queryResultToString("createVoucher", params)
	if err != nil {
		return "", err
	}
	if result == "" {
		return "", errors.New("failed to create voucher, empty string in client response")
	}
	return result, nil
}

//SetWebSocketState - set WSS Notification state
func (c *UtopiaClient) SetWebSocketState(enabled bool, port int) error {
	var enabledStr string = "0"
	if enabled {
		enabledStr = "1"
	}
	params := map[string]interface{}{
		"enabled": enabledStr,
		"port":    port,
	}
	result, err := c.queryResultToString("setWebSocketState", params)
	if err != nil {
		return err
	}
	if result == "" {
		return errors.New("failed to set websocker state")
	}
	return nil
}

//GetWebSocketState - returns WSS Notifications state, 0 - disabled or active listening port number.
func (c *UtopiaClient) GetWebSocketState() (int64, error) {
	result, err := c.queryResultToInt("getWebSocketState", nil)
	if err != nil {
		return 0, err
	}
	return result, nil
}

//ServeWs ..
func (c *UtopiaClient) ServeWs() error {
	//TODO
	return nil
}
