# 🚀 Ultra-Fast File Uploader

Ultra-Fast File Uploader adalah aplikasi web berbasis Go yang memungkinkan upload file dengan enkripsi tingkat tinggi ke Google Drive. File dienkripsi dengan AES-256-GCM dan disimpan dengan nama palsu untuk keamanan maksimal.

## ✨ Fitur Utama

- **Enkripsi Tingkat Tinggi**: AES-256-GCM dengan padding dan obfuscation
- **Upload Paralel**: Mendukung multiple upload simultan dengan pool management
- **Caching Intelligent**: Memory cache dengan batch writing ke MongoDB
- **Keamanan Maksimal**: File tersimpan dengan fake names dan enkripsi penuh
- **Performa Tinggi**: Optimasi memory dengan buffer pools dan goroutine workers
- **RESTful API**: Clean API untuk upload dan download file
- **Media Preview**: Dukungan preview untuk file gambar, video, dan audio
- **Quota Monitoring**: Real-time monitoring storage Google Drive

## 🛠️ Teknologi yang Digunakan

- **Backend**: Go (Golang) dengan Fiber framework
- **Database**: MongoDB
- **Storage**: Google Drive API
- **Enkripsi**: AES-256-GCM
- **Cache**: In-memory caching dengan sync.Map

## 📋 Prasyarat

Sebelum memulai instalasi, pastikan Anda memiliki:

- Go 1.19 atau lebih baru
- MongoDB 4.4 atau lebih baru  
- Node.js dan npm (untuk PM2)
- Google Drive API credentials
- VPS atau server (kami rekomendasikan **Atlantic Server** untuk performa optimal)

### 🌟 Rekomendasi VPS: Atlantic Server

