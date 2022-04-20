FROM golang:1.16-alpine AS build
RUN apk -U upgrade && \
    apk add gcc musl-dev
RUN wget -O aws-iam-authenticator https://amazon-eks.s3.us-west-2.amazonaws.com/1.19.6/2021-01-05/bin/linux/amd64/aws-iam-authenticator && \
    echo "fe958eff955bea1499015b45dc53392a33f737630efd841cd574559cc0f41800  aws-iam-authenticator" | sha256sum -c - && \
    install -o root -g root -m 0755 aws-iam-authenticator /usr/bin/aws-iam-authenticator
WORKDIR /src
RUN go get github.com/cespare/reflex
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN go build -o /app/im-inspector ./cmd/inspect

FROM alpine:3.13
RUN apk --no-cache -U upgrade
COPY --from=build /usr/bin/aws-iam-authenticator /usr/bin/aws-iam-authenticator
WORKDIR /app
COPY --from=build /app/im-inspector .
USER guest
ENTRYPOINT ["/app/im-inspector"]
