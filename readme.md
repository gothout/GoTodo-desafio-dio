
# 🧩 API REST de Tarefas e Tags em Go

Este projeto foi desenvolvido como parte do desafio **"Criando a sua Primeira API Rest com Go"** da plataforma [DIO](https://www.dio.me/). Ele representa uma aplicação simples de gerenciamento de tarefas (`Tasks`) e suas respectivas categorias (`Tags`), utilizando exclusivamente a linguagem Go com manipulação de dados em memória.

## 📌 Funcionalidades

- ✅ Criar, atualizar, listar e deletar **Tarefas**
- ✅ Criar, atualizar, listar e deletar **Tags**
- ✅ Filtrar tarefas por Tag
- ✅ UI simples com HTML + Tailwind para interação com a API
- ✅ Atualizações em tempo real com feedbacks visuais
- ❌ Sem banco de dados (dados em memória durante execução)

## 🚀 Tecnologias Utilizadas

- [Golang](https://golang.org/)
- [Gorilla Mux](https://github.com/gorilla/mux) – roteador HTTP
- [TailwindCSS](https://tailwindcss.com/) – estilização frontend
- HTML + JavaScript puro (sem frameworks)

---

## 🛠️ Como Rodar o Projeto

### Pré-requisitos

- Go 1.18+ instalado
- Git (opcional)
- Navegador

### Passo a Passo

1. Clone este repositório ou baixe os arquivos manualmente:

```bash
git clone https://github.com/seu-usuario/api-tarefas-go.git
cd api-tarefas-go
````

2. Execute o projeto:

```bash
go run main.go
```

3. Acesse no navegador:

```
http://localhost:3000
```

A interface será exibida com opções para **gerenciar tarefas e tags** visualmente.

---

## 📂 Estrutura do Projeto

```
.
├── main.go          # Código principal da API REST
├── static/
│   └── index.html   # Interface frontend com TailwindCSS
└── README.md        # Documentação do projeto
```

---

## 🧪 Exemplos de Requisições API (Opcional)

Você pode testar os endpoints via Postman, Insomnia ou `curl`:

### 🔹 Criar Tag

```http
POST /tag
Content-Type: application/json

{
  "nome": "Urgente"
}
```

### 🔹 Criar Tarefa

```http
POST /task
Content-Type: application/json

{
  "nome": "Estudar Go",
  "descricao": "Finalizar o projeto da DIO",
  "tag": {
    "nome": "Estudo"
  }
}
```

---

## 🧠 Aprendizados

Este projeto reforça os seguintes conceitos do curso:

* Criação de rotas RESTful com Go
* Manipulação de structs e slices
* Codificação e decodificação JSON
* Separação de responsabilidades
* Comunicação entre frontend e backend
* Boas práticas de tratamento de erros