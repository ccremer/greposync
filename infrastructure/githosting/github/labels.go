package github

import (
	"github.com/ccremer/greposync/domain"
	"github.com/google/go-github/v39/github"
)

// FetchLabels implements githosting.Remote.
func (r *GhRemote) FetchLabels(repository *domain.GitRepository) (domain.LabelSet, error) {
	ghLabels, err := r.fetchAllLabels(repository)
	if err == nil {
		r.labelCache[repository.URL] = ghLabels
	}
	return LabelSetConverter{}.ConvertToEntity(ghLabels), err
}

// EnsureLabels implements githosting.Remote.
func (r *GhRemote) EnsureLabels(repository *domain.GitRepository, labels domain.LabelSet) error {
	for _, label := range labels {
		cached, exists := r.findCachedLabel(repository.URL, label)
		if exists {
			if r.hasLabelChanged(cached, label) {
				err := r.updateLabel(repository, cached, label)
				if err != nil {
					return err
				}
			}
			continue
		}
		err := r.createLabel(repository, label)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteLabels implements githosting.Remote.
func (r *GhRemote) DeleteLabels(repository *domain.GitRepository, labels domain.LabelSet) error {
	for _, label := range labels {
		var converted *github.Label
		cached, exists := r.findCachedLabel(repository.URL, label)
		if exists {
			converted = cached
		} else {
			converted = LabelConverter{}.ConvertFromEntity(label)
		}
		_, err := r.deleteLabel(repository, converted)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *GhRemote) findCachedLabel(url *domain.GitURL, label domain.Label) (*github.Label, bool) {
	cachedSet, exists := r.labelCache[url]
	if !exists {
		return nil, false
	}
	for _, cached := range cachedSet {
		if *cached.Name == label.Name {
			return cached, true
		}
	}
	return nil, false
}

func (r *GhRemote) updateLabelCache(url *domain.GitURL, label *github.Label) {
	cachedSet, exists := r.labelCache[url]
	if !exists {
		r.labelCache[url] = []*github.Label{label}
		return
	}
	for i, cached := range cachedSet {
		if cached.GetName() == label.GetName() {
			cachedSet[i] = label
			return
		}
	}
	cachedSet = append(cachedSet, label)
	r.labelCache[url] = cachedSet
}

func (r *GhRemote) removeLabelFromCache(url *domain.GitURL, label *github.Label) {
	cachedSet, exists := r.labelCache[url]
	if !exists {
		return
	}
	for i, cached := range cachedSet {
		if cached.GetName() == label.GetName() {
			// replace the existing index with the last element
			cachedSet[i] = cachedSet[len(cachedSet)-1]
			// remove the (duplicated) last element
			newSet := cachedSet[:len(cachedSet)-1]
			r.labelCache[url] = newSet
			return
		}
	}
}

func (r *GhRemote) createLabel(repository *domain.GitRepository, label domain.Label) error {
	r.m.Lock()
	defer r.delayedUnlock()
	converted := LabelConverter{}.ConvertFromEntity(label)
	newLabel, _, err := r.client.Issues.CreateLabel(r.ctx, repository.URL.GetNamespace(), repository.URL.GetRepositoryName(), converted)
	r.updateLabelCache(repository.URL, newLabel)
	return r.instrumentation.createdLabel(repository, label, err)
}

func (r *GhRemote) updateLabel(repository *domain.GitRepository, ghLabel *github.Label, label domain.Label) error {
	r.m.Lock()
	defer r.delayedUnlock()
	ghLabel.Description = &label.Description
	color := ColorConverter{}.ConvertFromEntity(label.GetColor())
	ghLabel.Color = &color
	updatedLabel, _, err := r.client.Issues.EditLabel(r.ctx, repository.URL.GetNamespace(), repository.URL.GetRepositoryName(), label.Name, ghLabel)
	r.updateLabelCache(repository.URL, updatedLabel)
	return r.instrumentation.updatedLabel(repository, label, err)
}

func (r *GhRemote) deleteLabel(repository *domain.GitRepository, label *github.Label) (bool, error) {
	r.m.Lock()
	defer r.delayedUnlock()
	resp, err := r.client.Issues.DeleteLabel(r.ctx, repository.URL.GetNamespace(), repository.URL.GetRepositoryName(), label.GetName())
	if resp != nil && resp.StatusCode == 404 {
		// Not an error
		return false, nil
	}
	if err == nil {
		r.removeLabelFromCache(repository.URL, label)
	}
	return err == nil, r.instrumentation.deletedLabel(repository, label, err)
}

func (r *GhRemote) fetchAllLabels(repository *domain.GitRepository) ([]*github.Label, error) {
	r.m.Lock()
	defer r.delayedUnlock()
	nextPage := 1
	var allLabels []*github.Label
	for repeat := true; repeat; repeat = nextPage > 0 {
		labels, resp, err := r.client.Issues.ListLabels(r.ctx, repository.URL.GetNamespace(), repository.URL.GetRepositoryName(), &github.ListOptions{
			Page:    nextPage,
			PerPage: 100,
		})
		if err != nil {
			return nil, err
		}
		allLabels = append(allLabels, labels...)
		// On the last page, the NextPage is 0 again, we can use that to exit the loop
		nextPage = resp.NextPage
	}
	r.instrumentation.fetchedAllLabels(repository, allLabels)
	return allLabels, nil
}

func (r *GhRemote) hasLabelChanged(ghLabel *github.Label, repoLabel domain.Label) bool {
	converted := ColorConverter{}.ConvertFromEntity(repoLabel.GetColor())
	return ghLabel.GetDescription() != repoLabel.Description || ghLabel.GetColor() != converted
}
