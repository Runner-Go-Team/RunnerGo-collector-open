FROM runnergo/debian:stable-slim

ADD  collector  /data/collector/collector

#ADD wait-for-it.sh  /bin/

CMD ["/data/collector/collector","-m", "1"]
#ENTRYPOINT  ["/bin/entrypoint.sh"]
