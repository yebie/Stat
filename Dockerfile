FROM golang:1.6
ENV DIST /go/src/github.com/reyoung/PaddleStat
ENV POSTGRES_DB paddle_stat
ENV POSTGRES_USER user
ENV POSTGRES_PASSWORD passwd
ENV POSTGRES_HOST localhost
ENV PADDLE_VERSION 0.9.0a0
ENV COOKIE_KEY SOMEKEY

COPY . $DIST
RUN cd $DIST && go get ./... && go get .
EXPOSE 3000
ENTRYPOINT $DIST/run.sh
