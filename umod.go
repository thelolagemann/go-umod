package umod

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	gamesURL  = "https://assets.umod.org/games.json"
	searchURL = "https://umod.org/plugins/search.json"
)

type client interface {
	Get(url string) (*http.Response, error)
}

var (
	httpClient client = &http.Client{
		Timeout: time.Second * 5,
	}
)

// SearchResponse is the structure of the response returned by the search endpoint.
// It contains a paginated list of plugins found, the total number of pages of plugins
// found, the URL of the previous page, the URL of the next page, and the URL of the
// last page.
//
// See: https://umod.org/plugins/search.json
type SearchResponse struct {
	CurrentPage  int      `json:"current_page"`
	Data         []Plugin `json:"data"`
	FirstPageURL string   `json:"first_page_url"`
	From         int      `json:"from"`
	LastPageNum  int      `json:"last_page"`
	LastPageURL  string   `json:"last_page_url"`
	NextPageURL  string   `json:"next_page_url"`
	Path         string   `json:"path"`
	PerPage      int      `json:"per_page"`
	PrevPageURL  string   `json:"prev_page_url"`
	To           int      `json:"to"`
	Total        int      `json:"total"`

	search string
}

// PrevPage returns the previous page of the search response, or an error if
// there is no previous page.
func (s SearchResponse) PrevPage() (SearchResponse, error) {
	if s.PrevPageURL == "" {
		return s, fmt.Errorf("no previous page")
	}
	return search(s.search, Page(s.CurrentPage-1))
}

// NextPage returns the next page of the search response, or an error if
// there is no next page.
func (s SearchResponse) NextPage() (SearchResponse, error) {
	if s.NextPageURL == "" {
		return s, fmt.Errorf("no next page")
	}
	return search(s.search, Page(s.CurrentPage+1))
}

// LastPage returns the last page of the search response, or an error if
// there is no last page (although technically that should never happen).
func (s SearchResponse) LastPage() (SearchResponse, error) {
	if s.LastPageURL == "" {
		return s, fmt.Errorf("no last page") // this should never happen ?
	}
	return search(s.search, Page(s.LastPageNum))
}

// SearchOption is a function that can be used to modify the parameters when using
// the Search function.
type SearchOption func(*url.Values)

// Categories specifies the categories that the plugins must be compatible with.
func Categories(categories ...Category) SearchOption {
	return func(v *url.Values) {
		for i, g := range categories {
			v.Set("categories["+strconv.Itoa(i)+"]", string(g))
		}
	}
}

// Category is a simple string wrapper that wraps a category slug. It is used
// with the Categories functional option to specify the categories that the
// plugins must be compatible with.
type Category string

const (
	// CategoryUniversal is the category for plugins that are not specific to any game.
	CategoryUniversal Category = "universal"
	// Category7DaysToDie is the category for plugins that are compatible with 7 days to die
	Category7DaysToDie Category = "7-days-to-die"
	// CategoryHurtworld is the category for plugins that are compatible with Hurtworld
	CategoryHurtworld Category = "hurtworld"
	// CategoryReignOfKings is the category for plugins that are compatible with Reign of Kings
	CategoryReignOfKings Category = "reign-of-kings"
	// CategoryRust is the category for plugins that are compatible with Rust
	CategoryRust Category = "rust"
	// CategoryTheForest is the category for plugins that are compatible with The Forest
	CategoryTheForest Category = "the-forest"
)

// TODO: add description
func Tags(tags ...string) SearchOption {
	return func(v *url.Values) {
		for i, tag := range tags {
			v.Set("tags["+strconv.Itoa(i)+"]", tag)
		}
	}
}

// Page specifies the page of a SearchResponse.
func Page(page int) SearchOption {
	return func(v *url.Values) {
		v.Set("page", strconv.Itoa(page))
	}
}

// SortAscending specifies that the SearchResponse should be sorted in ascending order,
// according to the column specified as by the sort parameter.
func SortAscending(sort string) SearchOption {
	return func(v *url.Values) {
		v.Set("sort", sort)
		v.Set("sortdir", "asc")
	}
}

// SortDescending specifies that the SearchResponse should be sorted in descending order,
// according to the column specified as by the sort parameter.
func SortDescending(sort string) SearchOption {
	return func(v *url.Values) {
		v.Set("sort", sort)
		v.Set("sortdir", "desc")
	}
}

// Query specifies the query to search for, alternatively use the Search
// function.
func Query(query string) SearchOption {
	return func(v *url.Values) {
		v.Set("query", query)
	}
}

// Search searches the umod.org API for plugins matching the given query.
func Search(title string, opts ...SearchOption) (SearchResponse, error) {
	return search(title, opts...)
}

// Latest returns the latest plugins published on umod.org. It is equivalent to
// calling Search("", SortDescending("latest_release_at"), Page(1)).
func Latest() (SearchResponse, error) {
	return search("", SortDescending("latest_release_at"), Page(1))
}

// Oldest returns the oldest plugins published on umod.org. It is equivalent to
// calling Search("", SortAscending("latest_release_at"), Page(1)).
func Oldest() (SearchResponse, error) {
	return search("", SortAscending("latest_release_at"), Page(1))
}

func search(query string, opts ...SearchOption) (SearchResponse, error) {
	var res SearchResponse
	if query != "" {
		opts = append(opts, Query(query))
	}
	err := doRequest(searchURL, &res, opts...)
	res.search = query
	return res, err
}

func doRequest(path string, v interface{}, opts ...SearchOption) error {
	values := &url.Values{}
	for _, opt := range opts {
		opt(values)
	}

	res, err := httpClient.Get(fmt.Sprintf("%s?%s", path, values.Encode()))
	if err != nil {
		return err
	}
	if res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("non ok http status code: %v", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}

	return nil
}
