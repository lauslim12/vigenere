# Changelog

Changelog is used to keep track of version changes. The versioning scheme used is [SemVer](https://semver.org/). First integer is used for breaking change, second integer is used for major patches, and third integer is used for minor bug fixes.

## Version 1.0.3 (09/05/2022)

- Improve CI by introducing three new checks: `go mod verify`, `go vet ./...`, and `go build -v ./...`.
- Additional check for race conditions in unit-tests: `go test -race -v -cover ./... ./...`.
- Create new CI job specific for paralellized linting purposes with [`golangci-lint`](https://github.com/golangci/golangci-lint).

## Version 1.0.2 (04/05/2022)

- Fix formatting of the algorithm steps in Go doc.

## Version 1.0.1 (03/05/2022)

- Elaborate documentation in both repository and Go doc.
- Upgrade GitHub Action: `actions/github-script@v5` -> `actions/github-script@v6`.

## Version 1.0.0 (02/05/2022)

- Official initial release of the library.
