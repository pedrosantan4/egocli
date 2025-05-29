# egocli 🦫📦🍵

**egocli** é um terminal interativo, modular e extensível escrito em **Go** 🦫, utilizando os frameworks **Cobra** 📦 e **Bubble Tea** 🍵. Este projeto visa fomentar a automação, padronização e produtividade na criação de comandos que geram rapidamente snippets de infraestrutura como código (IaC) para módulos AWS, como IAM, EKS, S3, RDS e muitos outros.

---

## ✨ Principais Recursos

- 🏗️ **Geração automática de snippets**: IAM, EKS, S3, RDS, EC2 e outros módulos AWS.
- ⚡ **Interface interativa**: via **Bubble Tea** 🍵 com prompts dinâmicos.
- 🧱 **Arquitetura modular**: novos comandos facilmente adicionáveis com **Cobra** 📦.
- 📊 **Exibição de métricas locais**: CPU, Memória, Disco e Tempo de execução.
- 🔄 **Evita repetição de código**: uso de constantes globais para mensagens e métricas.
- ✅ **Padrões de código**: clean code, DRY (Don't Repeat Yourself), KISS (Keep It Simple, Stupid).
- 🖥️ **Experiência CLI rica**: autocompletes, prompts, menus e navegação fluida.
- 📦 **Templates prontos**: modulos AWS com parametrização rápida e segura.
- 🌍 **Testes Multi-ambiente**: facilidade para gerar e validar código para ambientes `prod`, `homolog` e `dev`.

---

## ⚙️ Tecnologias Utilizadas

- **Go** 🦫 — linguagem robusta e rápida.
- **Cobra** 📦 — criação de comandos CLI.
- **Bubble Tea** 🍵 — interface TUI elegante e reativa.
- **Viper** — gerenciamento de configurações.
- **Go Templates** — geração dinâmica de arquivos e módulos.
- **OS/Runtime Packages** — coleta de métricas do sistema.

---

## 🚀 Como pode ajudar desenvolvedores e equipes

✅ **Agilidade** na criação de módulos e recursos AWS.  
✅ **Padronização** no provisionamento de infraestrutura.  
✅ **Evita erros manuais**: reduz a necessidade de escrever código repetitivo.  
✅ **Integração rápida** com pipelines DevOps e GitOps.  
✅ **Ambientes isolados**: facilita testes entre `prod` e `homolog`.  
✅ **Métricas em tempo real**: saiba como está o seu terminal enquanto trabalha!

---

## 📦 Estrutura Técnica

- `cmd/` — comandos do terminal organizados por módulo.
- `internal/templates/` — templates Go para geração de módulos AWS.
- `pkg/ui/` — componentes de interface usando Bubble Tea.
- `pkg/metrics/` — coleta e exibição de métricas do sistema.
- `constants.go` — mensagens e métricas centralizadas para evitar duplicações.
- `moduleTemplate.go` — template engine global reutilizável.
- `main.go` — ponto de entrada com inicialização da CLI Cobra.

---

## 🧭 Roadmap: Melhorias Futuras

🔧 **Agente Healer**: monitoramento contínuo do terminal para detectar inconsistências ou erros, com autocorreção.  
📊 **Observabilidade**: coleta avançada de métricas e logs, monitoramento de clusters/pods Kubernetes, com informações personalizadas.  
🧠 **Integração com IA/LLMs**: para sugestões inteligentes de comandos, aprendizado de padrões de uso e melhoria no relacionamento entre Terminal e Desenvolvedor.  
🌐 **Gerenciamento Multi-ambiente**: ampliar suporte e testes entre `prod`, `homolog`, `dev` e `staging` de forma segura e padronizada.  
🔌 **Plugins**: arquitetura para que terceiros criem e integrem novos módulos.  
📄 **Documentação interativa**: comandos com auto-ajuda detalhada e exemplos práticos.

---

## 📈 Métricas exibidas no terminal

- 🔋 Uso de CPU.
- 💾 Consumo de memória.
- 🗄️ Espaço em disco.
- ⏱️ Tempo de execução da aplicação.

Tudo isso ajuda você a entender o impacto e saúde da sua máquina enquanto utiliza o terminal.

---

## 🧑‍💻 Padrões de Código Seguidos

- ✅ Clean Code.
- ✅ DRY (Don't Repeat Yourself).
- ✅ KISS (Keep It Simple, Stupid).
- ✅ Modularização.
- ✅ Reutilização de templates e constantes globais.

---

## 🤝 Contribuindo

1. Fork o repositório.
2. Crie sua feature branch: `git checkout -b minha-feature`.
3. Commit suas alterações: `git commit -m 'Minha nova feature'`.
4. Push para a branch: `git push origin minha-feature`.
5. Abra um Pull Request!

---

## 🛠️ Instalação e uso

```bash
git clone https://github.com/seuusuario/egocli.git
cd egocli
go build -o egocli
./egocli
