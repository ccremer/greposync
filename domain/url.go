package domain

import (
	"net/url"
	"path"
	"strings"
)

// GitURL is the same as url.URL but with additional helper methods.
type GitURL url.URL

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

// Redacted returns the same as url.URL:Redacted().
func (u *GitURL) Redacted() string {
	plain := url.URL(*u)
	return plain.Redacted()
}

// String returns the same as url.URL:String().
func (u *GitURL) String() string {
	plain := url.URL(*u)
	return plain.String()
}

// FromURL converts the given url.URL into a GitURL.
func FromURL(url *url.URL) *GitURL {
	g := GitURL(*url)
	return &g
}
