package models_test

import (
	"os"
	"testing"

	"github.com/psj2867/hsns/config"
)

func TestMain(m *testing.M) {
	config.SetTestDb("ramsql", "test")
	defer config.DeferDb()
	os.Exit(m.Run())
}
