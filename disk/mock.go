package disk

import (
	"fmt"
	"net/http"
)

type MockProvider struct {
}

func (m MockProvider) SignedPutURL(key, filetype string, public bool) (url string, headers http.Header, err error) {
	headers = make(http.Header)
	headers.Add("", "")
	return fmt.Sprintf("http://dl.test.com/%s", key), headers, nil
}

func (m MockProvider) SignedGetURL(key, filename string) (url string, err error) {
	return fmt.Sprintf("http://dl.test.com/%s", key), nil
}

func (m MockProvider) PublicURL(key string) (url string) {
	return fmt.Sprintf("http://dl.test.com/%s", key)
}

func (m MockProvider) ObjectDelete(key string) error {
	return nil
}

func (m MockProvider) ObjectsDelete(keys []string) error {
	return nil
}
