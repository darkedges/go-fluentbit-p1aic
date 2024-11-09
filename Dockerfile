FROM golang:1.23 AS gobuilder

WORKDIR /root

ENV GOOS=linux\
    GOARCH=amd64

COPY / /root/

ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go mod download & make all

FROM fluent/fluent-bit:3.1.10

COPY --from=gobuilder /root/in_p1aic.so /fluent-bit/bin/
COPY --from=gobuilder /root/fluent-bit.conf /fluent-bit/etc/fluent-bit.conf 
COPY --from=gobuilder /root/plugins.conf /fluent-bit/etc/

EXPOSE 2020
VOLUME /fluentbit

CMD ["/fluent-bit/bin/fluent-bit", "--config", "/fluent-bit/etc/fluent-bit.conf"]