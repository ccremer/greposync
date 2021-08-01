package core

import (
	"io/fs"
	"net/url"
	"path"
	"strings"
)

type (
	// GitURL is the same as url.URL but with additional helper methods.
	GitURL url.URL
	// GitHostingProvider is the identification key for a Git hosting service.
	GitHostingProvider string
	// GitRepositoryConfig holds all the relevant Git properties.
	GitRepositoryConfig struct {
		// URL is the repository location on the remote hosting provider.
		URL *GitURL
		// Provider returns the GitHostingProvider identity string.
		// Mainly used to identify the remote API implementation.
		Provider GitHostingProvider
		// RootDir is the local root path to the Git repository.
		RootDir string
	}
	// Values is a key-value construct with arbitrary hierarchy.
	Values map[string]interface{}
	// Template contains meta information about the source template.
	Template struct {
		// RelativePath is the path to a template file relative to the template root directory.
		// The path is delimited with a forward slash ("/") and not OS-specific.
		RelativePath string
		// FileMode is the mode of the RelativePath.
		// When rendering the template in a Git repository, the implementation must write the file with these file permissions.
		FileMode fs.FileMode
	}
	// Output contains meta information how the template should be processed.
	Output struct {
		// TargetPath is the actual file path relative to the Git repository where the resulting template should be written to.
		// The file permissions are copied from the Template property.
		TargetPath string
		// Template contains the information about the source template file.
		Template Template
		// Values contains the variables and parameters which the template should replace placeholders with.
		Values Values
		// Git contains the settings for the Git repository.
		Git GitRepositoryConfig
	}
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

// FromURL converts the given url.URL into a GitURL.
func FromURL(url *url.URL) *GitURL {
	g := GitURL(*url)
	return &g
}
