package generator

import (
	"os"
	"path/filepath"
	"strings"
)

func CreateCRUD(base, name string) error {
	lower := strings.ToLower(name)

	// handler
	handler := `package handler

import "fmt"

func Create` + name + `() {
	fmt.Println("Create ` + lower + `")
}

func GetAll` + name + `() {
	fmt.Println("Get all ` + lower + `")
}

func Get` + name + `() {
	fmt.Println("Get ` + lower + ` by id")
}

func Update` + name + `() {
	fmt.Println("Update ` + lower + `")
}

func Delete` + name + `() {
	fmt.Println("Delete ` + lower + `")
}
`

	// service
	service := `package service

type ` + name + `Service interface {
	Create()
	GetAll()
	GetByID()
	Update()
	Delete()
}
`

	// repository
	repo := `package repository

type ` + name + `Repository interface {
	Create()
	FindAll()
	FindByID()
	Update()
	Delete()
}
`

	// write file
	os.WriteFile(filepath.Join(base, "handler", lower+"_handler.go"), []byte(handler), os.ModePerm)
	os.WriteFile(filepath.Join(base, "service", lower+"_service.go"), []byte(service), os.ModePerm)
	os.WriteFile(filepath.Join(base, "repository", lower+"_repository.go"), []byte(repo), os.ModePerm)

	return nil
}
