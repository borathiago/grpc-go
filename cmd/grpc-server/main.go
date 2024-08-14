package main

import (
	"database/sql"
	"net"

	"github.com/devfullcycle/14-gRPC/internal/database"
	"github.com/devfullcycle/14-gRPC/internal/pb"
	"github.com/devfullcycle/14-gRPC/internal/services"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, error := sql.Open("sqlite3", "./db.sqlite")
	if error != nil {
		panic(error)
	}
	defer db.Close()
	categoryDb := database.NewCategory(db)
	categoryService := services.NewCategoryService(*categoryDb)
	/* Criando servidor gRPC */ grpcServer := grpc.NewServer()
	/* Registrando serviços no servidor */ pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	/* Acrescentando reflection para ler e processar informação */ reflection.Register(grpcServer)
	listen, error := net.Listen("tcp", ":50051")
	if error != nil {
		panic(error)
	}
	if error := grpcServer.Serve(listen); error != nil {
		panic(error)
	}
}
