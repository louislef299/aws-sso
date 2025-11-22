package provider

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Registry State
//
// The provider registry manages three distinct states for each provider:
//
//  1. Registered: The provider is known to the system (via Register). This
//     happens at init time when provider packages are imported. A registered
//     provider is available for inspection (schema, validation) but is not
//     yet configured for use.
//
//  2. Active: The provider has valid configuration that has been validated
//     (via Activate). This is a transient state during activation.
//
//  3. Initialized: The provider has been fully initialized and is ready for
//     authentication operations. This is the final state after successful
//     Activate, which validates config AND calls provider.Initialize().
//
// State transitions:
//
//	              ┌─────────────────────────────────────────┐
//	              │                                         │
//	              ▼                                         │
//	┌──────────────────┐    Activate()    ┌─────────────────────┐
//	│    Registered    │ ───────────────▶ │ Active + Initialized │
//	└──────────────────┘                  └─────────────────────┘
//	              ▲                                         │
//	              │         Deactivate()                    │
//	              └─────────────────────────────────────────┘
//
// Note: If Initialize() fails during Activate(), the provider is rolled back
// to Registered state. A provider is never left in an "active but not
// initialized" state, ensuring callers can safely use any active provider.
//
// This separation allows the CLI to:
//   - List all available providers (registered)
//   - Show which providers the user has enabled (active + initialized)
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
//	// Activate validates config AND calls Initialize - one step
//	err := provider.Activate(ctx, "aws-sso", userConfig)
//	if err != nil {
//	    // Config validation failed OR Initialize failed
//	    // Provider remains in Registered state, not broken
//	}
//
//	// Now the provider is ready for authentication
//	for _, p := range provider.ActiveProviders() {
//	    p.Authenticate(ctx, opts)  // Safe - Initialize already called
//	}

var (
	// mu protects providers, active, and initialized maps for concurrent access.
	// All exported functions acquire this lock appropriately.
	mu sync.RWMutex

	// providers holds all registered providers, keyed by name.
	// Registration typically happens at init time via blank imports.
	// A provider being registered does NOT mean it's configured or ready to use.
	providers = make(map[string]Provider)

	// active tracks which providers have validated configuration.
	// The map value is the validated config passed to Activate.
	// Note: A provider in active should also be in initialized (see below).
	active = make(map[string]map[string]any)

	// initialized tracks which providers have had Initialize() called successfully.
	// This is set atomically with active during Activate() - a provider is never
	// left in active without also being initialized.
	initialized = make(map[string]bool)
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

	// ErrProviderNotInitialized is returned when attempting to use a provider
	// that has not been initialized. This should not occur in normal use since
	// Activate() calls Initialize(), but may occur if Initialize() failed.
	ErrProviderNotInitialized = errors.New("provider: not initialized")

	// ErrMissingRequiredConfig is returned when activation fails due to
	// missing required configuration fields.
	ErrMissingRequiredConfig = errors.New("provider: missing required configuration")

	// ErrInitializationFailed is returned when provider.Initialize() fails
	// during Activate(). The underlying error is wrapped for inspection.
	ErrInitializationFailed = errors.New("provider: initialization failed")
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
// be active - use IsActive to check activation status, or IsInitialized to
// verify it's ready for use.
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
// Use ActiveProviders to get only providers that are configured and ready to use.
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

// Activate marks a provider as user-enabled, validates its configuration,
// and calls Initialize to prepare it for use. This is the single entry point
// for enabling a provider - after successful Activate, the provider is fully
// ready for Authenticate, Refresh, and Revoke operations.
//
// Activate performs the following steps atomically:
//  1. Validates all required configuration fields are present
//  2. Calls provider.ValidateConfig() for provider-specific validation
//  3. Calls provider.Initialize(ctx, config) to set up the provider
//  4. Marks the provider as active and initialized
//
// If any step fails, the provider remains in Registered state (not Active).
// This ensures a provider is never left in a broken "active but not initialized"
// state.
//
// Use MissingConfig to determine what fields are missing before calling Activate.
// Use Deactivate to disable a provider and allow reconfiguration.
//
// Returns:
//   - ErrProviderNotFound if the provider is not registered
//   - ErrProviderAlreadyActive if the provider is already active
//   - ErrMissingRequiredConfig if required fields are missing (wrapped with details)
//   - ErrInitializationFailed if provider.Initialize() fails (wrapped with cause)
//   - Provider-specific validation errors from ValidateConfig
//
// Example:
//
//	config := map[string]any{
//	    "sso_start_url": "https://my-sso.awsapps.com/start",
//	    "sso_region":    "us-east-1",
//	}
//	if err := provider.Activate(ctx, "aws-sso", config); err != nil {
//	    if errors.Is(err, provider.ErrMissingRequiredConfig) {
//	        // Prompt user for missing fields
//	    } else if errors.Is(err, provider.ErrInitializationFailed) {
//	        // Provider setup failed (network, credentials, etc.)
//	    }
//	}
func Activate(ctx context.Context, name string, config map[string]any) error {
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

	// Call Initialize - this is where the provider sets up connections,
	// loads credentials, etc. If this fails, we do NOT mark as active.
	if err := p.Initialize(ctx, config); err != nil {
		return fmt.Errorf("%w: %s: %w", ErrInitializationFailed, name, err)
	}

	// Success - mark as both active and initialized atomically
	active[name] = config
	initialized[name] = true
	return nil
}

// Deactivate removes a provider from the active and initialized sets. The
// provider remains registered and can be reactivated with new configuration.
//
// This is useful when:
//   - User wants to disable a provider without uninstalling
//   - Configuration needs to be changed (deactivate, then reactivate)
//   - Cleaning up during logout or reset operations
//
// Note: Deactivate does NOT call provider.Revoke() - that is for credential
// cleanup and should be called separately before Deactivate if needed.
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
	delete(initialized, name)
	return nil
}

