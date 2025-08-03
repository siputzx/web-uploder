package main

import (
    "bytes"
    "context"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "mime"
    "net/http"
    "os"
    "os/signal"
    "path/filepath"
    "runtime"
    "strings"
    "sync"
    "syscall"
    "time"
    "unsafe"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/compress"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/google/uuid"
    "go.mongodb.org/mongo-driver/bson"
    "google.golang.org/api/googleapi"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/drive/v3"
    "google.golang.org/api/option"
)

const (
    ChunkSize     = 1024 * 1024 * 16 // 16MB chunks
    BufferSize    = 1024 * 1024 * 8  // 8MB buffer
    MaxWorkers    = 16
    MaxUploads    = 4
    PaddingFactor = 3
)

type Config struct {
    ClientID     string `json:"client_id"`
    ClientSecret string `json:"client_secret"`
    RedirectURI  string `json:"redirect_uri"`
    RefreshToken string `json:"refresh_token"`
    FolderID     string `json:"folder_id"`
    MongoURI     string `json:"mongo_uri"`
}

type FileRecord struct {
    ID         string     `bson:"_id"`
    DriveID    string     `bson:"drive_id"`
    Password   string     `bson:"password"`
    FileInfo   FileInfo   `bson:"file_info"`
    UploadedAt time.Time  `bson:"uploaded_at"`
    FakeName   string     `bson:"fake_name"`
    RealSize   int64      `bson:"real_size"`
}

type FileInfo struct {
    Name      string `bson:"name"`
    Extension string `bson:"extension"`
    Size      int64  `bson:"size"`
    MimeType  string `bson:"mime_type"`
}

type UploadStats struct {
    StartTime   time.Time `json:"start_time"`
    EndTime     time.Time `json:"end_time"`
    Duration    string    `json:"duration"`
    TotalSize   int64     `json:"total_size"`
    UploadSpeed string    `json:"upload_speed"`
}

type DriveQuota struct {
    Used  int64 `json:"used"`
    Total int64 `json:"total"`
    Free  int64 `json:"free"`
}

type FileUploader struct {
    service       *drive.Service
    config        Config
    collection    *mongo.Collection
    memoryCache   sync.Map
    writeBuffer   []FileRecord
    writeMutex    sync.RWMutex
    ctx           context.Context
    cancel        context.CancelFunc
    uploadPool    chan struct{}
    workerPool    chan struct{}
    bufferPool    sync.Pool
    encryptPool   sync.Pool
}

