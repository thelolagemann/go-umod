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

// Game represents a game published on umod.org.
type Game struct {
	Name                          string      `json:"name"`
	Slug                          string      `json:"slug"`
	Description                   string      `json:"description"`
	Aliases                       string      `json:"aliases"`
	GameURL                       string      `json:"game_url"`
	SnapshotURL                   string      `json:"snapshot_url"`
	IconURL                       string      `json:"icon_url"`
	Repository                    string      `json:"repository"`
	ServerAppid                   string      `json:"server_appid"`
	ClientAppid                   string      `json:"client_appid"`
	Buildable                     int         `json:"buildable"`
	UmodBuildable                 int         `json:"umod_buildable"`
	InstallationPaths             string      `json:"installation_paths"`
	TargetFramework               string      `json:"target_framework"`
	TargetSdk                     string      `json:"target_sdk"`
	PublicBranchName              string      `json:"public_branch_name"`
	PublicBranchDescription       interface{} `json:"public_branch_description"`
	PreprocessorSymbol            string      `json:"preprocessor_symbol"`
	SteamAuthenticated            int         `json:"steam_authenticated"`
	FilesInstall                  interface{} `json:"files_install"`
	FilesUpdate                   interface{} `json:"files_update"`
	SkipInstall                   string      `json:"skip_install"`
	SkipUpdate                    interface{} `json:"skip_update"`
	Whitelist                     string      `json:"whitelist"`
	Blacklist                     string      `json:"blacklist"`
	UpdateCheckFrequency          string      `json:"update_check_frequency"`
	DownloadURL                   string      `json:"download_url"`
	URL                           string      `json:"url"`
	PluginCount                   int         `json:"plugin_count"`
	ExtensionCount                int         `json:"extension_count"`
	ProductCount                  int         `json:"product_count"`
	LatestReleaseVersion          string      `json:"latest_release_version"`
	LatestReleaseVersionFormatted string      `json:"latest_release_version_formatted"`
	LatestReleaseVersionChecksum  string      `json:"latest_release_version_checksum"`
	LatestReleaseAt               string      `json:"latest_release_at"`
	LatestReleaseAtAtom           time.Time   `json:"latest_release_at_atom"`
	Watchers                      int         `json:"watchers"`
	WatchersShortened             string      `json:"watchers_shortened"`
	Channels                      []struct {
		ChannelID string `json:"channel_id"`
		BotName   string `json:"bot_name"`
		BotSlug   string `json:"bot_slug"`
	} `json:"channels"`
	SteamBranches []struct {
		Name        string `json:"name"`
		Pwdrequired int    `json:"pwdrequired"`
		Timeupdated string `json:"timeupdated"`
		Buildid     int    `json:"buildid"`
	} `json:"steam_branches"`
}

// Games returns the list of games published on umod.org.
//
// See: https://assets.umod.org/games.json
func Games() ([]Game, error) {
	var games []Game
	res, err := httpClient.Get(gamesURL)
	if err != nil {
		return games, err
	}
	if res.StatusCode >= http.StatusBadRequest {
		return games, fmt.Errorf("non ok http status code: %v", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&games); err != nil {
		return games, err
	}

	return games, nil
}

func (g Game) Search(title string) (SearchResponse, error) {
	// link := fmt.Sprintf("%v&sort=latest_release_at&sortdir=desc&page=1&game=%v", url.QueryEscape(title), g.Slug)

	return search(title)
}

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

// SearchOption is the type used to specify the page of the search response.
type SearchOption func(*url.Values)

func Categories(categories ...Category) SearchOption {
	return func(v *url.Values) {
		for i, g := range categories {
			v.Set("categories"+genIndex(i), string(g))
		}
	}
}

type Category string

const (
	CategoryUniversal    Category = "universal"
	Category7DaysToDie   Category = "7-days-to-die"
	CategoryHurtworld    Category = "hurtworld"
	CategoryReignOfKings Category = "reign-of-kings"
	CategoryRust         Category = "rust"
	CategoryTheForest    Category = "the-forest"
)

func genIndex(i int) string {
	return url.QueryEscape(fmt.Sprintf("[%v]", i))
}

func Tags(tags ...string) SearchOption {
	return func(v *url.Values) {
		for i, tag := range tags {
			v.Set("tags"+genIndex(i), tag)
		}

	}
}

// Page specifies the page of the search response.
func Page(page int) SearchOption {
	return func(v *url.Values) {
		v.Set("page", strconv.Itoa(page))
	}
}

// SortAscending specifies that the search response should be sorted in ascending order,
// according to the column specified as by the sort parameter.
func SortAscending(sort string) SearchOption {
	return func(v *url.Values) {
		v.Set("sort", sort)
		v.Set("sortdir", "asc")
	}
}

// SortDescending specifies that the search response should be sorted in descending order,
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

// Latest returns the latest plugins published on umod.org.
func Latest() (SearchResponse, error) {
	return search("", SortDescending("latest_release_at"), Page(1))
}

// Oldest returns the oldest plugins published on umod.org.
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

	dec, err := url.QueryUnescape(values.Encode())
	if err != nil {
		return err
	}
	res, err := httpClient.Get(fmt.Sprintf("%s?%s", path, dec))
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
