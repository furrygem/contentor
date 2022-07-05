FROM golang:1.18

WORKDIR /app
COPY . /app/

RUN apt update -y; apt install -y python3-pip


RUN pip3 install -r requirements.txt
RUN go build -o build/app cmd/app/main.go

CMD ["./start.py", "use-private-key", "-f", "key.pem"]

