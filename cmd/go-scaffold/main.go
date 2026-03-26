package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/isaydiev/go-scaffold/internal/generator"
)

func main() {
	root := &cobra.Command{
		Use:   "go-scaffold",
		Short: "Go service scaffolding tool",
	}

	root.AddCommand(newCmd())
	root.AddCommand(versionCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func newCmd() *cobra.Command {
	var moduleName string
	var outputDir string

	cmd := &cobra.Command{
		Use:   "new [service-name]",
		Short: "Generate a new Go service from the Clean Architecture template",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var serviceName string

			if len(args) == 1 {
				serviceName = strings.ToLower(args[0])
			}

			// Interactive mode: ask missing values
			reader := bufio.NewReader(os.Stdin)

			if serviceName == "" {
				serviceName = prompt(reader, "Service name", "")
				if serviceName == "" {
					return fmt.Errorf("service name is required")
				}
				serviceName = strings.ToLower(serviceName)
			}

			if moduleName == "" {
				defaultModule := "github.com/myorg/" + serviceName
				moduleName = prompt(reader, "Go module path", defaultModule)
				if moduleName == "" {
					moduleName = defaultModule
				}
			}

			if outputDir == "" {
				outputDir = prompt(reader, "Output directory", serviceName)
				if outputDir == "" {
					outputDir = serviceName
				}
			}

			if _, err := os.Stat(outputDir); err == nil {
				return fmt.Errorf("directory %q already exists", outputDir)
			}

			vars := generator.NewVars(serviceName, moduleName)

			fmt.Printf("\nGenerating service %q → %s\n", serviceName, outputDir)

			if err := generator.Generate(vars, outputDir); err != nil {
				return fmt.Errorf("generate: %w", err)
			}

			fmt.Printf("\nDone! Next steps:\n")
			fmt.Printf("  cd %s\n", outputDir)
			fmt.Printf("  go mod tidy\n")
			fmt.Printf("  make up          # start postgres + redis\n")
			fmt.Printf("  make migrate-up  # apply migrations\n")
			fmt.Printf("  make run         # start the service\n")

			return nil
		},
	}

	cmd.Flags().StringVar(&moduleName, "module", "", "Go module path (default: github.com/myorg/<service-name>)")
	cmd.Flags().StringVar(&outputDir, "output", "", "Output directory (default: service name)")

	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("go-scaffold v0.1.0")
		},
	}
}

// prompt prints a question with an optional default value and reads the user's answer.
// If the user presses Enter without typing anything, the default is returned.
func prompt(reader *bufio.Reader, question, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("? %s [%s]: ", question, defaultVal)
	} else {
		fmt.Printf("? %s: ", question)
	}

	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(answer)

	if answer == "" {
		return defaultVal
	}

	return answer
}
