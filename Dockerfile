FROM kio.ee/base/go:1.19 as builder
WORKDIR /go/src/app
COPY cmd cmd
COPY Makefile Makefile
COPY go.mod go.mod
RUN  make small-binary
RUN chmod +x bin/hsr

FROM kio.ee/base/abi:edge as runner
COPY --from=builder --chown=appuser:appgroup /go/src/app/bin/hsr /hsr
COPY  --chown=appuser:appgroup entrypoint.sh /entrypoint.sh
RUN chmod u+x /entrypoint.sh
USER appuser
ENTRYPOINT ["/entrypoint.sh"]