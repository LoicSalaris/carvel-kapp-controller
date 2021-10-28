// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageavailable

import (
	"context"
	"fmt"
	"strings"

	"github.com/cppforlife/go-cli-ui/ui"
	uitable "github.com/cppforlife/go-cli-ui/ui/table"
	cmdcore "github.com/k14s/kapp/pkg/kapp/cmd/core"
	"github.com/k14s/kapp/pkg/kapp/logger"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GetOptions struct {
	ui          ui.UI
	depsFactory cmdcore.DepsFactory
	logger      logger.Logger

	NamespaceFlags cmdcore.NamespaceFlags
	PackageName    string
}

func NewGetOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *GetOptions {
	return &GetOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewGetCmd(o *GetOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "Get details for an available package or the openAPI schema of a package with a specific version",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVarP(&o.PackageName, "package", "P", "", "List all available versions of package")
	return cmd
}

func (o *GetOptions) Run() error {
	var pkgName, pkgVersion string
	pkgNameVersion := strings.Split(o.PackageName, "/")
	if len(pkgNameVersion) == 2 {
		pkgName = pkgNameVersion[0]
		pkgVersion = pkgNameVersion[1]
	} else if len(pkgNameVersion) == 1 {
		pkgName = pkgNameVersion[0]
	} else {
		return fmt.Errorf("package should be of the format name or name/version")
	}

	tableTitle := fmt.Sprintf("Package details for '%s'", pkgName)
	headers := []uitable.Header{
		uitable.NewHeader("name"),
		uitable.NewHeader("display-name"),
		uitable.NewHeader("short-description"),
		uitable.NewHeader("package-provider"),
		uitable.NewHeader("long-description"),
		uitable.NewHeader("maintainers"),
		uitable.NewHeader("support"),
		uitable.NewHeader("category"),
	}

	client, err := o.depsFactory.PackageClient()
	if err != nil {
		return err
	}

	pkgMetadata, err := client.DataV1alpha1().PackageMetadatas(
		o.NamespaceFlags.Name).Get(context.Background(), pkgName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	row := []uitable.Value{
		uitable.NewValueString(pkgMetadata.Name),
		uitable.NewValueString(pkgMetadata.Spec.DisplayName),
		uitable.NewValueString(pkgMetadata.Spec.ShortDescription),
		uitable.NewValueString(pkgMetadata.Spec.ProviderName),
		uitable.NewValueString(pkgMetadata.Spec.LongDescription),
		uitable.NewValueInterface(pkgMetadata.Spec.Maintainers),
		uitable.NewValueString(pkgMetadata.Spec.SupportDescription),
		uitable.NewValueStrings(pkgMetadata.Spec.Categories),
	}

	if pkgVersion != "" {

		pkg, err := client.DataV1alpha1().Packages(
			o.NamespaceFlags.Name).Get(context.Background(), fmt.Sprintf("%s.%s", pkgName, pkgVersion), metav1.GetOptions{})
		if err != nil {
			return err
		}
		headers = append(headers, []uitable.Header{
			uitable.NewHeader("version"),
			uitable.NewHeader("released-at"),
			uitable.NewHeader("minimum-capacity-requirements"),
			uitable.NewHeader("release-notes"),
			uitable.NewHeader("license"),
		}...)

		row = append(row, []uitable.Value{
			uitable.NewValueString(pkg.Spec.Version),
			uitable.NewValueString(pkg.Spec.ReleasedAt.String()),
			uitable.NewValueString(pkg.Spec.CapactiyRequirementsDescription),
			uitable.NewValueString(pkg.Spec.ReleaseNotes),
			uitable.NewValueStrings(pkg.Spec.Licenses),
		}...)
	}

	table := uitable.Table{
		Title:   tableTitle,
		Content: "Package Details",

		Header: headers,

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
			{Column: 1, Asc: true},
		},
	}

	table.Rows = append(table.Rows, row)
	table.Transpose = true
	o.ui.PrintTable(table)

	return nil
}