// cmd/gen.go
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate infrastructure components",
}

func init() {
	// Subcomandos
	genCmd.AddCommand(vpcCmd)
	genCmd.AddCommand(eksCmd)
	genCmd.AddCommand(iamCmd)
	genCmd.AddCommand(rdsCmd)
	genCmd.AddCommand(s3Cmd)
	genCmd.AddCommand(lambdaCmd)
	rootCmd.AddCommand(genCmd)
}

// Subcomando para VPC
var vpcCmd = &cobra.Command{
	Use:   "vpc",
	Short: "Generate VPC configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if err := generateInfra("vpc", genDir); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ VPC generated successfully")
	},
}

// Subcomando para EKS
var eksCmd = &cobra.Command{
	Use:   "eks",
	Short: "Generate EKS configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if err := generateInfra("eks", genDir); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ EKS generated successfully")
	},
}

// Subcomando para IAM - CORRIGIDO
var iamCmd = &cobra.Command{
	Use:   "iam",
	Short: "Generate IAM roles configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if err := generateInfra("iam", genDir); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ IAM roles generated successfully")
	},
}

// Subcomando para RDS - CORRIGIDO
var rdsCmd = &cobra.Command{
	Use:   "rds",
	Short: "Generate RDS database configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if err := generateInfra("rds", genDir); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ RDS database generated successfully")
	},
}

// Subcomando para S3 - CORRIGIDO
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "Generate S3 bucket configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if err := generateInfra("s3", genDir); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ S3 bucket generated successfully")
	},
}

// Subcomando para Lambda
var lambdaCmd = &cobra.Command{
	Use:   "lambda",
	Short: "Generate lambda configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if err := generateInfra("lambda", genDir); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Lambda function generated successfully")
	},
}

// Lógica unificada - CORRIGIDA para retornar error
func generateInfra(module string, outputDir string) error {
	template, exists := Templates[module]
	if !exists {
		return fmt.Errorf("unknown module: %s", module)
	}

	modulePath := filepath.Join(outputDir, template.DirName)
	outputPath := filepath.Join(modulePath, template.FileName)

	if err := os.MkdirAll(modulePath, 0755); err != nil {
		return fmt.Errorf("couldn't create directory: %w", err)
	}

	if !confirmOverwrite(outputPath) {
		return fmt.Errorf("operation cancelled by user")
	}

	if err := os.WriteFile(outputPath, []byte(template.Content), 0644); err != nil {
		return fmt.Errorf("failed to generate %s: %w", module, err)
	}

	fmt.Printf("\n✅ Generated %s module\n📁 Location: %s\n", module, outputPath)
	return nil
}

// Função utilitária compartilhada
func confirmOverwrite(path string) bool {
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("⚠️  %s exists. Overwrite? (y/n): ", filepath.Base(path))
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		return strings.EqualFold(strings.TrimSpace(input), "y")
	}
	return true
}
