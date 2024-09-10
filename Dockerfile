FROM golang:1.22.4-alpine3.20 AS build

# https://github.com/kubernetes-sigs/aws-iam-authenticator/releases/
ARG AWS_IAM_AUTHENTICATOR_VERSION=0.6.11
ARG AWS_IAM_AUTHENTICATOR_CHECKSUM=8593d0c5125f8fba4589008116adf12519cdafa56e1bfa6b11a277e2886fc3c8

# https://github.com/cespare/reflex/releases
ARG REFLEX_VERSION=v0.3.1

RUN apk add gcc musl-dev git && \
\
    wget -O aws-iam-authenticator https://github.com/kubernetes-sigs/aws-iam-authenticator/releases/download/v${AWS_IAM_AUTHENTICATOR_VERSION}/aws-iam-authenticator_${AWS_IAM_AUTHENTICATOR_VERSION}_linux_amd64 && \
    echo "${AWS_IAM_AUTHENTICATOR_CHECKSUM}  aws-iam-authenticator" | sha256sum -c - && \
    install -o root -g root -m 0755 aws-iam-authenticator /usr/bin/aws-iam-authenticator
WORKDIR /src
RUN go install github.com/cespare/reflex@${REFLEX_VERSION}
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN go build -o /app/im-inspector ./cmd/inspect

FROM alpine:3.20
RUN apk --no-cache -U upgrade
COPY --from=build /usr/bin/aws-iam-authenticator /usr/bin/aws-iam-authenticator
WORKDIR /app
COPY --from=build /app/im-inspector .
USER guest
ENTRYPOINT ["/app/im-inspector"]
