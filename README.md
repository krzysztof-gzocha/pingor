[![Build Status](https://travis-ci.org/krzysztof-gzocha/pingor.svg?branch=master)](https://travis-ci.org/krzysztof-gzocha/pingor)
[![Go Report Card](https://goreportcard.com/badge/github.com/krzysztof-gzocha/pingor)](https://goreportcard.com/report/github.com/krzysztof-gzocha/pingor)
[![codecov](https://codecov.io/gh/krzysztof-gzocha/pingor/branch/master/graph/badge.svg)](https://codecov.io/gh/krzysztof-gzocha/pingor)

# pinGOr
Logs for connection monitoring with Golang.
Run pinGOr and see it's logs to know if your internet connection was interrupted or not.
It's not supporting any database or reporting mechanism yet, but it's architecture is easy to add new features.

# Usage
In order to build the executable simply run:
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
go test ./... -tags unit -cover
```

# Config
```
success_rate_threshold: 0.74  # rate of successfull sub-checks to mark whole check as successfull
success_time_threshold: 5s    # Max average time of sub-checks to mark whole check as successfull
single_check_timeout: 10s     # Timeout for single sub-check
minimal_checking_period: 1m   # Minimal, starting period for periodic checks. Will double after success
maximal_checking_period: 30m  # Maximal period for periodic checks
dns:
  hosts:  # Hosts to resolve in order to confirm connection to DNS is working
    - wp.pl
    - onet.pl
    - google.com
    - upc.pl
    - mbank.pl
ping:
  ips:   # IPs to ping in order to confirm connection is working
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

### Contributing
All ideas and pull requests are welcomed and appreciated :)
If you have any problem with usage don't hesitate to create an issue, we can figure your problem out together.

# Author
Krzysztof Gzocha  
[![](https://img.shields.io/badge/Twitter-%40kgzocha-blue.svg)](https://twitter.com/kgzocha)
