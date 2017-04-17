FROM scratch

ADD netcan /netcan

ADD ui /ui

ENTRYPOINT ["/netcan"]