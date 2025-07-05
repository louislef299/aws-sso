package dlogin

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"sync"

	"github.com/spf13/cobra"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]ILogin)
)

// Register makes a login driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, driver ILogin) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("login: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("login: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

// Drivers returns a sorted list of the names of the registered drivers.
func Drivers() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()
	return slices.Sorted(maps.Keys(drivers))
}

type Config struct {
	Role   string
	Secret string

	Cluster string
}

type ConfigOptions any
type ConfigOptionsFunc func(*ConfigOptions) error
type ILogin interface {
	Init(cmd *cobra.Command) error
	Login(ctx context.Context, config any, opts ...ConfigOptionsFunc) error
}

func Init(driverName string, cmd *cobra.Command) error {
	driversMu.RLock()
	driverLogin, ok := drivers[driverName]
	driversMu.RUnlock()
	if !ok {
		return fmt.Errorf("login: unknown driver %q (forgotten import?)", driverName)
	}

	return driverLogin.Init(cmd)
}

func DLogin(driverName string, config any) error {
	driversMu.RLock()
	driverLogin, ok := drivers[driverName]
	driversMu.RUnlock()
	if !ok {
		return fmt.Errorf("login: unknown driver %q (forgotten import?)", driverName)
	}

	return driverLogin.Login(context.TODO(), config)
}
