FROM --platform=$BUILDPLATFORM public.ecr.aws/docker/library/golang:1.27 AS builder

ARG BUILD_DATE
ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src/github.com/Freelance-launchpad/backend-interview/services/payroll/payroll-cronjobs
COPY ./ /go/src/github.com/Freelance-launchpad/backend-interview

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -mod=vendor -o payroll-cronjobs

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app/
COPY --from=builder /go/src/github.com/Freelance-launchpad/backend-interview/services/payroll/payroll-cronjobs/payroll-cronjobs /app/payroll-cronjobs

EXPOSE 8080
USER 1000

ENTRYPOINT ["./payroll-cronjobs", "--config=/etc/payroll-cronjobs/configuration.yaml"]
