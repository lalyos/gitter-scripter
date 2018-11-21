FROM golang:1.11-alpine
RUN apk add -U git musl-dev gcc
ADD . /app
WORKDIR /app
RUN go build -v -o gitter-scripter-linux .

FROM alpine:3.8

RUN apk add -U curl bash jq
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl \
   && chmod +x ./kubectl \
   && mv ./kubectl /usr/local/bin/kubectl

COPY --from=0 /app/gitter-scripter-linux /gitter-scripter
COPY ./getsession.sh /
ENV PORT 8080

ENTRYPOINT ["/gitter-scripter"]
