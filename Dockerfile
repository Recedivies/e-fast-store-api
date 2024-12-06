# Build stage
FROM golang:1.23.4-alpine AS builder
WORKDIR /app

COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.16
WORKDIR /app

# set env variables
ARG DB_SOURCE
ARG ENVIRONMENT

ENV DB_SOURCE ${DB_SOURCE}
ENV ENVIRONMENT ${ENVIRONMENT}

COPY --from=builder /app/main .

RUN echo ${DB_SOURCE}

# Create .env file
RUN echo DB_SOURCE=${DB_SOURCE} >> .env && \
    echo ENVIRONMENT=${ENVIRONMENT} >> .env && \
    cat .env

COPY start.sh ./

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
