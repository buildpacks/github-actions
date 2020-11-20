# GitHub Actions
`github-actions` is a collection of end-user [GitHub Actions][gha] that integrate with Cloud Native Buildpacks projects.

[gha]: https://docs.github.com/en/free-pro-team@latest/actions

- [GitHub Actions](#github-actions)
  - [Buildpack](#buildpack)
    - [Compute Mmetadata Action](#compute-mmetadata-action)
      - [Inputs](#inputs)
      - [Outputs](#outputs)
  - [Buildpackage](#buildpackage)
    - [Verify Metadata Action](#verify-metadata-action)
      - [Inputs](#inputs-1)
  - [Registry](#registry)
    - [Add Entry Action](#add-entry-action)
      - [Inputs](#inputs-2)
    - [Yank Entry Action](#yank-entry-action)
      - [Inputs](#inputs-3)
  - [Setup pack CLI Action](#setup-pack-cli-action)
      - [Inputs](#inputs-4)
  - [License](#license)

## Buildpack

### Compute Mmetadata Action
The `buildpack/compute-metadata` action parses a `buildpack.toml` and exposes the contents of the `[buildpack]` block as output parameters.

```yaml
uses: docker://ghcr.io/buildpacks/actions/buildpack/compute-metadata
```

#### Inputs
| Parameter | Description
| :-------- | :----------
| `path` | Optional path to `buildpack.toml`. Defaults to `<working-dir>/buildpack.toml` 

#### Outputs
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

#### Inputs
| Parameter | Description
| :-------- | :----------
| `id` | The expected `id` for the buildpackage
| `version` | The expected `version` for the buildpackage
| `address` | The digest-style address of the buildpackage to verify

## Registry

### Add Entry Action 
The `registry/add-entry` action adds an entry to the [Buildpack Registry Index][bri].

[bri]: https://github.com/buildpacks/registry-index

```yaml
uses: docker://ghcr.io/buildpacks/actions/registry/add-entry
with:
  token:   ${{ secrets.IMPLEMENTATION_PAT }}
  id:      $buildpacksio/test-buildpack
  version: ${{ steps.deploy.outputs.version }}
  address: index.docker.io/buildpacksio/test-buildpack@${{ steps.deploy.outputs.digest }}
```

#### Inputs
| Parameter | Description
| :-------- | :----------
| `token` | A GitHub token with `public_repo` scope to open an issue against [`buildpacks/registry-index`][bri].
| `id` | A buildpack id that your user is allowed to manage.  This is must be in `{namespace}/{name}` format.
| `version` | The version of the buildpack that is being added to the registry.
| `address` | The Docker URI of the buildpack artifact.  This is must be in `{host}/{repo}@{digest}` form.

### Yank Entry Action
The `registry/yank-entry` action yanks an entry from the [Buildpack Registry Index][bri].

```yaml
uses: docker://ghcr.io/buildpacks/actions/registry/yank-entry
with:
  token:   ${{ secrets.IMPLEMENTATION_PAT }}
  id:      buildpacksio/test-buildpack
  version: ${{ steps.deploy.outputs.version }}
```

#### Inputs
| Parameter | Description
| :-------- | :----------
| `token` | A GitHub token with `public_repo` scope to open an issue against [`buildpacks/registry-index`][bri].
| `id` | A buildpack id that your user is allowed to manage.  This is must be in `{namespace}/{name}` format.
| `version` | The version of the buildpack that is being added to the registry.
| `yank` | `true` if this version should be yanked.

## Setup pack CLI Action
The setup-pack action adds [crane][crane], [`jq`][jq], [`pack`][pack], and [`yj`][yj] to the environment.

[crane]: https://github.com/google/go-containerregistry/tree/master/cmd/crane
[jq]:    https://stedolan.github.io/jq/
[pack]:  https://github.com/buildpacks/pack
[yj]:    https://github.com/sclevine/yj

```yaml
uses: buildpacks/github-actions/setup-pack
```

#### Inputs
| Parameter | Description
| :-------- | :----------
| `crane-version` | Optional version of [`crane`][crane] to install. Defaults to latest release.
| `jq-version` | Optional version of [`jq`][jq] to install. Defaults to latest release.
| `pack-version` | Optional version of [`pack`][pack] to install. Defaults to latest release.
| `yj-version` | Optional version of [`yj`][yj] to install. Defaults to latest release. 

## License
This library is released under version 2.0 of the [Apache License][a].

[a]: https://www.apache.org/licenses/LICENSE-2.0
