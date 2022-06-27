// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package common

const (
	FetchContentAnnotationKey = "fetch-content-from"
	LocalFetchAnnotationKey   = "kctrl.carvel.dev/local-fetch-0"
	PackageBuildFileName      = "package-build.yml"
)

const (
	FetchReleaseArtifactFromGithub string = "Release artifact from Github Repository"
	FetchManifestFromGithub        string = "Git Repository(Not supported)"
	FetchChartFromHelmRepo         string = "Helm Chart from Helm Repository"
	FetchChartFromGithub           string = "Helm Chart from Github repository(Not supported)"
	FetchFromLocalDirectory        string = "Local Directory"
)
