package jenv

import (
	"os"
	"testing"
)

var (
	inLocal       bool
	inDevelopment bool
	inStaging     bool
	inDemo        bool
	inProduction  bool
)

type Environment string

const (
	Local       Environment = "local"
	Development Environment = "development"
	Staging     Environment = "staging"
	Demo        Environment = "demo"
	Production  Environment = "production"
)

func init() {
	env := os.Getenv("JUMP_ENVIRONMENT")
	switch env {
	case string(Development):
		inDevelopment = true
	case string(Staging):
		inStaging = true
	case string(Demo):
		inDemo = true
	case string(Production):
		inProduction = true
	default:
		inLocal = true
	}
}

func Get() Environment {
	switch {
	case inDevelopment:
		return Development
	case inStaging:
		return Staging
	case inDemo:
		return Demo
	case inProduction:
		return Production
	default:
		return Local
	}
}

func Name() string {
	return os.Getenv("JUMP_ENVIRONMENT")
}

func InLocal() bool {
	return inLocal
}

func InDevelopment() bool {
	return inDevelopment
}

func InStaging() bool {
	return inStaging
}

func InDemo() bool {
	return inDemo
}

func InProd() bool {
	return inProduction
}

func InTest() bool {
	return testing.Testing()
}
