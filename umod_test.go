package umod

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type mockClient struct {
	responses  map[string]string
	err        error
	statusCode int
}

func (m *mockClient) Get(url string) (*http.Response, error) {
	res, ok := m.responses[url]
	if !ok {
		fmt.Println("first call to", url, " creating mock response")
		r, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		m.responses[url] = string(b)
		if err := m.save(); err != nil {
			return nil, err
		}
		return &http.Response{
			StatusCode: m.statusCode,
			Body:       ioutil.NopCloser(bytes.NewReader(b)),
		}, m.err
	}
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(strings.NewReader(res)),
	}, m.err
}

func (m *mockClient) save() error {
	b, err := json.Marshal(m.responses)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("tests.json", b, 0644)
}

var (
	clientMock = &mockClient{}
)

func init() {
	results := map[string]string{}
	b, err := ioutil.ReadFile("tests.json")
	if err == nil {
		if err := json.Unmarshal(b, &results); err != nil {
			panic(err)
		}
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
	if resp.Plugins[0].LatestReleaseAtAtom.IsZero() {
		t.Errorf("Search() returned zero time")
	}

	t.Run("Category", func(t *testing.T) {
		categories := []Category{Universal, _7DaysToDie, Hurtworld, ReignOfKings, Rust, TheForest}

		for _, c := range categories {
			resp, err := Search("", Categories(c))
			if err != nil {
				t.Errorf("Search() returned error: %v", err)
			}
			if len(resp.Plugins) == 0 {
				t.Errorf("Search() returned no results %v", c)
			}
			for _, p := range resp.Plugins {
				hasSupport := false
				for _, g := range p.GamesDetail {
					if g.Slug == string(c) {
						hasSupport = true
					}
				}
				if !hasSupport {
					t.Errorf("Search() returned wrong games, expecting: %v, got: %v", string(c), p.GamesDetail)
				}
			}
		}
	})
	t.Run("Page", func(t *testing.T) {
		resp, err := Search("", Page(2))
		if err != nil {
			t.Errorf("Search() returned error: %v", err)
		}
		if resp.Plugins[0].LatestReleaseAtAtom.IsZero() {
			t.Errorf("Search() returned zero time")
		}
		if resp.CurrentPage != 2 {
			t.Errorf("Search() returned wrong page, expecting: 2, got: %v", resp.CurrentPage)
		}
	})
	t.Run("Tag", func(t *testing.T) {
		tags := []string{"fun", "mechanics"}
		resp, err = Search("", Tags(tags...))
		if err != nil {
			t.Errorf("Search() returned error: %v", err)
		}
		if len(resp.Plugins) == 0 {
			t.Errorf("Search() returned no results")
		}
		for _, p := range resp.Plugins {
			for _, tag := range tags {
				if !strings.Contains(p.TagsAll, tag) {
					t.Errorf("Search() returned wrong tags, expecting: %v, got: %v", tags, p.TagsAll)
				}
			}
		}
	})
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
			if next.CurrentPage != resp.CurrentPage+1 {
				t.Errorf("NextPage() returned wrong page number, expecting: %v, got: %v", resp.CurrentPage+1, next.CurrentPage)
			}
			if next.Total != resp.Total {
				t.Errorf("NextPage() returned wrong total, expecting: %v, got: %v", resp.Total, next.Total)
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
			prev, err := next.PrevPage()
			if err != nil {
				t.Errorf("PrevPage() returned error: %v", err)
			}
			if prev.CurrentPage != next.CurrentPage-1 {
				t.Errorf("PrevPage() returned wrong page number, expecting: %v, got: %v", resp.CurrentPage-1, prev.CurrentPage)
			}
			if prev.Total != next.Total {
				t.Errorf("PrevPage() returned wrong total, expecting: %v, got: %v", next.Total, prev.Total)
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
	if resp.Plugins[0].LatestReleaseAtAtom.IsZero() {
		t.Errorf("Latest() returned zero time")
	} // TODO test is in ascending order
}

func TestOldest(t *testing.T) {
	resp, err := Oldest()
	if err != nil {
		t.Errorf("Oldest() returned error: %v", err)
	}
	if resp.Plugins[0].LatestReleaseAtAtom.IsZero() {
		t.Errorf("Oldest() returned zero time")
	}
}

func TestGames(t *testing.T) {
	resp, err := Games()
	if err != nil {
		t.Errorf("Games() returned error: %v", err)
	}
	if resp[0].LatestReleaseAtAtom.IsZero() {
		t.Errorf("Games() returned zero time")
	}

	t.Run("Search", func(t *testing.T) {
		for _, g := range resp {
			if g.Slug == "heat" {
				continue // has no plugins
			}
			resp, err := g.Search("")
			if err != nil {
				t.Errorf("Search() returned error: %v", err)
			}
			if len(resp.Plugins) == 0 {
				t.Errorf("Search() returned no results")
			}
			if !searchHasGame(resp, g.Slug) {
				t.Errorf("Search() returned wrong game, expecting: %v, got: %v", g.Slug, resp)
			}
		}
	})
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
		clientMock.responses["https://umod.org/plugins/search.json?query=test"] = "invalid json"
		_, err := Search("test")
		if err == nil {
			t.Errorf("Search() should return error")
		}
	})
}

func searchHasGame(resp SearchResponse, gameSlug string) bool {
	for _, p := range resp.Plugins {
		for _, g := range p.GamesDetail {
			if g.Slug == gameSlug {
				return true
			}
		}
	}
	return false
}
