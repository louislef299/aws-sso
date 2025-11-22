package provider

import (
	"errors"
	"fmt"
	"sync"
)

// Registry State
//
// The provider registry manages two distinct states for each provider:
//
//  1. Registered: The provider is known to the system (via Register). This
//     happens at init time when provider packages are imported. A registered
//     provider is available for inspection (schema, validation) but is not
//     yet configured for use.
//
//  2. Active: The provider has been configured and activated by the user
//     (via Activate). An active provider has passed configuration validation
//     and is ready for authentication operations.
//
// This separation allows the CLI to:
//   - List all available providers (registered)
//   - Show which providers the user has enabled (active)
//   - Detect missing configuration and prompt/fail helpfully before activation
//
// Example flow:
//
//	// At init time (automatic via import)
//	provider.Register(awsSSOProvider)
//
//	// At runtime, check what config is missing
//	missing, _ := provider.MissingConfig("aws-sso", userConfig)
//	if len(missing) > 0 {
//	    // Prompt user or print "run: knot config set ..."
//	}
//
//	// Once configured, activate for use
//	provider.Activate("aws-sso", userConfig)
//
//	// Now the provider is ready
//	for _, p := range provider.Active() {
//	    p.Initialize(ctx, config)
//	    p.Authenticate(ctx, opts)
//	}

var (
	// mu protects both providers and active maps for concurrent access.
	// All exported functions acquire this lock appropriately.
	mu sync.RWMutex

	// providers holds all registered providers, keyed by name.
	// Registration typically happens at init time via blank imports.
	// A provider being registered does NOT mean it's configured or ready to use.
	providers = make(map[string]Provider)

	// active tracks which providers the user has explicitly activated.
	// Only active providers should be used for authentication operations.
	// The map value is the validated config passed to Activate.
	active = make(map[string]map[string]any)
)

// Sentinel errors for registry operations. These enable callers to use
// errors.Is for specific error handling without string matching.
var (
	// ErrProviderNotFound is returned when attempting to access a provider
	// that has not been registered.
	ErrProviderNotFound = errors.New("provider: not found")

	// ErrProviderAlreadyRegistered is returned when attempting to register
	// a provider with a name that is already in use.
	ErrProviderAlreadyRegistered = errors.New("provider: already registered")

	// ErrProviderNotActive is returned when attempting to use a provider
	// that is registered but has not been activated.
	ErrProviderNotActive = errors.New("provider: not active")

	// ErrProviderAlreadyActive is returned when attempting to activate a
	// provider that is already active. Use Deactivate first to reconfigure.
	ErrProviderAlreadyActive = errors.New("provider: already active")

	// ErrMissingRequiredConfig is returned when activation fails due to
	// missing required configuration fields.
	ErrMissingRequiredConfig = errors.New("provider: missing required configuration")
)

// Register adds a provider to the registry, making it available for
// activation. This is typically called from init functions in provider
// packages, triggered by blank imports.
//
// Registration does NOT activate the provider - it only makes it known
// to the system. Users must call Activate with valid configuration to
// enable a provider for authentication.
//
// Returns ErrProviderAlreadyRegistered if a provider with the same name
// is already registered. This typically indicates a configuration error
// (duplicate imports or naming collision).
//
// Example:
//
//	func init() {
//	    provider.Register(&AWSProvider{})
//	}
func Register(p Provider) error {
	mu.Lock()
	defer mu.Unlock()
	name := p.Name()
	if _, exists := providers[name]; exists {
		return ErrProviderAlreadyRegistered
	}
	providers[name] = p
	return nil
}

// Get retrieves a registered provider by name. The provider may or may not
// be active - use IsActive to check activation status.
//
// This is useful for inspecting provider capabilities (schema, type) before
// activation, or for accessing an active provider for authentication.
//
// Returns ErrProviderNotFound if no provider with the given name is registered.
func Get(name string) (Provider, error) {
	mu.RLock()
	defer mu.RUnlock()
	p, ok := providers[name]
	if !ok {
		return nil, ErrProviderNotFound
	}
	return p, nil
}

// All returns all registered providers, regardless of activation status.
// Use Active to get only providers that are configured and ready to use.
//
// The returned slice is a copy; modifications do not affect the registry.
func All() []Provider {
	mu.RLock()
	defer mu.RUnlock()
	result := make([]Provider, 0, len(providers))
	for _, p := range providers {
		result = append(result, p)
	}
	return result
}

