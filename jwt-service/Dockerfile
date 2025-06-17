FROM ubuntu

RUN apt update && \
    DEBIAN_FRONTEND=noninteractive apt install -y tzdata && \
    rm -rf /var/lib/apt/lists/*

RUN apt update && \
    apt install -y wget build-essential

WORKDIR /jwt-service

COPY config/config.yaml config/
COPY jwt-service /jwt-service

EXPOSE 8009

CMD ["./jwt-service"]