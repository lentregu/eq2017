package face

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lentregu/Equinox/HumanIdentification/comms"
	"github.com/lentregu/Equinox/oxford"
)

const (
	// PrimaryKey for FaceAPI
	PrimaryKey = "567c560aa85245418459b82634bc7a98"
	// SecondaryKey for FaceAPI
	SecondaryKey = "4c1a4e7a02104577b045a2d046b20d29"
)

type findSimilarRequestType struct {
	URL        string `json:"url"`
	FaceListID string `json:"faceListID"`
}

type findWhoisRequestType struct {
	Base64Image string `json:"img"`
	FaceListID  string `json:"faceListID"`
}

// Index is the welcome handler
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

// FindSimilar is a handler to detect faces
func FindSimilar(w http.ResponseWriter, r *http.Request) {
	requestBody := findSimilarRequestType{}

	json.NewDecoder(r.Body).Decode(&requestBody)

	faceService := oxford.NewFace("567c560aa85245418459b82634bc7a98")
	faceID, _ := faceService.Detect(requestBody.URL)

	fmt.Printf("El faceID es: %s\n", faceID)

	if similarList, err := faceService.FindSimilar(faceID, requestBody.FaceListID); err != nil {
		fmt.Printf("Error %v", err)
	} else if isSimilar(similarList) {
		sms := comms.NewSMS(comms.Smppadapter)
		sms.SendSMS("PERSONA AUTORIZADA")
		fmt.Println("PERSONA AUTORIZADA")
	} else {
		fmt.Println("PERSONA NO AUTORIZADA")
	}

}

func Whois(w http.ResponseWriter, r *http.Request) {
	requestBody := findWhoisRequestType{}

	json.NewDecoder(r.Body).Decode(&requestBody)

	fmt.Println("WHOIS------")
	fmt.Print(requestBody)

	faceService := oxford.NewFace("567c560aa85245418459b82634bc7a98")

	var bestMatch *oxford.FaceSimilarResponseType = nil
	var similarList []oxford.FaceSimilarResponseType = nil
	binaryImg, err := oxford.Base64ToByteArray(requestBody.Base64Image)

	if err == nil {
		faceID, _ := faceService.DetectBin(binaryImg)

		fmt.Printf("El faceID es: %s\n", faceID)

		similarList, err = faceService.FindSimilar(faceID, requestBody.FaceListID)

		bestMatch = getBestMatch(similarList)
	}

	var jsonValue []byte
	if bestMatch != nil {
		jsonValue, err = json.MarshalIndent(bestMatch, "", "\t")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		fmt.Fprintf(w, fmt.Sprintf("%s", jsonValue))
	} else {
		err = fmt.Errorf("User Not Found")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(404)
		fmt.Fprintf(w, err.Error())
	}
}

func isSimilar(similarList []oxford.FaceSimilarResponseType) bool {
	found := false
	for _, similar := range similarList {
		if similar.Confidence > 0.6 {
			found = true
		}
	}
	return found
}

func getBestMatch(similarList []oxford.FaceSimilarResponseType) *oxford.FaceSimilarResponseType {
	var bestMatch *oxford.FaceSimilarResponseType = nil
	for _, similar := range similarList {
		if bestMatch == nil || similar.Confidence > bestMatch.Confidence {
			bestMatch = &similar
		}
	}
	if bestMatch.Confidence <= 0.6 {
		bestMatch = nil
	}
	return bestMatch
}
