package oxford

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	apiURL string = "https://api.projectoxford.ai"
	// V1 is the v1.0 version
	V1 string = "v1.0"
	// AzureSubscriptionID is my azure subscription
	AzureSubscriptionID string = "70306775-8047-4d29-9540-679cc5412f0f"
)

type M map[string]string

// APIType is a type for the different apis
type APIType int

const (
	// Face represents the face api
	Face APIType = iota
	// SpeakerRecognition represents the SpeakerRecognition api
	SpeakerRecognition
)

var apis = map[APIType]string{
	Face:               "face",
	SpeakerRecognition: "spid",
}

type oxfordError struct {
	StatusCode string `json:"code"`
	Message    string `json:"message"`
}

// Error interface
func (err oxfordError) Error() string {
	return err.Message
}

// APIErrorResponse is ...
type APIErrorResponse struct {
	Err oxfordError `json:"error"`
}

// GetResource builds a resource
func GetResource(apiType APIType, version string, resource string) string {
	u, _ := url.ParseRequestURI(apiURL)
	u.Path = apis[apiType] + "/" + version + "/" + resource
	urlStr := fmt.Sprintf("%v", u)
	return urlStr
}

func parseError(body io.Reader) APIErrorResponse {
	err := APIErrorResponse{}
	json.NewDecoder(body).Decode(&err)
	return err
}

type printOption int

const (
	pretty printOption = iota
	normal
)

func toJSON(value interface{}, option printOption) string {

	var jsonValue []byte

	switch option {
	case pretty:
		jsonValue, _ = json.MarshalIndent(value, "", "\t")
	case normal:
		jsonValue, _ = json.Marshal(value)
	}

	return fmt.Sprintf("%s", jsonValue)
}

type HTTPMethod int

const (
	// GET represents the HTTP GET method
	HTTP_GET HTTPMethod = iota
	HTTP_PUT
	HTTP_POST
)

var (
	requestGET  = createHTTPRequest(HTTP_GET)
	requestPOST = createHTTPRequest(HTTP_POST)
	requestPUT  = createHTTPRequest(HTTP_PUT)
)

func GET(url string, apiKey string, queryParams map[string]string, headers map[string]string) (*http.Response, error) {

	r := requestGET(url, queryParams, apiKey, headers, "", nil)
	client := getClient()
	return client.Do(r)
}

func PUT(url string, queryParams map[string]string, apiKey string, headers map[string]string, contentType string, body interface{}) (*http.Response, error) {

	r := requestPUT(url, queryParams, apiKey, headers, contentType, body)
	client := getClient()
	return client.Do(r)
}

func POST(url string, queryParams map[string]string, apiKey string, headers map[string]string, contentType string, body interface{}) (*http.Response, error) {

	r := requestPOST(url, queryParams, apiKey, headers, contentType, body)
	client := getClient()
	fmt.Printf("--> %s\n\n", formatRequest(r))
	return client.Do(r)
}

func createBody(body interface{}, contentType string) (bodyReader io.Reader) {

	switch {
	case contentType == "application/octet-stream":
		bodyReader = bytes.NewBuffer(body.([]byte))
	case contentType == "application/json":
		bodyReader = bytes.NewBufferString(toJSON(body, normal))
	}

	return bodyReader
}

type createRequest func(url string, queryParams map[string]string, apiKey string, headers map[string]string, contentType string, body interface{}) *http.Request

func createHTTPRequest(method HTTPMethod) createRequest {
	return func(url string, queryParams map[string]string, apiKey string, headers map[string]string, contentType string, body interface{}) *http.Request {
		var req *http.Request
		switch {
		case method == HTTP_GET:
			req, _ = http.NewRequest("GET", url, nil)
		case method == HTTP_POST:
			req, _ = http.NewRequest("POST", url, createBody(body, contentType))
		case method == HTTP_PUT:
			req, _ = http.NewRequest("PUT", url, createBody(body, contentType))

		}
		req.Header.Add("Content-Type", contentType)
		req.Header.Add("Ocp-Apim-Subscription-Key", apiKey)

		for k, v := range headers {
			req.Header.Add(k, v)
		}

		q := req.URL.Query()
		for k, v := range queryParams {
			q.Add(k,v)
			req.URL.RawQuery = q.Encode()
			//req.URL.Query().Add(k, v)
		}

		return req
	}
}

func getClient() *http.Client {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	return client
}

func FileToByteArray(imageFileName string) ([]byte, error) {
	file, err := os.Open(imageFileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()
	size := fileInfo.Size()
	bytes := make([]byte, size)

	// read file into bytes
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)

	// fileOutput, _ := os.Create("/tmp/image.jpg")
	// defer fileOutput.Close()
	// imageOutput := bufio.NewWriter(fileOutput)
	// imageOutput.Write(bytes)
	return bytes, err
}

func ByteArrayToBase64(binaryByteArray []byte) string {
	imgBase64Str := base64.StdEncoding.EncodeToString(binaryByteArray)
	return imgBase64Str
}

func Base64ToByteArray(imgBase64Str string) ([]byte, error) {
	fmt.Printf("-------------------\n")
	fmt.Print(imgBase64Str)
	fmt.Printf("\n\n-------------------")
	binaryByteArray, err := base64.StdEncoding.DecodeString(imgBase64Str)
	return binaryByteArray, err
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
 // Create return string
 var request []string
 // Add the request string
 url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
 request = append(request, url)
 // Add the host
 request = append(request, fmt.Sprintf("Host: %v", r.Host))
 // Loop through headers
 for name, headers := range r.Header {
   name = strings.ToLower(name)
   for _, h := range headers {
     request = append(request, fmt.Sprintf("%v: %v", name, h))
   }
 }
 
 // If this is a POST, add post data
 if r.Method == "POST" {
    r.ParseForm()
    request = append(request, "\n")
    request = append(request, r.Form.Encode())
 } 
  // Return the request as a string
  return strings.Join(request, "\n")
}