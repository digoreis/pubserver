# pubserver

Servidor privado de pacotes Dart/Flutter integrado ao GitLab, com análise automática de interfaces públicas via tree-sitter-dart.

---

## Visão geral

O **pubserver** é um servidor de publicação de pacotes Dart/Flutter que:

- **Publica pacotes** em releases no GitLab, anexando o tarball e metadados YAML.
- **Gera automaticamente** um arquivo `interfaces.json` para cada pacote, contendo todas as interfaces públicas (classes, mixins, extensions, typedefs, métodos e atributos) extraídas via [tree-sitter-dart](https://github.com/UserNobody14/tree-sitter-dart).
- **Disponibiliza uma API REST** para upload, download e consulta dessas interfaces.
- **Protege todas as rotas por autenticação** (token).

---

## Funcionalidades principais

- **Publicação de pacotes:** Recebe uploads, valida, gera metadados, cria release no GitLab.
- **Análise de interfaces:** Usa tree-sitter-dart para extrair e publicar `interfaces.json`.
- **API REST:** Endpoints para publicação, download do pacote, download de `interfaces.json`, listagem de versões etc.
- **Autenticação:** JWT ou token de API.
- **Armazenamento:** Usa releases e assets do GitLab como backend.

---

## Como instalar e rodar localmente

Veja [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) para instruções completas.

### Resumo rápido

```bash
git clone https://github.com/seu-usuario/pubserver.git
cd pubserver

# Instale dependências Go
go mod tidy

# (Opcional) Instale tree-sitter CLI para debug
npm install -g tree-sitter-cli
git clone https://github.com/UserNobody14/tree-sitter-dart.git vendor/tree-sitter-dart
cd vendor/tree-sitter-dart
git checkout <commit-usado-no-go.mod>
cd ../..

# Configure as variáveis de ambiente (ou .env)
export GITLAB_TOKEN=seu_token
export GITLAB_HOST=https://gitlab.com
export GITLAB_PROJECT=grupo/repo
export PUBSERVER_PORT=8080

# Rode o servidor
go run cmd/pubserver/main.go
```

---

## Exemplos de uso da API

### 1. Publicar um pacote

Faça um POST com os campos e o arquivo tar.gz para `/api/packages/versions/newUploadFinish` (ver exemplos no [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md)).

### 2. Baixar interfaces públicas do pacote

```bash
curl -H "Authorization: Bearer SEU_TOKEN" \
  http://localhost:8080/api/packages/NOME/versions/VERSAO/interfaces
```

---

## Estrutura do JSON de interfaces

```json
[
  {
    "name": "MyService",
    "type": "class",
    "superclass": "BaseService",
    "interfaces": ["MyInterface"],
    "fields": [
      {
        "name": "value",
        "type": "int",
        "isStatic": false,
        "isFinal": false,
        "isConst": false
      }
    ],
    "methods": [
      {
        "name": "doSomething",
        "returnType": "void",
        "isStatic": false,
        "isAbstract": false,
        "parameters": [
          { "name": "x", "type": "int", "required": true }
        ]
      }
    ]
  }
]
```

---

## Roadmap

- [ ] Suporte a outros linguagens além de Dart.
- [ ] Melhorias de performance para grandes pacotes.
- [ ] Integração com sistemas de permissão do GitLab.

---

## Contribua

Pull requests e sugestões são sempre bem-vindos! Veja [CONTRIBUTING.md](CONTRIBUTING.md) se existir.

---

## Licença

MIT

---

## Créditos

- [tree-sitter-dart](https://github.com/UserNobody14/tree-sitter-dart)
- [go-tree-sitter](https://github.com/smacker/go-tree-sitter)
