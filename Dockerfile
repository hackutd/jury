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

FROM rust as builder
WORKDIR /usr/src/jury

RUN apt update
RUN apt install libgsl-dev -y

COPY src src
COPY --from=client-builder /client/build public
COPY *.toml .
COPY Cargo.lock .

ARG MONGODB_URI=$MONGODB_URI
ARG JURY_ADMIN_PASSWORD=$JURY_ADMIN_PASSWORD
ARG FILESERVER="/public"

RUN cargo install --locked --path .

FROM debian:bullseye-slim

ENV FILESERVER="/public"
ENV ROCKET_ADDRESS=0.0.0.0

EXPOSE $PORT

RUN apt-get update && apt-get install -y libgsl-dev && rm -rf /var/lib/apt/lists/*
COPY --from=builder /usr/local/cargo/bin/jury /usr/local/bin/jury
COPY --from=client-builder /client/build /public
CMD ["jury"]
