package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

// Defina as flags como variáveis globais
var (
	vpcNewFlag    bool
	eksNewFlag    bool
	rdsNewFlag    bool
	s3NewFlag     bool
	iamNewFlag    bool
	lambdaNewFlag bool
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Cria novos templates de código",
	Long: `Cria arquivos com estrutura inicial para:
- Lambda Functions
- Templates Terraform
- Configurações AWS`,
	Example: `  egocli new --lambda  # Cria template Lambda
  egocli new --vpc     # Cria template VPC`,
}

func init() {
	// Configure as flags
	newCmd.Flags().BoolVarP(&vpcNewFlag, "vpc", "v", false, "Template para VPC")
	newCmd.Flags().BoolVarP(&eksNewFlag, "eks", "e", false, "Template para EKS")
	newCmd.Flags().BoolVarP(&rdsNewFlag, "rds", "r", false, "Template para RDS")
	newCmd.Flags().BoolVarP(&s3NewFlag, "s3", "s", false, "Template para S3")
	newCmd.Flags().BoolVarP(&iamNewFlag, "iam", "i", false, "Template para IAM")
	newCmd.Flags().BoolVarP(&lambdaNewFlag, "lambda", "l", false, "Template para Lambda")

	// Registre o comando
	rootCmd.AddCommand(newCmd)
	newCmd.Run = newCommand
}

func newCommand(cmd *cobra.Command, args []string) {
	start := time.Now()
	memBefore := GetMemoryUsage()

	// Mapeamento de flags para módulos usando templates.go
	flagMappings := []struct {
		flag   *bool
		module string
	}{
		{&vpcNewFlag, "vpc"},
		{&eksNewFlag, "eks"},
		{&rdsNewFlag, "rds"},
		{&s3NewFlag, "s3"},
		{&iamNewFlag, "iam"},
		{&lambdaNewFlag, "lambda"},
	}

	// Processar todos os templates selecionados
	for _, mapping := range flagMappings {
		if *mapping.flag {
			template, exists := Templates[mapping.module]
			if !exists {
				fmt.Printf("❌ Template não encontrado: %s\n", mapping.module)
				continue
			}

			// Usar diretório específico para new (snippets)
			snippetDir := filepath.Join(newDir, template.DirName)
			CreateTemplate(snippetDir, template.FileName, template.Content)
		}
	}

	PrintOperationStats(start, memBefore)
}

func CreateTemplate(dir, filename, content string) {
	fullPath := filepath.Join(dir, filename)

	// Criar diretório se não existir
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("❌ Erro ao criar diretório: %v\n", err)
		return
	}

	// Verificar se arquivo já existe
	if _, err := os.Stat(fullPath); err == nil {
		fmt.Printf("⚠️  Arquivo já existe: %s\n", fullPath)
		return
	}

	// Criar arquivo
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		fmt.Printf("❌ Erro ao criar arquivo: %v\n", err)
		return
	}

	// Tentar abrir na IDE
	if err := openInEditor(fullPath); err != nil {
		fmt.Printf("✅ Arquivo criado: %s\n", fullPath)
		fmt.Printf("⚠️  Não foi possível abrir no editor: %v\n", err)
		return
	}

	fmt.Printf("✅ Template criado e aberto: %s\n", fullPath)
}

// Funções auxiliares mantidas:
func GetMemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

func PrintOperationStats(start time.Time, memBefore uint64) {
	duration := time.Since(start)
	memAfter := GetMemoryUsage()
	memUsed := float64(memAfter-memBefore) / 1024 / 1024

	fmt.Printf("\n📊 Estatísticas:\n")
	fmt.Printf("⏱️  Duração: %v\n", duration.Round(time.Millisecond))
	fmt.Printf("💾 Memória usada: %.2fMB\n", memUsed)
	fmt.Printf("💵 Custo estimado: $%.6f\n", CalculateCost(duration, memUsed))
	fmt.Println("🌐 Modo: Local (sem deploy na cloud)")
}

func CalculateCost(duration time.Duration, memMB float64) float64 {
	return (0.0000166667 * (memMB / 1024) * duration.Seconds())
}

func openInEditor(path string) error {
	if err := exec.Command("code", path).Start(); err == nil {
		return nil
	}

	editors := []string{"subl", "gedit", "nano", "vim"}
	for _, editor := range editors {
		if err := exec.Command(editor, path).Start(); err == nil {
			return nil
		}
	}
	return fmt.Errorf("editor não encontrado")
}
