# go-umod

![GitHub tag (latest SemVer pre-release)](https://img.shields.io/github/v/tag/thelolagemann/go-umod?include_prereleases&label=release&sort=semver&style=for-the-badge)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://pkg.go.dev/github.com/thelolagemann/go-umod)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/thelolagemann/go-umod?style=for-the-badge)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/thelolagemann/go-umod/Test?label=tests&style=for-the-badge)
![GitHub Workflow Status](https://img.shields.io/github/workflow/status/thelolagemann/go-umod/CodeQL?label=CodeQL&style=for-the-badge)
[![Go ReportCard](https://goreportcard.com/badge/github.com/thelolagemann/go-umod?style=for-the-badge)](https://goreportcard.com/report/thelolagemann/go-umod)

> A small golang library to access the [uMod API](https://umod.org/plugins/search.json).

## Getting started

### Installation
```shell
go get -u github.com/thelolagemann/go-umod
```

### Usage

#### Plugin Searching
```go
// Search for a plugin
var names []string
results, _ := umod.Search("heli")
for _, result := range results.Plugins {
	names = append(names, result.Name)
}
fmt.Println(strings.Join(names, ", "))

// Output: HeliScrap, HeliSams, HelicopterHover, PersonalHeli, NoHeliFire, HelicopterProtection, HeliControl, HeliEditor, NudistHeli, NoHeliFlyhack

// Get the most recently updated plugins
var names []string
results, _ := umod.Latest()
for _, result := range results.Plugins {
	names = append(names, result.Name)
}
fmt.Println(strings.Join(names, ", "))

// Output: Kits, OptimalBurn, NightZombies, NeverWear, MountComputerStation, SafeRecycler, ComputersPlus, CarRadio, DiscordLinker, XPerience

// Search with parameters
var names []string
results, _ := umod.Search("heli", umod.Tags("fun", "mechanics"), umod.SortDescending("latest_release_at"))
for _, result := range results.Plugins {
	names = append(names, result.Name)
}
fmt.Println(strings.Join(names, ", "))

// Output: AdvancedArrows, SuicideBomber

// Search for plugins specific to a game
var names []string
results, _ := umod.Search("", umod.Categories(umod.SevenDaysToDie))
for _, result := range results.Plugins {
	names = append(names, result.Name)
}
fmt.Println(strings.Join(names, ", "))

// Output: AutomaticPluginUpdater, DegreeTags, CraftingStore, DaySeven, DiscordLinker, Give, ImgurApi, Punish, SDTeleportation, SleeperGroup
```

#### Games

```go
// Get all games
var names []string
games, _ := umod.Games()
for _, game := range games {
	names = append(names, game.Name)
}
fmt.Println(strings.Join(names, ", "))

// Output: Rust, Hurtworld, 7 Days To Die, Reign Of Kings, The Forest, Heat

// You can also search for plugins specific to a game
// using the Game.Search method, although this is simply
// a shortcut for calling umod.Search with the game's category slug.
var names []string
results, _ := games[0].Search("kits")
for _, result := range results.Plugins {
	names = append(names, result.Name)
}
fmt.Println(strings.Join(names, ", "))

// Output:  CustomAutoKits, WipeKits, Kits, Murderers, Factions, EasyTeams, FactionsCore, CustomSets, Battlefield, Loadoutless
```

### TODO 
- [x] Implement games endpoint
- [ ] Add search parameters (game, author, etc.)
- [ ] Pagination Examples