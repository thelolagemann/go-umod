package umod

import (
	"time"
)

// Game represents a game published on umod.org. The slug can be used
// to filter plugins when searching, or alternatively you can call
// the Game.Search method, which is simply a shortcut for
// Search(title, Categories(Category...))
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

// Games returns the list of Game(s) published on umod.org.
//
// See: https://assets.umod.org/games.json
func Games() ([]Game, error) {
	var games []Game
	err := doRequest(gamesURL, &games)
	return games, err
}

// Search is a shortcut for Search(title, Categories(Category(Game.Slug))).
func (g Game) Search(title string, opts ...SearchOption) (SearchResponse, error) {
	opts = append(opts, Categories(Category(g.Slug)))
	return search(title, opts...)
}

// Category is a simple string wrapper that wraps a category slug. It is used
// with the Categories functional option to specify the categories that the
// plugins must be compatible with.
type Category string

const (
	// Universal is the category for plugins that are not specific to any game.
	Universal Category = "universal"
	// SevenDaysToDie is the category for plugins that are compatible with 7 days to die
	SevenDaysToDie Category = "7-days-to-die"
	// Hurtworld is the category for plugins that are compatible with Hurtworld
	Hurtworld Category = "hurtworld"
	// ReignOfKings is the category for plugins that are compatible with Reign of Kings
	ReignOfKings Category = "reign-of-kings"
	// Rust is the category for plugins that are compatible with Rust
	Rust Category = "rust"
	// TheForest is the category for plugins that are compatible with The Forest
	TheForest Category = "the-forest"
)
