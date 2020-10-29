# GitHub Actions
`github-actions` is a collection of end-user [GitHub Actions][gha] that integrate with Cloud Native Buildpacks projects.

[gha]: https://docs.github.com/en/free-pro-team@latest/actions

## Registry Action
The registry action adds and yanks buildpack releases in the [Buildpack Registry Index][bri].

[bri]: https://github.com/buildpacks/registry-index

### Add
```yaml
uses: docker://ghcr.io/buildpacks/actions/registry
with:
  token:   ${{ secrets.IMPLEMENTATION_PAT }}
  id:      $buildpacksio/test-buildpack
  version: {{ steps.deploy.outputs.version }}
  address: index.docker.io/buildpacksio/test-buildpack@${{ steps.deploy.outputs.digest }}
```

| Parameter | Description
| :-------- | :----------
| `token` | A GitHub token with `public_repo` scope to open an issue against [`buildpacks/registry-index`][bri].
| `id` | A buildpack id that your user is allowed to manage.  This is must be in `{namespace}/{name}` format.
| `version` | The version of the buildpack that is being added to the registry.
| `address` | The Docker URI of the buildpack artifact.  This is must be in `{host}/{repo}@{digest}` form.

### Yank
```yaml
uses: docker://ghcr.io/buildpacks/actions/registry
with:
  token:   ${{ secrets.IMPLEMENTATION_PAT }}
  id:      buildpacksio/test-buildpack
  version: ${{ steps.deploy.outputs.version }}
  yank:    true
```

| Parameter | Description
| :-------- | :----------
| `token` | A GitHub token with `public_repo` scope to open an issue against [`buildpacks/registry-index`][bri].
| `id` | A buildpack id that your user is allowed to manage.  This is must be in `{namespace}/{name}` format.
| `version` | The version of the buildpack that is being added to the registry.
| `yank` | `true` if this version should be yanked.

## License
This library is released under version 2.0 of the [Apache License][a].

[a]: https://www.apache.org/licenses/LICENSE-2.0

