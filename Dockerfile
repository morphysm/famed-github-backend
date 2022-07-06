# https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/base
COPY famed-github-backend /usr/bin/famed-github-backend
COPY config.json /usr/bin/config.json
ENTRYPOINT ["/usr/bin/famed-github-backend"]