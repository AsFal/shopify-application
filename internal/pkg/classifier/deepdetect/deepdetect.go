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
)

type DeepDetectClassifier struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func NewDeepDetectClassifier(host string) *DeepDetectClassifier {
	return &DeepDetectClassifier{
		baseURL: &url.URL{
			Scheme: "https",
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
		uri     string
		loss    float32
		classes []PredictClass
	}
}

func (c *DeepDetectClassifier) Classify(file multipart.File) (string, error) {
	reader := bufio.NewReader(file)
	content, _ := ioutil.ReadAll(reader)
	encoded := base64.StdEncoding.EncodeToString(content)

	// TODO: Properly set the service name
	body := PredictReqBody{
		service: "service",
		data: []string{
			encoded,
		},
	}

	b, _ := json.Marshal(body)
	// TODO: Should never happen given I control the input?

	req, err := http.NewRequest(
		http.MethodPost,
		c.baseURL.ResolveReference(&url.URL{Path: "/predict"}).String(),
		bytes.NewReader(b),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	var resBody *PredictResBody
	err = json.Unmarshal(buf.Bytes(), resBody)

	// Transform classess to a more standard format

	allClasses := ""
	for _, class := range resBody.predictions.classes {
		allClasses += " " + class.cat
	}
	return allClasses, err
}
