FROM golang:1.23
WORKDIR /jury

COPY tests/go.mod tests/go.sum ./

RUN go mod download

ENV MONGODB_URI=$MONGODB_URI

CMD ["go", "run", "."]
