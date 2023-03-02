FROM runnergo/debian:stable-slim

ADD  collector  /data/collector/collector


CMD ["/data/collector/collector","-m", "1"]
