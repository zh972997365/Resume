package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Host        string
	Port        string
	Debug       bool
	StoragePath string
	MaxFileSize int64
	BaseURL     string
	DBDriver    string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBCharset   string
	DBParseTime string
	DBLoc       string
}

func LoadConfig() *Config {
	envPath := findEnvFile()
	if envPath == "" {
		log.Fatal("错误: 未找到 .env 配置文件")
	}

	log.Printf("加载配置文件: %s", envPath)

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("错误: 无法加载 .env 文件: %v", err)
	}

	config := validateAndParseConfig()

	logConfig(config)

	return config
}

func findEnvFile() string {
	workDir, err := os.Getwd()
	if err != nil {
		log.Printf("警告: 无法获取工作目录: %v", err)
		return ""
	}

	execDir := getExecutableDir()

	possiblePaths := []string{
		filepath.Join(workDir, ".env"),
		filepath.Join(execDir, ".env"),
		filepath.Join(workDir, "..", ".env"),
		filepath.Join(execDir, "..", ".env"),
		os.Getenv("ENV_FILE_PATH"),
	}

	for _, path := range possiblePaths {
		if path == "" {
			continue
		}

		cleanPath := filepath.Clean(path)

		if info, err := os.Stat(cleanPath); err == nil && !info.IsDir() {
			absPath, err := filepath.Abs(cleanPath)
			if err == nil {
				return absPath
			}
		}
	}

	return ""
}

func getExecutableDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		return filepath.Dir(filename)
	}

	execPath, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(execPath)
}

func validateAndParseConfig() *Config {
	requiredVars := map[string]string{
		"HOST":         "",
		"PORT":         "",
		"STORAGE_PATH": "",
		"BASE_URL":     "",
		"DB_DRIVER":    "",
		"DB_HOST":      "",
		"DB_PORT":      "",
		"DB_USER":      "",
		"DB_PASSWORD":  "",
		"DB_NAME":      "",
	}

	for key := range requiredVars {
		value := os.Getenv(key)
		if value == "" {
			log.Fatalf("错误: 必需配置项 %s 未设置", key)
		}
		requiredVars[key] = value
	}

	storagePath := requiredVars["STORAGE_PATH"]
	absStoragePath, err := filepath.Abs(storagePath)
	if err != nil {
		log.Fatalf("错误: 无法解析存储路径 '%s': %v", storagePath, err)
	}

	debug := parseBoolConfig("DEBUG", false)

	maxFileSize := parseIntConfig("MAX_FILE_SIZE", 50*1024*1024)

	return &Config{
		Host:        requiredVars["HOST"],
		Port:        requiredVars["PORT"],
		Debug:       debug,
		StoragePath: absStoragePath,
		MaxFileSize: maxFileSize,
		BaseURL:     requiredVars["BASE_URL"],
		DBDriver:    requiredVars["DB_DRIVER"],
		DBHost:      requiredVars["DB_HOST"],
		DBPort:      requiredVars["DB_PORT"],
		DBUser:      requiredVars["DB_USER"],
		DBPassword:  requiredVars["DB_PASSWORD"],
		DBName:      requiredVars["DB_NAME"],
		DBCharset:   os.Getenv("DB_CHARSET"),
		DBParseTime: os.Getenv("DB_PARSE_TIME"),
		DBLoc:       os.Getenv("DB_LOC"),
	}
}

func parseBoolConfig(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	switch value {
	case "true", "TRUE", "1", "yes", "YES", "on", "ON":
		return true
	case "false", "FALSE", "0", "no", "NO", "off", "OFF":
		return false
	default:
		log.Printf("警告: 配置项 %s 的值 '%s' 无效，使用默认值 %v", key, value, defaultValue)
		return defaultValue
	}
}

func parseIntConfig(key string, defaultValue int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Printf("警告: 配置项 %s 的值 '%s' 无效，使用默认值 %d", key, value, defaultValue)
		return defaultValue
	}

	return intValue
}

func logConfig(config *Config) {
	log.Printf("========== 配置信息 ==========")
	log.Printf("服务器地址: %s:%s", config.Host, config.Port)
	log.Printf("调试模式: %v", config.Debug)
	log.Printf("存储路径: %s", config.StoragePath)
	log.Printf("基础URL: %s", config.BaseURL)
	log.Printf("最大文件大小: %d 字节 (%.2f MB)",
		config.MaxFileSize,
		float64(config.MaxFileSize)/(1024*1024))
	log.Printf("数据库驱动: %s", config.DBDriver)
	log.Printf("数据库地址: %s:%s", config.DBHost, config.DBPort)
	log.Printf("数据库名称: %s", config.DBName)
	log.Printf("数据库用户: %s", config.DBUser)
	log.Printf("=============================")
}
