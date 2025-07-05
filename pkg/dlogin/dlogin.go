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

// ILogin defines the interface that all dynamic login plugins must implement
// to be compatible with the dlogin plugin system.
//
// Each plugin is responsible for registering itself (typically in an init
// function) and providing implementations for initialization, login, and logout
// behaviors.
//
// The interface is designed to be flexible and support a variety of
// authentication
// mechanisms (e.g., AWS, Kubernetes, custom SSO).
//
// Usage:
//   - Init is called to allow the plugin to register CLI flags and perform setup.
//   - Login is called to perform authentication and acquire credentials.
//   - Logout is called to clean up or revoke credentials.
//
// The config parameter for Login and Logout is plugin-specific and should be
// documented by the plugin implementation.
type ILogin interface {
	// Init allows the plugin to register CLI flags or perform setup on the
	// provided command.
	Init(cmd *cobra.Command) error

	// Login performs the authentication logic for the plugin.
	// The config parameter is plugin-specific and may contain credentials or
	// options.
	Login(ctx context.Context, config any) error

	// Logout performs any necessary cleanup or credential revocation.
	// The config parameter is plugin-specific and may contain credentials or
	// options.
	Logout(ctx context.Context, config any) error
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

func DLogin(ctx context.Context, driverName string, config any) error {
	driversMu.RLock()
	driverLogin, ok := drivers[driverName]
	driversMu.RUnlock()
	if !ok {
		return fmt.Errorf("login: unknown driver %q (forgotten import?)", driverName)
	}

	return driverLogin.Login(ctx, config)
}

func DLogout(ctx context.Context, driverName string, config any) error {
	driversMu.RLock()
	driverLogin, ok := drivers[driverName]
	driversMu.RUnlock()
	if !ok {
		return fmt.Errorf("login: unknown driver %q (forgotten import?)", driverName)
	}

	return driverLogin.Logout(ctx, config)
}
