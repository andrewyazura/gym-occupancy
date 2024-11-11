# Well Fitness parser

Parse Well Fitness client portal to get number of people at the given time in selected clubs.
Data is written to InfluxDB upon running the `main.go` script. It's intended to be scheduled
using `cron` or other schedulers.

![image](https://github.com/user-attachments/assets/66396d5c-00a4-423f-b727-74837b4c7f9e)
