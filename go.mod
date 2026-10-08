module github.com/buildpacks/github-actions

// This must use the full `1.N.0` form, since that's what Go's tooling (and so Dependabot)
// rewrites it to. The GitHub Actions workflows use only the major.minor from this line for
// `setup-go` and the `golang` image tag, so that they get the latest patch release.
go 1.27.0

require (
	github.com/buildpacks/libcnb v1.30.4
	github.com/google/go-containerregistry v0.22.1
	github.com/google/go-github/v89 v89.0.0
	github.com/onsi/gomega v1.44.0
	github.com/pelletier/go-toml/v2 v2.4.3
	github.com/sclevine/spec v1.4.0
	github.com/stretchr/testify v1.12.1
	gopkg.in/retry.v1 v1.0.3
)

require (
	github.com/BurntSushi/toml v1.6.0 // indirect
	github.com/Masterminds/semver/v3 v3.5.0 // indirect
	github.com/docker/cli v29.8.2+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.9.9 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/go-querystring v1.2.0 // indirect
	github.com/klauspost/compress v1.20.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/sirupsen/logrus v1.10.2 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/net v0.59.0 // indirect
	golang.org/x/sync v0.23.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	gotest.tools/v3 v3.5.2 // indirect
)
