package oxford

import (
	"fmt"
	"net/http"
	"strconv"

	"encoding/json"

	"github.com/TDAF/gologops"
)

type face struct {
	apiKey string
}

// NewFace creates a face client
func NewFace(key string) face {

	f := face{}
	f.apiKey = key
	return f
}

type faceList struct {
	FaceListID string `json:"faceListId,omitempty"`
	Name       string `json:"name"`
	UserData   string `json:"userData"`
}

type faceType struct {
	PersistedFaceID string `json:"persistedFaceId"`
	UserData        string `json:"userData"`
}
type FaceListContent struct {
	FaceListID     string `json:"faceListId"`
	Name           string `json:"name"`
	UserData       string `json:"userData"`
	PersistedFaces []faceType
}

func (f face) Verify(img string) bool {
	url := GetResource(Face, V1, "detect")
	gologops.Infof("url: %s", url)
	return true
}

// PhotoURLType is..
type PhotoURLType struct {
	URL string `json:"url"`
}

// PhotoLocalType  is ...
type PhotoLocalType struct {
	URL string `json:"url"`
}

type faceResponseType struct {
	PersistedFaceID string `json:"persistedFaceId"`
}

type faceSimilarRequestType struct {
	FaceID                     string `json:"faceId"`
	FaceListID                 string `json:"faceListId"`
	MaxNumOfCandidatesReturned int    `json:"maxNumOfCandidatesReturned"`
}

type FaceSimilarResponseType struct {
	PersistedFaceID string  `json:"persistedFaceId"`
	Confidence      float64 `json:"confidence"`
}

type FaceDetectInfo struct {
	FaceID string `json:"faceId"`
	//El resto no me interesan
}

func (f face) Detect(photoURL string) (string, error) {
	url := GetResource(Face, V1, "detect")
	photo := PhotoURLType{URL: photoURL}

	resp, err := POST(url, M{"returnFaceId": "true"}, f.apiKey, nil, "application/json", photo)

	var faceID string

	if err != nil {
		return faceID, err
	}

	var faceDetectResponse []FaceDetectInfo

	switch resp.StatusCode {
	case http.StatusOK:
		json.NewDecoder(resp.Body).Decode(&faceDetectResponse)
		gologops.InfoC(gologops.C{"op": "Detect", "result": "OK"}, "%s", resp.Status)
		faceID = faceDetectResponse[0].FaceID
	default:
		var faceErrorResponse APIErrorResponse
		json.NewDecoder(resp.Body).Decode(&faceErrorResponse)
		gologops.InfoC(gologops.C{"op": "Detect", "result": "NOK"}, "%s", resp.Status)
		//gologops.Info("Status:%s|Request:%s", resp.Status, req.URL.RequestURI())
		fmt.Print(toJSON(faceErrorResponse, pretty))
	}

	return faceID, nil

}

func (f face) DetectBin(imageByteArray []byte) (persistedFaceID string, err error) {
	url := GetResource(Face, V1, "detect")
	if err != nil {
		return "", err
	}

	resp, err := POST(url, M{"returnFaceId": "true"}, f.apiKey, M{"Content-Length": strconv.Itoa(len(imageByteArray))}, "application/octet-stream", imageByteArray)

	var faceID string

	if err != nil {
		return faceID, err
	}

	var faceDetectResponse []FaceDetectInfo

	switch resp.StatusCode {
	case http.StatusOK:
		json.NewDecoder(resp.Body).Decode(&faceDetectResponse)
		gologops.InfoC(gologops.C{"op": "DetectBin", "result": "OK"}, "%s", resp.Status)
		faceID = faceDetectResponse[0].FaceID
	default:
		var faceErrorResponse APIErrorResponse
		json.NewDecoder(resp.Body).Decode(&faceErrorResponse)
		gologops.InfoC(gologops.C{"op": "DetectBin", "result": "NOK"}, "%s", resp.Status)
		//gologops.Info("Status:%s|Request:%s", resp.Status, req.URL.RequestURI())
		fmt.Print(toJSON(faceErrorResponse, pretty))
	}

	return faceID, nil
}

