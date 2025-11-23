package provider

import (
	"context"
	"errors"
	"testing"
)

// mockProvider implements the Provider interface for testing.
type mockProvider struct {
	name            string
	providerType    Type
	schema          ConfigSchema
	validateErr     error
	initializeErr   error
	initializeCalls int
}

func (m *mockProvider) Name() string { return m.name }
func (m *mockProvider) Type() Type   { return m.providerType }

func (m *mockProvider) Initialize(ctx context.Context, config map[string]any) error {
	m.initializeCalls++
	return m.initializeErr
}

func (m *mockProvider) Authenticate(ctx context.Context, opts AuthOptions) (*Credentials, error) {
	return nil, nil
}

func (m *mockProvider) Refresh(ctx context.Context, creds *Credentials) (*Credentials, error) {
	return nil, nil
}

func (m *mockProvider) Revoke(ctx context.Context, creds *Credentials) error {
	return nil
}

func (m *mockProvider) GetConfigSchema() ConfigSchema {
	return m.schema
}

func (m *mockProvider) ValidateConfig(config map[string]any) error {
	return m.validateErr
}

func newMockProvider(name string) *mockProvider {
	return &mockProvider{
		name:         name,
		providerType: TypeOIDC,
		schema: ConfigSchema{
			Fields: []ConfigField{
				{Name: "required_field", Type: "string", Required: true},
				{Name: "optional_field", Type: "string", Required: false},
			},
		},
	}
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name        string
		providers   []string
		wantErr     error
		description string
	}{
		{
			name:        "register single provider",
			providers:   []string{"test-provider"},
			wantErr:     nil,
			description: "should successfully register a new provider",
		},
		{
			name:        "register duplicate provider",
			providers:   []string{"dup-provider", "dup-provider"},
			wantErr:     ErrProviderAlreadyRegistered,
			description: "should return error when registering duplicate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reset()

			var lastErr error
			for _, name := range tt.providers {
				lastErr = Register(newMockProvider(name))
			}

			if !errors.Is(lastErr, tt.wantErr) {
				t.Errorf("Register() error = %v, wantErr %v", lastErr, tt.wantErr)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name         string
		registerName string
		getName      string
		wantErr      error
	}{
		{
			name:         "get existing provider",
			registerName: "existing",
			getName:      "existing",
			wantErr:      nil,
		},
		{
			name:         "get non-existent provider",
			registerName: "other",
			getName:      "missing",
			wantErr:      ErrProviderNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reset()
			Register(newMockProvider(tt.registerName))

			p, err := Get(tt.getName)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && p == nil {
				t.Error("Get() returned nil provider when expecting valid provider")
			}
			if tt.wantErr == nil && p.Name() != tt.getName {
				t.Errorf("Get() returned provider with name %q, want %q", p.Name(), tt.getName)
			}
		})
	}
}

func TestAll(t *testing.T) {
	tests := []struct {
		name      string
		providers []string
		wantCount int
	}{
		{
			name:      "empty registry",
			providers: nil,
			wantCount: 0,
		},
		{
			name:      "single provider",
			providers: []string{"p1"},
			wantCount: 1,
		},
		{
			name:      "multiple providers",
			providers: []string{"p1", "p2", "p3"},
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reset()
			for _, name := range tt.providers {
				Register(newMockProvider(name))
			}

			got := All()

			if len(got) != tt.wantCount {
				t.Errorf("All() returned %d providers, want %d", len(got), tt.wantCount)
			}
		})
	}
}

