package services

import (
	"context"
	"io"

	"github.com/devfullcycle/14-gRPC/internal/database"
	"github.com/devfullcycle/14-gRPC/internal/pb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{
		CategoryDB: categoryDB,
	}
}

func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
	category, error := c.CategoryDB.Create(in.Name, in.Description)
	if error != nil {
		return nil, error
	}
	categoryResponse := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}
	return &pb.CategoryResponse{
		Category: categoryResponse,
	}, nil
}

func (c *CategoryService) ListCategories(ctx context.Context, in *pb.Blank) (*pb.CategoryList, error) {
	categories, error := c.CategoryDB.FindAll()
	if error != nil {
		return nil, error
	}
	var categoriesResponse []*pb.Category
	for _, category := range categories {
		categoryResponse := &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		}
		categoriesResponse = append(categoriesResponse, categoryResponse)
	}
	return &pb.CategoryList{Categories: categoriesResponse}, nil
}

func (c *CategoryService) GetCategory(ctx context.Context, in *pb.CategoryGetRequest) (*pb.Category, error) {
	category, error := c.CategoryDB.Find(in.Id)
	if error != nil {
		return nil, error
	}
	categoryResponse := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}
	return categoryResponse, nil
}

func (c *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.CategoryList{}
	for {
		category, error := stream.Recv()
		if error == io.EOF {
			return stream.SendAndClose(categories)
		}
		if error != nil {
			return error
		}
		categoryResult, error := c.CategoryDB.Create(category.Name, category.Description)
		if error != nil {
			return error
		}
		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		})
	}
}

func (c *CategoryService) CreateCategoryStreamBidirectional(stream pb.CategoryService_CreateCategoryStreamBidirectionalServer) error {
	for {
		categoty, error := stream.Recv()
		if error == io.EOF {
			return nil
		}
		if error != nil {
			return error
		}
		categoryResult, error := c.CategoryDB.Create(categoty.Name, categoty.Description)
		if error != nil {
			return error
		}
		error = stream.Send(&pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		})
		if error != nil {
			return error
		}
	}
}
