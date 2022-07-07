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
```go
// Search for a plugin
var names []string
results, _ := umod.Search("heli")
for _, result := range results.Data {
	names = append(names, result.Name)
}
fmt.Println(strings.Join(names, ", "))

// Output: HeliScrap, HeliSams, HelicopterHover, PersonalHeli, NoHeliFire, HelicopterProtection, HeliControl, HeliEditor, NudistHeli, NoHeliFlyhack

// Get the most recently updated plugins
var names []string
results, _ := umod.Latest()
for _, result := range results.Data {
	names = append(names, result.Name)
}
fmt.Println(strings.Join(names, ", "))

// Output: Kits, OptimalBurn, NightZombies, NeverWear, MountComputerStation, SafeRecycler, ComputersPlus, CarRadio, DiscordLinker, XPerience
```

### TODO 
- [ ] Implement games endpoint
- [ ] Add search parameters (game, author, etc.)
- [ ] Pagination Examples