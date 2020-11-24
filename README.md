# GitHub Actions
`github-actions` is a collection of end-user [GitHub Actions][gha] that integrate with Cloud Native Buildpacks projects.

[gha]: https://docs.github.com/en/free-pro-team@latest/actions

- [GitHub Actions](#github-actions)
  - [Buildpack](#buildpack)
    - [Compute Mmetadata Action](#compute-mmetadata-action)
  - [Buildpackage](#buildpackage)
    - [Verify Metadata Action](#verify-metadata-action)
  - [Registry](#registry)
    - [Add Entry Action](#add-entry-action)
    - [Compute Metadata Action](#compute-metadata-action)
    - [Request Add Entry Action](#request-add-entry-action)
    - [Request Yank Entry Action](#request-yank-entry-action)
    - [Verify Namespace Owner Action](#verify-namespace-owner-action)
    - [Yank Entry Action](#yank-entry-action)
  - [Setup pack CLI Action](#setup-pack-cli-action)
  - [License](#license)

## Buildpack

### Compute Mmetadata Action
The `buildpack/compute-metadata` action parses a `buildpack.toml` and exposes the contents of the `[buildpack]` block as output parameters.

```yaml
uses: docker://ghcr.io/buildpacks/actions/buildpack/compute-metadata
```

#### Inputs <!-- omit in toc -->
| Parameter | Description
| :-------- | :----------
| `path` | Optional path to `buildpack.toml`. Defaults to `<working-dir>/buildpack.toml`

#### Outputs <!-- omit in toc -->
| Parameter | Description
| :-------- | :----------
| `id` | The contents of `buildpack.id`
| `name` | The contents of `buildpack.name`
| `version` | The contents of `buildpack.version`
| `homepage` | The contents of `buildpack.homepage`

## Buildpackage

### Verify Metadata Action
The `buildpackage/verify-metadata` action parses the metadata on a buildpackage and verifies that the `id` and `version` match expected values.

```yaml
uses: docker://ghcr.io/buildpacks/actions/buildpackage/verify-metadata
with:
  id:      test-buildpack
  version: "1.0.0"
  address: ghcr.io/example/test-buildpack@sha256:04ba2d17480910bd340f0305d846b007148dafd64bc6fc2626870c174b7c7de7
```

#### Inputs <!-- omit in toc -->
| Parameter | Description
| :-------- | :----------
| `id` | The expected `id` for the buildpackage
| `version` | The expected `version` for the buildpackage
| `address` | The digest-style address of the buildpackage to verify

## Registry
[bri]: https://github.com/buildpacks/registry-index

### Add Entry Action
The `registry/add-entry` adds an entry to the [Buildpack Registry Index][bri].

```yaml
uses: docker://ghcr.io/buildpacks/actions/registry/add-entry
with:
  token: ${{ secrets.BOT_TOKEN }}
  owner: ${{ env.INDEX_OWNER }}
  repository: ${{ env.INDEX_REPOSITORY }}
  namespace: ${{ steps.metadata.outputs.namespace }}
  name: ${{ steps.metadata.outputs.name }}
  version: ${{ steps.metadata.outputs.version }}
  address: ${{ steps.metadata.outputs.address }}
```

#### Inputs <!-- omit in toc -->
| Parameter | Description
| :-------- | :----------
| `token` | A GitHub token with permissions to commit to the registry index repository.
| `owner` | The owner name of the registry index repository.
| `repository` | The repository name of the registry index repository.
| `namespace` | The namespace of the buildpack to register.
| `name` | The name of the buildpack to register.
| `version` | The version of the buildpack to register.
| `address` | The address of the buildpack to register.

### Compute Metadata Action
The `registry/compute-metadata` action parses a [`buildpacks/registry-index`][bri] issue and exposes the contents as output parameters.

```yaml
uses: docker://ghcr.io/buildpacks/actions/registry/add-entry
with:
  issue:   ${{ toJSON(github.events.issue) }}
```

#### Inputs <!-- omit in toc -->
| Parameter | Description
| :-------- | :----------
| `issue` | The GitHub issue payload.

#### Outputs <!-- omit in toc -->
| Parameter | Description
| :-------- | :----------
| `id` | The contents of `id`
| `version` | The contents of `version`
| `address` | The contents of `addr`
| `namespace` | The namespace portion of `id`
| `name` | The name portion of `id`

### Request Add Entry Action
The `registry/request-add-entry` action adds an entry to the [Buildpack Registry Index][bri].

```yaml
uses: docker://ghcr.io/buildpacks/actions/registry/request-add-entry
with:
  token:   ${{ secrets.IMPLEMENTATION_PAT }}
  id:      $buildpacksio/test-buildpack
  version: ${{ steps.deploy.outputs.version }}
  address: index.docker.io/buildpacksio/test-buildpack@${{ steps.deploy.outputs.digest }}
```

#### Inputs <!-- omit in toc -->
| Parameter | Description
| :-------- | :----------
| `token` | A GitHub token with `public_repo` scope to open an issue against [`buildpacks/registry-index`][bri].
| `id` | A buildpack id that your user is allowed to manage.  This is must be in `{namespace}/{name}` format.
| `version` | The version of the buildpack that is being added to the registry.
| `address` | The Docker URI of the buildpack artifact.  This is must be in `{host}/{repo}@{digest}` form.

### Request Yank Entry Action
The `registry/request-yank-entry` action yanks an entry from the [Buildpack Registry Index][bri].

```yaml
uses: docker://ghcr.io/buildpacks/actions/registry/request-yank-entry
with:
  token:   ${{ secrets.IMPLEMENTATION_PAT }}
  id:      buildpacksio/test-buildpack
  version: ${{ steps.deploy.outputs.version }}
```

#### Inputs <!-- omit in toc -->
| Parameter | Description
| :-------- | :----------
| `token` | A GitHub token with `public_repo` scope to open an issue against [`buildpacks/registry-index`][bri].
| `id` | A buildpack id that your user is allowed to manage.  This is must be in `{namespace}/{name}` format.
| `version` | The version of the buildpack that is being added to the registry.

### Verify Namespace Owner Action
The `registry/verify-namespace-owner` action verifies that a user is an owner of a namespace in the [Buildpack Registry Index][bri].

```yaml
uses: docker://ghcr.io/buildpacks/actions/registry/verify-namespace-owner
with:
  token: ${{ secrets.BOT_TOKEN }}
  owner: ${{ env.NAMESPACES_OWNER }}
  repository: ${{ env.NAMESPACES_REPOSITORY }}
  namespace: ${{ steps.metadata.outputs.namespace }}
  user: ${{ toJSON(github.event.issue.user) }}
  add-if-missing: true
```

#### Inputs <!-- omit in toc -->
| Parameter | Description
| :-------- | :----------
| `token` | A GitHub token with permissions to commit to the registry namespaces repository.
| `owner` | The owner name of the registry namespaces repository.
| `repository` | The repository name of the registry namespaces repository.
| `namespace` | The namespace to check ownership for.
| `user` | The Github user payload.
| `add-if-missing` | Whether to add the current user as the owner of the namespace if that namespace does not exist. (Optional. Default `false`)

### Yank Entry Action
The `registry/yank-entry` yanks an entry from the [Buildpack Registry Index][bri].

```yaml
uses: docker://ghcr.io/buildpacks/actions/registry/yank-entry
with:
  token: ${{ secrets.BOT_TOKEN }}
  owner: ${{ env.INDEX_OWNER }}
  repository: ${{ env.INDEX_REPOSITORY }}
  namespace: ${{ steps.metadata.outputs.namespace }}
  name: ${{ steps.metadata.outputs.name }}
  version: ${{ steps.metadata.outputs.version }}
```

#### Inputs <!-- omit in toc -->
| Parameter | Description
| :-------- | :----------
| `token` | A GitHub token with permissions to commit to the registry index repository.
| `owner` | The owner name of the registry index repository.
| `repository` | The repository name of the registry index repository.
| `namespace` | The namespace of the buildpack to register.
| `name` | The name of the buildpack to register.
| `version` | The version of the buildpack to register.
| `address` | The address of the buildpack to register.

## Setup pack CLI Action
The setup-pack action adds [crane][crane], [`jq`][jq], [`pack`][pack], and [`yj`][yj] to the environment.

[crane]: https://github.com/google/go-containerregistry/tree/master/cmd/crane
[jq]:    https://stedolan.github.io/jq/
[pack]:  https://github.com/buildpacks/pack
[yj]:    https://github.com/sclevine/yj

```yaml
uses: buildpacks/github-actions/setup-pack
```

#### Inputs <!-- omit in toc -->
| Parameter | Description
| :-------- | :----------
| `crane-version` | Optional version of [`crane`][crane] to install. Defaults to latest release.
| `jq-version` | Optional version of [`jq`][jq] to install. Defaults to latest release.
| `pack-version` | Optional version of [`pack`][pack] to install. Defaults to latest release.
| `yj-version` | Optional version of [`yj`][yj] to install. Defaults to latest release.

## License
This library is released under version 2.0 of the [Apache License][a].

[a]: https://www.apache.org/licenses/LICENSE-2.0
