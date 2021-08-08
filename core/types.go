package core

import (
	"net/url"
	"path"
	"strings"
)

type (
	// GitURL is the same as url.URL but with additional helper methods.
	GitURL url.URL
	// GitHostingProvider is the identification key for a Git hosting service.
	GitHostingProvider string
)

// GetRepositoryName returns the last element of the Git URL.
// Strips the name from any .git extensions in the URL.
func (u *GitURL) GetRepositoryName() string {
	return strings.TrimSuffix(path.Base(u.Path), ".git")
}

// GetNamespace returns the middle element(s) of the Git URL.
// Depending on the Git hosting service, this name may contain multiple slashes.
// Any leading "/" is removed.
func (u *GitURL) GetNamespace() string {
	return strings.TrimPrefix(path.Dir(u.Path), "/")
}

// Redacted returns the same as url.URL::Redacted.
func (u *GitURL) Redacted() string {
	plain := url.URL(*u)
	return plain.Redacted()
}

// FromURL converts the given url.URL into a GitURL.
func FromURL(url *url.URL) *GitURL {
	g := GitURL(*url)
	return &g
}
