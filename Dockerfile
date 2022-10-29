FROM golang:1.18-alpine AS build

WORKDIR /app
COPY . ./
RUN go mod download
RUN cd /app/cmd/lunch-go && go build -o /lunch-go

FROM alpine

RUN apk update && apk add poppler-utils
COPY --from=build /lunch-go /lunch-go

EXPOSE 8080

CMD [ "/lunch-go" ]
