package face

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lentregu/eq2017/HumanIdentification/comms"
	"github.com/lentregu/eq2017/oxford"
	"io"
	"strings"
)
type printOption int

const (
	pretty printOption = iota
	normal
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

type findWhoisURLRequestType struct {
	ImageFileName string `json:"img_path"`
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

	//faceService := oxford.NewFace("567c560aa85245418459b82634bc7a98")	
	faceService := oxford.NewFace("83dc246bac2b447782b5aab70604bc97")
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

func WhoisURL(w http.ResponseWriter, r *http.Request) {
	requestBody := findWhoisURLRequestType{}
	json.NewDecoder(r.Body).Decode(&requestBody)
	fmt.Println("WHOIS------")
	fmt.Print(formatRequest(r))

	fmt.Print(toJSON(requestBody, pretty))

	//faceService := oxford.NewFace("567c560aa85245418459b82634bc7a98")	
	faceService := oxford.NewFace("83dc246bac2b447782b5aab70604bc97")
	
	var bestMatch *oxford.FaceSimilarResponseType = nil

	var email string = ""


	faceID, _ := faceService.DetectBinFromFile(requestBody.ImageFileName)

	fmt.Printf("El faceID es: %s\n", faceID)

	similarList, err := faceService.FindSimilar(faceID, requestBody.FaceListID)

	if err == nil {
		bestMatch = getBestMatch(similarList)

		if bestMatch != nil {
				faces, _:= faceService.GetObjectFacesInAList(requestBody.FaceListID)
				for _, face := range faces.PersistedFaces {
					if face.PersistedFaceID == bestMatch.PersistedFaceID {
						email = face.UserData
						fmt.Printf("Todo ok, el email es: %s", email)
					}
				}
				fmt.Printf("El FaceID detectado es: %s\n", bestMatch.PersistedFaceID)
			} else {
				err = fmt.Errorf("User Not Found")
			}
	}
	fmt.Printf("--------------------------Whois Response-----------------------------------")

	if email != "" {
		w.Header().Set("Content-Type", "application/json")
    	result, _ := json.Marshal(map[string]string{"email": email})
    	io.WriteString(w, string(result))
		w.WriteHeader(200)
	} else {
		err = fmt.Errorf("User Not Found")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(404)
		fmt.Fprintf(w, err.Error())
	}
}


func Whois(w http.ResponseWriter, r *http.Request) {
	requestBody := findWhoisRequestType{}
	json.NewDecoder(r.Body).Decode(&requestBody)
	fmt.Println("WHOIS------")
	fmt.Print(formatRequest(r))

	fmt.Print(toJSON(requestBody, pretty))

	//faceService := oxford.NewFace("567c560aa85245418459b82634bc7a98")	
	faceService := oxford.NewFace("83dc246bac2b447782b5aab70604bc97")

	var bestMatch *oxford.FaceSimilarResponseType = nil
	var similarList []oxford.FaceSimilarResponseType = nil
	binaryImg, err := oxford.Base64ToByteArray(requestBody.Base64Image)

	var email string = ""

	if err == nil {
		faceID, _ := faceService.DetectBin(binaryImg)

		fmt.Printf("El faceID es: %s\n", faceID)

		similarList, err = faceService.FindSimilar(faceID, requestBody.FaceListID)
        
		if err == nil {
			bestMatch = getBestMatch(similarList)

			if bestMatch != nil {
				faces, _:= faceService.GetObjectFacesInAList(requestBody.FaceListID)
				for _, face := range faces.PersistedFaces {
					if face.PersistedFaceID == bestMatch.PersistedFaceID {
						email = face.UserData
						fmt.Printf("Todo ok, el email es: %s", email)
					}
				}
				fmt.Printf("El FaceID detectado es: %s\n", bestMatch.PersistedFaceID)
			} else {
				err = fmt.Errorf("User Not Found")
			}
		}
	}

	fmt.Printf("--------------------------Whois Response-----------------------------------")

	if email != "" {
		w.Header().Set("Content-Type", "application/json")
    	result, _ := json.Marshal(map[string]string{"email": email})
    	io.WriteString(w, string(result))
		w.WriteHeader(200)
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
	fmt.Println("---------getBestMatch-------------")
	fmt.Print(similarList)
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