var fakeNames = []string{
    "system_backup.tmp", "cache_data.bin", "temp_file.dat",
    "log_archive.tmp", "buffer_swap.bin", "memory_dump.dat",
    "config_cache.tmp", "process_data.bin", "kernel_dump.dat",
    "registry_backup.tmp", "thread_dump.bin", "heap_data.dat",
    "driver_cache.tmp", "service_log.bin", "event_data.dat",
    "boot_config.tmp", "system_info.bin", "device_dump.dat",
    "network_cache.tmp", "socket_data.bin", "protocol_dump.dat",
    "security_log.tmp", "auth_cache.bin", "session_data.dat",
    "resource_dump.tmp", "handle_cache.bin", "object_data.dat",
    "scheduler_log.tmp", "queue_data.bin", "task_dump.dat",

    "app_config.xml", "user_settings.json", "preferences.ini",
    "database_backup.sql", "migration_log.txt", "schema_dump.db",
    "api_response.json", "request_log.txt", "endpoint_cache.xml",
    "session_store.dat", "cookie_cache.tmp", "token_data.bin",
    "upload_temp.file", "download_cache.tmp", "media_buffer.dat",
    "thumbnail_cache.jpg", "preview_data.png", "image_temp.bmp",
    "video_buffer.mp4", "audio_cache.wav", "stream_data.tmp",
    "plugin_config.dll", "extension_data.so", "module_cache.dylib",
    "theme_data.css", "style_cache.scss", "layout_temp.html",
    "script_cache.js", "library_data.min.js", "framework_temp.ts",

    "document_backup.docx", "spreadsheet_temp.xlsx", "presentation_cache.pptx",
    "pdf_buffer.pdf", "text_backup.txt", "note_cache.md",
    "template_data.dotx", "form_backup.xltx", "slide_temp.potx",
    "report_cache.doc", "analysis_temp.xls", "summary_data.ppt",
    "contract_backup.pdf", "invoice_temp.xlsx", "receipt_cache.doc",
    "proposal_data.docx", "budget_temp.xls", "forecast_cache.xlsx",
    "manual_backup.pdf", "guide_temp.doc", "tutorial_cache.txt",
    "readme_data.md", "changelog_temp.txt", "version_cache.log",
    "license_backup.txt", "terms_temp.pdf", "policy_cache.doc",
    "specification_data.docx", "requirement_temp.xlsx", "design_cache.ppt",

    "source_backup.cpp", "header_cache.h", "library_temp.lib",
    "binary_data.exe", "object_cache.o", "assembly_temp.asm",
    "makefile_backup.mk", "script_cache.sh", "batch_temp.bat",
    "config_data.cmake", "project_cache.vcxproj", "solution_temp.sln",
    "package_backup.json", "dependency_cache.lock", "module_temp.py",
    "class_data.java", "interface_cache.kt", "struct_temp.rs",
    "function_backup.js", "method_cache.ts", "procedure_temp.sql",
    "query_data.sql", "view_cache.sql", "trigger_temp.sql",
    "index_backup.sql", "constraint_cache.sql", "procedure_temp.sql",
    "migration_data.rb", "seed_cache.sql", "fixture_temp.yml",

    "image_backup.jpg", "photo_cache.png", "picture_temp.gif",
    "video_data.mp4", "movie_cache.avi", "clip_temp.mov",
    "audio_backup.mp3", "music_cache.wav", "sound_temp.ogg",
    "recording_data.m4a", "podcast_cache.mp3", "voice_temp.wav",
    "graphics_backup.svg", "vector_cache.ai", "design_temp.psd",
    "icon_data.ico", "logo_cache.png", "banner_temp.jpg",
    "animation_backup.gif", "sprite_cache.png", "texture_temp.tga",
    "model_data.obj", "mesh_cache.fbx", "scene_temp.blend",
    "font_backup.ttf", "typeface_cache.otf", "glyph_temp.woff",
    "document_scan.pdf", "page_cache.tiff", "archive_temp.zip",

    "dataset_backup.csv", "records_cache.json", "entries_temp.xml",
    "database_dump.sql", "table_cache.db", "collection_temp.bson",
    "index_data.idx", "search_cache.lucene", "query_temp.elasticsearch",
    "log_backup.log", "trace_cache.txt", "debug_temp.out",
    "metrics_data.json", "stats_cache.csv", "analytics_temp.xml",
    "config_backup.yaml", "settings_cache.toml", "properties_temp.ini",
    "registry_data.reg", "preference_cache.plist", "option_temp.conf",
    "certificate_backup.crt", "key_cache.pem", "token_temp.jwt",
    "hash_data.md5", "checksum_cache.sha256", "signature_temp.sig",
    "backup_archive.tar.gz", "compressed_cache.zip", "packed_temp.rar",
    
    "webpage_backup.html", "style_cache.css", "script_temp.js",
    "component_data.jsx", "template_cache.vue", "module_temp.ts",
    "api_backup.json", "response_cache.xml", "request_temp.http",
    "session_data.cookie", "storage_cache.localStorage", "temp_sessionStorage",
    "manifest_backup.json", "worker_cache.js", "service_temp.sw.js",
    "bundle_data.js", "chunk_cache.js", "vendor_temp.js",
    "asset_backup.css", "resource_cache.scss", "theme_temp.less",
    "font_data.woff2", "icon_cache.svg", "image_temp.webp",
    "config_backup.webpack.js", "build_cache.rollup.js", "compile_temp.babel.js",
    "package_data.npm", "dependency_cache.yarn.lock", "module_temp.node_modules",

    "backup_full.zip", "archive_daily.tar", "compressed_weekly.rar",
    "snapshot_data.7z", "backup_incremental.gz", "archive_monthly.bz2",
    "system_backup.tar.gz", "data_archive.zip", "file_backup.rar",
    "config_snapshot.tar", "setting_backup.7z", "preference_archive.gz",
    "database_backup.sql.gz", "log_archive.tar.bz2", "temp_backup.zip",
    "user_data.backup", "profile_archive.tar", "session_backup.gz",
    "cache_snapshot.zip", "temp_archive.rar", "buffer_backup.7z",
    "memory_dump.tar.gz", "crash_backup.zip", "error_archive.tar",
    "debug_snapshot.gz", "trace_backup.bz2", "log_archive.7z",
    "metrics_backup.zip", "stats_archive.tar", "report_backup.gz",

    "processing_temp.tmp", "workflow_cache.wf", "pipeline_data.pipe",
    "batch_process.batch", "queue_item.queue", "job_data.job",
    "task_temp.task", "worker_cache.work", "thread_data.thread",
    "process_buffer.proc", "execution_temp.exec", "runtime_cache.run",
    "compile_temp.build", "link_cache.link", "deploy_temp.deploy",
    "test_data.test", "mock_cache.mock", "stub_temp.stub",
    "fixture_data.fixture", "sample_cache.sample", "demo_temp.demo",
    "prototype_data.proto", "template_cache.tmpl", "pattern_temp.pattern",
    "schema_data.schema", "model_cache.model", "entity_temp.entity",
    "service_data.service", "handler_cache.handler", "controller_temp.ctrl",
    "middleware_data.middleware", "filter_cache.filter", "interceptor_temp.int",

    "utility_backup.util", "helper_cache.help", "tool_temp.tool",
    "script_data.script", "command_cache.cmd", "function_temp.func",
    "library_backup.lib", "module_cache.mod", "package_temp.pkg",
    "plugin_data.plugin", "extension_cache.ext", "addon_temp.addon",
    "component_backup.comp", "widget_cache.widget", "control_temp.ctrl",
    "service_data.svc", "daemon_cache.daemon", "agent_temp.agent",
    "monitor_backup.mon", "watcher_cache.watch", "observer_temp.obs",
    "listener_data.listen", "handler_cache.handle", "processor_temp.proc",
    "parser_backup.parse", "validator_cache.valid", "formatter_temp.fmt",
    "converter_data.conv", "transformer_cache.trans", "mapper_temp.map",
    "serializer_backup.serial", "encoder_cache.encode", "decoder_temp.decode",

    "network_config.net", "connection_cache.conn", "socket_temp.sock",
    "protocol_data.proto", "packet_cache.packet", "frame_temp.frame",
    "message_backup.msg", "request_cache.req", "response_temp.resp",
    "header_data.header", "payload_cache.payload", "body_temp.body",
    "session_backup.session", "token_cache.token", "auth_temp.auth",
    "certificate_data.cert", "key_cache.key", "signature_temp.sig",
    "encryption_backup.enc", "hash_cache.hash", "digest_temp.digest",
    "checksum_data.chk", "validation_cache.valid", "verification_temp.verify",
    "firewall_backup.fw", "security_cache.sec", "access_temp.access",
    "permission_data.perm", "privilege_cache.priv", "right_temp.right",
    "policy_backup.policy", "rule_cache.rule", "condition_temp.cond",

    "security_log.sec", "audit_cache.audit", "compliance_temp.comp",
    "vulnerability_data.vuln", "threat_cache.threat", "risk_temp.risk",
    "scan_backup.scan", "analysis_cache.analysis", "report_temp.report",
    "incident_data.incident", "alert_cache.alert", "warning_temp.warn",
    "forensic_backup.forensic", "evidence_cache.evidence", "trace_temp.trace",
    "malware_data.malware", "virus_cache.virus", "trojan_temp.trojan",
    "quarantine_backup.quar", "sandbox_cache.sandbox", "isolation_temp.iso",
    "backup_encrypted.enc", "archive_secured.sec", "data_protected.prot",
    "file_locked.lock", "content_sealed.seal", "information_hidden.hide",
    "secret_data.secret", "private_cache.private", "confidential_temp.conf",
    "classified_backup.class", "restricted_cache.restrict", "limited_temp.limit",

    "monitor_log.monitor", "performance_cache.perf", "benchmark_temp.bench",
    "metrics_data.metrics", "statistics_cache.stats", "analytics_temp.analytics",
    "profiling_backup.profile", "tracing_cache.trace", "debugging_temp.debug",
    "diagnostic_data.diag", "health_cache.health", "status_temp.status",
    "uptime_backup.uptime", "availability_cache.avail", "reliability_temp.rel",
    "load_data.load", "capacity_cache.capacity", "usage_temp.usage",
    "resource_backup.resource", "memory_cache.memory", "cpu_temp.cpu",
    "disk_data.disk", "network_cache.network", "io_temp.io",
    "bandwidth_backup.bandwidth", "latency_cache.latency", "throughput_temp.through",
    "response_data.response", "request_cache.request", "transaction_temp.trans",

    "app_config.config", "system_settings.settings", "user_preferences.prefs",
    "environment_vars.env", "runtime_config.runtime", "deployment_settings.deploy",
    "database_config.db.config", "server_settings.server", "client_config.client",
    "api_settings.api", "service_config.service", "module_settings.module",
    "plugin_config.plugin", "theme_settings.theme", "locale_config.locale",
    "cache_settings.cache", "session_config.session", "cookie_settings.cookie",
    "security_config.security", "auth_settings.auth", "permission_config.permission",
    "logging_settings.logging", "debug_config.debug", "error_settings.error",
    "backup_config.backup", "archive_settings.archive", "restore_config.restore",
    "sync_settings.sync", "replication_config.replication", "mirror_settings.mirror",
    "cluster_config.cluster", "node_settings.node", "shard_config.shard",
}

