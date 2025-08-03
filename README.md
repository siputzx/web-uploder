# 🚀 Ultra-Fast Web Uploader

<div align="center">
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/MongoDB-4EA94B?style=for-the-badge&logo=mongodb&logoColor=white" alt="MongoDB">
  <img src="https://img.shields.io/badge/Google_Drive-4285F4?style=for-the-badge&logo=googledrive&logoColor=white" alt="Google Drive">
  <img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge" alt="MIT License">
</div>

<div align="center">
  <h3>⚡ Blazing fast file uploader with advanced encryption and secure storage</h3>
  <p>Upload files securely to Google Drive with AES-256 encryption, steganography, and blazing fast performance</p>
  
  <a href="https://cdn.arvan.my.id" target="_blank">🌟 Live Demo</a> •
  <a href="#features">✨ Features</a> •
  <a href="#installation">🛠️ Installation</a> •
  <a href="#deployment">🚀 Deployment</a>
</div>

---

## ✨ Features

### 🔒 **Advanced Security**
- **AES-256-GCM Encryption** - Military-grade encryption for all uploaded files
- **Steganography Protection** - Files disguised with fake headers and random padding
- **Secure Password Generation** - Cryptographically secure 24-character passwords
- **Memory-Safe Operations** - Zero-copy operations and secure memory handling

### ⚡ **High Performance**
- **Parallel Processing** - Multi-threaded upload/download with worker pools
- **Memory Optimization** - Advanced buffer pooling and garbage collection
- **Streaming I/O** - Efficient streaming for large files (up to 50MB)
- **Intelligent Caching** - In-memory cache with MongoDB persistence

### 🎯 **User Experience**
- **Drag & Drop Interface** - Modern, responsive web interface
- **Real-time Progress** - Live upload progress with speed metrics
- **Instant Access** - Direct file access via unique URLs
- **Media Preview** - Inline preview for images, videos, and audio

### 🌐 **Cloud Integration**
- **Google Drive Storage** - Unlimited storage via Google Drive API
- **MongoDB Database** - Fast metadata storage and retrieval
- **Global CDN Ready** - Optimized for content delivery networks
- **API First Design** - RESTful API for easy integration

---

## 🛠️ Installation

### Prerequisites

