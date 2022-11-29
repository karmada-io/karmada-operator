FROM golang:1.19.3 as builder

WORKDIR /go/src/github.com/karmada-io/karmada-operator

COPY . /go/src/github.com/karmada-io/karmada-operator

RUN cd /go/src/github.com/karmada-io/karmada-operator && \
    go build -o /bin/karmada-operator cmd/operator/operator.go


# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/base:nonroot

COPY --from=builder /bin/karmada-operator  /bin/karmada-operator

USER 65532:65532
