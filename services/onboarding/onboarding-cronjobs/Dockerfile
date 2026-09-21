FROM --platform=$BUILDPLATFORM public.ecr.aws/docker/library/golang:1.27 AS builder

ARG BUILD_DATE
ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src/github.com/Freelance-launchpad/backend-interview/services/onboarding/onboarding-cronjobs
COPY ./ /go/src/github.com/Freelance-launchpad/backend-interview

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -mod=vendor -o onboarding-cronjobs

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app/
COPY --from=builder /go/src/github.com/Freelance-launchpad/backend-interview/services/onboarding/onboarding-cronjobs/onboarding-cronjobs /app/onboarding-cronjobs

EXPOSE 8080
USER 1000

ENTRYPOINT ["./onboarding-cronjobs", "--config=/etc/onboarding-cronjobs/configuration.yaml"]