// Activate marks a provider as user-enabled and validates its configuration.
// After successful activation, the provider is ready for Initialize and
// subsequent authentication operations.
//
// Activation validates the provided config against the provider's schema:
//   - All required fields must be present
//   - Field types must match the schema
//   - Provider-specific validation rules must pass
//
// If validation fails, the provider remains inactive and an error is returned.
// Use MissingConfig to determine what fields are missing before calling Activate.
//
// Returns:
//   - ErrProviderNotFound if the provider is not registered
//   - ErrProviderAlreadyActive if the provider is already active
//   - ErrMissingRequiredConfig if required fields are missing (wrapped with details)
//   - Provider-specific validation errors from ValidateConfig
//
// Example:
//
//	config := map[string]any{
//	    "sso_start_url": "https://my-sso.awsapps.com/start",
//	    "sso_region":    "us-east-1",
//	}
//	if err := provider.Activate("aws-sso", config); err != nil {
//	    // Handle missing config, validation errors, etc.
//	}
func Activate(name string, config map[string]any) error {
	mu.Lock()
	defer mu.Unlock()

	// Verify provider is registered
	p, ok := providers[name]
	if !ok {
		return ErrProviderNotFound
	}

	// Check if already active
	if _, isActive := active[name]; isActive {
		return ErrProviderAlreadyActive
	}

	// Check for missing required configuration
	missing := getMissingConfigLocked(p, config)
	if len(missing) > 0 {
		// Build a helpful error message listing missing fields
		names := make([]string, len(missing))
		for i, f := range missing {
			names[i] = f.Name
		}
		return fmt.Errorf("%w: %v", ErrMissingRequiredConfig, names)
	}

	// Run provider-specific validation
	if err := p.ValidateConfig(config); err != nil {
		return fmt.Errorf("provider %s: %w", name, err)
	}

	// Store config and mark as active
	active[name] = config
	return nil
}

// Deactivate removes a provider from the active set. The provider remains
// registered and can be reactivated with new configuration.
//
// This is useful when:
//   - User wants to disable a provider without uninstalling
//   - Configuration needs to be changed (deactivate, then reactivate)
//   - Cleaning up during logout or reset operations
//
// Deactivate is idempotent - calling it on an inactive provider is a no-op.
//
// Returns ErrProviderNotFound if the provider is not registered.
func Deactivate(name string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := providers[name]; !ok {
		return ErrProviderNotFound
	}

	delete(active, name)
	return nil
}

// IsActive returns whether the named provider is currently active.
// Returns false if the provider is not registered or not active.
func IsActive(name string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, isActive := active[name]
	return isActive
}

// Active returns all providers that are currently active and ready for use.
// These providers have passed configuration validation via Activate.
//
// The returned slice is a copy; modifications do not affect the registry.
func Active() []Provider {
	mu.RLock()
	defer mu.RUnlock()
	result := make([]Provider, 0, len(active))
	for name := range active {
		if p, ok := providers[name]; ok {
			result = append(result, p)
		}
	}
	return result
}

// GetActiveConfig returns the configuration that was used to activate the
// named provider. This is useful for passing to Initialize.
//
// Returns ErrProviderNotFound if not registered, ErrProviderNotActive if
// registered but not active.
func GetActiveConfig(name string) (map[string]any, error) {
	mu.RLock()
	defer mu.RUnlock()

	if _, ok := providers[name]; !ok {
		return nil, ErrProviderNotFound
	}

	config, isActive := active[name]
	if !isActive {
		return nil, ErrProviderNotActive
	}

	return config, nil
}

// MissingConfig returns the required configuration fields that are not
// present in the provided config map. This enables the CLI to:
//
//  1. Detect missing configuration before activation fails
//  2. Generate interactive prompts for missing fields
//  3. Print helpful commands like "run: knot config set <field> <value>"
//
// The returned ConfigFields include descriptions and types, providing
// everything needed to prompt the user or generate help text.
//
// Returns an empty slice if all required fields are present.
// Returns ErrProviderNotFound if the provider is not registered.
//
// Example:
//
//	missing, err := provider.MissingConfig("aws-sso", userConfig)
//	if err != nil {
//	    return err
//	}
//	for _, field := range missing {
//	    if interactive {
//	        value := prompt("Enter %s (%s): ", field.Name, field.Description)
//	        userConfig[field.Name] = value
//	    } else {
//	        fmt.Printf("Run: knot config set %s.%s <value>\n", name, field.Name)
//	        fmt.Printf("  %s\n", field.Description)
//	    }
//	}
func MissingConfig(name string, config map[string]any) ([]ConfigField, error) {
	mu.RLock()
	defer mu.RUnlock()

	p, ok := providers[name]
	if !ok {
		return nil, ErrProviderNotFound
	}

	return getMissingConfigLocked(p, config), nil
}

// getMissingConfigLocked returns required fields not present in config.
// Caller must hold mu (read or write lock).
func getMissingConfigLocked(p Provider, config map[string]any) []ConfigField {
	schema := p.GetConfigSchema()
	var missing []ConfigField

	for _, field := range schema.Fields {
		if !field.Required {
			continue
		}
		if _, exists := config[field.Name]; !exists {
			missing = append(missing, field)
		}
	}

	return missing
}
