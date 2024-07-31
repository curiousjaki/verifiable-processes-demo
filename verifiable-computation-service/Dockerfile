FROM alpine:3.18 AS final

#RUN apk add --no-cache zip
# Create a non-privileged user that the app will run under.
# See https://docs.docker.com/go/dockerfile-user-best-practices/
# ARG UID=10001
# RUN adduser \
#     --disabled-password \
#     --gecos "" \
#     --home "/nonexistent" \
#     --shell "/sbin/nologin" \
#     --no-create-home \
#     --uid "${UID}" \
#     appuser
# USER appuser
ADD README.md /bin/proving-server
# Copy the executable from the "build" stage.
#COPY ./target /bin/target
#COPY ./target/release/proving-server.zip /bin/proving-server.zip
#COPY ./target/riscv-guest/riscv32im-risc0-zkvm-elf/docker/carbon_accounting/carbon_accounting.zip /bin/carbon_accounting.zip
#RUN unzip /bin/proving-server.zip -d /bin
#RUN unzip /bin/carbon_accounting.zip -d /bin

# Expose the port that the application listens on.
EXPOSE 50051

# What the container should run when it is started.
CMD ["/bin/proving-server"]