func (f face) DetectBinFromFile(imageFileName string) (persistedFaceID string, err error) {
	url := GetResource(Face, V1, "detect")
	imageByteArray, err := FileToByteArray(imageFileName)
	if err != nil {
		return "", err
	}

	resp, err := POST(url, M{"returnFaceId": "true"}, f.apiKey, M{"Content-Length": strconv.Itoa(len(imageByteArray))}, "application/octet-stream", imageByteArray)

	var faceID string

	if err != nil {
		return faceID, err
	}

	var faceDetectResponse []FaceDetectInfo

	switch resp.StatusCode {
	case http.StatusOK:
		json.NewDecoder(resp.Body).Decode(&faceDetectResponse)
		gologops.InfoC(gologops.C{"op": "DetectBinFromFile", "result": "OK"}, "%s", resp.Status)
		faceID = faceDetectResponse[0].FaceID
	default:
		var faceErrorResponse APIErrorResponse
		json.NewDecoder(resp.Body).Decode(&faceErrorResponse)
		gologops.InfoC(gologops.C{"op": "DetectBinFromFile", "result": "NOK"}, "%s", resp.Status)
		//gologops.Info("Status:%s|Request:%s", resp.Status, req.URL.RequestURI())
		fmt.Print(toJSON(faceErrorResponse, pretty))
	}

	fmt.Printf("DetectBinFromFile return......->: %s\n.......", faceID)
	return faceID, nil
}

func (f face) FindSimilar(faceID string, faceListID string) ([]FaceSimilarResponseType, error) {
	url := GetResource(Face, V1, "findsimilars")
	faceSimilarBody := faceSimilarRequestType{FaceID: faceID, FaceListID: faceListID, MaxNumOfCandidatesReturned: 5}

	var similarList []FaceSimilarResponseType
	resp, err := POST(url, nil, f.apiKey, nil, "application/json", faceSimilarBody)

	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		json.NewDecoder(resp.Body).Decode(&similarList)
		fmt.Print(toJSON(similarList, pretty))
		gologops.InfoC(gologops.C{"op": "FindSimilar", "result": "OK"}, "%s", resp.Status)
	default:
		var similarErrorResponse APIErrorResponse
		json.NewDecoder(resp.Body).Decode(&similarErrorResponse)
		gologops.InfoC(gologops.C{"op": "FindSimilar", "result": "NOK"}, "%s", resp.Status)
		fmt.Print(toJSON(similarErrorResponse, pretty))
	}

	if err != nil {
		fmt.Println(err)
	}
	return similarList, err
}

func (f face) AddFace(faceListID string, imageFileName string, email string) (persistedFaceID string, err error) {
	url := GetResource(Face, V1, "facelists")
	url = url + "/" + faceListID + "/persistedFaces"
	imageByteArray, err := FileToByteArray(imageFileName)

	if err != nil {
		return "", err
	}

	fmt.Printf("Add Face -------->\n")
	fmt.Printf("User data is: %s\n", email)
	resp, err := POST(url, M{"userData": email}, f.apiKey, M{"Content-Length": strconv.Itoa(len(imageByteArray))}, "application/octet-stream", imageByteArray)
 
	if err != nil {
		return "", err
	}

	var faceResponse faceResponseType
	switch resp.StatusCode {
	case http.StatusOK:
		json.NewDecoder(resp.Body).Decode(&faceResponse)
		gologops.InfoC(gologops.C{"op": "AddFace", "result": "OK"}, "%s", resp.Status)
	default:
		var faceErrorResponse APIErrorResponse
		json.NewDecoder(resp.Body).Decode(&faceErrorResponse)
		gologops.InfoC(gologops.C{"op": "AddFace", "result": "NOK"}, "%s", resp.Status)
		//gologops.Info("Status:%s|Request:%s", resp.Status, req.URL.RequestURI())
		fmt.Print(toJSON(faceErrorResponse, pretty))
	}

	if err != nil {
		fmt.Println(err)
	}

	return toJSON(faceResponse, pretty), err

}

