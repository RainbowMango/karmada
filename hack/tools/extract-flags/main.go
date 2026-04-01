/*
Copyright 2026 The Karmada Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	cliflag "k8s.io/component-base/cli/flag"
	_ "k8s.io/component-base/logs/json/register"
	controllerruntime "sigs.k8s.io/controller-runtime"

	controllermanagerapp "github.com/karmada-io/karmada/cmd/controller-manager/app"
	"github.com/karmada-io/karmada/pkg/util/names"
)

const (
	// argvOutputDirOptional is os.Args[1], the optional output directory when the user passes one argument after the program name.
	argvOutputDirOptional = 1
	// maxArgCount is the maximum number of os.Args entries we accept: program name plus at most one path.
	maxArgCount = 2

	// defaultOutputDir is used when no output directory is given; it matches the repo layout for generated flag docs.
	defaultOutputDir = "docs/command-flags"

	dirPerm  = 0o755
	filePerm = 0o600
)

// componentSpec wires a canonical component name (see pkg/util/names) to its Cobra root command.
// Add new entries here when documenting flags for additional binaries.
var componentSpecs = []struct {
	name string
	new  func(context.Context) *cobra.Command
}{
	{
		name: names.KarmadaControllerManagerComponentName,
		new:  controllermanagerapp.NewControllerManagerCommand,
	},
}

// formatDeprecatedFlags extracts and formats deprecated flags from a command and its subcommands.
func formatDeprecatedFlags(cmd *cobra.Command) string {
	var collect func(*cobra.Command) map[string]string
	collect = func(c *cobra.Command) map[string]string {
		deprecated := make(map[string]string)
		c.Flags().VisitAll(func(flag *pflag.Flag) {
			if flag.Deprecated != "" {
				deprecated[flag.Name] = flag.Deprecated
			}
		})
		for _, sub := range c.Commands() {
			maps.Copy(deprecated, collect(sub))
		}
		return deprecated
	}
	deprecated := collect(cmd)
	if len(deprecated) == 0 {
		return ""
	}
	names := make([]string, 0, len(deprecated))
	for name := range deprecated {
		names = append(names, name)
	}
	sort.Strings(names)
	var lines []string
	lines = append(lines, "", "Deprecated flags:", "")
	for _, flagName := range names {
		// Label line, message line with fixed padding (cobra help–like alignment), then "" for a blank line between entries.
		lines = append(lines, fmt.Sprintf("      [DEPRECATED] --%s", flagName), fmt.Sprintf("                           %s", deprecated[flagName]), "")
	}
	return strings.Join(lines, "\n")
}

func writeComponentHelp(outputDir string, componentName string, cmd *cobra.Command) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultCompletionCmd()
	cmd.SetGlobalNormalizationFunc(cliflag.WordSepNormalizeFunc)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	if err := cmd.Help(); err != nil {
		return fmt.Errorf("help for %s: %w", componentName, err)
	}
	content := buf.Bytes()
	if s := formatDeprecatedFlags(cmd); s != "" {
		content = append(content, []byte(s)...)
	}

	outputPath := filepath.Join(outputDir, componentName+".txt")
	if err := os.WriteFile(outputPath, content, filePerm); err != nil {
		return fmt.Errorf("write %s: %w", outputPath, err)
	}
	fmt.Printf("Wrote %s\n", outputPath)
	return nil
}

func main() {
	if len(os.Args) > maxArgCount {
		fmt.Fprintln(os.Stderr, "Usage: extract-flags [output-dir]")
		fmt.Fprintf(os.Stderr, "  default output-dir: %s\n", defaultOutputDir)
		os.Exit(1)
	}

	outputDir := defaultOutputDir
	// len(os.Args) is 1 when only argvProgram is present; otherwise the user supplied an explicit output directory.
	if len(os.Args) > argvOutputDirOptional {
		outputDir = os.Args[argvOutputDirOptional]
	}

	if err := os.MkdirAll(outputDir, dirPerm); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir %s: %v\n", outputDir, err)
		os.Exit(1)
	}

	ctx := controllerruntime.SetupSignalHandler()

	var genErrs []error
	for _, spec := range componentSpecs {
		cmd := spec.new(ctx)
		if err := writeComponentHelp(outputDir, spec.name, cmd); err != nil {
			genErrs = append(genErrs, fmt.Errorf("%s: %w", spec.name, err))
		}
	}
	if err := errors.Join(genErrs...); err != nil {
		fmt.Fprintf(os.Stderr, "extract-flags: finished with errors:\n%v\n", err)
		os.Exit(1)
	}
}
