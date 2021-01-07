# build from scratch image
FROM scratch
# set working directory
WORKDIR /app/
# add compiled binary to root path
COPY ./bin/fumeping* fumeping
# start application located at /app/fumeping
ENTRYPOINT ["./fumeping", "-config", "./config.toml"]