func NewFileUploader(configPath string) (*FileUploader, error) {
    configFile, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }

    var config Config
    if err := json.Unmarshal(configFile, &config); err != nil {
        return nil, err
    }

    ctx, cancel := context.WithCancel(context.Background())
    oauth2Config := &oauth2.Config{
        ClientID:     config.ClientID,
        ClientSecret: config.ClientSecret,
        RedirectURL:  config.RedirectURI,
        Scopes:       []string{drive.DriveScope},
        Endpoint:     google.Endpoint,
    }

    token := &oauth2.Token{RefreshToken: config.RefreshToken}
    client := oauth2Config.Client(ctx, token)
    client.Timeout = 5 * time.Minute

    service, err := drive.NewService(ctx, option.WithHTTPClient(client))
    if err != nil {
        cancel()
        return nil, err
    }

    clientOpts := options.Client().
        ApplyURI(config.MongoURI).
        SetMaxPoolSize(20).
        SetMinPoolSize(5).
        SetMaxConnIdleTime(30 * time.Second).
        SetServerSelectionTimeout(5 * time.Second).
        SetSocketTimeout(10 * time.Second)

    mongoClient, err := mongo.Connect(ctx, clientOpts)
    if err != nil {
        cancel()
        return nil, err
    }

    collection := mongoClient.Database("Uploader").Collection("files")
    
    fu := &FileUploader{
        service:       service,
        config:        config,
        collection:    collection,
        writeBuffer:   make([]FileRecord, 0, 100),
        ctx:           ctx,
        cancel:        cancel,
        uploadPool:    make(chan struct{}, MaxUploads),
        workerPool:    make(chan struct{}, MaxWorkers),
        bufferPool: sync.Pool{
            New: func() interface{} {
                return make([]byte, BufferSize)
            },
        },
        encryptPool: sync.Pool{
            New: func() interface{} {
                return &bytes.Buffer{}
            },
        },
    }

    runtime.GC()
    go fu.batchWrite()
    go fu.handleShutdown()

    return fu, nil
}

