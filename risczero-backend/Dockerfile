FROM rust:1.80 AS builder

# Install Risc0
ARG RISC0_VERSION=1.0.5
RUN cargo install cargo-binstall
RUN cargo binstall -y --force cargo-risczero@${RISC0_VERSION}
RUN cargo risczero install

# Install the Protobuf-Compiler
RUN apt-get update
RUN apt-get install protobuf-compiler -y

#Copy the source code to the container
COPY ./host/ /app/host
COPY ./methods/ /app/methods
COPY ./Cargo.toml /app/Cargo.toml
#COPY ./target/release/deps /app/target/release/deps
WORKDIR /app

#Fetch and build the server
RUN cargo fetch
RUN cargo build --release

#Copy the binary to /bin/server
RUN cp ./target/release/proving-server /bin/server

FROM rust:1.80 AS final

# Install Risc0
ARG RISC0_VERSION=1.0.5
RUN cargo install cargo-binstall
RUN cargo binstall -y --force cargo-risczero@${RISC0_VERSION}
RUN cargo risczero install

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
