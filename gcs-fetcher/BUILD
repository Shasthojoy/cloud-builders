load("@io_bazel_rules_go//go:def.bzl", "gazelle", "go_prefix")

go_prefix("github.com/GoogleCloudPlatform/cloud-builders/gcs-fetcher")

gazelle(
    name = "gazelle",
    external = "vendored",
)

load("@io_bazel_rules_docker//docker:docker.bzl", "container_push")

container_push(
    name = "push-gcs-fetcher",
    format = "Docker",
    image = "//cmd/gcs-fetcher:image",
    registry = "gcr.io",
    repository = "$(project)/gcs-fetcher",
    stamp = True,
)

container_push(
    name = "push-upload-manifest",
    format = "Docker",
    image = "//cmd/upload-manifest:image",
    registry = "gcr.io",
    repository = "$(project)/upload-manifest",
    stamp = True,
)
