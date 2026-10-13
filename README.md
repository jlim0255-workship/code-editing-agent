mkdir code-editing-agent
cd code-editing-agent
go mod init agent
touch main.go

export ANTHROPIC_API_KEY="this is the last time i will tell you to set this"

# Download the dependencies
go mod tidy
# Run it
go run main.go