func (fu *FileUploader) batchWrite() {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            fu.flushBuffer()
        case <-fu.ctx.Done():
            fu.flushBuffer()
            return
        }
    }
}

func (fu *FileUploader) flushBuffer() {
    fu.writeMutex.Lock()
    if len(fu.writeBuffer) == 0 {
        fu.writeMutex.Unlock()
        return
    }

    batch := make([]FileRecord, len(fu.writeBuffer))
    copy(batch, fu.writeBuffer)
    fu.writeBuffer = fu.writeBuffer[:0]
    fu.writeMutex.Unlock()

    go func(records []FileRecord) {
        operations := make([]mongo.WriteModel, len(records))
        for i, record := range records {
            op := mongo.NewReplaceOneModel()
            op.SetFilter(bson.M{"_id": record.ID})
            op.SetReplacement(record)
            op.SetUpsert(true)
            operations[i] = op
            fu.memoryCache.Store(record.ID, record)
        }

        opts := options.BulkWrite().SetOrdered(false)
        ctx, cancel := context.WithTimeout(fu.ctx, 30*time.Second)
        defer cancel()
        
        _, err := fu.collection.BulkWrite(ctx, operations, opts)
        if err != nil {
            log.Printf("Batch write error: %v", err)
        }
    }(batch)
}

