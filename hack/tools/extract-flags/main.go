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

// extract-flags captures the --help output of every Karmada component and its
// subcommands, writing one .txt file per command into an output directory.
//
// Usage:
//
//	go run hack/tools/extract-flags/main.go [output-dir]
//
// The default output directory is docs/command-flags.
// Files are named after the command path with spaces replaced by underscores,
// e.g. "karmada-controller-manager.txt" and "karmada-controller-manager_version.txt".
//
// The generated files are checked into the repository and verified in CI by
// hack/verify-command-flags.sh to detect unintentional flag changes.
package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
	_ "k8s.io/component-base/logs/json/register"

	agentapp "github.com/karmada-io/karmada/cmd/agent/app"
	aaapp "github.com/karmada-io/karmada/cmd/aggregated-apiserver/app"
	cmapp "github.com/karmada-io/karmada/cmd/controller-manager/app"
	deschapp "github.com/karmada-io/karmada/cmd/descheduler/app"
	searchapp "github.com/karmada-io/karmada/cmd/karmada-search/app"
	adapterapp "github.com/karmada-io/karmada/cmd/metrics-adapter/app"
	estiapp "github.com/karmada-io/karmada/cmd/scheduler-estimator/app"
	schapp "github.com/karmada-io/karmada/cmd/scheduler/app"
	webhookapp "github.com/karmada-io/karmada/cmd/webhook/app"
	"github.com/karmada-io/karmada/pkg/util/names"
)

const (
	// defaultOutputDir matches the repo layout for generated flag docs.
	defaultOutputDir = "docs/command-flags"

	dirPerm  = 0o755
	filePerm = 0o644
)

// componentSpec maps a component name to its cobra root command constructor.
// Add new entries here when documenting flags for additional binaries.
var componentSpecs = []struct {
	name string
	cmd  func(context.Context) *cobra.Command
}{
	{names.KarmadaControllerManagerComponentName, func(ctx context.Context) *cobra.Command {
		return cmapp.NewControllerManagerCommand(ctx)
	}},
	{names.KarmadaSchedulerComponentName, func(ctx context.Context) *cobra.Command {
		return schapp.NewSchedulerCommand(ctx)
	}},
	{names.KarmadaAgentComponentName, func(ctx context.Context) *cobra.Command {
		return agentapp.NewAgentCommand(ctx)
	}},
	{names.KarmadaAggregatedAPIServerComponentName, func(ctx context.Context) *cobra.Command {
		return aaapp.NewAggregatedApiserverCommand(ctx)
	}},
	{names.KarmadaDeschedulerComponentName, func(ctx context.Context) *cobra.Command {
		return deschapp.NewDeschedulerCommand(ctx)
	}},
	{names.KarmadaSearchComponentName, func(ctx context.Context) *cobra.Command {
		return searchapp.NewKarmadaSearchCommand(ctx)
	}},
	{names.KarmadaSchedulerEstimatorComponentName, func(ctx context.Context) *cobra.Command {
		return estiapp.NewSchedulerEstimatorCommand(ctx)
	}},
	{names.KarmadaWebhookComponentName, func(ctx context.Context) *cobra.Command {
		return webhookapp.NewWebhookCommand(ctx)
	}},
	{names.KarmadaMetricsAdapterComponentName, func(ctx context.Context) *cobra.Command {
		return adapterapp.NewMetricsAdapterCommand(ctx)
	}},
}

// captureHelp captures the --help output of a cobra command by redirecting
// stdout to a buffer. The returned bytes are exactly what a user would see
// when running "<binary> --help" or "<binary> <subcommand> --help".
func captureHelp(cmd *cobra.Command) ([]byte, error) {
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	if err := cmd.Help(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// buildFilename converts a command path like "karmada-controller-manager version"
// to a filename like "karmada-controller-manager_version.txt".
// The root command (single word) produces "karmada-controller-manager.txt".
func buildFilename(cmdPath string) string {
	return strings.ReplaceAll(cmdPath, " ", "_") + ".txt"
}

// writeCommandTree recursively captures help text for cmd and all its
// available subcommands, writing one file per command.
func writeCommandTree(outputDir string, cmd *cobra.Command) error {
	content, err := captureHelp(cmd)
	if err != nil {
		return fmt.Errorf("help for %q: %w", cmd.CommandPath(), err)
	}

	outputPath := filepath.Join(outputDir, buildFilename(cmd.CommandPath()))
	if err := os.WriteFile(outputPath, content, filePerm); err != nil {
		return fmt.Errorf("write %s: %w", outputPath, err)
	}
	fmt.Printf("Wrote %s\n", outputPath)

	for _, sub := range cmd.Commands() {
		if !sub.IsAvailableCommand() && !sub.IsAdditionalHelpTopicCommand() {
			continue
		}
		if err := writeCommandTree(outputDir, sub); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "Usage: extract-flags [output-dir]")
		fmt.Fprintf(os.Stderr, "  default output-dir: %s\n", defaultOutputDir)
		os.Exit(1)
	}

	outputDir := defaultOutputDir
	if len(os.Args) == 2 {
		outputDir = os.Args[1]
	}

	if err := os.MkdirAll(outputDir, dirPerm); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir %s: %v\n", outputDir, err)
		os.Exit(1)
	}

	ctx := context.TODO()

	var genErrs []error
	for _, spec := range componentSpecs {
		cmd := spec.cmd(ctx)
		// Ensure default subcommands (help, completion) are registered before
		// we walk the command tree.
		cmd.InitDefaultHelpCmd()
		cmd.InitDefaultCompletionCmd()
		cmd.SetGlobalNormalizationFunc(cliflag.WordSepNormalizeFunc)

		if err := writeCommandTree(outputDir, cmd); err != nil {
			genErrs = append(genErrs, fmt.Errorf("%s: %w", spec.name, err))
		}
	}
	if err := errors.Join(genErrs...); err != nil {
		fmt.Fprintf(os.Stderr, "extract-flags: finished with errors:\n%v\n", err)
		os.Exit(1)
	}
}
