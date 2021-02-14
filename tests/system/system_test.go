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
	PICTURE_FOLDER = "./pictures"
)

var STAGING_API_URL = &url.URL{
	Scheme: "http",
	Host:   "157.245.246.219:8080",
}

type SystemTestSuite struct {
	suite.Suite
	localToRepoURI map[string]string
	httpClient     *http.Client
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *SystemTestSuite) SetupTest() {
	images, err := ioutil.ReadDir(PICTURE_FOLDER)
	if err != nil {
		panic(err)
	}
	suite.httpClient = &http.Client{}

	postImageUrl := STAGING_API_URL.ResolveReference(&url.URL{Path: "image"})
	for _, image := range images {
		fmt.Println(image.Name())
		body := buildMultipartFormDataBody(image.Name())
		req, err := http.NewRequest(
			http.MethodPost,
			postImageUrl.String(),
			body,
		)

		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("Content-Type", "multipart/form-data")

		res, err := suite.httpClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		b, err := ioutil.ReadAll(res.Body)

		if err != nil {
			panic(err)
		}
		uri := string(b)
		if res.StatusCode != 201 {
			suite.Failf("Fail", "Received status %s. Body: %s", res.Status, uri)
		}
		fmt.Println(uri)

		suite.localToRepoURI[image.Name()] = uri
	}
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *SystemTestSuite) TestImageSearchReflexivity() {
	searchByImageUrl := STAGING_API_URL.ResolveReference(&url.URL{Path: "_/search/_image"})

	SAMPLE_IMAGE_NAME := "cat_sky.jpg"

	body := buildMultipartFormDataBody(SAMPLE_IMAGE_NAME)
	req, err := http.NewRequest(
		http.MethodPost,
		searchByImageUrl.String(),
		body,
	)

	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "multipart/form-data")
	res, err := suite.httpClient.Do(req)

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	foundImageUris := make([]string, 0)
	if err := json.NewDecoder(res.Body).Decode(&foundImageUris); err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
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

	searchImageUrl := STAGING_API_URL.ResolveReference(&url.URL{Path: "_search"})
	searchImageUrl.Query().Add("tags", "[cat]")

	req, err := http.NewRequest(
		http.MethodGet,
		searchImageUrl.String(),
		nil,
	)

	if err != nil {
		panic(err)
	}

	res, err := suite.httpClient.Do(req)

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	foundImageUris := make([]string, 0)
	if err := json.NewDecoder(res.Body).Decode(&foundImageUris); err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
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

	searchImageUrl := STAGING_API_URL.ResolveReference(&url.URL{Path: "_search"})
	searchImageUrl.Query().Add("text", "A cat with a blue sky")

	req, err := http.NewRequest(
		http.MethodGet,
		searchImageUrl.String(),
		nil,
	)

	if err != nil {
		panic(err)
	}

	res, err := suite.httpClient.Do(req)

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	foundImageUris := make([]string, 0)
	if err := json.NewDecoder(res.Body).Decode(&foundImageUris); err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
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
	if err != nil {
		panic(err)
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	fi, err := file.Stat()
	if err != nil {
		panic(err)
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", fi.Name())
	if err != nil {
		panic(err)
	}
	part.Write(fileContents)
	return body
}