func (f face) AddFaceURL(faceListID string, photoURL string) (list string, err error) {
	url := GetResource(Face, V1, "facelists")
	url = url + "/" + faceListID + "/persistedFaces"
	photo := PhotoURLType{URL: photoURL}

	resp, err := POST(url, nil, f.apiKey, nil, "application/json", photo)

	fmt.Printf("AddPhoto----->")

	if err != nil {
		return "", err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		gologops.InfoC(gologops.C{"op": "AddPhoto", "result": "OK"}, "%s", resp.Status)
	default:
		gologops.InfoC(gologops.C{"op": "AddPhoto", "result": "NOK"}, "%s", resp.Status)
		//gologops.Info("Status:%s|Request:%s", resp.Status, req.URL.RequestURI())
	}

	if err != nil {
		fmt.Println(err)
	}

	var faceResponse faceResponseType
	switch resp.StatusCode {
	case http.StatusOK:
		json.NewDecoder(resp.Body).Decode(&faceResponse)
		gologops.InfoC(gologops.C{"op": "AddFace", "result": "OK"}, "%s", resp.Status)
	default:
		var faceErrorResponse APIErrorResponse
		json.NewDecoder(resp.Body).Decode(&faceErrorResponse)
		gologops.InfoC(gologops.C{"op": "AddFace", "result": "NOK"}, "%s", resp.Status)
		//gologops.Info("Status:%s|Request:%s", resp.Status, req.URL.RequestURI())
		fmt.Print(toJSON(faceErrorResponse, pretty))
	}

	if err != nil {
		fmt.Println(err)
	}

	return toJSON(faceResponse, pretty), err

}

func (f face) CreateFaceList(faceListID string) (id string, err error) {
	url := GetResource(Face, V1, "facelists")
	url = url + "/" + faceListID
	fl := faceList{Name: faceListID, UserData: "Face List for Equinox"}

	resp, err := PUT(url, nil, f.apiKey, nil, "application/json", fl)

	if err != nil {
		return "", err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		gologops.InfoC(gologops.C{"op": "createFaceList", "result": "OK"}, "%s", resp.Status)
	default:
		gologops.InfoC(gologops.C{"op": "createFaceList", "result": "NOK"}, "%s", resp.Status)
		//gologops.Info("Status:%s|Request:%s", resp.Status, req.URL.RequestURI())
	}

	if err != nil {
		fmt.Println(err)
	}

	return faceListID, err
}

func (f face) GetFaceList() (list string, err error) {
	url := GetResource(Face, V1, "facelists")

	resp, err := GET(url, f.apiKey, nil, nil)

	if err != nil {
		return "", err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		gologops.InfoC(gologops.C{"op": "GetFaceList", "result": "OK"}, "%s", resp.Status)
	default:
		gologops.InfoC(gologops.C{"op": "GetFaceList", "result": "NOK"}, "%s", resp.Status)
		//gologops.Info("Status:%s|Request:%s", resp.Status, req.URL.RequestURI())
	}

	if err != nil {
		fmt.Println(err)
	}

	var listOfFacesList []faceList
	json.NewDecoder(resp.Body).Decode(&listOfFacesList)
	return toJSON(listOfFacesList, pretty), err

}

func (f face) GetFacesInAList(faceListID string) (list string, err error) {
	url := GetResource(Face, V1, "facelists")
	url = url + "/" + faceListID

	resp, err := GET(url, f.apiKey, nil, nil)

	if err != nil {
		return "", err
	}

	//gologops.Info("Status:%s|Request:%s", resp.Status, req.URL.RequestURI())
	switch resp.StatusCode {
	case http.StatusOK:
		gologops.InfoC(gologops.C{"op": "GetFaceList", "result": "OK"}, "%s", resp.Status)
	default:
		gologops.InfoC(gologops.C{"op": "GetFaceList", "result": "NOK"}, "%s", resp.Status)
	}

	if err != nil {
		fmt.Println(err)
	}

	facesInAList := FaceListContent{}
	json.NewDecoder(resp.Body).Decode(&facesInAList)
	return toJSON(facesInAList, pretty), err

}

func (f face) GetObjectFacesInAList(faceListID string) (list FaceListContent, err error) {
	url := GetResource(Face, V1, "facelists")
	url = url + "/" + faceListID

	resp, err := GET(url, f.apiKey, nil, nil)

	if err != nil {
		return FaceListContent{}, err
	}

	//gologops.Info("Status:%s|Request:%s", resp.Status, req.URL.RequestURI())
	switch resp.StatusCode {
	case http.StatusOK:
		gologops.InfoC(gologops.C{"op": "GetFaceList", "result": "OK"}, "%s", resp.Status)
	default:
		gologops.InfoC(gologops.C{"op": "GetFaceList", "result": "NOK"}, "%s", resp.Status)
	}

	if err != nil {
		fmt.Println(err)
	}

	facesInAList := FaceListContent{}
	json.NewDecoder(resp.Body).Decode(&facesInAList)
	return facesInAList, nil

}