func (fu *FileUploader) addToBuffer(record FileRecord) {
    fu.writeMutex.Lock()
    fu.writeBuffer = append(fu.writeBuffer, record)
    if len(fu.writeBuffer) >= 50 {
        fu.writeMutex.Unlock()
        fu.flushBuffer()
    } else {
        fu.writeMutex.Unlock()
    }
}

func (fu *FileUploader) handleShutdown() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
    fu.flushBuffer()
    time.Sleep(1 * time.Second)
    fu.cancel()
    os.Exit(0)
}

func generateSecurePassword() string {
    b := make([]byte, 24)
    rand.Read(b)
    const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
    for i, v := range b {
        b[i] = chars[v%byte(len(chars))]
    }
    return *(*string)(unsafe.Pointer(&b))
}

func getFakeName() string {
    return fakeNames[time.Now().UnixNano()%int64(len(fakeNames))]
}

func (fu *FileUploader) advancedEncrypt(data []byte, password string) ([]byte, error) {
    key := sha256.Sum256([]byte(password))
    block, err := aes.NewCipher(key[:])
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonce := make([]byte, gcm.NonceSize())
    rand.Read(nonce)

    encrypted := gcm.Seal(nonce, nonce, data, nil)
    
    paddingSize := len(encrypted) * PaddingFactor
    padding := make([]byte, paddingSize)
    rand.Read(padding)
    
    fakeHeaders := []byte{
        0x50, 0x4B, 0x03, 0x04, 0x14, 0x00, 0x00, 0x00, 0x08, 0x00,
        0x21, 0x0C, 0x4B, 0x50, 0x28, 0xB5, 0x2F, 0xFD, 0x87, 0x00,
        0x00, 0x00, 0x75, 0x00, 0x00, 0x00, 0x08, 0x00, 0x1C, 0x00,
    }
    
    result := make([]byte, 0, len(fakeHeaders)+len(padding)+len(encrypted)+1024)
    result = append(result, fakeHeaders...)
    result = append(result, padding...)
    result = append(result, encrypted...)
    
    extraNoise := make([]byte, 1024)
    rand.Read(extraNoise)
    result = append(result, extraNoise...)
    
    return result, nil
}

func (fu *FileUploader) advancedDecrypt(data []byte, password string) ([]byte, error) {
    if len(data) < 1024+30 {
        return nil, fmt.Errorf("invalid data")
    }

    headerSize := 30
    originalLen := len(data) - headerSize - 1024
    paddingSize := originalLen * PaddingFactor / (PaddingFactor + 1)
    encryptedStart := headerSize + paddingSize
    encryptedEnd := encryptedStart + (originalLen - paddingSize)

    if encryptedEnd > len(data)-1024 {
        return nil, fmt.Errorf("data corruption")
    }

    encrypted := data[encryptedStart:encryptedEnd]

    key := sha256.Sum256([]byte(password))
    block, err := aes.NewCipher(key[:])
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    if len(encrypted) < gcm.NonceSize() {
        return nil, fmt.Errorf("invalid encrypted data")
    }

    nonce, ciphertext := encrypted[:gcm.NonceSize()], encrypted[gcm.NonceSize():]
    return gcm.Open(nil, nonce, ciphertext, nil)
}

