FROM node as client-builder
WORKDIR /client
COPY client/public public
COPY client/src src
COPY ["client/package.json", "client/yarn.lock", "client/tailwind.config.js", "client/tsconfig.json", ".env", "./"]

ENV NODE_ENV=production
ENV REACT_APP_JURY_NAME=$JURY_NAME
ENV REACT_APP_JURY_URL=$JURY_URL

RUN yarn install --frozen-lockfile
RUN yarn build

FROM rust as builder
WORKDIR /usr/src/jury

RUN apt update
RUN apt install libgsl-dev -y
RUN cargo install cargo-watch

COPY src src
COPY --from=client-builder /client/build public
COPY *.toml .

ARG MONGODB_URI=$MONGODB_URI
ARG JURY_ADMIN_PASSWORD=$JURY_ADMIN_PASSWORD
ARG FILESERVER="/public"

RUN cargo install --path .

FROM debian:bullseye-slim

ENV FILESERVER="/public"

ENV ROCKET_ADDRESS=0.0.0.0
EXPOSE $PORT

RUN apt-get update && apt-get install -y libgsl-dev && rm -rf /var/lib/apt/lists/*
COPY --from=builder /usr/local/cargo/bin/jury /usr/local/bin/jury
COPY --from=client-builder /client/build /public
CMD ["jury"]
