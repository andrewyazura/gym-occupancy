# Well Fitness parser

Parse Well Fitness client portal to get number of people at the given time in selected clubs.
Data is written to InfluxDB upon running the `main.go` script. It's intended to be scheduled
using `cron` or other schedulers.
