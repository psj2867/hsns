package config_test

import (
	"testing"

	"github.com/psj2867/hsns/config"
)

func TestDb(t *testing.T) {
	config.SetTestDb("sqlite", "../"+config.DbConfig)
}
