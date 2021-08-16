// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package secretgencontroller

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/carvel-kapp-controller/test/e2e"
)

func Test_PlaceholderSecrets_DeletedWhenPackageInstallDeleted(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	name := "placeholder-garbage-collection"
	sas := e2e.ServiceAccounts{env.Namespace}

	pkgiYaml := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: pkg.test.carvel.dev
  namespace: %[1]s
spec:
  # This is the name we want to reference in resources such as PackageInstall.
  displayName: "Test PackageMetadata in repo"
  shortDescription: "PackageMetadata used for testing"
  longDescription: "A longer, more detailed description of what the package contains and what it is for"
  providerName: Carvel
  maintainers:
  - name: carvel
  categories:
  - testing
  supportDescription: "Description of support provided for the package"
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: pkg.test.carvel.dev.1.0.0
  namespace: %[1]s
spec:
  refName: pkg.test.carvel.dev
  version: 1.0.0
  licenses:
  - Apache 2.0
  capactiyRequirementsDescription: "cpu: 1,RAM: 2, Disk: 3"
  releaseNotes: |
    - Introduce simple-app package
  releasedAt: 2021-05-05T18:57:06Z
  template:
    spec:
      fetch:
      - imgpkgBundle:
          image: k8slt/kctrl-example-pkg:v1.0.0
      template:
      - ytt: {}
      - kbld:
          paths:
          - "-"
          - ".imgpkg/images.yml"
      deploy:
      - kapp: {}
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  name: %[2]s
  namespace: %[1]s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/packageinstalls
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  packageRef:
    refName: pkg.test.carvel.dev
    versionSelection:
      constraints: 1.0.0
  values:
  - secretRef:
      name: pkg-demo-values
---
apiVersion: v1
kind: Secret
metadata:
  name: pkg-demo-values
stringData:
  values.yml: |
    hello_msg: "hi"
`, env.Namespace, name) + sas.ForNamespaceYAML()

	cleanUp := func() {
		// Delete App with kubectl first since kapp
		// deletes ServiceAccount before App
		kubectl.RunWithOpts([]string{"delete", "apps/" + name}, e2e.RunOpts{AllowError: true})
		kapp.Run([]string{"delete", "-a", name})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("Create PackageInstall", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgiYaml)})

		kubectl.Run([]string{"wait", "--for=condition=ReconcileSucceeded", "pkgi/" + name, "--timeout", "1m"})
	})

	logger.Section("Check placeholder secret created", func() {
		kubectl.Run([]string{"get", "secret", name + "-fetch0"})
	})

	logger.Section("Check placeholder secret deleted after PackageInstall deleted", func() {
		cleanUp()
		out, err := kubectl.RunWithOpts([]string{"get", "secret", name + "fetch0"}, e2e.RunOpts{AllowError: true})
		assert.NotNil(t, err, "expected error from not finding placeholder secret.\nGot: "+out)
	})
}

func Test_PackageInstall_CanAuthenticateToPrivateRepository_UsingPlaceholderSecret(t *testing.T) {
	env := e2e.BuildEnv(t)
	logger := e2e.Logger{}
	kapp := e2e.Kapp{t, env.Namespace, logger}
	kubectl := e2e.Kubectl{t, env.Namespace, logger}
	name := "placeholder-private-auth"
	sas := e2e.ServiceAccounts{env.Namespace}

	// If this changes, the skip-tls-verify domain must be updated to match
	registryNamespace := "registry"
	registryName := "test-registry"

	pkgiYaml := fmt.Sprintf(`---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: PackageMetadata
metadata:
  name: pkg.test.carvel.dev
  namespace: %[1]s
spec:
  # This is the name we want to reference in resources such as PackageInstall.
  displayName: "Test PackageMetadata in repo"
  shortDescription: "PackageMetadata used for testing"
  longDescription: "A longer, more detailed description of what the package contains and what it is for"
  providerName: Carvel
  maintainers:
  - name: carvel
  categories:
  - testing
  supportDescription: "Description of support provided for the package"
---
apiVersion: data.packaging.carvel.dev/v1alpha1
kind: Package
metadata:
  name: pkg.test.carvel.dev.1.0.0
  namespace: %[1]s
spec:
  refName: pkg.test.carvel.dev
  version: 1.0.0
  template:
    spec:
      fetch:
      - imgpkgBundle:
          image: registry-svc.%[3]s.svc.cluster.local:443/my-repo/image
      template:
      - ytt: {}
      - kbld:
          paths:
          - "-"
          - ".imgpkg/images.yml"
      deploy:
      - kapp: {}
---
apiVersion: packaging.carvel.dev/v1alpha1
kind: PackageInstall
metadata:
  name: %[2]s
  namespace: %[1]s
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/packageinstalls
spec:
  syncPeriod: 30s
  serviceAccountName: kappctrl-e2e-ns-sa
  packageRef:
    refName: pkg.test.carvel.dev
    versionSelection:
      constraints: 1.0.0
`, env.Namespace, name, registryNamespace) + sas.ForNamespaceYAML()

	secretYaml := fmt.Sprintf(`
---
apiVersion: v1
kind: Secret
metadata:
  name: regcred
type: kubernetes.io/dockerconfigjson
stringData:
  .dockerconfigjson: |
    {
      "auths": {
        "registry-svc.%s.svc.cluster.local": {
          "username": "testuser",
          "password": "testpassword",
          "auth": ""
        }
      }
    }
---
apiVersion: secretgen.k14s.io/v1alpha1
kind: SecretExport
metadata:
  name: regcred
spec:
  toNamespaces:
  - %s
`, registryNamespace, env.Namespace)

	cleanUp := func() {
		// Delete App with kubectl first since kapp
		// deletes ServiceAccount before App
		kubectl.RunWithOpts([]string{"delete", "secret", name + "-fetch0"}, e2e.RunOpts{AllowError: true})
		kubectl.RunWithOpts([]string{"delete", "apps/" + name}, e2e.RunOpts{AllowError: true})
		kapp.Run([]string{"delete", "-a", registryName, "-n", registryNamespace})
		kapp.Run([]string{"delete", "-a", name})
		kapp.Run([]string{"delete", "-a", "secret-export"})
	}
	cleanUp()
	defer cleanUp()

	logger.Section("deploy registry with self signed certs", func() {
		kapp.Run([]string{"deploy", "-f", "../assets/registry/registry2.yml", "-f", "../assets/registry/certs-for-skip-tls.yml", "-f", "../assets/registry/htpasswd-auth", "-a", registryName})
	})

	logger.Section("Create Docker Registry Secret", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", "secret-export", "-f", "-"},
			e2e.RunOpts{StdinReader: strings.NewReader(secretYaml)})
	})

	logger.Section("Create PackageInstall", func() {
		kapp.RunWithOpts([]string{"deploy", "-a", name, "-f", "-"},
			e2e.RunOpts{StdinReader: strings.NewReader(pkgiYaml)})
		kubectl.Run([]string{"wait", "--for=condition=ReconcileSucceeded", "pkgi/" + name, "--timeout", "10m"})
	})

	logger.Section("Check PackageInstall/App succeed", func() {
		kubectl.Run([]string{"wait", "--for=condition=ReconcileSucceeded", "pkgi/" + name, "--timeout", "5m"})
		kubectl.Run([]string{"wait", "--for=condition=ReconcileSucceeded", "app/" + name, "--timeout", "5m"})
		kubectl.Run([]string{"get", "configmap", "simple-app-values"})
	})
}
