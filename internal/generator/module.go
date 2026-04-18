package generator

import (
	"os"
	"path/filepath"
	"strings"
)

func GenerateFullModule(base, moduleName, modulePath string) error {
	name := capitalize(moduleName)
	lower := strings.ToLower(moduleName)

	// ===== MODEL =====
	model := `package model

type ` + name + ` struct {
	ID   int    ` + "`json:\"id\"`" + `
	Name string ` + "`json:\"name\"`" + `
}
`

	// ===== REPOSITORY =====
	repo := `package repository

import "` + modulePath + `/modules/` + lower + `/model"

type ` + name + `Repository interface {
	FindAll() []model.` + name + `
	FindByID(id int) *model.` + name + `
	Create(data model.` + name + `)
}

type ` + lower + `Repository struct {
	data []model.` + name + `
}

func New` + name + `Repository() ` + name + `Repository {
	return &` + lower + `Repository{}
}

func (r *` + lower + `Repository) FindAll() []model.` + name + ` {
	return r.data
}

func (r *` + lower + `Repository) FindByID(id int) *model.` + name + ` {
	for _, v := range r.data {
		if v.ID == id {
			return &v
		}
	}
	return nil
}

func (r *` + lower + `Repository) Create(data model.` + name + `) {
	r.data = append(r.data, data)
}
`

	// ===== SERVICE =====
	service := `package service

import (
	"` + modulePath + `/modules/` + lower + `/model"
	"` + modulePath + `/modules/` + lower + `/repository"
)

type ` + name + `Service interface {
	GetAll() []model.` + name + `
	GetByID(id int) *model.` + name + `
	Create(data model.` + name + `)
}

type ` + lower + `Service struct {
	repo repository.` + name + `Repository
}

func New` + name + `Service(r repository.` + name + `Repository) ` + name + `Service {
	return &` + lower + `Service{repo: r}
}

func (s *` + lower + `Service) GetAll() []model.` + name + ` {
	return s.repo.FindAll()
}

func (s *` + lower + `Service) GetByID(id int) *model.` + name + ` {
	return s.repo.FindByID(id)
}

func (s *` + lower + `Service) Create(data model.` + name + `) {
	s.repo.Create(data)
}
`

	// ===== HANDLER =====
	handler := `package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"` + modulePath + `/modules/` + lower + `/model"
	"` + modulePath + `/modules/` + lower + `/service"
)

type ` + name + `Handler struct {
	service service.` + name + `Service
}

func New` + name + `Handler(s service.` + name + `Service) *` + name + `Handler {
	return &` + name + `Handler{service: s}
}

func (h *` + name + `Handler) GetAll(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.GetAll())
}

func (h *` + name + `Handler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data := h.service.GetByID(id)

	if data == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *` + name + `Handler) Create(c *gin.Context) {
	var req model.` + name + `

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.service.Create(req)

	c.JSON(http.StatusOK, req)
}
`

	// ===== ROUTE =====
	route := `package route

import (
	"github.com/gin-gonic/gin"

	"` + modulePath + `/modules/` + lower + `/handler"
	"` + modulePath + `/modules/` + lower + `/repository"
	"` + modulePath + `/modules/` + lower + `/service"
)

func Register` + name + `Routes(r *gin.Engine) {
	repo := repository.New` + name + `Repository()
	svc := service.New` + name + `Service(repo)
	h := handler.New` + name + `Handler(svc)

	group := r.Group("/` + lower + `")
	{
		group.GET("/", h.GetAll)
		group.GET("/:id", h.GetByID)
		group.POST("/", h.Create)
	}
}
`

	// write files
	if err := write(filepath.Join(base, "model", lower+".go"), model); err != nil {
		return err
	}
	if err := write(filepath.Join(base, "repository", lower+"_repository.go"), repo); err != nil {
		return err
	}
	if err := write(filepath.Join(base, "service", lower+"_service.go"), service); err != nil {
		return err
	}
	if err := write(filepath.Join(base, "handler", lower+"_handler.go"), handler); err != nil {
		return err
	}
	if err := write(filepath.Join(base, "route", lower+"_route.go"), route); err != nil {
		return err
	}

	return nil
}

func write(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(content), os.ModePerm); err != nil {
		return err
	}
	return nil
}

func capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}
