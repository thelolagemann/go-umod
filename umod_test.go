package umod

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

var (
	testURLs = []string{
		"https://umod.org/plugins/search.json?query=heli&sort=latest_release_at&sortdir=desc&page=1",
		"https://umod.org/plugins/search.json?query=heli&sort=latest_release_at&sortdir=desc&page=2",
		"https://umod.org/plugins/search.json?sort=latest_release_at&sortdir=desc&page=1",
	}
)

func generateTestData() error {
	results := map[string]string{}
	for _, url := range testURLs {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		results[url] = string(b)
		resp.Body.Close()
	}

	// for programmatically acquiring the last page
	var r SearchResponse
	if err := json.Unmarshal([]byte(results["https://umod.org/plugins/search.json?query=heli&sort=latest_release_at&sortdir=desc&page=1"]), &r); err != nil {
		return err
	}

	res, err := http.Get(fmt.Sprintf("%s%s", baseURL, r.LastPageURL))
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	results[fmt.Sprintf("%s%s", baseURL, r.LastPageURL)] = string(b)

	b, err = json.Marshal(results)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile("tests.json", b, 0644); err != nil {
		return err
	}

	return nil
}

type mockClient struct {
	responses  map[string]string
	err        error
	statusCode int
}

func (m *mockClient) Get(url string) (*http.Response, error) {
	res, ok := m.responses[url]
	if !ok {
		return nil, fmt.Errorf("not found: %s", url)
	}
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(strings.NewReader(res)),
	}, m.err
}

var (
	clientMock = &mockClient{}
)

func init() {
	if _, err := os.Stat("tests.json"); err != nil {
		if err := generateTestData(); err != nil {
			panic(err)
		}
	}
	b, err := ioutil.ReadFile("tests.json")
	if err != nil {
		panic(err)
	}
	var results map[string]string
	if err := json.Unmarshal(b, &results); err != nil {
		panic(err)
	}
	clientMock.responses = results
	clientMock.statusCode = http.StatusOK
	httpClient = clientMock
}

func TestSearch(t *testing.T) {
	resp, err := Search("heli")
	if err != nil {
		t.Fatalf("Search() returned error: %v", err)
	}
	if resp.Data[0].LatestReleaseAtAtom.IsZero() {
		t.Errorf("Search() returned zero time")
	}

}

func TestPagination(t *testing.T) {
	resp, err := Search("heli")
	if err != nil {
		t.Fatalf("Search() returned error: %v", err)
	}

	var next SearchResponse
	t.Run("Next", func(t *testing.T) {
		t.Run("HasNext", func(t *testing.T) {
			var err error
			next, err = resp.NextPage()
			if err != nil {
				t.Fatalf("NextPage() returned error: %v", err)
			}
			if next.Data[0].LatestReleaseAtAtom.IsZero() {
				t.Errorf("NextPage() returned zero time")
			}
		})
		t.Run("NoNext", func(t *testing.T) {
			resp.NextPageURL = ""
			_, err := resp.NextPage()
			if err == nil {
				t.Errorf("NextPage() should return error")
			}
		})
	})
	t.Run("Prev", func(t *testing.T) {
		t.Run("HasPrev", func(t *testing.T) {
			_, err := next.PrevPage()
			if err != nil {
				t.Errorf("PrevPage() returned error: %v", err)
			}
		})
		t.Run("NoPrev", func(t *testing.T) {
			_, err := resp.PrevPage()
			if err == nil {
				t.Errorf("Search() returned no error")
			}
		})
	})
	t.Run("Last", func(t *testing.T) {
		t.Run("HasLast", func(t *testing.T) {
			_, err := resp.LastPage()
			if err != nil {
				t.Errorf("LastPage() returned error: %v", err)
			}
		})
		t.Run("NoLast", func(t *testing.T) { // should technically never happen
			resp.LastPageURL = ""
			_, err := resp.LastPage()
			if err == nil {
				t.Errorf("LastPage() returned no error")
			}
		})
	})
}

func TestLatest(t *testing.T) {
	resp, err := Latest()
	if err != nil {
		t.Errorf("Latest() returned error: %v", err)
	}
	if resp.Data[0].LatestReleaseAtAtom.IsZero() {
		t.Errorf("Latest() returned zero time")
	}
}

func TestRequest(t *testing.T) {
	t.Run("ClientError", func(t *testing.T) {
		clientMock.err = fmt.Errorf("client error")
		defer func() { clientMock.err = nil }()
		_, err := Search("heli")
		if err == nil {
			t.Errorf("Search() should return error")
		}
	})
	t.Run("ServerError", func(t *testing.T) {
		clientMock.statusCode = http.StatusInternalServerError
		defer func() { clientMock.statusCode = http.StatusOK }()
		_, err := Search("heli")
		if err == nil {
			t.Errorf("Search() should return error")
		}
	})
	t.Run("InvalidJSON", func(t *testing.T) {
		clientMock.responses["https://umod.org/plugins/search.json?query=test&sort=latest_release_at&sortdir=desc&page=1"] = "invalid json"
		_, err := Search("test")
		if err == nil {
			t.Errorf("Search() should return error")
		}
	})
}
