# KafkaBiller

## Installation
Get the Repository
```bash
git clone github.com/j03hanafi/KafkaBiller
cd KafkaBiller
```
Prepare Package
```bash
go get github.com/confluentinc/confluent-kafka-go
go get github.com/gorilla/mux
go get github.com/mofax/iso8583
go get github.com/rivo/uniseg
```
Run the Program
```bash
go run .
```
Build and Execute the Program
```bash
go build {program_name}
./{program_name}
```

## Preparation
Make sure to add ```storage/request/``` and ```storage/response/``` directory to store request/response ISO file
