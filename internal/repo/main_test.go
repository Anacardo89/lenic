package repo

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var (
	TestDSN  string
	SeedPath string
	Repo     DBRepo
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	var (
		closeDB func()
		err     error
	)
	repo, testDSN, closeDB, seedPath, err := BuildTestDBEnv(ctx)
	if err != nil {
		fmt.Println("Failed to start test environment:", err)
		os.Exit(1)
	}
	defer closeDB()
	TestDSN = testDSN
	SeedPath = filepath.Join(seedPath, "repo_tests.sql")
	Repo = repo
	code := m.Run()
	os.Exit(code)
}
