package core

import (
	"net/url"
	"path"
	"strings"
)

type (
	// GitUrl is the same as url.URL but with additional helper methods.
	GitUrl url.URL
	// GitHostingProvider is the provider
	GitHostingProvider string
	// GitRepositoryConfig holds all the relevant Git properties.
	GitRepositoryConfig struct {
		// Url is the repository location on the remote hosting provider.
		Url *GitUrl
		// Provider returns the GitHostingProvider identity string.
		// Mainly used to identify the remote API implementation.
		Provider GitHostingProvider
	}
)

// GetRepositoryName returns the last element of the Git URL.
// Strips the name from any .git extensions in the URL.
func (u *GitUrl) GetRepositoryName() string {
	return strings.TrimSuffix(path.Base(u.Path), ".git")
}

// GetNamespace returns the middle element(s) of the Git URL.
// Depending on the Git hosting service, this name may contain multiple slashes.
// Any leading "/" is removed.
func (u *GitUrl) GetNamespace() string {
	return strings.TrimPrefix(path.Dir(u.Path), "/")
}
