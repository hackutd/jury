# STEP 0: Statically build node client
FROM node as client-builder
WORKDIR /client
COPY client/public public
COPY client/src src
COPY ["client/package.json", "client/yarn.lock", "client/tailwind.config.js", "client/tsconfig.json", "./"]

ARG NODE_ENV=production
ARG REACT_APP_JURY_NAME
ARG REACT_APP_JURY_URL

RUN yarn install --frozen-lockfile
RUN yarn build

# STEP 1: Compile backend
FROM rust as builder
WORKDIR /usr/src/jury

RUN apt update
RUN apt install libgsl-dev -y

# Jank way to cache built rust dependencies
RUN mkdir src
RUN echo "fn main() {}" > ./src/main.rs
COPY ["Cargo.toml", "Cargo.lock",  "./"]
RUN cargo build --locked --release

# Then actually copy over the app and build it
COPY src src
COPY --from=client-builder /client/build public
COPY Rocket.toml .

ARG MONGODB_URI=$MONGODB_URI
ARG JURY_ADMIN_PASSWORD=$JURY_ADMIN_PASSWORD
ARG FILESERVER="/public"

RUN cargo install --locked --path .

# STEP 2: Main running container
FROM debian:bullseye-slim

ENV FILESERVER="/public"
ENV ROCKET_ADDRESS=0.0.0.0

EXPOSE $PORT

RUN apt-get update && apt-get install -y libgsl-dev && rm -rf /var/lib/apt/lists/*
COPY --from=builder /usr/local/cargo/bin/jury /usr/local/bin/jury
COPY --from=client-builder /client/build /public
CMD ["jury"]
