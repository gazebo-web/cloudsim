package main

import (
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gitlab.com/ignitionrobotics/web/ign-go/testhelpers"
	"bytes"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

// Test utilities and helper functions (note: which are only compiled with `go test`)

const (
	apiVersion  string = "1.0"
	ctTextPlain string = "text/plain; charset=utf-8"
	ctJSON      string = "application/json"
	ctZip       string = "application/zip"
)

type uriTest struct {
	// description of the test
	testDesc string
	// a url (eg. /1.0/simulations)
	URL string
	// an optional JWT definition (can contain a plain jwt or a claims map)
	jwtGen *testJWT
	// optional expected ign.ErrMsg response. If the test case represents an error case
	// in such case, content type text/plain will be used
	expErrMsg *ign.ErrMsg
	// in case of error response, whether to parse the response body to get an ign.ErrMsg struct
	ignoreErrorBody bool
	// should we ignore invoking the OPTIONS method by default
	ignoreOptionsCall bool
}

type okCallback func(bslice *[]byte, resp *igntest.AssertResponse)

var tokenGeneratorPrivateKey string

// sptr returns a pointer to a given string.
// This function is especially useful when using string literals as argument.
func sptr(s string) *string {
	return &s
}

// intptr returns a pointer to a given int.
// This function is especially useful when using ints as arguments.
func intptr(n int) *int {
	return &n
}

// boolptr returns a pointer to a given bool.
// This function is especially useful when using bools as arguments.
func boolptr(b bool) *bool {
	return &b
}

// timeptr returns a pointer to a given Time.
// This function is especially useful when using Time as arguments.
func timeptr(t time.Time) *time.Time {
	return &t
}

// invokeURITest is a helper function invoked by tests to make the usual steps
// of invoking a GET route with a JWT and comparing against expected status and content
// type. It also optionally invokes the OPTIONS route, and optionally tests against
// an ign.ErrMsg.
// If the route invocation returns Status OK (200), then this function will invoke
// the okCallback to the calling test so specific comparisons can be done
// there.
func invokeURITest(t *testing.T, test uriTest, cb okCallback) {
	invokeURITestWithArgs(t, test, "GET", nil, cb)
}

// invokeURITestPOST is a helper function invoked by tests to make the usual steps
// of invoking a POST route with a JWT and Form fields, and comparing against expected
// status and content type. It also optionally invokes the OPTIONS route, and optionally tests against
// an ign.ErrMsg in case of error.
// If the route invocation returns Status OK (200), then this function will invoke
// the okCallback to the calling test so specific comparisons can be done
// there.
func invokeURITestMultipartPOST(t *testing.T, test uriTest, params map[string]string, cb okCallback) {
	jwt := getJWTToken(t, test.jwtGen)
	expEm, _ := errMsgAndContentType(test.expErrMsg, ctJSON)
	expStatus := expEm.StatusCode
	// first, check the OPTIONS method work
	if !test.ignoreOptionsCall {
		igntest.AssertRoute("OPTIONS", test.URL, http.StatusOK, t)
	}
	code, bslice, ok := igntest.SendMultipartPOST(t.Name(), t, test.URL, jwt, params, nil)
	require.Equal(t, expStatus, code)
	if expStatus != http.StatusOK && !test.ignoreErrorBody {
		igntest.AssertBackendErrorCode(t.Name(), bslice, expEm.ErrCode, t)
	} else if expStatus == http.StatusOK {
		var resp igntest.AssertResponse
		resp.Ok = ok
		resp.BodyAsBytes = bslice
		cb(bslice, &resp)
	}
}

// invokeURITestWithArgs is a helper function invoked by tests to make the usual steps
// of invoking a route with a JWT and optionally a Body, and comparing against expected
// status and content type. It also optionally invokes the OPTIONS route, and optionally tests against
// an ign.ErrMsg in case of error.
// If the route invocation returns Status OK (200), then this function will invoke
// the okCallback to the calling test so specific comparisons can be done
// there.
func invokeURITestWithArgs(t *testing.T, test uriTest, method string, b *bytes.Buffer, cb okCallback) {
	jwt := getJWTToken(t, test.jwtGen)
	expEm, expCt := errMsgAndContentType(test.expErrMsg, ctJSON)
	expStatus := expEm.StatusCode
	// first, check the OPTIONS method work
	if !test.ignoreOptionsCall {
		igntest.AssertRoute("OPTIONS", test.URL, http.StatusOK, t)
	}
	reqArgs := igntest.RequestArgs{Method: method, Route: test.URL, Body: b, SignedToken: jwt}
	resp := igntest.AssertRouteMultipleArgsStruct(reqArgs, expStatus, expCt, t)
	bslice := resp.BodyAsBytes
	require.Equal(t, expStatus, resp.RespRecorder.Code)
	if expStatus != http.StatusOK && !test.ignoreErrorBody {
		igntest.AssertBackendErrorCode(t.Name(), bslice, expEm.ErrCode, t)
	} else if expStatus == http.StatusOK {
		cb(bslice, resp)
	}
}

func getDefaultTestJWT() *testJWT {
	// Check for auth0 environment variables.
	return newJWT(os.Getenv("IGN_TEST_JWT"))
}

// testJWT is either a explicit jwt token , or a map of jwtClaims
// used to generate a jwt token (using the TOKEN_GENERATOR_PRIVATE_RSA256_KEY env var)
type testJWT struct {
	jwt       *string
	jwtClaims *jwt.MapClaims
}

// newClaimsJWT creates a testJWT definition using a map of claims
func newClaimsJWT(cl *jwt.MapClaims) *testJWT {
	return &testJWT{jwtClaims: cl}
}

// newJWT creates a new testJWT definition based on a given string token.
func newJWT(tk string) *testJWT {
	return &testJWT{jwt: &tk}
}

// getTestJWT - given an optional testJWT it creates and returns a token (or nil).
func getJWTToken(t *testing.T, jwtDef *testJWT) *string {
	if jwtDef != nil {
		s := generateJWT(*jwtDef, t)
		return &s
	}
	return nil
}

// generateJWT creates a JWT given a testJWT struct.
func generateJWT(jwt testJWT, t *testing.T) string {
	if jwt.jwt != nil {
		return *jwt.jwt
	}
	if tokenGeneratorPrivateKey == "" {
		tokenGeneratorPrivateKey = os.Getenv("TOKEN_GENERATOR_PRIVATE_RSA256_KEY")
	}
	testPrivateKeyAsPEM := []byte("-----BEGIN RSA PRIVATE KEY-----\n" + tokenGeneratorPrivateKey + "\n-----END RSA PRIVATE KEY-----")
	token, err := GenerateTokenRSA256(t, testPrivateKeyAsPEM, *jwt.jwtClaims)
	assert.NoError(t, err, "Error while generating token")
	return token
}

// GenerateTokenRSA256 generates an RSA256 token containing the given claims,
// the returns the token signed with the given PEM private key.
// Used with public - private keys.
func GenerateTokenRSA256(t *testing.T, pemPrivKey []byte, jwtClaims jwt.MapClaims) (signedToken string, err error) {
	signingKey, err := jwt.ParseRSAPrivateKeyFromPEM(pemPrivKey)
	if err != nil {
		require.NoError(t, err, "error while parsing private key")
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtClaims)
	ss, err := token.SignedString(signingKey)
	if err != nil {
		require.NoError(t, err, "error while encoding token into signed string")
		return
	}
	// All OK !
	signedToken = ss
	return
}

// Generate a new test JWT token with the given identity.
func createJWTForIdentity(t *testing.T, identity string) string {
	return generateJWT(testJWT{jwtClaims: &jwt.MapClaims{"sub": identity}}, t)
}

// errMsgAndContentType is a helper that given an optional errMsg and a content type to use
// when OK (ie. http status code 200), it returns a tuple with the ErrMsg and contentType to use
// in a subsequent call to 'igntest.AssertRouteMultipleArgs'.
// It was created to reduce LOC.
func errMsgAndContentType(em *ign.ErrMsg, successCT string) (ign.ErrMsg, string) {
	if em != nil {
		return *em, ctTextPlain
	}
	return ign.ErrorMessageOK(), successCT
}
