FROM golang:1.6
ENV DIST /go/src/github.com/reyoung/PaddleStat
ENV POSTGRES_DBNAME paddle_stat
ENV POSTGRES_DBUSER user
ENV POSTGRES_DBPASSWD passwd
ENV PADDLE_VERSION 0.9.0a0
ENV COOKIE_KEY SOMEKEY

COPY . $DIST
RUN cd $DIST && go get ./... && go get .
EXPOSE 3000
ENTRYPOINT ["PaddleStat"]
CMD ["-dbconnect=${POSTGRES_DBUSER}:${POSTGRES_DBPASSWD}@localhost/${POSTGRES_DBNAME}?sslmode=disable", \
     "-version=${PADDLE_VERSION}", \
     "-key=${COOKIE_KEY}"]
