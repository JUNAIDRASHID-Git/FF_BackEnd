package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"funfillers/backend/config"
	"funfillers/backend/models"
	"github.com/gin-gonic/gin"
)

type CategoryController struct{}

func NewCategoryController() *CategoryController {
	return &CategoryController{}
}

var mockCategories = []models.Category{
	{
		ID:          "cat_action",
		Name:        "Action Figures",
		Slug:        "action-figures",
		Icon:        "https://images.unsplash.com/photo-1608889825205-eebdb9fc5806?auto=format&fit=crop&w=300&q=80",
		Image:       "https://images.unsplash.com/photo-1608889825205-eebdb9fc5806?auto=format&fit=crop&w=300&q=80",
		Description: "Superhero action figures and collectible toys",
		ItemCount:   124,
		SubCategories: []models.SubCategory{
			{ID: "sub_1", CategoryID: "cat_action", Name: "Superheroes", Slug: "superheroes", Icon: "https://images.unsplash.com/photo-1608889825205-eebdb9fc5806?auto=format&fit=crop&w=300&q=80", Description: "Avengers, Spider-Man & Batman figures"},
			{ID: "sub_2", CategoryID: "cat_action", Name: "Anime Collectibles", Slug: "anime", Icon: "https://images.unsplash.com/photo-1563089145-599997674d42?auto=format&fit=crop&w=300&q=80", Description: "Popular anime and manga statues"},
		},
		CreatedAt: time.Now(),
	},
	{
		ID:          "cat_dolls",
		Name:        "Dolls & Plush",
		Slug:        "dolls-plush",
		Icon:        "https://images.unsplash.com/photo-1566576721346-d4a3b4eaeb55?auto=format&fit=crop&w=300&q=80",
		Image:       "https://images.unsplash.com/photo-1566576721346-d4a3b4eaeb55?auto=format&fit=crop&w=300&q=80",
		Description: "Soft cuddly teddy bears and fashion princess dolls",
		ItemCount:   98,
		SubCategories: []models.SubCategory{
			{ID: "sub_3", CategoryID: "cat_dolls", Name: "Teddy Bears", Slug: "teddy-bears", Icon: "https://images.unsplash.com/photo-1558060370-d644479be6f7?auto=format&fit=crop&w=300&q=80", Description: "Ultra-soft cuddly plush teddy bears"},
			{ID: "sub_4", CategoryID: "cat_dolls", Name: "Fashion Dolls", Slug: "fashion-dolls", Icon: "https://images.unsplash.com/photo-1566576721346-d4a3b4eaeb55?auto=format&fit=crop&w=300&q=80", Description: "Royal princess dolls with accessories"},
		},
		CreatedAt: time.Now(),
	},
	{
		ID:          "cat_building",
		Name:        "Building Blocks",
		Slug:        "building-blocks",
		Icon:        "https://images.unsplash.com/photo-1587654780291-39c9404d746b?auto=format&fit=crop&w=300&q=80",
		Image:       "https://images.unsplash.com/photo-1587654780291-39c9404d746b?auto=format&fit=crop&w=300&q=80",
		Description: "Creative wooden blocks and brick construction sets",
		ItemCount:   76,
		SubCategories: []models.SubCategory{
			{ID: "sub_5", CategoryID: "cat_building", Name: "Wooden Blocks", Slug: "wooden-blocks", Icon: "https://images.unsplash.com/photo-1587654780291-39c9404d746b?auto=format&fit=crop&w=300&q=80", Description: "Rainbow educational wooden blocks"},
			{ID: "sub_6", CategoryID: "cat_building", Name: "Construction Sets", Slug: "construction", Icon: "https://images.unsplash.com/photo-1513151233558-d860c5398176?auto=format&fit=crop&w=300&q=80", Description: "Lego-style building brick sets"},
		},
		CreatedAt: time.Now(),
	},
	{
		ID:          "cat_vehicles",
		Name:        "Vehicles & RC",
		Slug:        "vehicles-rc",
		Icon:        "https://images.unsplash.com/photo-1594787318286-3d835c1d207f?auto=format&fit=crop&w=300&q=80",
		Image:       "https://images.unsplash.com/photo-1594787318286-3d835c1d207f?auto=format&fit=crop&w=300&q=80",
		Description: "High-speed remote control cars, stunt trucks & diecast models",
		ItemCount:   112,
		SubCategories: []models.SubCategory{
			{ID: "sub_7", CategoryID: "cat_vehicles", Name: "RC Stunt Cars", Slug: "rc-cars", Icon: "https://images.unsplash.com/photo-1594787318286-3d835c1d207f?auto=format&fit=crop&w=300&q=80", Description: "360-flip remote control vehicles"},
			{ID: "sub_8", CategoryID: "cat_vehicles", Name: "Diecast Models", Slug: "diecast", Icon: "https://images.unsplash.com/photo-1581235720704-06d3acfcb36f?auto=format&fit=crop&w=300&q=80", Description: "1:24 scale GT-R diecast cars"},
		},
		CreatedAt: time.Now(),
	},
	{
		ID:          "cat_educational",
		Name:        "Educational Toys",
		Slug:        "educational",
		Icon:        "https://images.unsplash.com/photo-1500995617113-cf789362a3e1?auto=format&fit=crop&w=300&q=80",
		Image:       "https://images.unsplash.com/photo-1500995617113-cf789362a3e1?auto=format&fit=crop&w=300&q=80",
		Description: "STEM learning kits, puzzles and musical instruments",
		ItemCount:   64,
		SubCategories: []models.SubCategory{
			{ID: "sub_9", CategoryID: "cat_educational", Name: "STEM & Puzzles", Slug: "stem", Icon: "https://images.unsplash.com/photo-1500995617113-cf789362a3e1?auto=format&fit=crop&w=300&q=80", Description: "Interactive logic puzzles and science kits"},
			{ID: "sub_10", CategoryID: "cat_educational", Name: "Musical Toys", Slug: "musical", Icon: "https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?auto=format&fit=crop&w=300&q=80", Description: "Interactive toy guitars and keypads"},
		},
		CreatedAt: time.Now(),
	},
}