Untuk performa terbaik, kami merekomendasikan menggunakan VPS dari **[Atlantic Server](https://atlantic-server.com)**:

- ⚡ SSD NVMe ultra-cepat
- 🔒 Keamanan tingkat enterprise
- 🌍 Multiple lokasi server
- 💬 Support 24/7 dalam bahasa Indonesia
- 💰 Harga kompetitif dengan performa premium
- 🚀 Bandwidth unlimited untuk traffic tinggi

**Spesifikasi minimal yang direkomendasikan:**
- CPU: 2 vCPU
- RAM: 4GB
- Storage: 20GB SSD NVMe
- Bandwidth: Unlimited

## 🚀 Instalasi dan Setup

### 1. Clone Repository

```bash
git clone https://github.com/yourusername/ultra-fast-uploader.git
cd ultra-fast-uploader
```

### 2. Install Dependencies Go

```bash
go mod tidy
```

### 3. Setup Google Drive API

1. Buka [Google Cloud Console](https://console.cloud.google.com/)
2. Buat project baru atau pilih project yang sudah ada
3. Enable Google Drive API
4. Buat credentials OAuth 2.0
5. Download file JSON credentials

### 4. Setup MongoDB

#### Instalasi MongoDB di Ubuntu/Debian:

```bash
# Import MongoDB GPG key
wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | sudo apt-key add -

# Add MongoDB repository
echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/6.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-6.0.list

# Update dan install
sudo apt update
sudo apt install -y mongodb-org

# Start MongoDB service
sudo systemctl start mongod
sudo systemctl enable mongod
```

#### Untuk CentOS/RHEL:

```bash
# Buat file repository
sudo tee /etc/yum.repos.d/mongodb-org-6.0.repo <<EOF
[mongodb-org-6.0]
name=MongoDB Repository
baseurl=https://repo.mongodb.org/yum/redhat/8/mongodb-org/6.0/x86_64/
gpgcheck=1
enabled=1
gpgkey=https://www.mongodb.org/static/pgp/server-6.0.asc
EOF

# Install MongoDB
sudo yum install -y mongodb-org

# Start service
sudo systemctl start mongod
sudo systemctl enable mongod
```

### 5. Konfigurasi Aplikasi

Buat file `config.json`:

```json
{
    "client_id": "your-google-client-id",
    "client_secret": "your-google-client-secret", 
    "redirect_uri": "http://localhost:5000/oauth/callback",
    "refresh_token": "your-refresh-token",
    "folder_id": "your-google-drive-folder-id",
    "mongo_uri": "mongodb://localhost:27017"
}
```

### 6. Mendapatkan Refresh Token

Jalankan script helper untuk mendapatkan refresh token:

```bash
go run scripts/get_token.go
```

Ikuti instruksi yang muncul untuk authorize aplikasi dan dapatkan refresh token.

### 7. Build Aplikasi

```bash
go build -o uploader main.go
```

## 🏃‍♂️ Menjalankan Aplikasi

### Development Mode

```bash
go run main.go
```

### Production Mode

```bash
./uploader
```

Aplikasi akan berjalan di `http://localhost:5000`

## 🔧 Setup PM2 untuk Production

### 1. Install PM2

```bash
npm install -g pm2
```

### 2. Buat file ecosystem.config.js

```javascript
module.exports = {
  apps: [{
    name: 'ultra-fast-uploader',
    script: './uploader',
    instances: 'max',
    exec_mode: 'cluster',
    env: {
      NODE_ENV: 'production',
      PORT: 5000
    },
    error_file: './logs/err.log',
    out_file: './logs/out.log',
    log_file: './logs/combined.log',
    time: true,
    max_memory_restart: '1G',
    restart_delay: 4000,
    max_restarts: 10,
    min_uptime: '10s'
  }]
};
```

### 3. Jalankan dengan PM2

```bash
# Buat direktori logs
mkdir logs

# Start aplikasi
pm2 start ecosystem.config.js

# Save PM2 configuration
pm2 save

# Setup auto-startup
pm2 startup
```

### 4. Monitoring PM2

```bash
# Monitor real-time
pm2 monit

# Lihat logs
pm2 logs ultra-fast-uploader

# Restart aplikasi
pm2 restart ultra-fast-uploader

# Stop aplikasi
pm2 stop ultra-fast-uploader

# Reload (zero-downtime restart)
pm2 reload ultra-fast-uploader
```

## 📡 API Documentation

### Upload File

```http
POST /upload
Content-Type: multipart/form-data

Body: files (max 50MB per file)
```

**Response:**
```json
{
  "success": true,
  "id": "uuid-file-id",
  "url": "/file/uuid-file-id",
  "stats": {
    "start_time": "2024-01-01T10:00:00Z",
    "end_time": "2024-01-01T10:00:05Z", 
    "duration": "5.2s",
    "total_size": 1048576,
    "upload_speed": "200.0 KB/s"
  }
}
```

### Download/View File

```http
GET /file/:id
```

### Check Drive Quota

```http
GET /quota
```

**Response:**
```json
{
  "used": 5368709120,
  "total": 16106127360,
  "free": 10737418240
}
```

## 🔒 Keamanan

### Enkripsi File

- **Algorithm**: AES-256-GCM
- **Key Derivation**: SHA-256
- **Padding**: 3x random padding untuk obfuscation
- **Fake Headers**: ZIP-like headers untuk menyamarkan file
- **Random Noise**: 1KB random data tambahan

### Keamanan Database

- Password unik untuk setiap file
- Fake filename untuk menyembunyikan file asli
- Metadata terenkripsi

### Rekomendasi Keamanan Production

1. **Gunakan HTTPS**: Setup SSL/TLS dengan Let's Encrypt
2. **Firewall**: Batasi akses port yang tidak perlu
3. **Rate Limiting**: Implementasi rate limiting untuk API
4. **Backup**: Setup backup otomatis database dan konfigurasi

## ⚙️ Optimasi Performance

### Konfigurasi MongoDB untuk Production

```javascript
// mongo shell configuration
db.adminCommand({
  "setParameter": 1,
  "wiredTigerCacheSizeGB": 2
})

// Index untuk performa query
db.files.createIndex({"_id": 1})
db.files.createIndex({"uploaded_at": -1})
```

### Optimasi Sistem Linux

```bash
# Increase file descriptor limits
echo "* soft nofile 65536" >> /etc/security/limits.conf
echo "* hard nofile 65536" >> /etc/security/limits.conf

# Optimize network settings
echo "net.core.rmem_max = 16777216" >> /etc/sysctl.conf
echo "net.core.wmem_max = 16777216" >> /etc/sysctl.conf
echo "net.ipv4.tcp_rmem = 4096 12582912 16777216" >> /etc/sysctl.conf
echo "net.ipv4.tcp_wmem = 4096 12582912 16777216" >> /etc/sysctl.conf

# Apply changes
sysctl -p
```

## 🐳 Docker Deployment

### Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o uploader main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/uploader .
COPY --from=builder /app/config.json .
COPY --from=builder /app/index.html .

EXPOSE 5000
CMD ["./uploader"]
```

### docker-compose.yml

```yaml
version: '3.8'

services:
  uploader:
    build: .
    ports:
      - "5000:5000"
    volumes:
      - ./config.json:/root/config.json
      - ./logs:/root/logs
    depends_on:
      - mongodb
    restart: unless-stopped

  mongodb:
    image: mongo:6.0
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db
    restart: unless-stopped

volumes:
  mongodb_data:
```

## 📊 Monitoring dan Logging

### Setup Log Rotation

```bash
# Buat file logrotate
sudo tee /etc/logrotate.d/uploader <<EOF
/path/to/uploader/logs/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 0644 www-data www-data
    postrotate
        pm2 reload ultra-fast-uploader
    endscript
}
EOF
```

### Monitoring dengan Grafana (Opsional)

1. Install Prometheus dan Grafana
2. Setup metrics endpoint di aplikasi
3. Configure Grafana dashboard untuk monitoring

## 🔧 Troubleshooting

### Common Issues

**1. MongoDB Connection Error**
```bash
# Check MongoDB status
sudo systemctl status mongod

# Check logs
sudo tail -f /var/log/mongodb/mongod.log
```

**2. Google Drive API Quota Exceeded**
```bash
# Check quota usage
curl http://localhost:5000/quota
```

**3. Memory Issues**
```bash
# Monitor memory usage
ps aux | grep uploader
free -h

# Adjust PM2 max memory restart
pm2 start ecosystem.config.js --max-memory-restart 2G
```

### Debug Mode

Set environment variable untuk debug:

```bash
export DEBUG=true
go run main.go
```

## 🤝 Contributing

1. Fork repository ini
2. Buat feature branch (`git checkout -b feature/amazing-feature`)
3. Commit perubahan (`git commit -m 'Add amazing feature'`)
4. Push ke branch (`git push origin feature/amazing-feature`)
5. Buat Pull Request

## 📄 License

Distributed under the MIT License. See `LICENSE` file for more information.

## 🙏 Acknowledgments

- [Fiber Framework](https://gofiber.io/) - Fast HTTP framework untuk Go
- [MongoDB](https://www.mongodb.com/) - Database NoSQL yang powerful
- [Google Drive API](https://developers.google.com/drive) - Cloud storage solution
- [Atlantic Server](https://atlantic-server.com) - Premium VPS hosting partner

## 📞 Support

Jika Anda mengalami masalah atau memiliki pertanyaan:

1. Buka issue di GitHub repository
2. Join Discord community kami
3. Email: support@yourproject.com

---

**Dibuat dengan ❤️ untuk komunitas developer Indonesia**

> 💡 **Tips**: Untuk performa optimal, gunakan VPS dari [Atlantic Server](https://atlantic-server.com) dengan spesifikasi SSD NVMe dan bandwidth unlimited!
