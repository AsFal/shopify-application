package deepdetect

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"log"
)

type DeepDetectClassifier struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func NewDeepDetectClassifier(host string) *DeepDetectClassifier {
	return &DeepDetectClassifier{
		baseURL: &url.URL{
			Scheme: "http",
			Host:   host,
		},
		httpClient: &http.Client{},
	}

}

type PredictReqBody struct {
	service string   `json:"service"`
	data    []string `json:"data"`
}

type PredictClass struct {
	prob float64
	cat  string
}
type PredictResBody struct {
	predictions struct {
		uri     string `json:"uri"`
		loss    float32 `json:"loss"`
		classes []PredictClass `json:"classes"`
	} `json:""prediction"`
}

func (c *DeepDetectClassifier) Classify(file multipart.File) (string, error) {
	reader := bufio.NewReader(file)
	content, _ := ioutil.ReadAll(reader)
	encoded := base64.StdEncoding.EncodeToString(content)

	// TODO: Properly set the service name
	buf := new(bytes.Buffer)
	body := PredictReqBody{
		service: "service",
		data: []string{
			encoded,
		},
	}

	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		c.baseURL.ResolveReference(&url.URL{Path: "/predict"}).String(),
		buf,
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil{
		return "", err
	}

	defer res.Body.Close()
	resBody := new(PredictResBody)
	if err := json.NewDecoder(res.Body).Decode(resBody); err != nil {
		return "", err
	}
	log.Println(*resBody)

	allClasses := ""
	for _, class := range resBody.predictions.classes {
		allClasses += " " + class.cat
	}
	return allClasses, err
}
