package version

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestFindValidVersion(t *testing.T) {
	tests := []struct {
		name     string
		versions []string
		want     string
		wantErr  bool
	}{
		{
			name:     "valid semver found",
			versions: []string{"v1.2.3", "Release 1.2.3"},
			want:     "v1.2.3",
			wantErr:  false,
		},
		{
			name:     "valid semver with prefix",
			versions: []string{"version-1.2.3", "v1.2.3"},
			want:     "v1.2.3",
			wantErr:  false,
		},
		{
			name:     "no valid semver",
			versions: []string{"version-abc", "release-xyz"},
			want:     "",
			wantErr:  true,
		},
		{
			name:     "empty input",
			versions: []string{},
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findValidVersion(tt.versions...)
			if (err != nil) != tt.wantErr {
				t.Errorf("findValidVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("findValidVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckForUpdate(t *testing.T) {
	// Save original version and HTTP client
	originalVersion := Version
	originalClient := http.DefaultClient
	defer func() {
		Version = originalVersion
		http.DefaultClient = originalClient
	}()

	tests := []struct {
		name           string
		currentVersion string
		mockResponse   *latestVersion
		mockStatusCode int
		wantErr        bool
	}{
		{
			name:           "newer version available",
			currentVersion: "1.0.0",
			mockResponse: &latestVersion{
				Name:    "Release 2.0.0",
				TagName: "v2.0.0",
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "current version up to date",
			currentVersion: "2.0.0",
			mockResponse: &latestVersion{
				Name:    "Release 2.0.0",
				TagName: "v2.0.0",
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "invalid version format",
			currentVersion: "1.0.0",
			mockResponse: &latestVersion{
				Name:    "Invalid Release",
				TagName: "invalid-tag",
			},
			mockStatusCode: http.StatusOK,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the current version for this test
			Version = tt.currentVersion

			// Create a custom HTTP client that returns our mock responses
			http.DefaultClient = &http.Client{
				Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
					// Check request headers
					if req.Header.Get("X-GitHub-Api-Version") != "2022-11-28" {
						t.Errorf("Expected X-GitHub-Api-Version header to be '2022-11-28', got %s",
							req.Header.Get("X-GitHub-Api-Version"))
					}

					// Prepare mock response
					var respBody []byte
					if tt.mockResponse != nil {
						respBody, _ = json.Marshal(tt.mockResponse)
					}

					// Create a mock response
					return &http.Response{
						StatusCode: tt.mockStatusCode,
						Body:       io.NopCloser(strings.NewReader(string(respBody))),
					}, nil
				}),
			}

			// Call CheckForUpdate
			err := CheckForUpdate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckForUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// RoundTripFunc is a helper type to mock HTTP responses
type RoundTripFunc func(req *http.Request) (*http.Response, error)

// RoundTrip implements the http.RoundTripper interface
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// Test for SemverString function
func TestSemverString(t *testing.T) {
	originalVersion := Version
	defer func() { Version = originalVersion }()

	tests := []struct {
		version string
		want    string
	}{
		{"1.2.3", "v1.2.3"},
		{"0.1.0", "v0.1.0"},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			Version = tt.version
			if got := SemverString(); got != tt.want {
				t.Errorf("SemverString() = %v, want %v", got, tt.want)
			}
		})
	}
}
