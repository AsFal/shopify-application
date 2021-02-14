// +build system

package system_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	PICTURE_FOLDER  = "tests/system/pictures"
	STAGING_API_URL = ""
)

type SystemTestSuite struct {
	suite.Suite
	localToRepoURI map[string]string
	httpClient     *http.Client
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *SystemTestSuite) SetupTest() {
	images, _ := ioutil.ReadDir(PICTURE_FOLDER)
	suite.httpClient = &http.Client{}
	STAGING_URL_IMAGE := path.Join(STAGING_API_URL, "image")
	for _, image := range images {
		fmt.Println(image.Name())
		body := buildMultipartFormDataBody(image.Name())
		req, err := http.NewRequest(
			http.MethodPost,
			STAGING_URL_IMAGE,
			body,
		)
		req.Header.Set("Content-Type", "multipart/form-data")

		res, err := suite.httpClient.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()
		b, err := ioutil.ReadAll(res.Body)
		uri := string(b)
		fmt.Println(uri)

		suite.localToRepoURI[image.Name()] = uri
	}
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *SystemTestSuite) TestImageSearchReflexivity() {
	STAGING_URL_SEARCH_IMAGE := path.Join(STAGING_API_URL, "_search/_image")

	SAMPLE_IMAGE_NAME := "cat_sky.jpg"

	body := buildMultipartFormDataBody(SAMPLE_IMAGE_NAME)
	req, err := http.NewRequest(
		http.MethodPost,
		STAGING_URL_SEARCH_IMAGE,
		body,
	)
	req.Header.Set("Content-Type", "multipart/form-data")
	res, err := suite.httpClient.Do(req)
	defer res.Body.Close()

	foundImageUris := make([]string, 0)
	if err := json.NewDecoder(res.Body).Decode(&foundImageUris); err != nil {
		fmt.Println(err)
	}
	suite.True(
		contains(foundImageUris, suite.localToRepoURI[SAMPLE_IMAGE_NAME]),
		"%s (uri: %s) missing from found uris %s",
		SAMPLE_IMAGE_NAME, suite.localToRepoURI[SAMPLE_IMAGE_NAME],
		foundImageUris,
	)
}

func (suite *SystemTestSuite) TestBasicSearchByTag() {
	SAMPLE_IMAGE_NAME := "cat_sky.jpg"

	searchImageUrl := &url.URL{
		Scheme: "http",
		Host:   STAGING_API_URL,
		Path:   "_search",
	}
	searchImageUrl.Query().Add("tags", "[cat]")

	req, err := http.NewRequest(
		http.MethodGet,
		searchImageUrl.String(),
		nil,
	)

	res, err := suite.httpClient.Do(req)
	defer res.Body.Close()

	foundImageUris := make([]string, 0)
	if err := json.NewDecoder(res.Body).Decode(&foundImageUris); err != nil {
		fmt.Println(err)
	}
	suite.True(
		contains(foundImageUris, suite.localToRepoURI[SAMPLE_IMAGE_NAME]),
		"%s (uri: %s) missing from found uris %s",
		SAMPLE_IMAGE_NAME, suite.localToRepoURI[SAMPLE_IMAGE_NAME],
		foundImageUris,
	)
}

func (suite *SystemTestSuite) TestBasicFullTextSearch() {
	SAMPLE_IMAGE_NAME := "cat_sky.jpg"

	searchImageUrl := &url.URL{
		Scheme: "http",
		Host:   STAGING_API_URL,
		Path:   "_search",
	}
	searchImageUrl.Query().Add("text", "A cat with a blue sky")

	req, err := http.NewRequest(
		http.MethodGet,
		searchImageUrl.String(),
		nil,
	)

	res, err := suite.httpClient.Do(req)
	defer res.Body.Close()

	foundImageUris := make([]string, 0)
	if err := json.NewDecoder(res.Body).Decode(&foundImageUris); err != nil {
		fmt.Println(err)
	}
	suite.True(
		contains(foundImageUris, suite.localToRepoURI[SAMPLE_IMAGE_NAME]),
		"%s (uri: %s) missing from found uris %s",
		SAMPLE_IMAGE_NAME, suite.localToRepoURI[SAMPLE_IMAGE_NAME],
		foundImageUris,
	)
}

func contains(xs []string, y string) bool {
	for _, x := range xs {
		if x == y {
			return true
		}
	}
	return false
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSystemTestSuite(t *testing.T) {
	suite.Run(t, new(SystemTestSuite))
}

func buildMultipartFormDataBody(imageName string) io.Reader {
	// TODO: Properly set the service name

	file, err := os.Open(path.Join(PICTURE_FOLDER, imageName))
	fileContents, err := ioutil.ReadAll(file)
	fi, err := file.Stat()
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", fi.Name())
	part.Write(fileContents)
	return body
}

func TestSomething() {
	// TODO: Properly set the service name

	PICTURE_FOLDER := "tests/system/pictures"
	file, err := os.Open(path.Join(PICTURE_FOLDER, "cat_sky.jpg"))
	fileContents, err := ioutil.ReadAll(file)
	fi, err := file.Stat()
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", fi.Name())
	part.Write(fileContents)

	STAGING_URL := ""
	STAGING_URL_ROOT := STAGING_URL
	req, err := http.NewRequest(
		http.MethodPost,
		STAGING_URL_ROOT,
		body,
	)
	req.Header.Set("Content-Type", "multipart/form-data")

	httpClient := &http.Client{}
	_, err := httpClient.Do(req)
}
