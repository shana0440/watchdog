# WatchDog

[![cover.run](https://cover.run/go/https:/github.com/shana0440/watchdog.svg?style=flat&tag=golang-1.10)](https://cover.run/go?tag=golang-1.10&repo=https%3A%2Fgithub.com%2Fshana0440%2Fwatchdog)

execute script when file change

[![asciicast](https://asciinema.org/a/3vsKIrda45uXwrBXsYluxMgi0.png)](https://asciinema.org/a/3vsKIrda45uXwrBXsYluxMgi0)

Useage example:

```bash
watchdog -c "go test ./..." -ignore "*.swp" -ignore "vendor"
```

## Installation

```bash
go get github.com/shana0440/watchdog
```
