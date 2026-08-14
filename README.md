# GitInspector Go

A fast, modular Go port of the original [gitinspector](https://github.com/ejwa/gitinspector). It reads the Git object database directly to analyze repository history, intelligently filtering out noise (like comments and empty lines) to calculate real code impact.

## Features
* **Language-Aware Filtering:** Plugin system automatically ignores irrelevant lines for supported languages.
* **Rich Output:** Renders a colorized terminal table with impact progress bars.
* **Standalone:** Uses `go-git` to read the Git object database directly—no local Git binary required.

## Building

    go build -o gitinspector ./cmd/gitinspector/main.go

## Runnig 

    go run ./cmd/gitinspector/main.go <path/to/your/git/repo>

## Customizing Language Rules

Add custom language rules (e.g., filtering `#` in Python) by implementing the `LanguagePlugin` interface in `internal/plugins/` and registering it in `plugin.go`.