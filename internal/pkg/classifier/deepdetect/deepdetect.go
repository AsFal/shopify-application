package deepdetect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/AsFal/shopify-application/internal/pkg/imgrepo"
	"github.com/AsFal/shopify-application/internal/pkg/search"
)

const classifierService = "ggnet"

type DeepDetectClassifier struct {
	baseURL    *url.URL
	httpClient *http.Client
}

type HTTPClientError struct {
	status string
	msg    string
}

func (e *HTTPClientError) Error() string {
	return fmt.Sprintf("HTTP Client Error: %s\n %s", e.status, e.msg)
}

func isClientError(res *http.Response) bool {
	return res.StatusCode/100 == 4
}

func isClientConflictError(res *http.Response) bool {
	return res.StatusCode == 409
}

func newHTTPClientError(res *http.Response) *HTTPClientError {
	b, _ := ioutil.ReadAll(res.Body)
	return &HTTPClientError{
		status: res.Status,
		msg:    string(b),
	}
}

func NewDeepDetectClassifier(host string) (*DeepDetectClassifier, error) {
	c := &DeepDetectClassifier{
		baseURL: &url.URL{
			Scheme: "http",
			Host:   host,
		},
		httpClient: &http.Client{},
	}

	// Create Image Service
	body := map[string]interface{}{
		"mllib":       "caffe",
		"description": "image classification service",
		"type":        "supervised",
		"parameters": map[string]interface{}{
			"input": map[string]string{
				"connector": "image",
			},
			"mllib": map[string]interface{}{
				"nclasses": 1000,
			},
		},
		"model": map[string]string{
			"repository": "/opt/models/ggnet/",
		},
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		c.baseURL.ResolveReference(&url.URL{Path: "services/" + classifierService}).String(),
		buf,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	} else if isClientError(res) && !isClientConflictError(res) {
		defer res.Body.Close()
		return nil, newHTTPClientError(res)
	}

	return c, nil
}

type PredictRequest struct {
	Service    string                   `json:"service"`
	Parameters PredictRequestParameters `json:"parameters"`
	Data       []string                 `json:"data"`
}

type PredictRequestParameters struct {
	Output PredictRequestParametersOutput `json:"output"`
}

type PredictRequestParametersOutput struct {
	Best int `json:"best"`
}

type PredictResponse struct {
	Status interface{}         `json:"status"`
	Head   interface{}         `json:"head"`
	Body   PredictResponseBody `json:"body"`
}

type PredictResponseBody struct {
	Predictions []struct {
		Uri     string         `json:"uri"`
		Loss    float32        `json:"loss"`
		Classes []PredictClass `json:"classes"`
	} `json:"predictions"`
}

type PredictClass struct {
	Prob float64 `json:"prob"`
	Cat  string  `json:"cat"`
}

func (c *DeepDetectClassifier) Classify(imgURI imgrepo.ImgURI) (search.Tags, error) {
	// TODO: Properly set the service name
	buf := new(bytes.Buffer)
	body := PredictRequest{
		Service: classifierService,
		Parameters: PredictRequestParameters{
			Output: PredictRequestParametersOutput{
				Best: 3,
			},
		},
		Data: []string{
			string(imgURI),
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
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if err != nil {
		return "", err
	} else if isClientError(res) {
		return "", newHTTPClientError(res)
	}

	wrapper := new(PredictResponse)
	if err := json.NewDecoder(res.Body).Decode(wrapper); err != nil {
		return "", err
	}

	allClasses := ""
	for _, class := range wrapper.Body.Predictions[0].Classes {
		allClasses += " " + class.Cat
	}
	return search.Tags(allClasses), err
}
