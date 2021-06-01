package idp_test

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/elimity-com/scim"
	"io"
	"io/fs"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

//go:embed testdata
var testdata embed.FS

func TestIdP(t *testing.T) {
	idPs, _ := testdata.ReadDir("testdata")
	for _, idP := range idPs {
		var newServer func() scim.Server
		switch idP.Name() {
		case "okta":
			newServer = newOktaTestServer
		case "azuread":
			newServer = newAzureADTestServer
		}
		t.Run(idP.Name(), func(t *testing.T) {
			idpPath := fmt.Sprintf("testdata/%s", idP.Name())
			if err := fs.WalkDir(testdata, idpPath, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				raw, err := fs.ReadFile(testdata, path)
				if err != nil {
					return fmt.Errorf("%s: %v", path, err)
				}
				var test testCase
				if err := json.Unmarshal(raw, &test); err != nil {
					return fmt.Errorf("%s: %v", path, err)
				}
				fileName, _ := filepath.Rel(idpPath, path)
				t.Run(strings.TrimSuffix(fileName, ".json"), func(t *testing.T) {
					if err := testRequest(test, newServer); err != nil {
						t.Error(err)
					}
				})
				return nil
			}); err != nil {
				t.Error(err)
			}
		})
	}
}

func testRequest(t testCase, newServer func() scim.Server) error {
	rr := httptest.NewRecorder()
	var br io.Reader
	if len(t.Request) != 0 {
		br = bytes.NewReader(t.Request)
	}
	newServer().ServeHTTP(
		rr,
		httptest.NewRequest(t.Method, t.Path, br),
	)
	if code := rr.Code; code != t.StatusCode {
		return fmt.Errorf("expected %d, got %d\n", t.StatusCode, code)
	}
	if len(t.Response) != 0 {
		var response map[string]interface{}
		if err := unmarshal(rr.Body.Bytes(), &response); err != nil {
			return err
		}
		var reference map[string]interface{}
		if err := unmarshal(t.Response, &reference); err != nil {
			return err
		}
		if !reflect.DeepEqual(reference, response) {
			return fmt.Errorf("expected, got:\n%v\n%v", reference, response)
		}
	}
	return nil
}

type testCase struct {
	Request    json.RawMessage
	Response   json.RawMessage
	Method     string
	Path       string
	StatusCode int
}