func (cc *CategoryController) GetCategories(c *gin.Context) {
	if config.DB != nil {
		var categories []models.Category
		if err := config.DB.Preload("SubCategories").Find(&categories).Error; err == nil && len(categories) > 0 {
			c.JSON(http.StatusOK, categories)
			return
		}
	}
	c.JSON(http.StatusOK, mockCategories)
}

func (cc *CategoryController) CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := strings.TrimSpace(req.Name)
	if len(name) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category name must be at least 2 characters long"})
		return
	}

	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	}

	iconImg := strings.TrimSpace(req.Icon)
	if iconImg == "" {
		iconImg = strings.TrimSpace(req.Image)
	}

	// Update existing category if ID is provided
	if req.ID != "" {
		for i, existing := range mockCategories {
			if existing.ID == req.ID {
				mockCategories[i].Name = name
				mockCategories[i].Slug = slug
				mockCategories[i].Icon = iconImg
				mockCategories[i].Image = iconImg
				mockCategories[i].Description = strings.TrimSpace(req.Description)

				if config.DB != nil {
					config.DB.Model(&models.Category{}).Where("id = ?", req.ID).Updates(map[string]interface{}{
						"name":        name,
						"slug":        slug,
						"icon":        iconImg,
						"image":       iconImg,
						"description": strings.TrimSpace(req.Description),
					})
				}

				c.JSON(http.StatusOK, gin.H{
					"message":  "Category updated successfully",
					"category": mockCategories[i],
				})
				return
			}
		}
	}

	cat := models.Category{
		ID:          fmt.Sprintf("cat_%d", time.Now().UnixNano()/1e6),
		Name:        name,
		Slug:        slug,
		Icon:        iconImg,
		Image:       iconImg,
		Description: strings.TrimSpace(req.Description),
		CreatedAt:   time.Now(),
	}

	if config.DB != nil {
		config.DB.Create(&cat)
	}
	mockCategories = append(mockCategories, cat)

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Category created successfully",
		"category": cat,
	})
}

func (cc *CategoryController) CreateSubCategory(c *gin.Context) {
	var req models.CreateSubCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := strings.TrimSpace(req.Name)
	if len(name) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Subcategory name must be at least 2 characters long"})
		return
	}

	if strings.TrimSpace(req.CategoryID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parent category ID is required"})
		return
	}

	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	}

	iconImg := strings.TrimSpace(req.Icon)
	if iconImg == "" {
		iconImg = strings.TrimSpace(req.Image)
	}

	// Update existing subcategory if ID is provided
	if req.ID != "" {
		for i, cat := range mockCategories {
			if cat.ID == req.CategoryID {
				for j, sub := range cat.SubCategories {
					if sub.ID == req.ID {
						mockCategories[i].SubCategories[j].Name = name
						mockCategories[i].SubCategories[j].Slug = slug
						mockCategories[i].SubCategories[j].Icon = iconImg
						mockCategories[i].SubCategories[j].Image = iconImg
						mockCategories[i].SubCategories[j].Description = strings.TrimSpace(req.Description)

						if config.DB != nil {
							config.DB.Model(&models.SubCategory{}).Where("id = ?", req.ID).Updates(map[string]interface{}{
								"name":        name,
								"slug":        slug,
								"icon":        iconImg,
								"image":       iconImg,
								"description": strings.TrimSpace(req.Description),
							})
						}

						c.JSON(http.StatusOK, gin.H{
							"message":     "SubCategory updated successfully",
							"subCategory": mockCategories[i].SubCategories[j],
						})
						return
					}
				}
			}
		}
	}

	sub := models.SubCategory{
		ID:          fmt.Sprintf("sub_%d", time.Now().UnixNano()/1e6),
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Slug:        slug,
		Icon:        iconImg,
		Image:       iconImg,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	if config.DB != nil {
		config.DB.Create(&sub)
	}

	for i, cat := range mockCategories {
		if cat.ID == req.CategoryID {
			mockCategories[i].SubCategories = append(mockCategories[i].SubCategories, sub)
			break
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "SubCategory created successfully",
		"subCategory": sub,
	})
}

func (cc *CategoryController) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if config.DB != nil {
		config.DB.Where("id = ?", id).Delete(&models.Category{})
	}

	updated := make([]models.Category, 0)
	for _, cat := range mockCategories {
		if cat.ID != id {
			updated = append(updated, cat)
		}
	}
	mockCategories = updated

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
