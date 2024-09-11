FROM golang:1.22
WORKDIR /jury

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

COPY server/go.mod server/go.sum ./

RUN go mod download

ENV MONGODB_URI=$MONGODB_URI
ENV JURY_ADMIN_PASSWORD=$JURY_ADMIN_PASSWORD

CMD ["air"]
