package umod

import "time"

// Plugin is the structure of the plugins returned by the umod.org API. It contains
// all the information about the plugin, including amongst other things, a direct link
// to the plugin's download URL, alongside the checksum of the plugin, the latest release date,
// and the latest release version.
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
