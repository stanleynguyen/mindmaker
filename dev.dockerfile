FROM golang:1.13-alpine
RUN apk update && apk upgrade
RUN apk add --no-cache git jq bash sed && go get github.com/pilu/fresh
RUN set -x \
  # Install ngrok (latest official stable from https://ngrok.com/download).
  && apk add --no-cache curl \
  && curl -Lo /ngrok.zip https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-stable-linux-386.zip \
  && unzip -o /ngrok.zip -d /bin \
  && rm -f /ngrok.zip \
  # Create non-root user.
  && adduser -h /home/ngrok -D -u 6737 ngrok
WORKDIR $GOPATH/src/github.com/stanleynguyen/mindmaker
ENV PORT=8080
ENV GO_ENV=dev
EXPOSE 8080
COPY dev-entrypoint.sh .
CMD ["bash", "./dev-entrypoint.sh"]
