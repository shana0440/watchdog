# WatchDog

[![Build Status](https://travis-ci.com/shana0440/watchdog.svg?branch=master)](https://travis-ci.com/shana0440/watchdog)
[![codecov](https://codecov.io/gh/shana0440/watchdog/branch/master/graph/badge.svg)](https://codecov.io/gh/shana0440/watchdog)

execute script when file change

[![asciicast](https://asciinema.org/a/3vsKIrda45uXwrBXsYluxMgi0.png)](https://asciinema.org/a/3vsKIrda45uXwrBXsYluxMgi0)

Useage example:

```bash
watchdog -c "go test ./..." -i "*.swp" -i "vendor"
```

## Installation

```bash
go get github.com/shana0440/watchdog
```
