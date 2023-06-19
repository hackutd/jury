# STEP 0: Statically build node client
FROM node:lts-hydrogen as client-builder
WORKDIR /client
COPY client/public public
COPY client/src src
COPY ["client/package.json", "client/tailwind.config.js", "client/tsconfig.json", "./"]

ARG NODE_ENV=production
ARG REACT_APP_JURY_NAME
ARG REACT_APP_JURY_URL
ARG REACT_APP_HUB

RUN yarn install
RUN yarn build

# STEP 1: Compile backend
FROM rust:1-alpine3.18 as builder
WORKDIR /usr/src/jury

# This is important, see https://github.com/rust-lang/docker-rust/issues/85
ENV RUSTFLAGS="-C target-feature=-crt-static"

RUN apk add --no-cache musl-dev

# Jank way to cache built rust dependencies
RUN mkdir src
RUN echo "fn main() {}" > ./src/main.rs
COPY ["Cargo.toml", "./"]
RUN cargo build --release

# Then actually copy over the app and build it
COPY src src
COPY --from=client-builder /client/build public
COPY Rocket.toml .

ARG MONGODB_URI=$MONGODB_URI
ARG JURY_ADMIN_PASSWORD=$JURY_ADMIN_PASSWORD
ARG FILESERVER="/public"

RUN cargo install --locked --path .
RUN ls -la /usr/local/cargo/bin

# STEP 2: Main running container
FROM alpine:3.18

# Extra dependencies needed
RUN apk add --no-cache ca-certificates openssl-dev bash openssl libgcc libstdc++

ENV FILESERVER="/public"
ENV ROCKET_ADDRESS=0.0.0.0

EXPOSE $PORT

COPY --from=builder /usr/local/cargo/bin/jury /usr/local/bin/jury
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=client-builder /client/build /public

CMD ["jury"]
