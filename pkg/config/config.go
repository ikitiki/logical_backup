package config

import (
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx"
	"gopkg.in/yaml.v2"
)

type Config struct {
	TempDir                                string         `yaml:"tempDir"`
	Tables                                 []string       `yaml:"tables"`
	DB                                     pgx.ConnConfig `yaml:"db"`
	Slotname                               string         `yaml:"slotname"`
	PublicationName                        string         `yaml:"publication"`
	TrackNewTables                         bool           `yaml:"trackNewTables"`
	DeltasPerFile                          int            `yaml:"deltasPerFile"`
	BackupThreshold                        int            `yaml:"backupThreshold"`
	ConcurrentBasebackups                  int            `yaml:"concurrentBasebackups"`
	InitialBasebackup                      bool           `yaml:"initialBasebackup"`
	SendStatusOnCommit                     bool           `yaml:"sendStatusOnCommit"`
	Fsync                                  bool           `yaml:"fsync"`
	ArchiveDir                             string         `yaml:"archiveDir"`
	ForceBasebackupAfterInactivityInterval time.Duration  `yaml:"forceBasebackupAfterInactivityInterval"`
	ArchiverTimeout                        time.Duration  `yaml:"archiverTimeout"`
}

func New(filename string) (*Config, error) {
	var cfg Config

	configFp, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %v", err)
	}

	if err := yaml.NewDecoder(configFp).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("could not decode config file: %v", err)
	}
	// forcing backups with sub-minute inactivity period makes no sense.
	cfg.ForceBasebackupAfterInactivityInterval = cfg.ForceBasebackupAfterInactivityInterval.Truncate(1 * time.Minute)

	return &cfg, nil
}
