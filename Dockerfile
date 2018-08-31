FROM golang:alpine

WORKDIR /go/src/app
COPY . .

RUN go build main.go

RUN ls -la .

RUN chmod +x main


RUN apk add --no-cache curl

EXPOSE 8080

ENV SEWOBEUSER restuser

ENV SEWOBEPASSWORD 12345secret

ENV SEWOBEURL https://server30.der-moderne-verein.de/restservice/standard/v1.0/auswertungen/get_auswertung_data.php


COPY start.sh /go/src/app/start.sh
COPY query_sewobe.sh /go/src/app/query_sewobe.sh
COPY crontab /crontab
RUN chmod 755 /go/src/app/start.sh /go/src/app/query_sewobe.sh
RUN /usr/bin/crontab /crontab
RUN /usr/sbin/crond -f -l 8 &

WORKDIR /go/src/app/testdata
WORKDIR /go/src/app


CMD ["./start.sh"]
