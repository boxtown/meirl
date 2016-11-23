package main

import (
	"flag"
	"os"
	"strings"
)

type environment int

func (e environment) String() string {
	switch e {
	case prod:
		return "prod"
	default:
		return "dev"
	}
}

const (
	dev environment = iota
	prod
)

const signingKeyVar = "MEIRL_KEY"
const pgUserVar = "MEIRL_PG_USER"
const pgPassVar = "MERIRL_PG_PASS"

var appEnv environment
var signingKey []byte
var pgUser string
var pgPass string
var pgHost string
var pgPort string

const pgDBName = "meirldb"

func init() {
	appEnv = loadAppEnvironment()
	signingKey = loadSigningKey(appEnv)
	pgUser, pgPass = loadPostgresCredentials(appEnv)
	pgHost, pgPort = loadPostgresHostAndPort(appEnv)
}

func loadAppEnvironment() environment {
	env := strings.ToLower(*flag.String("env", "dev", "application environment type"))
	switch env {
	case "prod":
		return prod
	default:
		return dev
	}
}

func loadSigningKey(env environment) []byte {
	switch env {
	case prod:
		return []byte(os.Getenv(signingKeyVar))
	default:
		return []byte("dev-key")
	}
}

func loadPostgresCredentials(env environment) (string, string) {
	switch env {
	case prod:
		return os.Getenv(pgUserVar), os.Getenv(pgPassVar)
	default:
		return "tester", "test"
	}
}

func loadPostgresHostAndPort(env environment) (string, string) {
	switch env {
	case prod:
		return "localhost", "5432"
	default:
		return "localhost", "5432"
	}
}

func debug() bool {
	return appEnv == dev
}
