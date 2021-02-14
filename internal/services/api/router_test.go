package api_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/suite"
)

// func TestRegisterUnauthenticatedUnavailableUsername(t *testing.T) {
// 	saveLists()
// 	w := httptest.NewRecorder()
//
// 	r := getRouter(true)
//
// 	r.POST("/u/register", register)
//
// 	registrationPayload := getLoginPOSTPayload()
// 	req, _ := http.NewRequest("POST", "/u/register", strings.NewReader(registrationPayload))
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	req.Header.Add("Content-Length", strconv.Itoa(len(registrationPayload)))
//
// 	r.ServeHTTP(w, req)
//
// 	if w.Code != http.StatusBadRequest {
// 		t.Fail()
// 	}
// 	restoreLists()
// }

type PostRootSuite struct {
	suite.Suite
	router *gin.Engine
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *PostRootSuite) SetupTest() {
	// Setup router
	return
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *PostRootSuite) TestBasic() {
	assert.Equal(suite.T(), true, true)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPostRootSuite(t *testing.T) {
	suite.Run(t, new(PostRootSuite))
}
