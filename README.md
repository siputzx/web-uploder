# 🚀 Ultra-Fast Web Uploader

<div align="center">

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org)
[![MongoDB](https://img.shields.io/badge/MongoDB-4.4+-green.svg)](https://mongodb.com)
[![Google Drive API](https://img.shields.io/badge/Google%20Drive-API%20v3-red.svg)](https://developers.google.com/drive)

**Modern, secure, and lightning-fast file uploader with Google Drive integration**

[Demo](https://cdn.arvan.my.id) • [Atlantic Server VPS](https://atlantic-server.com) • [Documentation](#documentation)

</div>

## ✨ Features

🔐 **Advanced Security**
- AES-256 encryption with GCM mode
- Secure password generation
- File obfuscation with fake names
- Multi-layer data protection

⚡ **High Performance**
- Parallel upload processing
- Memory-efficient streaming
- Connection pooling
- Batch database operations
- Built-in caching system

🎯 **Modern Architecture**
- Clean Go codebase
- RESTful API design
- Graceful shutdown handling
- Resource optimization
- Error recovery mechanisms

🌐 **Cloud Integration**
- Google Drive storage backend
- MongoDB database
- Quota monitoring
- Automatic backups

## 🎮 Demo

Experience the uploader in action: **[https://cdn.arvan.my.id](https://cdn.arvan.my.id)**

## 📋 Requirements

- **Go 1.19+**
- **MongoDB 4.4+**
- **Google Drive API credentials**
- **VPS/Server** (Recommended: [Atlantic Server](https://atlantic-server.com))

## 🏗️ Installation

### 1. Clone Repository

```bash
git clone https://github.com/siputzx/web-uploder.git
cd web-uploder
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
```

### 3. Configure Google Drive API

#### Get Google Drive API Credentials

Follow this comprehensive guide: [Upload file to Google Drive with Node.js](https://blog.tericcabrel.com/upload-file-to-google-drive-with-nodejs/)

1. **Create Google Cloud Project**
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project or select existing one

2. **Enable Google Drive API**
   - Navigate to "APIs & Services" > "Library"
   - Search for "Google Drive API" and enable it

3. **Create OAuth 2.0 Credentials**
   - Go to "APIs & Services" > "Credentials"
   - Click "Create Credentials" > "OAuth 2.0 Client IDs"
   - Set application type to "Web application"
   - Add authorized redirect URIs: `http://localhost:3000/auth/callback`

4. **Get Refresh Token**
   ```bash
   # Use the OAuth 2.0 Playground or follow the blog tutorial
   # https://developers.google.com/oauthplayground/
   ```

### 4. Setup Configuration

Create `config.json`:

```json
{
    "client_id": "your-google-client-id",
    "client_secret": "your-google-client-secret",
    "redirect_uri": "http://localhost:3000/auth/callback",
    "refresh_token": "your-refresh-token",
    "folder_id": "your-google-drive-folder-id",
    "mongo_uri": "mongodb://localhost:27017"
}
```

### 5. Setup MongoDB

#### Local Installation
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install mongodb

# Start MongoDB
sudo systemctl start mongodb
sudo systemctl enable mongodb
```

#### Docker Installation
```bash
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

### 6. Build and Run

#### Development Mode
```bash
go run main.go
```

#### Production Build
```bash
go build -o uploader main.go
./uploader
```

## 🚀 Production Deployment

### Using PM2 (Recommended)

#### 1. Install PM2
```bash
npm install -g pm2
```

#### 2. Create PM2 Ecosystem File

Create `ecosystem.config.js`:

```javascript
module.exports = {
  apps: [{
    name: 'web-uploader',
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
    autorestart: true,
    watch: false,
    max_memory_restart: '1G'
  }]
}
```

#### 3. Start with PM2
```bash
# Create logs directory
mkdir logs

# Start application
pm2 start ecosystem.config.js

# Save PM2 configuration
pm2 save

# Setup auto-start on boot
pm2 startup
```

#### 4. PM2 Management Commands
```bash
pm2 status          # Check status
pm2 logs            # View logs
pm2 restart all     # Restart all apps
pm2 stop all        # Stop all apps
pm2 delete all      # Delete all apps
pm2 monit           # Monitor resources
```

### Nginx Reverse Proxy

Create `/etc/nginx/sites-available/uploader`:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    client_max_body_size 50M;

    location / {
        proxy_pass http://localhost:5000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        proxy_read_timeout 600s;
        proxy_send_timeout 600s;
    }
}
```

Enable the site:
```bash
sudo ln -s /etc/nginx/sites-available/uploader /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 🌟 Why Choose Atlantic Server?

<div align="center">

### 🏆 **Premium VPS Hosting Starting from 30,000 IDR/month**

</div>

**Atlantic Server** provides the perfect hosting environment for your web uploader:

✅ **High-Performance Infrastructure**
- SSD storage for lightning-fast file operations
- Optimized for Go applications
- Global CDN network integration

✅ **Reliability & Security**
- 99.9% uptime guarantee
- DDoS protection included
- SSL certificates provided
- Regular automated backups

✅ **Developer-Friendly**
- Easy server management panel
- API for automation
- Expert 24/7 support
- No hidden fees or long-term contracts

✅ **Perfect for File Uploaders**
- High bandwidth allocation
- Unlimited file transfers
- Scalable resources
- MongoDB hosting available

**Contact Atlantic Server:**
- 📧 admin@atlantic-server.com
- 📱 +62 882-2220-7701
- 🏢 Komplek BTN No.6A, Sekejati, Kec. Buahbatu, Kota Bandung, Jawa Barat 40286

[**Get Your VPS Now →**](https://atlantic-server.com)

## 🔧 API Endpoints

### Upload File
```http
POST /upload
Content-Type: multipart/form-data

# Form data:
files: [file]
```

**Response:**
```json
{
  "success": true,
  "id": "uuid-string",
  "url": "/file/uuid-string",
  "stats": {
    "start_time": "2024-01-01T00:00:00Z",
    "end_time": "2024-01-01T00:00:05Z",
    "duration": "5s",
    "total_size": 1048576,
    "upload_speed": "204.8 KB/s"
  }
}
```

### Download File
```http
GET /file/{id}
```

### Check Quota
```http
GET /quota
```

**Response:**
```json
{
  "used": 1073741824,
  "total": 107374182400,
  "free": 106300440576
}
```

## 🎯 Configuration Options

### Environment Variables

```bash
# Server Configuration
PORT=5000
GO_ENV=production

# MongoDB Configuration
MONGO_URI=mongodb://localhost:27017

# Google Drive Configuration
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REFRESH_TOKEN=your-refresh-token
GOOGLE_FOLDER_ID=your-folder-id
```

### Performance Tuning

```go
const (
    ChunkSize     = 1024 * 1024 * 16 // 16MB chunks
    BufferSize    = 1024 * 1024 * 8  // 8MB buffer
    MaxWorkers    = 16               // Concurrent workers
    MaxUploads    = 4                // Parallel uploads
    PaddingFactor = 3                // Encryption padding
)
```

## 🔒 Security Features

- **AES-256-GCM Encryption**: Military-grade encryption for all files
- **Secure Password Generation**: Cryptographically secure random passwords
- **File Obfuscation**: Files stored with fake system names
- **Multi-layer Protection**: Header spoofing and padding for additional security
- **Access Control**: Unique UUID-based file access
- **Data Integrity**: Built-in checksums and validation

## 📊 Performance Optimizations

- **Memory Pool**: Reusable buffer allocation
- **Connection Pooling**: Efficient database connections
- **Batch Operations**: Bulk database writes
- **Streaming**: Low memory footprint for large files
- **Caching**: In-memory record caching
- **Graceful Degradation**: Fallback mechanisms

## 🛠️ Development

### Project Structure
```
web-uploder/
├── main.go              # Main application
├── config.json          # Configuration file
├── index.html           # Frontend interface
├── ecosystem.config.js  # PM2 configuration
├── go.mod              # Go dependencies
├── go.sum              # Go checksums
├── logs/               # Application logs
└── README.md           # This file
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/siputzx/web-uploder/blob/main/LICENSE) file for details.

```
MIT License

Copyright (c) 2024 siputzx

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## 🤝 Support

Need help? Here are your options:

- 📖 Check the [documentation](#documentation)
- 🐛 Report bugs in [Issues](https://github.com/siputzx/web-uploder/issues)
- 💬 Join discussions in [Discussions](https://github.com/siputzx/web-uploder/discussions)
- 🌟 Star the project if you find it useful!

## 🎉 Acknowledgments

- Google Drive API for cloud storage
- MongoDB for database solutions
- Fiber framework for high-performance HTTP
- Atlantic Server for reliable hosting

---

<div align="center">

**Made with ❤️ by [siputzx](https://github.com/siputzx)**

[⭐ Star this project](https://github.com/siputzx/web-uploder) • [🐛 Report Bug](https://github.com/siputzx/web-uploder/issues) • [✨ Request Feature](https://github.com/siputzx/web-uploder/issues)

</div>
