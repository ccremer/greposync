# git-repo-sync

This is a PoC to reimplement modulesync in Go

## Project status

This project in its current status is merely more than a weekend's work.
It got triggered by the lack of deep-merging YAML in modulesync configs and limitations in ERB templates.
It got inspired by Helm's templating mechanisms.
Most features are untested yet ("it works on my machine").

## Features implemented

Feature | modulesync | git-repo-sync
---     | ---        | ---
Git Clone | ✔️ | ✔️
Git Commit | ✔️ | ✔️
Git Commit --amend | ✔️ | ✔️
Git Push | ✔️ | ✔️
Git Push --force | ✔️ | ✔️
Git Tags | ✔️ |
GitHub create PR | ✔️ | ✔️
GitHub update PR | ❌ | ✔️
GitLab create PR | ✔️ | ❌
GitLab update PR | ❌ | ❌
Pre-Commit script | ✔️ | ❌
Default git namespace | ✔️ | ✔️
Default git base url | ✔️ | ✔️
Template Defaults | ✔️ | ✔️
Per repository overrides | ✔️ | ✔️
Hooks | ✔️ | ❌
CLI help | ✔️ | ✔️
Filtering repositories | ✔️ |
Dry run | ✔️ | ✔️
Changelog | ✔️ | ❌

> ✔️ Feature implemented
>
> ❌ Feature not implemented (not planned)

Some features aren't planned yet resp. won't be added to git-repo-sync.

## Differences

Feature | modulesync | git-repo-sync
---     | ---        | ---
Template engine | ERB | gotemplate
Installation | Gemfile | single binary, Docker, gomodule, apk, deb, rpm
Deep-merge YAML | ❌ | ✔️
PullRequest template | ❌ | ✔️

## Note before migrating

> It's too soon yet!

Most notably, the template engine is using gotemplate with Sprig library.
This means that you would have to rewrite your templates when migrating to gsync.
To clearly differentiate between the two, the `template` dir is `moduleroot`'s pendent.

The configuration syntax in `managed_repos.yml` (`managed_modules.yml`) and `gitreposync.yml` (`modulesync.yml`) has changed as well compared to their predecessors.
Checkout the examples and documentation (coming later).
