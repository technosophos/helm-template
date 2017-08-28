package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/engine"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/strvals"
	"k8s.io/helm/pkg/timeconv"
)

const globalUsage = `
Render chart templates locally and display the output.

This does not require Tiller. However, any values that would normally be
looked up or retrieved in-cluster will be faked locally. Additionally, none
of the server-side testing of chart validity (e.g. whether an API is supported)
is done.

To render just one template in a chart, use '-x':

	$ helm template mychart -x mychart/templates/deployment.yaml
`

var (
	setVals     []string
	valsFiles   valueFiles
	flagVerbose bool
	showNotes   bool
	releaseName string
	namespace   string
	renderFiles []string
)

var version = "DEV"

func main() {
	cmd := &cobra.Command{
		Use:   "template [flags] CHART",
		Short: fmt.Sprintf("locally render templates (helm-template %s)", version),
		RunE:  run,
	}

	f := cmd.Flags()
	f.StringArrayVar(&setVals, "set", []string{}, "set values on the command line. See 'helm install -h'")
	f.VarP(&valsFiles, "values", "f", "specify one or more YAML files of values")
	f.BoolVarP(&flagVerbose, "verbose", "v", false, "show the computed YAML values as well.")
	f.BoolVar(&showNotes, "notes", false, "show the computed NOTES.txt file as well.")
	f.StringVarP(&releaseName, "release", "r", "RELEASE-NAME", "release name")
	f.StringVarP(&namespace, "namespace", "n", "NAMESPACE", "namespace")
	f.StringArrayVarP(&renderFiles, "execute", "x", []string{}, "only execute the given templates.")

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

	vv, err := vals()
	if err != nil {
		return err
	}

	config := &chart.Config{Raw: string(vv), Values: map[string]*chart.Value{}}

	if err = chartutil.ProcessRequirementsEnabled(c, config); err != nil {
		return err
	}

	if flagVerbose {
		fmt.Println("---\n# merged values")
		fmt.Println(string(vv))
	}

	options := chartutil.ReleaseOptions{
		Name:      releaseName,
		Time:      timeconv.Now(),
		Namespace: namespace,
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

	in := func(needle string, haystack []string) bool {
		for _, h := range haystack {
			if h == needle {
				return true
			}
		}
		return false
	}

	sortedKeys := make([]string, 0, len(out))
	for key := range out {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	// If renderFiles is set, we ONLY print those.
	if len(renderFiles) > 0 {
		for _, name := range sortedKeys {
			data := out[name]
			if in(name, renderFiles) {
				fmt.Printf("---\n# Source: %s\n", name)
				fmt.Println(data)
			}
		}
		return nil
	}

	for _, name := range sortedKeys {
		data := out[name]
		b := filepath.Base(name)
		if !showNotes && b == "NOTES.txt" {
			continue
		}
		if strings.HasPrefix(b, "_") {
			continue
		}
		fmt.Printf("---\n# Source: %s\n", name)
		fmt.Println(data)
	}
	return nil
}

// liberally borrows from Helm
func vals() ([]byte, error) {
	base := map[string]interface{}{}

	// User specified a values files via -f/--values
	for _, filePath := range valsFiles {
		currentMap := map[string]interface{}{}
		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return []byte{}, err
		}

		if err := yaml.Unmarshal(bytes, &currentMap); err != nil {
			return []byte{}, fmt.Errorf("failed to parse %s: %s", filePath, err)
		}
		// Merge with the previous map
		base = mergeValues(base, currentMap)
	}

	// User specified a value via --set
	for _, value := range setVals {
		if err := strvals.ParseInto(value, base); err != nil {
			return []byte{}, fmt.Errorf("failed parsing --set data: %s", err)
		}
	}

	return yaml.Marshal(base)
}

// Copied from Helm.

func mergeValues(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = v
			continue
		}
		nextMap, ok := v.(map[string]interface{})
		// If it isn't another map, overwrite the value
		if !ok {
			dest[k] = v
			continue
		}
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = nextMap
			continue
		}
		// Edge case: If the key exists in the destination, but isn't a map
		destMap, isMap := dest[k].(map[string]interface{})
		// If the source map has a map for this key, prefer it
		if !isMap {
			dest[k] = v
			continue
		}
		// If we got to this point, it is a map in both, so merge them
		dest[k] = mergeValues(destMap, nextMap)
	}
	return dest
}

type valueFiles []string

func (v *valueFiles) String() string {
	return fmt.Sprint(*v)
}

func (v *valueFiles) Type() string {
	return "valueFiles"
}

func (v *valueFiles) Set(value string) error {
	for _, filePath := range strings.Split(value, ",") {
		*v = append(*v, filePath)
	}
	return nil
}
