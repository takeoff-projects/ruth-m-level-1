FROM golang:1.17.0-alpine3.13
RUN mkdir /app
ADD . /app
WORKDIR /app
ENV EVENTS_API_URL=https://events-api-ex7otr565q-uc.a.run.app
RUN go build -o main .
CMD ["/app/main"]