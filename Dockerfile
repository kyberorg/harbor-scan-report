FROM kio.ee/base/go:1.18 as builder
WORKDIR /go/src/app
COPY cmd/main.go cmd/main.go
RUN  GO111MODULE=off CGO_ENABLED=0 go install ./...
RUN chmod +x /go/bin/cmd

FROM kio.ee/base/abi:edge as runner
COPY --from=builder --chown=appuser:appgroup /go/bin/cmd /hsr
COPY  --chown=appuser:appgroup entrypoint.sh /entrypoint.sh
RUN chmod u+x /entrypoint.sh
USER appuser
ENTRYPOINT ["/entrypoint.sh"]