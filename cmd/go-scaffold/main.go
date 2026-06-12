package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/isaydiev86/go-scaffold/internal/generator"
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
	var withRedis bool
	var withKafkaConsumer bool
	var withKafkaProducer bool
	var nonInteractive bool

	cmd := &cobra.Command{
		Use:          "new [service-name]",
		Short:        "Generate a new Go service from the Clean Architecture template",
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var serviceName string

			if len(args) == 1 {
				serviceName = strings.ToLower(args[0])
			}

			if nonInteractive {
				// CI mode: no prompts, missing required values are an error,
				// boolean flags keep their defaults (false) unless passed explicitly.
				if serviceName == "" {
					return fmt.Errorf("service name is required in non-interactive mode: go-scaffold new <service-name> --non-interactive")
				}
				if moduleName == "" {
					return fmt.Errorf("--module is required in non-interactive mode")
				}
				if outputDir == "" {
					outputDir = serviceName
				}
			} else {
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

				if !cmd.Flags().Changed("redis") {
					answer := prompt(reader, "Add Redis? (caching example for GET) y/n", "n")
					withRedis = isYes(answer)
				}

				if !cmd.Flags().Changed("kafka-consumer") {
					answer := prompt(reader, "Add Kafka consumer? y/n", "n")
					withKafkaConsumer = isYes(answer)
				}

				if !cmd.Flags().Changed("kafka-producer") {
					answer := prompt(reader, "Add Kafka producer? (transactional outbox + worker) y/n", "n")
					withKafkaProducer = isYes(answer)
				}
			}

			if err := generator.ValidateServiceName(serviceName); err != nil {
				return err
			}

			if err := generator.ValidateModulePath(moduleName); err != nil {
				return err
			}

			if _, err := os.Stat(outputDir); err == nil {
				return fmt.Errorf("directory %q already exists", outputDir)
			}

			vars := generator.NewVars(serviceName, moduleName, generator.Options{
				WithRedis:         withRedis,
				WithKafkaConsumer: withKafkaConsumer,
				WithKafkaProducer: withKafkaProducer,
			})

			fmt.Printf("\nGenerating service %q → %s\n", serviceName, outputDir)

			if err := generator.Generate(vars, outputDir); err != nil {
				return fmt.Errorf("generate: %w", err)
			}

			services := "postgres"
			if withRedis {
				services += " + redis"
			}
			if withKafkaConsumer || withKafkaProducer {
				services += " + kafka"
			}

			fmt.Printf("\nDone! Next steps:\n")
			fmt.Printf("  cd %s\n", outputDir)
			fmt.Printf("  make bootstrap   # go mod tidy + go generate (mocks, OpenAPI) — required before the first build\n")
			fmt.Printf("  go test ./...\n")
			fmt.Printf("  make up          # start %s\n", services)
			fmt.Printf("  make migrate-up  # apply migrations\n")
			fmt.Printf("  make run         # start the service\n")

			return nil
		},
	}

	cmd.Flags().StringVar(&moduleName, "module", "", "Go module path (default: github.com/myorg/<service-name>)")
	cmd.Flags().StringVar(&outputDir, "output", "", "Output directory (default: service name)")
	cmd.Flags().BoolVar(&withRedis, "redis", false, "Include Redis (cache adapter + caching example for GET)")
	cmd.Flags().BoolVar(&withKafkaConsumer, "kafka-consumer", false, "Include Kafka consumer controller")
	cmd.Flags().BoolVar(&withKafkaProducer, "kafka-producer", false, "Include Kafka producer (transactional outbox + worker)")
	cmd.Flags().BoolVar(&nonInteractive, "non-interactive", false, "Never prompt: fail on missing required values, boolean flags default to false (for CI)")

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

// isYes reports whether the answer is an affirmative "y"/"yes" (case-insensitive).
func isYes(answer string) bool {
	return strings.EqualFold(answer, "y") || strings.EqualFold(answer, "yes")
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
