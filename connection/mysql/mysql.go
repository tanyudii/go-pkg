package mysql

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
	"pkg.tanyudii.me/go-pkg/go-mon/logger"
	"strings"
	"time"
)

type config struct {
	Username string `envconfig:"MYSQL_USERNAME" required:"true"`
	Password string `envconfig:"MYSQL_PASSWORD" required:"true"`
	Host     string `envconfig:"MYSQL_HOST" default:"127.0.0.1"`
	Port     string `envconfig:"MYSQL_PORT" default:"3306"`
	Database string `envconfig:"MYSQL_DATABASE" required:"true"`

	Charset           string `envconfig:"MYSQL_CHARSET" default:"utf8mb4"`
	ParseTime         bool   `envconfig:"MYSQL_PARSE_TIME" default:"true"`
	MultiStatements   bool   `envconfig:"MYSQL_MULTI_STATEMENTS" default:"true"`
	DefaultStringSize uint   `envconfig:"DEFAULT_STRING_SIZE" default:"255"`

	Location        string `envconfig:"MYSQL_LOCATION" default:"Local"`
	LogMode         bool   `envconfig:"MYSQL_LOG_MODE" default:"false"`
	MaxOpenConns    int    `envconfig:"MYSQL_MAX_OPEN_CONNS" default:"100"`
	MaxIdleConns    int    `envconfig:"MYSQL_MAX_IDLE_CONNS" default:"10"`
	ConnMaxLifetime int    `envconfig:"MYSQL_CONN_MAX_LIFETIME" default:"10"`

	db *gorm.DB
}

func (c *config) dsn() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%+v&loc=%s&multiStatements=%+v&tls=skip-verify",
		c.Username, c.Password, c.Host, c.Port, c.Database, c.Charset, c.ParseTime, c.Location, c.MultiStatements,
	)
}

func (c *config) connect() *gorm.DB {
	if c.db != nil {
		return c.db
	}

	gormConfig := &gorm.Config{}
	if c.LogMode {
		gormConfig.Logger = gormlog.Default.LogMode(gormlog.Info)
	}

	dbCon, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               c.dsn(),
		DefaultStringSize: c.DefaultStringSize,
	}), gormConfig)
	if err != nil {
		logger.Fatalf("failed create connection to database: %v", err)
		return nil
	}

	sqlDB, _ := dbCon.DB()
	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.ConnMaxLifetime) * time.Minute)

	c.db = dbCon
	return c.db
}

var mapCfg = make(map[string]*config)

func Connect(name ...string) *gorm.DB {
	var prefix string
	if len(name) > 0 {
		prefix = strings.ToUpper(name[0])
	}
	if cfg, ok := mapCfg[prefix]; ok {
		return cfg.db
	}

	cfg := &config{}
	envconfig.MustProcess(prefix, cfg)
	if len(name) > 0 && !strings.EqualFold(cfg.Database, prefix) {
		logger.Panicf("database must be same with name: %s", name[0])
	}

	cfg.connect()
	mapCfg[prefix] = cfg
	return cfg.db
}

func Close(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		logger.Fatalf("failed close connection database: %v", err)
		return
	}
	if err = dbSQL.Close(); err != nil {
		logger.Fatalf("failed close connection database: %v", err)
	}
}
