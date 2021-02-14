package deepdetect_test

import (
	"fmt"
	"testing"

	. "github.com/AsFal/shopify-application/internal/pkg/classifier/deepdetect"
	"github.com/AsFal/shopify-application/internal/pkg/imgrepo"
	"github.com/AsFal/shopify-application/internal/pkg/search"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

type NewDeepDetectClassifierSuite struct {
	suite.Suite
	host               string
	serviceInitRequest string
}

func (suite *NewDeepDetectClassifierSuite) SetupTest() {
	suite.host = "deepdetect.com"
	suite.serviceInitRequest = fmt.Sprintf("POST http://%s/services/ggnet", suite.host)
}

func (suite *NewDeepDetectClassifierSuite) TestShouldInitializeDeepDetectService() {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://deepdetect.com/services/ggnet",
		httpmock.NewStringResponder(200, ""),
	)

	classifier, err := NewDeepDetectClassifier(suite.host)

	suite.NotNil(classifier)
	suite.Nil(err)

	info := httpmock.GetCallCountInfo()
	suite.Equal(1, info[suite.serviceInitRequest])
}

func (suite *NewDeepDetectClassifierSuite) TestShouldNotFailWhenInitializingExistingService() {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://deepdetect.com/services/ggnet",
		httpmock.NewStringResponder(409, ""),
	)

	classifier, err := NewDeepDetectClassifier(suite.host)
	suite.NotNil(classifier)
	suite.Nil(err)

	info := httpmock.GetCallCountInfo()
	suite.Equal(1, info[suite.serviceInitRequest])
}

func (suite *NewDeepDetectClassifierSuite) TEstShouldFailIfNon409ErrorReceived() {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://deepdetect.com/services/ggnet",
		httpmock.NewStringResponder(400, ""), // TODO change for real DD possible error codes
	)

	classifier, err := NewDeepDetectClassifier(suite.host)
	suite.Nil(classifier)
	suite.NotNil(err)

	info := httpmock.GetCallCountInfo()
	suite.Equal(1, info[suite.serviceInitRequest])
}

// TODO: Test to validate request body

func TestNewDeepDetectClassifierSuite(t *testing.T) {
	suite.Run(t, new(NewDeepDetectClassifierSuite))
}

type DeepDetectClassifierClassifySuite struct {
	suite.Suite
	host            string
	url             string
	classifyRequest string
	classifier      *DeepDetectClassifier
}

func (suite *DeepDetectClassifierClassifySuite) SetupTest() {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	suite.host = "deepdetect.com"
	suite.url = fmt.Sprintf("http://%s/predict", suite.host)
	suite.classifyRequest = fmt.Sprintf("POST %s", suite.url)

	httpmock.RegisterResponder("POST", "http://deepdetect.com/services/ggnet",
		httpmock.NewStringResponder(200, ""),
	)

	suite.classifier, _ = NewDeepDetectClassifier(suite.host)
}

func (suite *DeepDetectClassifierClassifySuite) TestSendsURIToDeepDetectPredictEndpoint() {
	suite.T().Skip()
}

func (suite *DeepDetectClassifierClassifySuite) TestParsesDeepDetectTagsIntoSearchTags() {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Response taken from https://www.deepdetect.com/server/docs/api/#predict
	sample_response := `{"status":{"code":200,"msg":"OK"},"head":{"method":"/predict","time":1591.0,"service":"imageserv"},"body":{"predictions":{"uri":"http://i.ytimg.com/vi/0vxOhd4qlnA/maxresdefault.jpg","loss":0.0,"classes":[{"prob":0.24278657138347627,"cat":"n03868863 oxygen mask"},{"prob":0.20703653991222382,"cat":"n03127747 crash helmet"},{"prob":0.07931024581193924,"cat":"n03379051 football helmet"}]}}}`

	httpmock.RegisterResponder("POST", suite.url,
		httpmock.NewStringResponder(200, sample_response),
	)

	SAMPLE_IMG_URI := "hello.png"
	tags, err := suite.classifier.Classify(imgrepo.ImgURI(SAMPLE_IMG_URI))

	suite.Nil(err)

	EXPECTED_TAGS := []string{
		"oxygen", "mask", "crash", "helmet", "football", "helmet",
	}
	for _, tag := range EXPECTED_TAGS {
		suite.True(tags.Contains(tag))
	}

	info := httpmock.GetCallCountInfo()
	suite.Equal(1, info[suite.classifyRequest])
}

func (suite *DeepDetectClassifierClassifySuite) TestReturnEmptyTagsWhenReceivingClientError() {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", suite.url,
		httpmock.NewStringResponder(409, ""),
	)

	SAMPLE_IMG_URI := "hello.png"
	tags, err := suite.classifier.Classify(imgrepo.ImgURI(SAMPLE_IMG_URI))

	suite.NotNil(err)
	suite.Equal(search.Tags(""), tags)

	info := httpmock.GetCallCountInfo()
	suite.Equal(1, info[suite.classifyRequest])
}

func TestDeepDetectClassifierClassifySuite(t *testing.T) {
	suite.Run(t, new(DeepDetectClassifierClassifySuite))
}
