FROM ubuntu

RUN apt update && \
    DEBIAN_FRONTEND=noninteractive apt install -y tzdata && \
    rm -rf /var/lib/apt/lists/*

RUN apt update && \
    apt install -y wget build-essential

WORKDIR /project-service

COPY config/config.yaml config/
COPY project-service /project-service

EXPOSE 8012

CMD ["./project-service"]