// IsActive returns whether the named provider is currently active.
// An active provider has valid configuration and has been initialized.
// Returns false if the provider is not registered or not active.
func IsActive(name string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, isActive := active[name]
	return isActive
}

// IsInitialized returns whether the named provider has been successfully
// initialized. In normal operation, this is equivalent to IsActive since
// Activate() calls Initialize(). This function is primarily useful for
// debugging and testing.
//
// Returns false if the provider is not registered or not initialized.
func IsInitialized(name string) bool {
	mu.RLock()
	defer mu.RUnlock()
	return initialized[name]
}

// ActiveProviders returns all providers that are currently active and
// initialized, ready for authentication operations.
//
// The returned slice is a copy; modifications do not affect the registry.
func ActiveProviders() []Provider {
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
// named provider.
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

// Reinitialize re-runs Initialize on an already-active provider with new
// or updated configuration. This is useful when configuration changes without
// fully deactivating the provider.
//
// Unlike Activate, this:
//   - Requires the provider to already be active
//   - Updates the stored config on success
//   - Leaves the provider active (with old config) on Initialize failure
//
// If you need to change configuration and want atomic behavior (either fully
// updated or fully rolled back), use Deactivate followed by Activate instead.
//
// Returns:
//   - ErrProviderNotFound if the provider is not registered
//   - ErrProviderNotActive if the provider is not currently active
//   - ErrMissingRequiredConfig if required fields are missing
//   - ErrInitializationFailed if provider.Initialize() fails
func Reinitialize(ctx context.Context, name string, config map[string]any) error {
	mu.Lock()
	defer mu.Unlock()

	p, ok := providers[name]
	if !ok {
		return ErrProviderNotFound
	}

	if _, isActive := active[name]; !isActive {
		return ErrProviderNotActive
	}

	// Validate new config
	missing := getMissingConfigLocked(p, config)
	if len(missing) > 0 {
		names := make([]string, len(missing))
		for i, f := range missing {
			names[i] = f.Name
		}
		return fmt.Errorf("%w: %v", ErrMissingRequiredConfig, names)
	}

	if err := p.ValidateConfig(config); err != nil {
		return fmt.Errorf("provider %s: %w", name, err)
	}

	// Re-initialize with new config
	if err := p.Initialize(ctx, config); err != nil {
		return fmt.Errorf("%w: %s: %w", ErrInitializationFailed, name, err)
	}

	// Update stored config
	active[name] = config
	return nil
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

// Reset clears all registry state, removing all registered, active, and
// initialized providers. This is primarily useful for testing.
//
// WARNING: This will break any code holding references to providers
// obtained before the reset. Use with caution.
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	providers = make(map[string]Provider)
	active = make(map[string]map[string]any)
	initialized = make(map[string]bool)
}
