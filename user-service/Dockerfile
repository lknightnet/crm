FROM ubuntu

RUN apt update && \
    DEBIAN_FRONTEND=noninteractive apt install -y tzdata && \
    rm -rf /var/lib/apt/lists/*

RUN apt update && \
    apt install -y wget build-essential

WORKDIR /user-service

COPY config/config.yaml config/
COPY user-service /user-service

EXPOSE 8010

CMD ["./user-service"]