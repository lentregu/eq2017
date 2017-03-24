package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/lentregu/eq2017/oxford"
)

type printOption int

const (
	pretty printOption = iota
	normal
)

type wordOption int

const (
	oneWord = iota
	multipleWords
)

//const oneWordRegExp ="^\S.*\S$"
const oneWordRegExp = `^[^\t\n\f\r ]*$`

const multipleWordsRegExp = `^.*$`

//var faceService = oxford.NewFace("567c560aa85245418459b82634bc7a98")
var	faceService = oxford.NewFace("83dc246bac2b447782b5aab70604bc97")

var speakService = oxford.NewSpeak("af90809f8a0d4430ba2aabd44785ebc4")

func getBestMatch(similarList []oxford.FaceSimilarResponseType) *oxford.FaceSimilarResponseType {
	var bestMatch *oxford.FaceSimilarResponseType = nil
	for _, similar := range similarList {
		var similar = similar
		if bestMatch == nil || similar.Confidence > bestMatch.Confidence {
			bestMatch = &similar
		}
	}

	if bestMatch.Confidence <= 0.6 {
		bestMatch = nil
	}

	return bestMatch
}

func whois() (string, error) {

	faceListID, err := readString("FaceList ID", oneWordRegExp)
	imageFileName, err := readString("Face", oneWordRegExp)
	if err != nil {
		return "", err
	}

	faceID, _ := faceService.DetectBinFromFile(imageFileName)
	similarList, err := faceService.FindSimilar(faceID, faceListID)
	if err != nil {
		return "", err
	}

	bestMatch := getBestMatch(similarList)

	var email string = ""
	if bestMatch != nil {
		faces, _:= faceService.GetObjectFacesInAList(faceListID)
		for _, face := range faces.PersistedFaces {
			if face.PersistedFaceID == bestMatch.PersistedFaceID {
				email = face.UserData
			}
		}

		fmt.Printf("El FaceID detectado es: %s", bestMatch.PersistedFaceID)
		fmt.Print(toJSON(faces, pretty))

	} else {
		err = fmt.Errorf("User Not Found")
	}

	return email, err
}

func addFace() (string, error) {

	faceListID, err := readString("FaceList ID", oneWordRegExp)
	imageFileName, err := readString("Face", oneWordRegExp)
	email, err := readString("Email", oneWordRegExp)
	if err != nil {
		return "", err
	}
	//return faceService.AddFaceURL(faceListID, imageFileName)
	return faceService.AddFace(faceListID, imageFileName, email)
}

func getFacesInAList() (string, error) {

	faceListID, _ := readString("FaceList ID", oneWordRegExp)
	return faceService.GetFacesInAList(faceListID)
}

func getFaceList() (string, error) {

	return faceService.GetFaceList()
}

func createFaceList() (string, error) {

	faceListID, err := readString("FaceList Name", oneWordRegExp)

	if err != nil {
		return "", err
	}

	return faceService.CreateFaceList(faceListID)
}

func createSpeakerProfile() (string, error) {

	locale, err := readString("Locale", oneWordRegExp)

	if err != nil {
		return "", err
	}

	return speakService.CreateProfile(locale)
}

//------------------------

func whoIsBase64() {

	faceListID, err := readString("FaceList ID", oneWordRegExp)
	imageFileName, err := readString("Face", oneWordRegExp)

	url := "http://localhost:8080/whoisb64"
	imgContent, _ := fileToString(imageFileName)

	//var body inteface{}
	var body = map[string]string{"img": imgContent, "faceListID": faceListID}
	bodyJson, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	bodyObject := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&bodyObject)
	fmt.Println("WhoIsBase64 Result--------->")
	fmt.Print(toJSON(bodyObject, pretty))
}

func fileToString(imageFileName string) (string, error) {
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

	return string(bytes), err
}

//------------------------

func readString(name string, wordRegExp string) (string, error) {
	var value string
	fmt.Print(name + ": ")

	validExpression := regexp.MustCompile(wordRegExp)

	line, _, err := stdin.ReadLine()
	if err != nil {
		err = fmt.Errorf("Error reading value for %s: %s", name, err.Error())
	} else {

		value = fmt.Sprintf("%s", line)

		if !validExpression.MatchString(value) && wordRegExp == oneWordRegExp {
			err = fmt.Errorf("ERROR Not spaces are allowed for %s field\n", name)
		}
	}

	return value, err
}

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