- **Go 1.19+** - [Download Go](https://golang.org/dl/)
- **MongoDB** - [MongoDB Atlas](https://www.mongodb.com/atlas) or local installation
- **Google Drive API** - [Setup Guide](#google-drive-setup)

### Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/siputzx/web-uploder.git
   cd web-uploder
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Configure the application**
   ```bash
   cp config.json.example config.json
   # Edit config.json with your credentials
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

5. **Access the application**
   ```
   Open http://localhost:5000 in your browser
   ```

---

## ⚙️ Configuration

Create a `config.json` file in the root directory:

```json
{
  "client_id": "your-google-oauth-client-id",
  "client_secret": "your-google-oauth-client-secret",
  "redirect_uri": "http://localhost:5000/auth/callback",
  "refresh_token": "your-google-oauth-refresh-token",
  "folder_id": "your-google-drive-folder-id",
  "mongo_uri": "mongodb://localhost:27017/uploader"
}
```

### 🔑 Google Drive API Setup

Follow this comprehensive guide to set up Google Drive API: [Upload File to Google Drive with Node.js](https://blog.tericcabrel.com/upload-file-to-google-drive-with-nodejs/)

**Quick Steps:**

1. **Create Google Cloud Project**
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project or select existing one

2. **Enable Google Drive API**
   - Navigate to "APIs & Services" > "Library"
   - Search for "Google Drive API" and enable it

3. **Create OAuth 2.0 Credentials**
   - Go to "APIs & Services" > "Credentials"
   - Click "Create Credentials" > "OAuth client ID"
   - Choose "Web application"
   - Add redirect URI: `http://localhost:5000/auth/callback`

4. **Get Refresh Token**
   - Use the OAuth 2.0 Playground or the tutorial above
   - Save the refresh token in your config.json

5. **Create Drive Folder**
   - Create a folder in Google Drive for uploads
   - Copy the folder ID from the URL

---

## 🚀 Deployment

### Production Deployment with PM2

1. **Install PM2 globally**
   ```bash
   npm install -g pm2
   ```

2. **Build the application**
   ```bash
   go build -o web-uploader main.go
   ```

3. **Create PM2 ecosystem file**
   ```bash
   # Create ecosystem.config.js
   cat > ecosystem.config.js << EOF
   module.exports = {
     apps: [{
       name: 'web-uploader',
       script: './web-uploader',
       instances: 'max',
       exec_mode: 'fork',
       env: {
         NODE_ENV: 'production',
         PORT: 5000
       },
       error_file: './logs/err.log',
       out_file: './logs/out.log',
       log_file: './logs/combined.log',
       time: true,
       max_memory_restart: '1G'
     }]
   }
   EOF
   ```

4. **Start with PM2**
   ```bash
   mkdir -p logs
   pm2 start ecosystem.config.js
   pm2 save
   pm2 startup
   ```

5. **Monitor your application**
   ```bash
   pm2 status
   pm2 logs web-uploader
   pm2 monit
   ```

### 🌐 Recommended VPS Hosting

<div align="center">
  <a href="https://atlantic-server.com" target="_blank">
    <img src="https://img.shields.io/badge/Atlantic_Server-Recommended_VPS-00ADD8?style=for-the-badge&logo=server&logoColor=white" alt="Atlantic Server">
  </a>
</div>

**Why Atlantic Server?**

- ⚡ **High Performance** - Optimized infrastructure with global CDN
- 🔒 **Advanced Security** - DDoS protection, SSL certificates, regular backups
- 💰 **Affordable Plans** - Starting from 15,000 IDR/month with no hidden fees
- 🚀 **Easy Scaling** - Seamless resource upgrades via control panel
- 📊 **Monitoring** - Detailed analytics and instant alerts
- 🗄️ **Managed Databases** - Automatic backups and high availability
- 🆘 **24/7 Support** - Expert support team available around the clock

**Contact Information:**
- 📧 Email: admin@atlantic-server.com
- 📞 Phone: +62 882-2220-7701
- 📍 Address: Komplek BTN No.6A, Sekejati, Kec. Buahbatu, Kota Bandung, Jawa Barat 40286

---

## 🗂️ Project Structure

```
web-uploder/
├── main.go              # Main application server
├── config.json          # Configuration file
├── index.html           # Web interface
├── go.mod              # Go module dependencies
├── go.sum              # Go module checksums
├── ecosystem.config.js  # PM2 configuration (after setup)
└── logs/               # Application logs (after setup)
    ├── err.log
    ├── out.log
    └── combined.log
```

---

## 🔧 API Documentation

### Upload File
```http
POST /upload
Content-Type: multipart/form-data

files: File to upload (max 50MB)
```

**Response:**
```json
{
  "success": true,
  "id": "unique-file-id",
  "url": "/file/unique-file-id",
  "stats": {
    "start_time": "2024-01-01T12:00:00Z",
    "end_time": "2024-01-01T12:00:05Z",
    "duration": "5s",
    "total_size": 1048576,
    "upload_speed": "209.7 KB/s"
  }
}
```

### Download File
```http
GET /file/{id}
```

### Check Storage Quota
```http
GET /quota
```

**Response:**
```json
{
  "used": 1073741824,
  "total": 16106127360,
  "free": 15032385536
}
```

---

## 🎯 Performance Optimization

### System Requirements

**Minimum:**
- 1GB RAM
- 1 CPU Core
- 10GB Storage
- 100Mbps Network

**Recommended:**
- 4GB RAM
- 4 CPU Cores
- 50GB SSD Storage
- 1Gbps Network

### Tuning Parameters

```go
const (
    ChunkSize     = 1024 * 1024 * 16 // 16MB chunks
    BufferSize    = 1024 * 1024 * 8  // 8MB buffer
    MaxWorkers    = 16               // Concurrent workers
    MaxUploads    = 4                // Parallel uploads
)
```

### Memory Management

- **Buffer Pooling** - Reuses memory buffers for optimal performance
- **Garbage Collection** - Automatic memory cleanup
- **Memory Limits** - PM2 automatic restart at 1GB memory usage

---

## 🛡️ Security Features

### Encryption Details

- **Algorithm**: AES-256-GCM (Galois/Counter Mode)
- **Key Derivation**: SHA-256 hash of secure password
- **Steganography**: Files disguised with fake ZIP headers
- **Padding**: Random padding (3x file size) to obscure real content
- **Nonce**: Cryptographically secure random nonce per file

### Security Best Practices

1. **Change default ports** in production
2. **Use HTTPS** with SSL certificates
3. **Implement rate limiting** for upload endpoints
4. **Regular security updates** for dependencies
5. **Monitor access logs** for suspicious activity

---

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. **Fork the repository**
2. **Create a feature branch** (`git checkout -b feature/amazing-feature`)
3. **Commit your changes** (`git commit -m 'Add amazing feature'`)
4. **Push to the branch** (`git push origin feature/amazing-feature`)
5. **Open a Pull Request**

### Development Setup

```bash
# Install air for hot reloading (optional)
go install github.com/cosmtrek/air@latest

# Run with hot reloading
air

# Run tests
go test -v ./...
```

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/siputzx/web-uploder/blob/main/LICENSE) file for details.

```
MIT License

Copyright (c) 2024 Web Uploader

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

---

## 📞 Support & Contact

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/siputzx/web-uploder/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/siputzx/web-uploder/discussions)
- 📧 **Email Support**: Contact via GitHub
- 🌟 **Demo**: [https://cdn.arvan.my.id](https://cdn.arvan.my.id)

---

<div align="center">
  <h3>⭐ If you found this project helpful, please give it a star!</h3>
  
  **Built with ❤️ using Go, MongoDB, and Google Drive API**
  
  <a href="https://github.com/siputzx/web-uploder/stargazers">
    <img src="https://img.shields.io/github/stars/siputzx/web-uploder?style=social" alt="GitHub stars">
  </a>
  <a href="https://github.com/siputzx/web-uploder/network/members">
    <img src="https://img.shields.io/github/forks/siputzx/web-uploder?style=social" alt="GitHub forks">
  </a>
</div>
