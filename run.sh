cd ..
export $(grep -v '^#' .env | xargs -d '\n')
cd chat-server-go
go run .