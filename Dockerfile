FROM golang:alpine as go

RUN apk add git nano curl

# # Enable cgo - using C in Go
# RUN apk add build-base


WORKDIR /app
RUN git config --global --add safe.directory /app

COPY ./scripts/buildstep.sh .
RUN ./buildstep.sh


FROM go AS dev

WORKDIR /app
COPY . .
ENTRYPOINT [ "air", "-c", "./configs/air.toml" ]