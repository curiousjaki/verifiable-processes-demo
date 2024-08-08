
# To build run: docker build -f Dockerfile.release --build-arg="RISC0_TOOLCHAIN_VERSION=r0.1.79.0" -t proving-server .
FROM rust:1.80 AS builder

ARG APP_NAME=proving-server
ARG RISC0_TOOLCHAIN_VERSION=r0.1.79.0

RUN apt-get update
RUN apt-get install protobuf-compiler -y

RUN cargo install cargo-binstall
RUN cargo binstall -y --force cargo-risczero
RUN cargo risczero install --version ${RISC0_TOOLCHAIN_VERSION}

COPY ./host/ /app/host
COPY ./methods/ /app/methods
COPY ./Cargo.toml /app/Cargo.toml
#COPY ./target /app/target
WORKDIR /app

RUN cargo fetch
RUN cargo build --release
RUN cp ./target/release/proving-server /bin/server

FROM rust:1.80 AS final

# Create a non-privileged user that the app will run under.
# See https://docs.docker.com/go/dockerfile-user-best-practices/
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

# Copy the executable from the "build" stage.
COPY --from=builder /bin/server /bin/server

# Expose the port that the application listens on.
EXPOSE 50051

# What the container should run when it is started.
CMD ["/bin/server"]
