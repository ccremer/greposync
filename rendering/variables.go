package rendering

// ConstructMetadata returns a map with metadata values.
// Included are "Repository" and "PullRequest", both deserialized from the config struct.
func (r *Renderer) ConstructMetadata() Values {
	d := Values{
		"Repository":  r.cfg.Git,
		"PullRequest": r.cfg.PullRequest,
	}
	return d
}
