FROM golang:1.21.6

WORKDIR /app
RUN apt-get update
RUN apt-get install -y maven openjdk-17-jdk unzip xz-utils zip ca-certificates

RUN export PATH=$PATH:/usr/share/maven/bin

RUN install -m 0755 -d /etc/apt/keyrings
RUN curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc
RUN chmod a+r /etc/apt/keyrings/docker.asc

RUN echo \
      "deb [arch=amd64 signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian \
      $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
      tee /etc/apt/sources.list.d/docker.list > /dev/null
RUN apt-get update
RUN apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

RUN wget https://storage.googleapis.com/flutter_infra_release/releases/stable/linux/flutter_linux_3.16.9-stable.tar.xz
RUN tar xf flutter_linux_3.16.9-stable.tar.xz
RUN export PATH="$PATH:`pwd`/flutter/bin"

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o ./app

EXPOSE 8080
CMD ["./app"]