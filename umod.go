package umod

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var (
	ParseError = errors.New("unable to parse")
)

const (
	baseURL = "https://umod.org/plugins/search.json"
)

type SearchResponse struct {
	CurrentPage  int       `json:"current_page"`
	Data         []*Plugin `json:"data"`
	FirstPageURL string    `json:"first_page_url"`
	From         int       `json:"from"`
	LastPage     int       `json:"last_page"`
	LastPageURL  string    `json:"last_page_url"`
	NextPageURL  string    `json:"next_page_url"`
	Path         string    `json:"path"`
	PerPage      int       `json:"per_page"`
	PrevPageURL  string    `json:"prev_page_url"`
	To           int       `json:"to"`
	Total        int       `json:"total"`
}

func (s *SearchResponse) PrevPage() (*SearchResponse, error) {
	if s.PrevPageURL == "" {
		return nil, fmt.Errorf("no previous page")
	}
	return doRequest(s.PrevPageURL)
}

func (s *SearchResponse) NextPage() (*SearchResponse, error) {
	if s.NextPageURL == "" {
		return nil, fmt.Errorf("no next page")
	}
	return doRequest(s.NextPageURL)
}

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

	location string
	version  string
}

func (p *Plugin) String() string {
	return p.Title
}

func Search(title string) (*SearchResponse, error) {
	link := fmt.Sprintf("%v?query=%v&page=1&sort=latest_release_at&sortdir=desc", baseURL, url.QueryEscape(title))

	return doRequest(link)
}

func Latest() (*SearchResponse, error) {
	return doRequest(fmt.Sprintf("%v?query=&page=1&sort=latest_release_at&sortdir=desc", baseURL))
}

func doRequest(url string) (*SearchResponse, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("non ok http status code: %v", res.StatusCode)
	}

	var search *SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&search); err != nil {
		return nil, err
	}

	return search, nil
}
