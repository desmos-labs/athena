# To build the DJuno image, just run:
# > docker build -t djuno .
#
# In order to work properly, this Docker container needs to have a volume that:
# - as source points to a directory which contains a config.toml and firebase-config.toml files
# - as destination it points to the /home folder
#
# Simple usage with a mounted data directory (considering ~/.djuno/config as the configuration folder):
# > docker run -it -v ~/.djuno/config:/home djuno djuno parse config.toml firebase-config.json
#
# If you want to run this container as a daemon, you can do so by executing
# > docker run -td -v ~/.djuno/config:/home --name djuno djuno
#
# Once you have done so, you can enter the container shell by executing
# > docker exec -it djuno bash
#
# To exit the bash, just execute
# > exit
FROM golang:1.19-alpine
ARG arch=x86_64

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3 ca-certificates build-base
RUN set -eux; apk add --no-cache $PACKAGES;

# Set working directory for the build
WORKDIR /code

# Add source files
COPY . /code/

# See https://github.com/CosmWasm/wasmvm/releases
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.1.2/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 77b41e65f6c3327d910a7f9284538570727e380ab49bc3c88c8d4053811d5209

ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.1.2/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep e0a0955815a23c139d42781f1cc70beffa916aa74fe649e5c69ee7e95ff13b6b

# Copy the library you want to the final location that will be found by the linker flag `-lwasmvm_muslc`
RUN cp /lib/libwasmvm_muslc.${arch}.a /usr/local/lib/libwasmvm_muslc.a

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN BUILD_TAGS=muslc GOOS=linux GOARCH=amd64 LEDGER_ENABLED=true LINK_STATICALLY=true make build
RUN cp /code/build/djuno /usr/bin/djuno
RUN echo "Ensuring binary is statically linked ..." && (file /usr/bin/djuno | grep "statically linked")

ENTRYPOINT ["djuno"]