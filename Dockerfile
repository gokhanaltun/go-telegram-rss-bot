FROM golang:1.22.3
WORKDIR /opt/go_rss_bot
COPY . .
ENV TOKEN=""
ENV ID=""
RUN go build -o goRssBot
CMD [ "./goRssBot" ]