func TestActivate(t *testing.T) {
	tests := []struct {
		name          string
		providerName  string
		config        map[string]any
		setupProvider func() *mockProvider
		wantErr       error
	}{
		{
			name:         "successful activation",
			providerName: "test",
			config:       map[string]any{"required_field": "value"},
			setupProvider: func() *mockProvider {
				return newMockProvider("test")
			},
			wantErr: nil,
		},
		{
			name:         "provider not found",
			providerName: "missing",
			config:       map[string]any{},
			setupProvider: func() *mockProvider {
				return newMockProvider("other")
			},
			wantErr: ErrProviderNotFound,
		},
		{
			name:         "missing required config",
			providerName: "test",
			config:       map[string]any{},
			setupProvider: func() *mockProvider {
				return newMockProvider("test")
			},
			wantErr: ErrMissingRequiredConfig,
		},
		{
			name:         "validation error",
			providerName: "test",
			config:       map[string]any{"required_field": "value"},
			setupProvider: func() *mockProvider {
				p := newMockProvider("test")
				p.validateErr = errors.New("invalid config value")
				return p
			},
			wantErr: nil, // wrapped error, check with Contains
		},
		{
			name:         "initialization error",
			providerName: "test",
			config:       map[string]any{"required_field": "value"},
			setupProvider: func() *mockProvider {
				p := newMockProvider("test")
				p.initializeErr = errors.New("init failed")
				return p
			},
			wantErr: ErrInitializationFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reset()
			p := tt.setupProvider()
			Register(p)

			err := Activate(context.Background(), tt.providerName, tt.config)

			if tt.name == "validation error" {
				// Special case: validation error is wrapped differently
				if err == nil {
					t.Error("Activate() expected error for validation failure")
				}
				return
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Activate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestActivateAlreadyActive(t *testing.T) {
	Reset()
	p := newMockProvider("test")
	Register(p)

	config := map[string]any{"required_field": "value"}
	Activate(context.Background(), "test", config)

	err := Activate(context.Background(), "test", config)

	if !errors.Is(err, ErrProviderAlreadyActive) {
		t.Errorf("Activate() on already active provider error = %v, want %v", err, ErrProviderAlreadyActive)
	}
}

func TestDeactivate(t *testing.T) {
	tests := []struct {
		name         string
		providerName string
		setup        func()
		wantErr      error
	}{
		{
			name:         "deactivate active provider",
			providerName: "test",
			setup: func() {
				p := newMockProvider("test")
				Register(p)
				Activate(context.Background(), "test", map[string]any{"required_field": "value"})
			},
			wantErr: nil,
		},
		{
			name:         "deactivate registered provider",
			providerName: "test",
			setup: func() {
				Register(newMockProvider("test"))
			},
			wantErr: nil,
		},
		{
			name:         "deactivate non-existent provider",
			providerName: "missing",
			setup:        func() {},
			wantErr:      ErrProviderNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reset()
			tt.setup()

			err := Deactivate(tt.providerName)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Deactivate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && IsActive(tt.providerName) {
				t.Error("Deactivate() provider is still active")
			}
		})
	}
}

func TestIsActive(t *testing.T) {
	tests := []struct {
		name         string
		providerName string
		setup        func()
		want         bool
	}{
		{
			name:         "active provider",
			providerName: "test",
			setup: func() {
				p := newMockProvider("test")
				Register(p)
				Activate(context.Background(), "test", map[string]any{"required_field": "value"})
			},
			want: true,
		},
		{
			name:         "registered but not active",
			providerName: "test",
			setup: func() {
				Register(newMockProvider("test"))
			},
			want: false,
		},
		{
			name:         "non-existent provider",
			providerName: "missing",
			setup:        func() {},
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reset()
			tt.setup()

			got := IsActive(tt.providerName)

			if got != tt.want {
				t.Errorf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInitialized(t *testing.T) {
	tests := []struct {
		name         string
		providerName string
		setup        func()
		want         bool
	}{
		{
			name:         "initialized provider",
			providerName: "test",
			setup: func() {
				p := newMockProvider("test")
				Register(p)
				Activate(context.Background(), "test", map[string]any{"required_field": "value"})
			},
			want: true,
		},
		{
			name:         "registered but not initialized",
			providerName: "test",
			setup: func() {
				Register(newMockProvider("test"))
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reset()
			tt.setup()

			got := IsInitialized(tt.providerName)

			if got != tt.want {
				t.Errorf("IsInitialized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInvalid(t *testing.T) {
	tests := []struct {
		name         string
		providerName string
		setup        func()
		want         bool
	}{
		{
			name:         "invalid provider",
			providerName: "test",
			setup: func() {
				p := newMockProvider("test")
				p.initializeErr = errors.New("init failed")
				Register(p)
				Activate(context.Background(), "test", map[string]any{"required_field": "value"})
			},
			want: true,
		},
		{
			name:         "valid active provider",
			providerName: "test",
			setup: func() {
				p := newMockProvider("test")
				Register(p)
				Activate(context.Background(), "test", map[string]any{"required_field": "value"})
			},
			want: false,
		},
		{
			name:         "registered provider",
			providerName: "test",
			setup: func() {
				Register(newMockProvider("test"))
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reset()
			tt.setup()

			got := IsInvalid(tt.providerName)

			if got != tt.want {
				t.Errorf("IsInvalid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInvalidReason(t *testing.T) {
	tests := []struct {
		name         string
		providerName string
		setup        func()
		wantReason   bool
		wantErr      error
	}{
		{
			name:         "invalid provider has reason",
			providerName: "test",
			setup: func() {
				p := newMockProvider("test")
				p.initializeErr = errors.New("init failed")
				Register(p)
				Activate(context.Background(), "test", map[string]any{"required_field": "value"})
			},
			wantReason: true,
			wantErr:    nil,
		},
		{
			name:         "valid provider has no reason",
			providerName: "test",
			setup: func() {
				p := newMockProvider("test")
				Register(p)
				Activate(context.Background(), "test", map[string]any{"required_field": "value"})
			},
			wantReason: false,
			wantErr:    nil,
		},
		{
			name:         "non-existent provider",
			providerName: "missing",
			setup:        func() {},
			wantReason:   false,
			wantErr:      ErrProviderNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reset()
			tt.setup()

			reason, err := GetInvalidReason(tt.providerName)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetInvalidReason() err = %v, wantErr %v", err, tt.wantErr)
			}

			hasReason := reason != nil
			if hasReason != tt.wantReason {
				t.Errorf("GetInvalidReason() hasReason = %v, want %v", hasReason, tt.wantReason)
			}
		})
	}
}

func TestStatus(t *testing.T) {
	tests := []struct {
		name         string
		providerName string
		setup        func()
		want         string
	}{
		{
			name:         "active provider",
			providerName: "test",
			setup: func() {
				p := newMockProvider("test")
				Register(p)
				Activate(context.Background(), "test", map[string]any{"required_field": "value"})
			},
			want: "active",
		},
		{
			name:         "invalid provider",
			providerName: "test",
			setup: func() {
				p := newMockProvider("test")
				p.initializeErr = errors.New("init failed")
				Register(p)
				Activate(context.Background(), "test", map[string]any{"required_field": "value"})
			},
			want: "invalid",
		},
		{
			name:         "registered provider",
			providerName: "test",
			setup: func() {
				Register(newMockProvider("test"))
			},
			want: "registered",
		},
		{
			name:         "non-existent provider",
			providerName: "missing",
			setup:        func() {},
			want:         "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reset()
			tt.setup()

			got := Status(tt.providerName)

			if got != tt.want {
				t.Errorf("Status() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestActiveProviders(t *testing.T) {
	tests := []struct {
		name      string
		setup     func()
		wantCount int
	}{
		{
			name:      "no active providers",
			setup:     func() {},
			wantCount: 0,
		},
		{
			name: "one active provider",
			setup: func() {
				p := newMockProvider("test")
				Register(p)
				Activate(context.Background(), "test", map[string]any{"required_field": "value"})
			},
			wantCount: 1,
		},
		{
			name: "multiple active providers",
			setup: func() {
				for _, name := range []string{"p1", "p2", "p3"} {
					p := newMockProvider(name)
					Register(p)
					Activate(context.Background(), name, map[string]any{"required_field": "value"})
				}
			},
			wantCount: 3,
		},
		{
			name: "mix of active and inactive",
			setup: func() {
				Register(newMockProvider("inactive"))
				p := newMockProvider("active")
				Register(p)
				Activate(context.Background(), "active", map[string]any{"required_field": "value"})
			},
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reset()
			tt.setup()

			got := ActiveProviders()

			if len(got) != tt.wantCount {
				t.Errorf("ActiveProviders() returned %d providers, want %d", len(got), tt.wantCount)
			}
		})
	}
}

func TestInvalidProviders(t *testing.T) {
	tests := []struct {
		name      string
		setup     func()
		wantCount int
	}{
		{
			name:      "no invalid providers",
			setup:     func() {},
			wantCount: 0,
		},
		{
			name: "one invalid provider",
			setup: func() {
				p := newMockProvider("test")
				p.initializeErr = errors.New("init failed")
				Register(p)
				Activate(context.Background(), "test", map[string]any{"required_field": "value"})
			},
			wantCount: 1,
		},
		{
			name: "mix of valid and invalid",
			setup: func() {
				validP := newMockProvider("valid")
				Register(validP)
				Activate(context.Background(), "valid", map[string]any{"required_field": "value"})

				invalidP := newMockProvider("invalid")
				invalidP.initializeErr = errors.New("init failed")
				Register(invalidP)
				Activate(context.Background(), "invalid", map[string]any{"required_field": "value"})
			},
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reset()
			tt.setup()

			got := InvalidProviders()

			if len(got) != tt.wantCount {
				t.Errorf("InvalidProviders() returned %d providers, want %d", len(got), tt.wantCount)
			}
		})
	}
}

func TestGetActiveConfig(t *testing.T) {
	tests := []struct {
		name         string
		providerName string
		setup        func() map[string]any
		wantErr      error
	}{
		{
			name:         "get config for active provider",
			providerName: "test",
			setup: func() map[string]any {
				config := map[string]any{"required_field": "value", "extra": "data"}
				p := newMockProvider("test")
				Register(p)
				Activate(context.Background(), "test", config)
				return config
			},
			wantErr: nil,
		},
		{
			name:         "get config for inactive provider",
			providerName: "test",
			setup: func() map[string]any {
				Register(newMockProvider("test"))
				return nil
			},
			wantErr: ErrProviderNotActive,
		},
		{
			name:         "get config for non-existent provider",
			providerName: "missing",
			setup: func() map[string]any {
				return nil
			},
			wantErr: ErrProviderNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reset()
			expectedConfig := tt.setup()

			config, err := GetActiveConfig(tt.providerName)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetActiveConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && expectedConfig != nil {
				if config["required_field"] != expectedConfig["required_field"] {
					t.Errorf("GetActiveConfig() config mismatch")
				}
			}
		})
	}
}

func TestReinitialize(t *testing.T) {
	tests := []struct {
		name         string
		providerName string
		newConfig    map[string]any
		setup        func() *mockProvider
		wantErr      error
	}{
		{
			name:         "successful reinitialize",
			providerName: "test",
			newConfig:    map[string]any{"required_field": "new_value"},
			setup: func() *mockProvider {
				p := newMockProvider("test")
				Register(p)
				Activate(context.Background(), "test", map[string]any{"required_field": "old_value"})
				return p
			},
			wantErr: nil,
		},
		{
			name:         "reinitialize non-existent provider",
			providerName: "missing",
			newConfig:    map[string]any{"required_field": "value"},
			setup: func() *mockProvider {
				return newMockProvider("test")
			},
			wantErr: ErrProviderNotFound,
		},
		{
			name:         "reinitialize inactive provider",
			providerName: "test",
			newConfig:    map[string]any{"required_field": "value"},
			setup: func() *mockProvider {
				p := newMockProvider("test")
				Register(p)
				return p
			},
			wantErr: ErrProviderNotActive,
		},
		{
			name:         "reinitialize with missing config",
			providerName: "test",
			newConfig:    map[string]any{},
			setup: func() *mockProvider {
				p := newMockProvider("test")
				Register(p)
				Activate(context.Background(), "test", map[string]any{"required_field": "value"})
				return p
			},
			wantErr: ErrMissingRequiredConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reset()
			tt.setup()

			err := Reinitialize(context.Background(), tt.providerName, tt.newConfig)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Reinitialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMissingConfig(t *testing.T) {
	tests := []struct {
		name         string
		providerName string
		config       map[string]any
		setup        func()
		wantMissing  int
		wantErr      error
	}{
		{
			name:         "no missing config",
			providerName: "test",
			config:       map[string]any{"required_field": "value"},
			setup: func() {
				Register(newMockProvider("test"))
			},
			wantMissing: 0,
			wantErr:     nil,
		},
		{
			name:         "missing required field",
			providerName: "test",
			config:       map[string]any{},
			setup: func() {
				Register(newMockProvider("test"))
			},
			wantMissing: 1,
			wantErr:     nil,
		},
		{
			name:         "provider not found",
			providerName: "missing",
			config:       map[string]any{},
			setup:        func() {},
			wantMissing:  0,
			wantErr:      ErrProviderNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reset()
			tt.setup()

			missing, err := MissingConfig(tt.providerName, tt.config)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("MissingConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(missing) != tt.wantMissing {
				t.Errorf("MissingConfig() returned %d missing fields, want %d", len(missing), tt.wantMissing)
			}
		})
	}
}

func TestReset(t *testing.T) {
	Reset()

	// Setup some state
	p := newMockProvider("test")
	Register(p)
	Activate(context.Background(), "test", map[string]any{"required_field": "value"})

	invalidP := newMockProvider("invalid")
	invalidP.initializeErr = errors.New("init failed")
	Register(invalidP)
	Activate(context.Background(), "invalid", map[string]any{"required_field": "value"})

	// Verify state exists
	if len(All()) != 2 {
		t.Fatal("setup failed: expected 2 providers")
	}

	// Reset
	Reset()

	// Verify all state is cleared
	if len(All()) != 0 {
		t.Errorf("Reset() did not clear providers, got %d", len(All()))
	}
	if len(ActiveProviders()) != 0 {
		t.Errorf("Reset() did not clear active providers, got %d", len(ActiveProviders()))
	}
	if len(InvalidProviders()) != 0 {
		t.Errorf("Reset() did not clear invalid providers, got %d", len(InvalidProviders()))
	}
}

func TestActivateAutoRetry(t *testing.T) {
	Reset()

	// First activation fails
	p := newMockProvider("test")
	p.initializeErr = errors.New("init failed")
	Register(p)

	err := Activate(context.Background(), "test", map[string]any{"required_field": "value"})
	if err == nil {
		t.Fatal("expected first activation to fail")
	}
	if !IsInvalid("test") {
		t.Error("provider should be invalid after failed activation")
	}

	// Fix the provider and retry - should auto-clear invalid state
	p.initializeErr = nil
	err = Activate(context.Background(), "test", map[string]any{"required_field": "value"})
	if err != nil {
		t.Errorf("Activate() retry error = %v, want nil", err)
	}
	if !IsActive("test") {
		t.Error("provider should be active after successful retry")
	}
	if IsInvalid("test") {
		t.Error("provider should not be invalid after successful retry")
	}
}

func TestDeactivateClearsInvalid(t *testing.T) {
	Reset()

	p := newMockProvider("test")
	p.initializeErr = errors.New("init failed")
	Register(p)

	// Activate fails, provider becomes invalid
	Activate(context.Background(), "test", map[string]any{"required_field": "value"})
	if !IsInvalid("test") {
		t.Fatal("provider should be invalid")
	}

	// Deactivate clears invalid state
	Deactivate("test")

	if IsInvalid("test") {
		t.Error("Deactivate() should clear invalid state")
	}
	if Status("test") != "registered" {
		t.Errorf("Status() = %q, want %q", Status("test"), "registered")
	}
}
