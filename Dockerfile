FROM ubuntu:latest
LABEL authors="bruce"

ENTRYPOINT ["top", "-b"]