func (fu *FileUploader) parallelUpload(data []byte, filename string) (string, error) {
    fu.uploadPool <- struct{}{}
    defer func() { <-fu.uploadPool }()

    file := &drive.File{
        Name:    filename,
        Parents: []string{fu.config.FolderID},
    }

    ctx, cancel := context.WithTimeout(fu.ctx, 10*time.Minute)
    defer cancel()

    reader := bytes.NewReader(data)
    call := fu.service.Files.Create(file).Context(ctx).Media(reader, googleapi.ChunkSize(ChunkSize))
    
    res, err := call.Do()
    if err != nil {
        return "", err
    }
    return res.Id, nil
}

func (fu *FileUploader) streamDownload(driveID string) ([]byte, error) {
    fu.uploadPool <- struct{}{}
    defer func() { <-fu.uploadPool }()

    ctx, cancel := context.WithTimeout(fu.ctx, 10*time.Minute)
    defer cancel()

    resp, err := fu.service.Files.Get(driveID).Context(ctx).Download()
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    buf := fu.bufferPool.Get().([]byte)
    defer fu.bufferPool.Put(buf)

    var result []byte
    for {
        n, err := resp.Body.Read(buf)
        if n > 0 {
            result = append(result, buf[:n]...)
        }
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }
    }
    return result, nil
}

func isMediaFile(mimeType string) bool {
    return strings.HasPrefix(mimeType, "video/") ||
           strings.HasPrefix(mimeType, "image/") ||
           strings.HasPrefix(mimeType, "audio/")
}

func (fu *FileUploader) handleUpload(c *fiber.Ctx) error {
    startTime := time.Now()

    form, err := c.MultipartForm()
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid form"})
    }

    fileHeaders := form.File["files"]
    if len(fileHeaders) != 1 {
        return c.Status(400).JSON(fiber.Map{"error": "Single file only"})
    }

    fileHeader := fileHeaders[0]
    if fileHeader.Size > 50*1024*1024 {
        return c.Status(413).JSON(fiber.Map{"error": "File too large"})
    }

    file, err := fileHeader.Open()
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Cannot open file"})
    }
    defer file.Close()

    data := make([]byte, fileHeader.Size)
    _, err = io.ReadFull(file, data)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Read failed"})
    }

    uniqueID := uuid.New().String()
    password := generateSecurePassword()
    fakeName := getFakeName()

    ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
    mimeType := mime.TypeByExtension(ext)
    if mimeType == "" {
        mimeType = http.DetectContentType(data[:min(512, len(data))])
    }

    fileInfo := FileInfo{
        Name:      fileHeader.Filename,
        Extension: ext,
        Size:      fileHeader.Size,
        MimeType:  mimeType,
    }

    encryptedData, err := fu.advancedEncrypt(data, password)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Encryption failed"})
    }

    driveID, err := fu.parallelUpload(encryptedData, fakeName)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Upload failed"})
    }

    endTime := time.Now()
    duration := endTime.Sub(startTime)

    record := FileRecord{
        ID:         uniqueID,
        DriveID:    driveID,
        Password:   password,
        FileInfo:   fileInfo,
        UploadedAt: startTime,
        FakeName:   fakeName,
        RealSize:   fileHeader.Size,
    }

    // Store in memory cache immediately
    fu.memoryCache.Store(uniqueID, record)
    
    // Store in database synchronously to ensure it's saved
    ctx, cancel := context.WithTimeout(fu.ctx, 10*time.Second)
    defer cancel()
    
    _, err = fu.collection.ReplaceOne(ctx, 
        bson.M{"_id": uniqueID}, 
        record, 
        options.Replace().SetUpsert(true))
    if err != nil {
        log.Printf("Database save error: %v", err)
        // Still continue as we have it in memory cache
    }

    stats := UploadStats{
        StartTime:   startTime,
        EndTime:     endTime,
        Duration:    duration.String(),
        TotalSize:   fileHeader.Size,
        UploadSpeed: formatSpeed(fileHeader.Size, duration),
    }

    return c.JSON(fiber.Map{
        "success": true,
        "id":      uniqueID,
        "url":     fmt.Sprintf("/file/%s", uniqueID),
        "stats":   stats,
    })
}

