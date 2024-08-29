run-api-service:
		reflex -r '\.go$$' -s -- sh -c "go run cmd/main.go"