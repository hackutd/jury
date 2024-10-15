# STEP 0: Statically build node client
FROM node:lts-hydrogen AS client-builder
WORKDIR /client
COPY client ./
COPY ["client/package.json", "client/tailwind.config.js", "client/tsconfig.json", "./"]

ARG VITE_JURY_NAME
ARG VITE_HUB
ARG VITE_JURY_URL

RUN yarn install

ARG NODE_ENV=production
RUN yarn build

# STEP 1: Compile backend
FROM golang:1.22 AS builder
WORKDIR /usr/src/jury

# Copy over the app
COPY server ./
RUN rm -rf public
COPY --from=client-builder /client/build public

# Install dependencies
RUN go mod download

ARG MONGODB_URI=$MONGODB_URI
ARG JURY_ADMIN_PASSWORD=$JURY_ADMIN_PASSWORD
ARG EMAIL_HOST=$EMAIL_HOST
ARG EMAIL_PORT=$EMAIL_PORT
ARG EMAIL_FROM=$EMAIL_FROM
ARG EMAIL_FROM_NAME=$EMAIL_FROM_NAME
ARG EMAIL_USERNAME=$EMAIL_USERNAME
ARG EMAIL_PASSWORD=$EMAIL_PASSWORD
ARG SENDGRID_API_KEY=$SENDGRID_API_KEY

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/jury

# STEP 2: Main running container
FROM scratch

# Extra dependencies needed
# RUN apk add --no-cache ca-certificates openssl-dev bash openssl libgcc libstdc++

EXPOSE $PORT

COPY --from=builder /go/bin/jury .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=client-builder /client/build /public
COPY ./server/email.html /email.html

ENTRYPOINT [ "./jury" ]
