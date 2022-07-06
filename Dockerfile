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
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates

# Install bash
RUN apk add --no-cache bash

# Copy over binaries from the build-env
COPY --from=desmoslabs/builder:latest /code/build/djuno /usr/bin/djuno

# Run djuno by default, omit entrypoint to ease using container with desmos
CMD ["djuno"]