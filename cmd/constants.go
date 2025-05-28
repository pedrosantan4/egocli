// cmd/constants.go
package cmd

import (
	"time"
)

// ============== DIRETÓRIOS ==============
const (
	// Diretório para comandos gen (infraestrutura)
	genDir = "infra"

	// Diretório para comandos new (snippets)
	newDir = "mySnippets"

	// Diretório para backup de arquivos
	backupDir = "backup"
)

// ============== CONFIGURAÇÕES DE TEMPO ==============
const (
	// Intervalo de atualização de métricas no terminal
	metricsInterval = 1 * time.Second

	// Intervalo de piscar do cursor
	cursorInterval = 500 * time.Millisecond

	// Timeout para operações de arquivo
	fileTimeout = 30 * time.Second
)

// ============== CONFIGURAÇÕES DE ARQUIVO ==============
const (
	// Permissões padrão para diretórios
	dirPermissions = 0755

	// Permissões padrão para arquivos
	filePermissions = 0644

	// Extensão padrão para templates Terraform
	terraformExt = ".tf"

	// Extensão padrão para funções Lambda
	lambdaExt = ".js"
)

// ============== MENSAGENS PADRÃO ==============
const (
	// Mensagem de sucesso
	successMsg = "✅ %s generated successfully"

	// Mensagem de erro
	errorMsg = "❌ Error: %v"

	// Mensagem de confirmação
	confirmMsg = "⚠️  %s exists. Overwrite? (y/n): "

	// Mensagem de localização
	locationMsg = "\n✅ Generated %s module\n📁 Location: %s\n"
)
