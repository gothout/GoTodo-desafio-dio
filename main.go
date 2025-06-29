package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// ==========================
// Structs
// ==========================

type Tasks struct {
	Task []Task `json:"tasks"`
}

type Task struct {
	Id        int64  `json:"id"`
	Nome      string `json:"nome" binding:"required"`
	Descricao string `json:"descricao" binding:"required"`
	Tag       Tag    `json:"tag"`
}

type Tag struct {
	Id   int64  `json:"id"`
	Nome string `json:"nome" binding:"required"`
}

// ==========================
// In-memory storage
// ==========================

var tasks Tasks
var tag []Tag

// ==========================
// Handlers de TAG
// ==========================

func GetTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(tag)
	if err != nil {
		http.Error(w, "Erro ao retornar as tags", http.StatusInternalServerError)
	}
}

func GetTagByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	for _, t := range tag {
		if t.Id == id {
			json.NewEncoder(w).Encode(t)
			return
		}
	}
	http.Error(w, "Tag não encontrada", http.StatusNotFound)
}

func CreateTag(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var novaTag Tag
	if err := json.NewDecoder(r.Body).Decode(&novaTag); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	if novaTag.Nome == "" {
		http.Error(w, "Campo 'nome' é obrigatório", http.StatusBadRequest)
		return
	}

	for _, t := range tag {
		if t.Nome == novaTag.Nome {
			http.Error(w, "Tag com esse nome já existe", http.StatusConflict)
			return
		}
	}

	var lastID int64
	for _, t := range tag {
		if t.Id > lastID {
			lastID = t.Id
		}
	}
	novaTag.Id = lastID + 1
	tag = append(tag, novaTag)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(novaTag)
}

func UpdateTag(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var novaTag Tag
	if err := json.NewDecoder(r.Body).Decode(&novaTag); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if novaTag.Nome == "" {
		http.Error(w, "Campo 'nome' é obrigatório", http.StatusBadRequest)
		return
	}

	// Verifica duplicidade ignorando o próprio ID
	for _, t := range tag {
		if strings.EqualFold(t.Nome, novaTag.Nome) && t.Id != id {
			http.Error(w, "Já existe uma tag com esse nome", http.StatusConflict)
			return
		}
	}

	// Atualiza a tag
	atualizada := false
	for i, t := range tag {
		if t.Id == id {
			tag[i].Nome = novaTag.Nome
			atualizada = true
			break
		}
	}

	if !atualizada {
		http.Error(w, "Tag não encontrada", http.StatusNotFound)
		return
	}

	// Atualiza nome da tag em todas as tarefas que a utilizam
	for i, t := range tasks.Task {
		if t.Tag.Id == id {
			tasks.Task[i].Tag.Nome = novaTag.Nome
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(novaTag)
}

func DeleteTag(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	for _, t := range tasks.Task {
		if t.Tag.Id == id {
			http.Error(w, "Não é possível deletar. Existem tarefas com essa tag.", http.StatusConflict)
			return
		}
	}

	for i, t := range tag {
		if t.Id == id {
			tag = append(tag[:i], tag[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Tag não encontrada", http.StatusNotFound)
}

// ==========================
// Handlers de TASK
// ==========================

func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
func GetTasksByTag(w http.ResponseWriter, r *http.Request) {
	tagNome := mux.Vars(r)["tag"]
	var filtradas []Task

	for _, t := range tasks.Task {
		if strings.EqualFold(t.Tag.Nome, tagNome) {
			filtradas = append(filtradas, t)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(filtradas); err != nil {
		http.Error(w, "Erro ao retornar as tarefas por tag", http.StatusInternalServerError)
	}
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	for _, t := range tasks.Task {
		if t.Id == id {
			json.NewEncoder(w).Encode(t)
			return
		}
	}
	http.Error(w, "Tarefa não encontrada", http.StatusNotFound)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var novaTask Task

	if err := json.NewDecoder(r.Body).Decode(&novaTask); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	if novaTask.Nome == "" || novaTask.Descricao == "" {
		http.Error(w, "Campos obrigatórios faltando", http.StatusBadRequest)
		return
	}

	if novaTask.Tag.Nome == "" {
		novaTask.Tag = Tag{Id: 0, Nome: ""}
	} else {
		novaTask.Tag = verificaOuCriaTagPorNome(novaTask.Tag.Nome)
	}

	var lastID int64
	for _, t := range tasks.Task {
		if t.Id > lastID {
			lastID = t.Id
		}
	}
	novaTask.Id = lastID + 1
	tasks.Task = append(tasks.Task, novaTask)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(novaTask)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var nova Task
	if err := json.NewDecoder(r.Body).Decode(&nova); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	for i, t := range tasks.Task {
		if t.Id == id {
			if nova.Nome != "" {
				tasks.Task[i].Nome = nova.Nome
			}
			if nova.Descricao != "" {
				tasks.Task[i].Descricao = nova.Descricao
			}
			if nova.Tag.Nome != "" {
				tasks.Task[i].Tag = verificaOuCriaTagPorNome(nova.Tag.Nome)
			}
			json.NewEncoder(w).Encode(tasks.Task[i])
			return
		}
	}
	http.Error(w, "Tarefa não encontrada", http.StatusNotFound)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	for i, t := range tasks.Task {
		if t.Id == id {
			tasks.Task = append(tasks.Task[:i], tasks.Task[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Tarefa não encontrada", http.StatusNotFound)
}

// ==========================
// Util
// ==========================

func verificaOuCriaTagPorNome(nome string) Tag {
	for _, t := range tag {
		if t.Nome == nome {
			return t
		}
	}
	var lastID int64
	for _, t := range tag {
		if t.Id > lastID {
			lastID = t.Id
		}
	}
	nova := Tag{Id: lastID + 1, Nome: nome}
	tag = append(tag, nova)
	return nova
}

// ==========================
// Main
// ==========================

func main() {
	router := mux.NewRouter()

	// Rotas de TAG
	router.HandleFunc("/tag", GetTags).Methods("GET")
	router.HandleFunc("/tag/{id}", GetTagByID).Methods("GET")
	router.HandleFunc("/tag", CreateTag).Methods("POST")
	router.HandleFunc("/tag/{id}", UpdateTag).Methods("PUT")
	router.HandleFunc("/tag/{id}", DeleteTag).Methods("DELETE")

	// Rotas de TASK
	router.HandleFunc("/task", GetTasks).Methods("GET")
	router.HandleFunc("/task/{id}", GetTaskByID).Methods("GET")
	router.HandleFunc("/task/tag/{tag}", GetTasksByTag).Methods("GET")
	router.HandleFunc("/task", CreateTask).Methods("POST")
	router.HandleFunc("/task/{id}", UpdateTask).Methods("PUT")
	router.HandleFunc("/task/{id}", DeleteTask).Methods("DELETE")

	// Serve static (index.html no /static)
	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/").Handler(fs)
	fmt.Println("Servidor rodando em http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
