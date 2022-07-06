package umod

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL = "https://umod.org/plugins/search.json"
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
// It contains the list of plugins found, the total number of plugins found, the
// total number of pages of plugins found, the URL of the previous page, the URL
// of the next page, and the URL of the last page. The list of plugins is sorted
// by the latest release date.
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
}

// PrevPage returns the previous page of the search response, or an error if
// there is no previous page.
func (s SearchResponse) PrevPage() (SearchResponse, error) {
	if s.PrevPageURL == "" {
		return s, fmt.Errorf("no previous page")
	}
	return doRequest(fmt.Sprintf("%v%v", baseURL, s.PrevPageURL))
}

// NextPage returns the next page of the search response, or an error if
// there is no next page.
func (s SearchResponse) NextPage() (SearchResponse, error) {
	if s.NextPageURL == "" {
		return s, fmt.Errorf("no next page")
	}
	return doRequest(fmt.Sprintf("%v%v", baseURL, s.NextPageURL))
}

// LastPage returns the last page of the search response, or an error if
// there is no last page (although technically that should never happen).
func (s SearchResponse) LastPage() (SearchResponse, error) {
	if s.LastPageURL == "" {
		return s, fmt.Errorf("no last page") // this should never happen ?
	}
	return doRequest(fmt.Sprintf("%v%v", baseURL, s.LastPageURL))
}

// Plugin is the structure of the plugins returned by the umod.org API. It contains
// all the information about the plugin.
type Plugin struct {
	LatestReleaseAtAtom           time.Time `json:"latest_release_at_atom"`
	LatestReleaseAt               string    `json:"latest_release_at"`
	LatestReleaseVersionFormatted string    `json:"latest_release_version_formatted"`
	CategoryTags                  string    `json:"category_tags"`
	Description                   string    `json:"description"`
	CreatedAt                     string    `json:"created_at"`
	Watchers                      int       `json:"watchers"`
	AuthorIconURL                 string    `json:"author_icon_url"`
	Title                         string    `json:"title"`
	Distribution                  string    `json:"distribution"`
	UpdatedAtAtom                 time.Time `json:"updated_at_atom"`
	UpdatedAt                     string    `json:"updated_at"`
	Downloads                     int       `json:"downloads"`
	JSONURL                       string    `json:"json_url"`
	WatchersShortened             string    `json:"watchers_shortened"`
	DonateURL                     string    `json:"donate_url"`
	DownloadURL                   string    `json:"download_url"`
	PublishedAt                   string    `json:"published_at"`
	CreatedAtAtom                 time.Time `json:"created_at_atom"`
	Slug                          string    `json:"slug"`
	IconURL                       string    `json:"icon_url"`
	LatestReleaseVersionChecksum  string    `json:"latest_release_version_checksum"`
	LatestReleaseVersion          string    `json:"latest_release_version"`
	Author                        string    `json:"author"`
	GamesDetail                   []struct {
		IconURL string `json:"icon_url"`
		Name    string `json:"name"`
		URL     string `json:"url"`
		Slug    string `json:"slug"`
	} `json:"games_detail"`
	DownloadsShortened string `json:"downloads_shortened"`
	URL                string `json:"url"`
	StatusDetail       struct {
		Icon    string `json:"icon"`
		Text    string `json:"text"`
		Message string `json:"message"`
		Value   int    `json:"value"`
		Class   string `json:"class"`
	} `json:"status_detail"`
	TagsAll  string `json:"tags_all"`
	Name     string `json:"name"`
	AuthorID string `json:"author_id"`
}

// Search searches the umod.org API for plugins matching the given query.
func Search(title string) (SearchResponse, error) {
	link := fmt.Sprintf("%v?query=%v&sort=latest_release_at&sortdir=desc&page=1", baseURL, url.QueryEscape(title))

	return doRequest(link)
}

// Latest returns the latest plugins published on umod.org.
func Latest() (SearchResponse, error) {
	return doRequest(fmt.Sprintf("%v?sort=latest_release_at&sortdir=desc&page=1", baseURL))
}

func doRequest(url string) (SearchResponse, error) {
	var search SearchResponse
	res, err := httpClient.Get(url)
	if err != nil {
		return search, err
	}
	if res.StatusCode >= http.StatusBadRequest {
		return search, fmt.Errorf("non ok http status code: %v", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&search); err != nil {
		return search, err
	}

	return search, nil
}
