# pingor
Create simple logs for connection monitoring with Golang

# Usage
```
go build && ./pingor -debug -config config.yaml
```

# Config
```
success_rate_threshold: 0.74
success_time_threshold: 20s
single_check_timeout: 10s
minimal_checking_period: 1m
maximal_checking_period: 30m
dns:
  hosts:
    - wp.pl
    - onet.pl
    - google.com
    - upc.pl
    - mbank.pl
ping:
  ips:
    - 8.8.8.8
    - 8.8.4.4
    - 9.9.9.9
    - 1.1.1.1
```

# Example systemd config
File `/etc/systemd/system/pingor.service`
```
[Unit]
Description=pinGOr: checking internet connectivity
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
