package comms

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/TDAF/gologops"
)

// SMSGatewayType is ....
type SMSGatewayType int

const (
	// Pigeon is ...
	Pigeon SMSGatewayType = iota
	// Smppadapter is ...
	Smppadapter
)

type smsType struct {
	smsgateway SMSGatewayType
}

// curl -X POST -d '{"to": ["tel:+34699218702"],
// "message": "Tu PIN es 8765", "from": "tel:22949;phone-context=+34"}'
// --header "Content-Type:application/json" http://81.45.59.59:8000/sms/v2/smsoutbound

type smsRequestType struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Message string   `json:"message"`
}

type smsResponseType struct {
	ID string `json:"id"`
}

// SendSMS is ...
func (s smsType) SendSMS(text string) (string, error) {

	var url string

	if s.smsgateway == Smppadapter {
		url = "https://dev.mobileconnect.pdi.tid.es/es/sms/v2/smsoutbound"
	} else {
		url = "https://pigeon.compilon.com/sms/v2/smsoutbound"
	}

	dst := []string{"tel:+34699218702"}

	smsBody := smsRequestType{From: "tel:22154;phone-context=+34", To: dst, Message: text}

	client, req := getClient(url, smsBody, "application/json")

	if s.smsgateway == Pigeon {
		req.SetBasicAuth("f3381b48-133e-4cd9-8a1a-5811f4dbcd61", "31b50c52-5dc8-4084-86fa-d70486f807f8")
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	var smsResponse smsResponseType
	switch resp.StatusCode {
	case http.StatusOK:
		json.NewDecoder(resp.Body).Decode(&smsResponse)
		gologops.Infof("SMS ID: %s", smsResponse.ID)
		gologops.InfoC(gologops.C{"op": "SendSMS", "result": "OK"}, "%s", resp.Status)
	default:
		gologops.InfoC(gologops.C{"op": "SendSMS", "result": "NOK"}, "")
	}

	if err != nil {
		fmt.Println(err)
	}

	return smsResponse.ID, err

}

func getClient(url string, body interface{}, contentType string) (*http.Client, *http.Request) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	bodyJSON, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(fmt.Sprintf("%s", bodyJSON)))
	req.Header.Add("Content-Type", contentType)

	return client, req
}

// NewSMS creates a face client
func NewSMS(gatewayType SMSGatewayType) smsType {

	sms := smsType{smsgateway: gatewayType}
	return sms
}
