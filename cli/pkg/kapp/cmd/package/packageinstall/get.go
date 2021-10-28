// Copyright 2020 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package packageinstall

import (
	"bufio"
	"context"
	"fmt"
	"os"

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
	pkgiName    string
	valuesFile  string

	NamespaceFlags cmdcore.NamespaceFlags
}

func NewGetOptions(ui ui.UI, depsFactory cmdcore.DepsFactory, logger logger.Logger) *GetOptions {
	return &GetOptions{ui: ui, depsFactory: depsFactory, logger: logger}
}

func NewGetCmd(o *GetOptions, flagsFactory cmdcore.FlagsFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"g"},
		Short:   "Get installed Package",
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
	}
	o.NamespaceFlags.Set(cmd, flagsFactory)
	cmd.Flags().StringVar(&o.pkgiName, "name", "", "Name of PackageInstall")
	cmd.Flags().StringVar(&o.valuesFile, "values-file", "", "File path for exporting configuration values file")
	return cmd
}

func (o *GetOptions) Run() error {

	client, err := o.depsFactory.KappCtrlClient()
	if err != nil {
		return err
	}

	pkgi, err := client.PackagingV1alpha1().PackageInstalls(
		o.NamespaceFlags.Name).Get(context.Background(), o.pkgiName, metav1.GetOptions{})
	if err != nil {
		//Handle IsNotFound error?
		return err
	}

	//TODO: Verify and enhance how we export values file
	if o.valuesFile != "" {
		f, err := os.Create(o.valuesFile)
		if err != nil {
			return err
		}
		defer f.Close()
		w := bufio.NewWriter(f)

		coreClient, err := o.depsFactory.CoreClient()
		if err != nil {
			return err
		}

		dataValue := ""
		for _, value := range pkgi.Spec.Values {
			if value.SecretRef == nil {
				continue
			}

			s, err := coreClient.CoreV1().Secrets(o.NamespaceFlags.Name).Get(context.Background(), "db-user-pass", metav1.GetOptions{})
			if err != nil {
				return err
			}

			var data []byte
			yamlSeperator := "---"
			for _, value := range s.Data {
				if len(string(value)) < 3 {
					data = append(data, yamlSeperator...)
					data = append(data, "\n"...)
				}
				if len(string(value)) >= 3 && string(value)[:3] != yamlSeperator {
					data = append(data, yamlSeperator...)
					data = append(data, "\n"...)
				}
				data = append(data, value...)
			}

			if len(string(data)) < 3 {
				dataValue += yamlSeperator
				dataValue += "\n"
			}
			if len(string(data)) >= 3 && string(data)[:3] != yamlSeperator {
				dataValue += yamlSeperator
				dataValue += "\n"
			}
			dataValue += string(data)
		}
		if _, err = fmt.Fprintf(w, "%s", dataValue); err != nil {
			return err
		}
		w.Flush()
		return nil
	}

	tableTitle := "Package Information"
	table := uitable.Table{
		Title:     tableTitle,
		Content:   "PackageInstalls",
		Transpose: true,

		Header: []uitable.Header{
			uitable.NewHeader("Name"),
			uitable.NewHeader("Package Name"),
			uitable.NewHeader("Package Version"),
			uitable.NewHeader("Status"),
			uitable.NewHeader("Conditions"),
			uitable.NewHeader("Useful Error Message"),
		},

		SortBy: []uitable.ColumnSort{
			{Column: 0, Asc: true},
			{Column: 1, Asc: true},
		},
	}

	table.Rows = append(table.Rows, []uitable.Value{
		uitable.NewValueString(pkgi.Name),
		uitable.NewValueString(pkgi.Spec.PackageRef.RefName),
		uitable.NewValueString(pkgi.Status.Version),
		uitable.NewValueString(pkgi.Status.FriendlyDescription),
		uitable.NewValueInterface(pkgi.Status.Conditions),
		uitable.NewValueString(pkgi.Status.UsefulErrorMessage),
	})

	o.ui.PrintTable(table)

	return nil
}