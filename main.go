package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/engine"
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

	// Set up engine.
	e := engine.New()
	out, err := e.Render(c, map[string]interface{}{})
	if err != nil {
		return err
	}

	for name, data := range out {
		fmt.Printf("---\n# $s\n", name)
		fmt.Println(data)
	}
	return nil
}
