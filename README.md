[![Build Status](https://travis-ci.org/krzysztof-gzocha/pingor.svg?branch=master)](https://travis-ci.org/krzysztof-gzocha/pingor)
[![Go Report Card](https://goreportcard.com/badge/github.com/krzysztof-gzocha/pingor)](https://goreportcard.com/report/github.com/krzysztof-gzocha/pingor)
[![codecov](https://codecov.io/gh/krzysztof-gzocha/pingor/branch/master/graph/badge.svg)](https://codecov.io/gh/krzysztof-gzocha/pingor)

# pinGOr
Logs for connection monitoring with Golang.
Run pinGOr and see it's logs to know if your internet connection was interrupted or not.
It's not supporting any database or reporting mechanism yet, but it's architecture is easy to add new features.

# How?
PinGOr will read provided config and try to:
- resolve provided host names to IPs,
- run ping command on the host to specified IPs.

It will start checking the connection after configured minimal checking period and if the connection will be ok the period will be doubled.
When connection checks will drop below configured success rate threshold, then the connection will be marked as "dropped" and proper log will be created.

# Why?
I have signed SLA agreement with my ISP, but didn't have any tool to actually know if there was any connection-related problem, while the PC was running and I was away.

# Usage
In order to build the executable run:
```
curl https://glide.sh/get | sh
glide install
go build
```
In order to run the executable:
```
./pingor -debug -config config.yaml
```
In order to test it:
```
go test ./... -tags unit
```

# Config
```
success_rate_threshold: 0.74  # rate of successfull sub-checks to mark whole check as successfull
success_time_threshold: 5s    # Max average time of sub-checks to mark whole check as successfull
single_check_timeout: 10s     # Timeout for single sub-check
minimal_checking_period: 1m   # Minimal, starting period for periodic checks. Will double after success
maximal_checking_period: 30m  # Maximal period for periodic checks
www:
  urls:     # URLs to check if HTTP status is 200 OK
    - https://wp.pl
    - https://www.onet.pl
    - https://www.google.com
    - https://www.upc.pl
dns:
  hosts:  # Hosts to resolve in order to confirm connection to DNS is working. Leave empty to skip DNS checks
    - wp.pl
    - onet.pl
    - google.com
    - upc.pl
    - mbank.pl
ping:
  ips:   # IPs to ping in order to confirm connection is working. Leave empty to skip ping checks
    - 8.8.8.8
    - 8.8.4.4
    - 9.9.9.9
    - 1.1.1.1
```

# Recommended usage
Recommended usage is by adding pinGOr to [systemd](https://www.tecmint.com/create-new-service-units-in-systemd/).
Below you can see example `pingor.service` file, that can be helpful
```
[Unit]
Description=pinGOr: logging internet connectivity
After=network.target

[Service]
Type=simple
User=somebody # EDIT THIS
WorkingDirectory=/full/path/to/pingor/dir # EDIT THIS
ExecStart=/full/path/to/pingor/exec # EDIT THIS
Restart=on-abort

[Install]
WantedBy=multi-user.target
```
After configuring it you can inspect the logs to check for connection disruption

# Known issues
- Was not tested on Windows

# Contributing
All ideas and pull requests are welcomed and appreciated.
If you have any problem with usage don't hesitate to create an issue, we can figure out your problem together.

# Author
Krzysztof Gzocha
[![](https://img.shields.io/badge/Twitter-%40kgzocha-blue.svg)](https://twitter.com/kgzocha)
