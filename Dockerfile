FROM golang:1.21
WORKDIR /build
COPY go.mod go.sum .
RUN go mod download
COPY main.go main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

FROM debian:bullseye
WORKDIR /app
ADD https://github.com/jgraph/drawio-desktop/releases/download/v24.4.0/drawio-amd64-24.4.0.deb drawio.deb
# install dependencies required for:
# - drawio installation (static linking go fuck yourself)
# - running (electron, fuck you)
# - reading tls
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y libgtk-3-0 libnotify4 libnss3 libxss1 libxtst6 xdg-utils libatspi2.0-0 libsecret-1-0 \
        libgbm-dev libasound-dev xvfb xorg gtk2-engines-pixbuf dbus-x11 xfonts-base xfonts-100dpi xfonts-75dpi xfonts-cyrillic xfonts-scalable \
        ca-certificates && \
    apt-get --fix-broken install && \
    dpkg -i drawio.deb
EXPOSE 8080
COPY start.sh .
RUN chmod +x start.sh
COPY --from=0 /build/app .
CMD ["./start.sh"]