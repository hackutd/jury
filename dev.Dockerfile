FROM rust
WORKDIR /jury

RUN apt update
RUN apt install libgsl-dev -y
RUN cargo install cargo-watch

COPY public /public

ENV MONGODB_URI=$MONGODB_URI
ENV JURY_ADMIN_PASSWORD=$JURY_ADMIN_PASSWORD
ENV FILESERVER="/public"

CMD ["cargo", "watch", "-w", "src", "-w", "Cargo.toml", "-w", "Rocket.toml", "-x", "run"]

# docker build -f dev.Dockerfile -t jury-dev:latest .
