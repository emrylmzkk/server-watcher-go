example .env file

ADMIN_USERNAME=admin
ADMIN_PASSWORD=admin123!
BACKEND_PORT=5058



pm2 start "go run main.go" --name server-watcher

pm2 logs server-watche






