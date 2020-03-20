package postgres

import (
    "github.com/go-pg/pg"
    "github.com/secmohammed/meetups/models"
)

// CategoriesRepo is used to contain the db driver.
type CategoriesRepo struct {
    DB *pg.DB
}

// GetCategories is used to get categories from database.
func (c *CategoriesRepo) GetCategories(limit, offset *int) ([]*models.Category, error) {
    var categories []*models.Category
    query := c.DB.Model(&categories).Order("id")
    if limit != nil {
        query.Limit(*limit)
    }
    if offset != nil {
        query.Offset(*offset)
    }
    err := query.Select()
    if err != nil {
        return nil, err
    }
    return categories, nil
}

//Create is used to create a comment using the passed struct.
func (c *CategoriesRepo) Create(category *models.Category) (*models.Category, error) {
    _, err := c.DB.Model(category).Returning("*").Insert()
    return category, err
}

// Update is used to update the passed meetup by id.
func (c *CategoriesRepo) Update(category *models.Category) (*models.Category, error) {
    _, err := c.DB.Model(category).Where("id = ?", category.ID).Update()
    return category, err
}

// GetByName is used to fetch meetup by name.
func (c *CategoriesRepo) GetByName(name string) (*models.Category, error) {
    category := models.Category{}
    err := c.DB.Model(&category).Where("name = ?", name).First()
    if err != nil {
        return nil, err
    }
    return &category, nil
}

// Delete is used to delete meetup by its id.
func (c *CategoriesRepo) Delete(category *models.Category) error {
    _, err := c.DB.Model(category).Where("id = ?", category.ID).Delete()
    return err
}