func (fu *FileUploader) handleFile(c *fiber.Ctx) error {
    id := c.Params("id")
    if id == "" {
        return c.Status(400).JSON(fiber.Map{"error": "ID required"})
    }

    // Try to get from memory cache first
    recordInterface, exists := fu.memoryCache.Load(id)
    var record FileRecord
    
    if exists {
        record = recordInterface.(FileRecord)
    } else {
        // Try to get from database if not in memory cache
        ctx, cancel := context.WithTimeout(fu.ctx, 5*time.Second)
        defer cancel()
        
        err := fu.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&record)
        if err != nil {
            if err == mongo.ErrNoDocuments {
                return c.Status(404).JSON(fiber.Map{"error": "File not found"})
            }
            return c.Status(500).JSON(fiber.Map{"error": "Database error"})
        }
        
        // Store back in memory cache for future requests
        fu.memoryCache.Store(id, record)
    }

    encryptedData, err := fu.streamDownload(record.DriveID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Download failed"})
    }

    data, err := fu.advancedDecrypt(encryptedData, record.Password)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Decryption failed"})
    }

    c.Set("Content-Type", record.FileInfo.MimeType)
    c.Set("Cache-Control", "public, max-age=31536000")
    
    if isMediaFile(record.FileInfo.MimeType) {
        c.Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", record.FileInfo.Name))
    } else {
        c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", record.FileInfo.Name))
    }
    
    return c.Send(data)
}

func (fu *FileUploader) handleQuota(c *fiber.Ctx) error {
    about, err := fu.service.About.Get().Fields("storageQuota").Do()
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Quota check failed"})
    }

    quota := &DriveQuota{
        Used:  about.StorageQuota.Usage,
        Total: about.StorageQuota.Limit,
        Free:  about.StorageQuota.Limit - about.StorageQuota.Usage,
    }
    return c.JSON(quota)
}

func formatSpeed(bytes int64, duration time.Duration) string {
    if duration.Seconds() == 0 {
        return "∞ B/s"
    }
    return formatBytes(int64(float64(bytes)/duration.Seconds())) + "/s"
}

func formatBytes(bytes int64) string {
    const unit = 1024
    if bytes < unit {
        return fmt.Sprintf("%d B", bytes)
    }
    div, exp := int64(unit), 0
    for n := bytes / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    
    uploader, err := NewFileUploader("config.json")
    if err != nil {
        log.Fatal(err)
    }

    app := fiber.New(fiber.Config{
        BodyLimit:       50 * 1024 * 1024,
        ReadTimeout:     10 * time.Minute,
        WriteTimeout:    10 * time.Minute,
        IdleTimeout:     120 * time.Second,
        Prefork:         false,
        ServerHeader:    "",
        DisableKeepalive: false,
        StreamRequestBody: true,
    })

    app.Use(compress.New(compress.Config{
        Level: compress.LevelBestSpeed,
    }))

    app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
        AllowMethods: "GET,POST,OPTIONS",
        AllowHeaders: "*",
    }))

    app.Post("/upload", uploader.handleUpload)
    app.Get("/file/:id", uploader.handleFile)
    app.Get("/quota", uploader.handleQuota)
    app.Get("/", func(c *fiber.Ctx) error {
    		return c.SendFile("./index.html")
    	})

    log.Printf("🚀 Ultra-Fast Uploader running on port 5000")
    log.Fatal(app.Listen(":5000"))
}
