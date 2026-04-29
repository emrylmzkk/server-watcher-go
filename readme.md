# 🚀 Server Watcher (Go)

[![Go Version](https://img.shields.io/github/go-mod/go-version/emrylmzkk/server-watcher-go?style=for-the-badge&logo=go)](https://golang.org/)
[![Fiber Framework](https://img.shields.io/badge/Fiber-v2.x-00ADED?style=for-the-badge&logo=gofiber)](https://gofiber.io/)
[![GORM](https://img.shields.io/badge/GORM-SQLite-blue?style=for-the-badge&logo=sqlite)](https://gorm.io/)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)

**Server Watcher** is a robust, lightweight monitoring and management dashboard for Linux servers. It provides a centralized interface to manage **PM2 processes** and **Docker containers**, while monitoring system-level resources in real-time.

---

## ✨ Key Features

### 🟢 PM2 Process Management
- **Live Monitoring:** Real-time CPU and Memory usage for each PM2 process.
- **Control Actions:** Start, Stop, Restart, and Delete processes via API.
- **Smart Sync:** Automatically synchronizes the database with your server's current PM2 state.
- **Project Configuration:** Define project paths, start commands, and runtime types (Go, Python, React, etc.).

### 🐳 Docker Container Management
- **Dashboard View:** List all containers with their current status and uptime.
- **Live Stats:** Real-time monitoring of CPU, RAM (MB), and Network (TX/RX) metrics.
- **Lifecycle Control:** Start and Stop containers directly from the dashboard.
- **History:** Paginated container statistics for performance analysis.

### 📊 System Resource Monitoring
- **CPU:** Percentage of total system CPU usage.
- **RAM:** Used, Total, and Percentage of system memory.
- **Disk:** Detailed disk space usage metrics.

### 🔐 Security & Architecture
- **JWT Authentication:** Secure API access with Access and Refresh tokens.
- **Role-Based Access:** Pre-defined roles (Admin/Member) for secure management.
- **Clean Architecture:** Built with a layered structure (Handler -> Service -> Repository) for maintainability.

---

## 🛠 Tech Stack

- **Language:** [Go (Golang)](https://golang.org/)
- **Web Framework:** [Fiber](https://gofiber.io/)
- **ORM:** [GORM](https://gorm.io/)
- **Database:** SQLite (local persistence)
- **Authentication:** JWT (JSON Web Tokens)
- **Process Management:** PM2 & Docker Integration

---

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- PM2 installed (`npm install pm2 -g`)
- Docker installed (optional, for container monitoring)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/emrylmzkk/server-watcher-go.git
   cd server-watcher-go
   ```

2. **Setup environment variables:**
   Create a `.env` file in the root directory:
   ```env
   ADMIN_USERNAME=admin
   ADMIN_PASSWORD=your_secure_password
   BACKEND_PORT=5058
   JWT_SECRET=your_secret_key
   ```

3. **Install dependencies:**
   ```bash
   go mod tidy
   ```

4. **Run the application:**
   ```bash
   go run main.go
   ```

### 📦 Running with PM2
To run the backend in the background using PM2:
```bash
pm2 start "go run main.go" --name server-watcher
```

---

## 📖 API Documentation

The project includes a detailed API specification for frontend integration. 
Check out [frontend_api_docs.md](./frontend_api_docs.md) for endpoint details, request/reponse DTOs, and enumeration values.

---

## 🤝 Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

---

<p align="center">
  Developed with ❤️ by <a href="https://github.com/emrylmzkk">Emirhan YILMAZKAYA</a>
</p>
