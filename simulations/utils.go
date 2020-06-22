package simulations

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/go-playground/form"
	"gitlab.com/ignitionrobotics/web/cloudsim/globals"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/go-playground/validator.v9"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Utils includes helper functions

// Sleep pauses the current goroutine for at least the duration d. A negative or
// zero duration causes Sleep to return immediately.
// Dev note: we redefine the Sleep function here to allow mocking during testing.
var Sleep = func(d time.Duration) {
	time.Sleep(d)
}

// sptr returns a pointer to a given string.
// This function is specially useful when using string literals as argument.
func sptr(s string) *string {
	return &s
}

// intptr returns a pointer to a given int.
func intptr(i int) *int {
	return &i
}

func int32ptr(i int32) *int32 {
	return &i
}

func int64ptr(i int64) *int64 {
	return &i
}

func boolptr(b bool) *bool {
	return &b
}

func timeptr(t time.Time) *time.Time {
	return &t
}

// cloneStringsMap creates a new strings map by cloning the given map. It creates
// a shallow copy of the input map.
func cloneStringsMap(toClone map[string]string) map[string]string {
	newMap := make(map[string]string)
	for k, v := range toClone {
		newMap[k] = v
	}
	return newMap
}

func logger(ctx context.Context) ign.Logger {
	return ign.LoggerFromContext(ctx)
}

func timeTrack(ctx context.Context, start time.Time, name string) {
	elapsed := time.Since(start)
	logger(ctx).Info(fmt.Sprintf("%s took %s", name, elapsed))
}

// StrSliceContains determines if val is contained within list.
func StrSliceContains(val string, list []string) bool {
	for _, s := range list {
		if s == val {
			return true
		}
	}
	return false
}

// SliceToStr joins the elements of a string array using comma as the separator.
func SliceToStr(slice []string) string {
	return strings.Join(slice, ",")
}

// EnvVarToSlice reads the contents of an environment variable and splits it into an array of strings using comma as
// the separator.
func EnvVarToSlice(envVar string) []string {
	s, _ := ign.ReadEnvVar(envVar)
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

// StrToIntSlice converts a string to an array of ints, where each element is a section of the complete number.
func StrToIntSlice(str string) ([]int, error) {
	if str == "" {
		return nil, nil
	}
	noSpaces := strings.TrimSpace(str)
	noSpaces = strings.TrimPrefix(noSpaces, ",")
	noSpaces = strings.TrimSuffix(noSpaces, ",")
	var result []int
	for _, numStr := range strings.Split(noSpaces, ",") {
		numStr = strings.TrimSpace(numStr)
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, err
		}
		result = append(result, num)
	}
	return result, nil
}

// ParseStruct reads the http request and decodes sent values
// into the given struct. It uses the isForm bool to know if the values comes
// as "request.Form" values or as "request.Body".
// It also calls validator to validate the struct fields.
func ParseStruct(s interface{}, r *http.Request, isForm bool) *ign.ErrMsg {
	// TODO: stop using globals. Move to own packages.
	if isForm {
		if errs := globals.FormDecoder.Decode(s, r.Form); errs != nil {
			return ign.NewErrorMessageWithArgs(ign.ErrorFormInvalidValue, errs,
				getDecodeErrorsExtraInfo(errs))
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(s); err != nil {
			return ign.NewErrorMessageWithBase(ign.ErrorUnmarshalJSON, err)
		}
	}
	// Validate struct values
	if em := ValidateStruct(s); em != nil {
		return em
	}
	return nil
}

// ValidateStruct Validate struct values using golang validator.v9
func ValidateStruct(s interface{}) *ign.ErrMsg {
	if errs := globals.Validate.Struct(s); errs != nil {
		return ign.NewErrorMessageWithArgs(ign.ErrorFormInvalidValue, errs,
			getValidationErrorsExtraInfo(errs))
	}
	return nil
}

// Builds the ErrMsg extra info from the given DecodeErrors
func getDecodeErrorsExtraInfo(err error) []string {
	errs := err.(form.DecodeErrors)
	extra := make([]string, 0, len(errs))
	for field, er := range errs {
		extra = append(extra, fmt.Sprintf("Field: %s. %v", field, er.Error()))
	}
	return extra
}

// Builds the ErrMsg extra info from the given ValidationErrors
func getValidationErrorsExtraInfo(err error) []string {
	validationErrors := err.(validator.ValidationErrors)
	extra := make([]string, 0, len(validationErrors))
	for _, fe := range validationErrors {
		extra = append(extra, fmt.Sprintf("%s:%v", fe.StructField(), fe.Value()))
	}
	return extra
}

// getLocalIPAddressString returns a local IP address of this host.
func getLocalIPAddressString() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	var ip string
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}
	return ip, nil
}

// Min returns the minimum value between two ints.
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// Max returns the maximum value between two ints.
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// generateToken returns a random hexadecimal token of the specified length.
// `size` is the length of the token in bytes. Defaults to 32 bytes.
func generateToken(size *int) (string, error) {
	// Set size default value
	if size == nil {
		size = intptr(32)
	}

	b := make([]byte, *size)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

// IsWebsocketAddress checks that the given address is a valid websocket address for cloudsim.
// If a groupd is provided, it will check that the given address includes a group ID.
func IsWebsocketAddress(addr string, groupID *string) bool {
	if !strings.Contains(addr, "cloudsim-ws.ignitionrobotics.org") {
		return false
	}

	if !strings.Contains(addr, "/simulations/") {
		return false
	}

	if groupID != nil {
		if !strings.Contains(addr, *groupID) {
			return false
		}
	}

	return true
}