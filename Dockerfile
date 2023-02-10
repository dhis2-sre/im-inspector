FROM golang:1.18-alpine3.17 AS build
ARG REFLEX_VERSION=v0.3.1
RUN apk add gcc musl-dev git && \
\
    wget -O aws-iam-authenticator https://amazon-eks.s3.us-west-2.amazonaws.com/1.19.6/2021-01-05/bin/linux/amd64/aws-iam-authenticator && \
    echo "fe958eff955bea1499015b45dc53392a33f737630efd841cd574559cc0f41800  aws-iam-authenticator" | sha256sum -c - && \
    install -o root -g root -m 0755 aws-iam-authenticator /usr/bin/aws-iam-authenticator
WORKDIR /src
RUN go install github.com/cespare/reflex@${REFLEX_VERSION}
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN go build -o /app/im-inspector ./cmd/inspect

FROM alpine:3.17
RUN apk --no-cache -U upgrade
COPY --from=build /usr/bin/aws-iam-authenticator /usr/bin/aws-iam-authenticator
WORKDIR /app
COPY --from=build /app/im-inspector .
USER guest
ENTRYPOINT ["/app/im-inspector"]
