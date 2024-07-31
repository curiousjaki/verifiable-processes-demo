
# To build run: docker build -f Dockerfile.release --build-arg="RISC0_TOOLCHAIN_VERSION=r0.1.79.0" -t risczero/risc0-guest-builder:r0.1.79.0 .
FROM ubuntu:20.04@sha256:3246518d9735254519e1b2ff35f95686e4a5011c90c85344c1f38df7bae9dd37 as builder

ARG RISC0_TOOLCHAIN_VERSION=r0.1.79.0

RUN apt-get update
RUN apt-get install -y --no-install-recommends ca-certificates clang curl libssl-dev pkg-config
RUN curl --proto '=https' --tlsv1.2 --retry 10 --retry-connrefused -fsSL 'https://sh.rustup.rs' | sh -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"
RUN cargo install cargo-binstall
RUN cargo binstall -y --force cargo-risczero
RUN cargo risczero install --version ${RISC0_TOOLCHAIN_VERSION}

#ENTRYPOINT [ "/bin/sh" ]
# syntax=docker/dockerfile:1

# Comments are provided throughout this file to help you get started.
# If you need more help, visit the Dockerfile reference guide at
# https://docs.docker.com/go/dockerfile-reference/

# Want to help us make this template better? Share your feedback here: https://forms.gle/ybq9Krt8jtBL3iCk7

#ARG RUST_VERSION=1.79.0
ARG APP_NAME=proving-server

################################################################################
# Create a stage for building the application.

#FROM rust:${RUST_VERSION}-alpine AS build
#ARG APP_NAME
#WORKDIR /app

# Install host build dependencies.
#RUN apk add --no-cache clang lld musl-dev git build-base
#RUN cargo install cargo-binstall
#RUN cargo binstall cargo-risczero --no-confirm                                                                                                                                                                                                                                                                
#RUN cargo risczero install   

# Build the application.
# Leverage a cache mount to /usr/local/cargo/registry/
# for downloaded dependencies, a cache mount to /usr/local/cargo/git/db
# for git repository dependencies, and a cache mount to /app/target/ for
# compiled dependencies which will speed up subsequent builds.
# Leverage a bind mount to the src directory to avoid having to copy the
# source code into the container. Once built, copy the executable to an
# output directory before the cache mounted /app/target is unmounted.
COPY . /app
WORKDIR /app
#RUN --mount=type=bind,source=host,target=/app/host \
#    --mount=type=bind,source=methods,target=/app/methods \
#    --mount=type=bind,source=Cargo.toml,target=Cargo.toml \
#    --mount=type=bind,source=Cargo.lock,target=Cargo.lock \
#    --mount=type=cache,target=/app/target/ \
#    --mount=type=cache,target=/usr/local/cargo/git/db \
#    --mount=type=cache,target=/usr/local/cargo/registry/ \
RUN cargo build --release && \
    cp ./target/release/$APP_NAME /bin/server

################################################################################
# Create a new stage for running the application that contains the minimal
# runtime dependencies for the application. This often uses a different base
# image from the build stage where the necessary files are copied from the build
# stage.
#
# The example below uses the alpine image as the foundation for running the app.
# By specifying the "3.18" tag, it will use version 3.18 of alpine. If
# reproducability is important, consider using a digest
# (e.g., alpine@sha256:664888ac9cfd28068e062c991ebcff4b4c7307dc8dd4df9e728bedde5c449d91).
FROM alpine:3.18 AS final

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
COPY --from=builder /bin/server /bin/

# Expose the port that the application listens on.
EXPOSE 50051

# What the container should run when it is started.
CMD ["/bin/server"]
