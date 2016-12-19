package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/engine"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/timeconv"
)

const globalUsage = `
Render chart templates locally and display the output.

This does not require Tiller. However, any values that would normally be
looked up or retrieved in-cluster will be faked locally. Additionally, none
of the server-side testing of chart validity (e.g. whether an API is supported)
is done.
`

func main() {
	cmd := &cobra.Command{
		Use:   "template [flags] CHART",
		Short: "locally render templates",
		Long:  globalUsage,
		RunE:  run,
	}
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("chart is required")
	}
	c, err := chartutil.Load(args[0])
	if err != nil {
		return err
	}

	config := &chart.Config{Raw: "", Values: map[string]*chart.Value{}}

	options := chartutil.ReleaseOptions{
		Name:      "RELEASE-NAME",
		Time:      timeconv.Now(),
		Namespace: "NAMESPACE",
		//Revision:  1,
		//IsInstall: true,
	}

	// Set up engine.
	renderer := engine.New()

	vals, err := chartutil.ToRenderValues(c, config, options)
	if err != nil {
		return err
	}

	out, err := renderer.Render(c, vals)
	if err != nil {
		return err
	}

	for name, data := range out {
		b := filepath.Base(name)
		if strings.HasPrefix(b, "_") {
			continue
		}
		fmt.Printf("---\n# %s\n", name)
		fmt.Println(data)
	}
	return nil
}
