FROM golang

WORKDIR /app

RUN go get github.com/BeHappych/lastPlease/docs

CMD ["lastplease", "-command=./app"]