FROM --platform=$BUILDPLATFORM public.ecr.aws/docker/library/golang:1.27 AS builder

ARG BUILD_DATE
ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src/github.com/Freelance-launchpad/backend-interview/services/expenses/expenses-api
COPY ./ /go/src/github.com/Freelance-launchpad/backend-interview

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags "-X 'main.buildDate=${BUILD_DATE}'" -mod=vendor -o expenses-api

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app/
COPY --from=builder /go/src/github.com/Freelance-launchpad/backend-interview/services/expenses/expenses-api/expenses-api /app/expenses-api

EXPOSE 8080
USER 1000

CMD ["./expenses-api", "--config=/etc/expenses-api/configuration.yaml"]
