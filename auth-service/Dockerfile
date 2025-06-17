FROM ubuntu

RUN apt update && \
    DEBIAN_FRONTEND=noninteractive apt install -y tzdata && \
    rm -rf /var/lib/apt/lists/*

RUN apt update && \
    apt install -y wget build-essential

WORKDIR /auth-service

COPY config/config.yaml config/
COPY auth-service /auth-service

EXPOSE 8008

CMD ["./auth-service"]