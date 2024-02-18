//go:build integration

package postgresdb

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	testDB *pgxpool.Pool
)

func TestMain(m *testing.M) {
	const (
		dbPassword = "tests"
		dbUser     = "integration"
		dbName     = "postgres"
	)

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	databaseURL := os.Getenv("POSTGRES_TEST_DSN")

	var resource *dockertest.Resource
	if databaseURL == "" {
		// pulls an image, creates a container based on it and runs it
		resource, err = pool.RunWithOptions(&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "16",
			Env: []string{
				"POSTGRES_PASSWORD=" + dbPassword,
				"POSTGRES_USER=" + dbUser,
				"POSTGRES_DB=" + dbName,
				"listen_addresses = '*'",
			},
		}, func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		})
		if err != nil {
			log.Fatalf("could not start resource: %s", err)
		}

		hostAndPort := resource.GetHostPort("5432/tcp")
		databaseURL = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPassword, hostAndPort, dbName)
		err = resource.Expire(300) // Tell docker to hard kill the container in 300 seconds
		if err != nil {
			log.Printf("could not set hard kill time for container: %s", err)
		}
	}

	dbConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("failed parsing database config: %s", err)
	}

	log.Printf("connecting to database on url: %s", databaseURL)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second

	if err = pool.Retry(func() error {
		mainDB, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
		if err != nil {
			return err
		}

		dbName := fmt.Sprintf("testdb%d", rand.Intn(10000))
		_, err = mainDB.Exec(context.Background(), `CREATE DATABASE `+dbName)
		if err != nil {
			return err
		}
		mainDB.Close()

		dbConfig.ConnConfig.Database = dbName
		testDB, err = pgxpool.NewWithConfig(context.Background(), dbConfig)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	log.Printf("starting integration tests: %s", dbConfig.ConnConfig.Database)

	initDB()

	//Run tests
	code := m.Run()

	testDB.Close()

	if os.Getenv("POSTGRES_TEST_DSN") == "" {
		// You can't defer this because os.Exit doesn't care for defer
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("could not purge resource: %s", err)
		}
	}

	os.Exit(code)
}

func initDB() {
	files := []string{
		"../../../dev/database/1_init.sql",
	}

	for _, file := range files {
		initQuery, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("error reading db fixtures file(%s): %s", file, err)
		}

		_, err = testDB.Exec(context.Background(), string(initQuery))
		if err != nil {
			log.Fatalf("error initializing database(%s): %s", file, err)
		}
	}
}

func truncate() {
	q := `
		DELETE FROM accounts;
	`
	_, err := testDB.Exec(context.Background(), q)
	if err != nil {
		log.Fatalf("error truncating database: %s", err)
	}
}
