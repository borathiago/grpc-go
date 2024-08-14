package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type Course struct {
	db          *sql.DB
	ID          string
	Name        string
	Description string
	CategoryID  string
}

func NewCourse(db *sql.DB) *Course {
	return &Course{db: db}
}

func (c *Course) Create(name, description, categoryID string) (*Course, error) {
	id := uuid.New().String()
	_, error := c.db.Exec("INSERT INTO courses (id, name, description, category_id) VALUES ($1, $2, $3, $4)", id, name, description, categoryID)
	if error != nil {
		return nil, error
	}
	return &Course{
		ID:          id,
		Name:        name,
		Description: description,
		CategoryID:  categoryID,
	}, nil
}

func (c *Course) FindAll() ([]Course, error) {
	rows, error := c.db.Query("SELECT id, name, description, category_id FROM courses")
	if error != nil {
		return nil, error
	}
	defer rows.Close()
	courses := []Course{}
	for rows.Next() {
		var id, name, description, categoryID string
		if error := rows.Scan(&id, &name, &description, &categoryID); error != nil {
			return nil, error
		}
		courses = append(courses, Course{ID: id, Name: name, Description: description, CategoryID: categoryID})
	}
	return courses, nil
}

func (c *Course) FindByCategoryID(categoryID string) ([]Course, error) {
	rows, error := c.db.Query("SELECT id, name, description, category_id FROM courses WHERE category_id = $1", categoryID)
	if error != nil {
		return nil, error
	}
	defer rows.Close()
	courses := []Course{}
	for rows.Next() {
		var id, name, description, categoryID string
		if error := rows.Scan(&id, &name, &description, &categoryID); error != nil {
			return nil, error
		}
		courses = append(courses, Course{ID: id, Name: name, Description: description, CategoryID: categoryID})
	}
	return courses, nil
}
