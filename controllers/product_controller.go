package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"funfillers/backend/config"
	"funfillers/backend/models"
	"github.com/gin-gonic/gin"
)

type ProductController struct{}

const productsFilePath = "./data/products.json"

func loadProductsFromFile() {
	if _, err := os.Stat(productsFilePath); os.IsNotExist(err) {
		saveProductsToFile()
		return
	}
	data, err := os.ReadFile(productsFilePath)
	if err != nil {
		return
	}
	var loaded []models.Product
	if err := json.Unmarshal(data, &loaded); err == nil && len(loaded) > 0 {
		mockProducts = loaded
	}
}

func saveProductsToFile() {
	os.MkdirAll("./data", 0755)
	data, err := json.MarshalIndent(mockProducts, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(productsFilePath, data, 0644)
}

func NewProductController() *ProductController {
	loadProductsFromFile()
	return &ProductController{}
}

var mockProducts = []models.Product{
	{
		ID:            "prod_1",
		Name:          "FunFillers Premium Wireless Gaming Headset",
		Description:   "Ultra-low latency 2.4GHz wireless audio with 50mm neodymium drivers and customizable RGB lighting.",
		Price:         129.99,
		RegularPrice:  159.99,
		ComparePrice:  159.99,
		CategoryID:    "cat_1",
		SubCategoryID: "sub_1",
		Category:      "Electronics",
		SubCategory:   "Audio & Headphones",
		ImageURL:      "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=600",
		Images: []string{
			"https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=600",
			"https://images.unsplash.com/photo-1546435770-a3e426bf472b?w=600",
			"https://images.unsplash.com/photo-1484704849700-f032a568e944?w=600",
		},
		Stock:        45,
		Rating:       4.8,
		ReviewsCount: 124,
		IsFeatured:   true,
		CreatedAt:    time.Now(),
	},
	{
		ID:            "prod_2",
		Name:          "Ergonomic RGB Mechanical Keyboard",
		Description:   "Hot-swappable mechanical switches, PBT keycaps, and customizable macro keys for maximum productivity.",
		Price:         89.95,
		RegularPrice:  119.00,
		ComparePrice:  119.00,
		CategoryID:    "cat_1",
		SubCategoryID: "sub_2",
		Category:      "Electronics",
		SubCategory:   "Keyboards & Mice",
		ImageURL:      "https://images.unsplash.com/photo-1587829741301-dc798b83add3?w=600",
		Images: []string{
			"https://images.unsplash.com/photo-1587829741301-dc798b83add3?w=600",
			"https://images.unsplash.com/photo-1618384887929-16ec33fab9ef?w=600",
		},
		Stock:        30,
		Rating:       4.9,
		ReviewsCount: 89,
		IsFeatured:   true,
		CreatedAt:    time.Now(),
	},
	{
		ID:            "prod_3",
		Name:          "Ultra-Comfort Smart Desk Chair",
		Description:   "Mesh breathable back support, dynamic lumbar adjustment, and 4D adjustable armrests.",
		Price:         249.00,
		RegularPrice:  299.00,
		ComparePrice:  299.00,
		CategoryID:    "cat_2",
		SubCategoryID: "sub_3",
		Category:      "Furniture",
		SubCategory:   "Office Chairs",
		ImageURL:      "https://images.unsplash.com/photo-1580481072645-022f9a6d83d0?w=600",
		Images: []string{
			"https://images.unsplash.com/photo-1580481072645-022f9a6d83d0?w=600",
		},
		Stock:        15,
		Rating:       4.7,
		ReviewsCount: 56,
		IsFeatured:   true,
		CreatedAt:    time.Now(),
	},
	{
		ID:            "prod_4",
		Name:          "Minimalist Aluminum Laptop Stand",
		Description:   "Precision-engineered aluminum alloy stand with heat ventilation slots and anti-slip silicone pads.",
		Price:         39.99,
		RegularPrice:  49.99,
		ComparePrice:  49.99,
		CategoryID:    "cat_3",
		SubCategoryID: "sub_5",
		Category:      "Accessories",
		SubCategory:   "Laptop Accessories",
		ImageURL:      "https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?w=600",
		Images: []string{
			"https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?w=600",
		},
		Stock:        60,
		Rating:       4.6,
		ReviewsCount: 42,
		IsFeatured:   false,
		CreatedAt:    time.Now(),
	},
}

func (pc *ProductController) GetProducts(c *gin.Context) {
	category := c.Query("category")
	search := strings.ToLower(c.Query("search"))
	featuredOnly := c.Query("featured") == "true"

	if config.DB != nil {
		var products []models.Product
		query := config.DB.Model(&models.Product{})

		if category != "" {
			query = query.Where("category = ? OR category_id = ?", category, category)
		}
		if search != "" {
			query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", "%"+search+"%", "%"+search+"%")
		}
		if featuredOnly {
			query = query.Where("is_featured = ?", true)
		}

		if err := query.Order("created_at desc").Find(&products).Error; err == nil {
			if len(products) == 0 && category == "" && search == "" && !featuredOnly {
				// Seed initial mock products if database is empty
				for _, mp := range mockProducts {
					config.DB.Create(&mp)
				}
				products = mockProducts
			}
			c.JSON(http.StatusOK, gin.H{
				"products": products,
				"total":    len(products),
			})
			return
		}
	}

	var filtered []models.Product
	for _, p := range mockProducts {
		if category != "" && !strings.EqualFold(p.Category, category) && p.CategoryID != category {
			continue
		}
		if search != "" && !strings.Contains(strings.ToLower(p.Name), search) && !strings.Contains(strings.ToLower(p.Description), search) {
			continue
		}
		if featuredOnly && !p.IsFeatured {
			continue
		}
		filtered = append(filtered, p)
	}

	if filtered == nil {
		filtered = []models.Product{}
	}

	c.JSON(http.StatusOK, gin.H{
		"products": filtered,
		"total":    len(filtered),
	})
}

func (pc *ProductController) GetProductByID(c *gin.Context) {
	id := c.Param("id")

	if config.DB != nil {
		var prod models.Product
		if err := config.DB.First(&prod, "id = ?", id).Error; err == nil {
			c.JSON(http.StatusOK, prod)
			return
		}
	}

	for _, p := range mockProducts {
		if p.ID == id {
			c.JSON(http.StatusOK, p)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}

func (pc *ProductController) CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	name := strings.TrimSpace(req.Name)
	if len(name) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product title must be at least 3 characters long"})
		return
	}

	if req.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product price must be greater than zero"})
		return
	}

	if req.Stock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product stock quantity cannot be negative"})
		return
	}

	imgs := req.Images
	if len(imgs) == 0 && req.ImageURL != "" {
		imgs = []string{req.ImageURL}
	}

	img := req.ImageURL
	if img == "" && len(imgs) > 0 {
		img = imgs[0]
	}
	if img == "" {
		img = "https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=600"
		imgs = []string{img}
	}

	regPrice := req.RegularPrice
	if regPrice <= 0 {
		regPrice = req.Price * 1.2
	}

	warningCopy := computeLegalWarningCopy(req.HasSmallParts, req.HasSmallBall, req.HasMarble, req.HasBalloon)
	suffocationTag := req.GrossWidthCm > 12.7

	ageGrading := req.AgeGrading
	if ageGrading == "" {
		ageGrading = "AGE_2_4Y"
	}

	status := "PENDING_COMPLIANCE"

	prod := models.Product{
		ID:                     fmt.Sprintf("prod_%d", time.Now().UnixNano()/1e6),
		SKU:                    req.SKU,
		Name:                   req.Name,
		Description:            req.Description,
		Price:                  req.Price,
		RegularPrice:           regPrice,
		ComparePrice:           regPrice,
		CategoryID:             req.CategoryID,
		SubCategoryID:          req.SubCategoryID,
		Category:               req.Category,
		SubCategory:            req.SubCategory,
		ImageURL:               img,
		Images:                 imgs,
		Stock:                  req.Stock,
		Rating:                 5.0,
		ReviewsCount:           1,
		IsFeatured:             req.IsFeatured,
		CreatedAt:              time.Now(),
		ListingStatus:          status,
		AgeGrading:             ageGrading,
		HasSmallParts:          req.HasSmallParts,
		HasSmallBall:           req.HasSmallBall,
		HasMarble:              req.HasMarble,
		HasBalloon:             req.HasBalloon,
		LegalWarningCopy:       warningCopy,
		CountryOfOrigin:        req.CountryOfOrigin,
		ManufacturerName:       req.ManufacturerName,
		ManufacturerAddress:    req.ManufacturerAddress,
		ImporterName:           req.ImporterName,
		ImporterAddress:        req.ImporterAddress,
		CustomerCareContact:    req.CustomerCareContact,
		NetQuantity:            req.NetQuantity,
		Materials:              req.Materials,
		IsNonToxic:             req.IsNonToxic,
		IsBpaFree:              req.IsBpaFree,
		IsWashable:             req.IsWashable,
		BatteryRequirement:     req.BatteryRequirement,
		BatteryChemistry:       req.BatteryChemistry,
		Un383ReportUrl:         req.Un383ReportUrl,
		HasScrewLockBattery:    req.HasScrewLockBattery,
		NetWeightGrams:         req.NetWeightGrams,
		GrossWeightGrams:       req.GrossWeightGrams,
		GrossWidthCm:           req.GrossWidthCm,
		RequiresSuffocationTag: suffocationTag,
		BatchNumber:            req.BatchNumber,
	}

	if config.DB != nil {
		if err := config.DB.Create(&prod).Error; err != nil {
			fmt.Printf("⚠️ DB insert product error: %v\n", err)
		}
	}
	mockProducts = append([]models.Product{prod}, mockProducts...)
	saveProductsToFile()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"product": prod,
	})
}

func (pc *ProductController) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" && len(strings.TrimSpace(req.Name)) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Updated product title must be at least 3 characters long"})
		return
	}

	if req.Price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Updated product price cannot be negative"})
		return
	}

	if req.Stock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Updated product stock quantity cannot be negative"})
		return
	}

	warningCopy := computeLegalWarningCopy(req.HasSmallParts, req.HasSmallBall, req.HasMarble, req.HasBalloon)
	suffocationTag := req.GrossWidthCm > 12.7

	updateFields := func(p *models.Product) {
		if req.Name != "" {
			p.Name = req.Name
		}
		if req.SKU != "" {
			p.SKU = req.SKU
		}
		if req.Description != "" {
			p.Description = req.Description
		}
		if req.Price > 0 {
			p.Price = req.Price
		}
		if req.RegularPrice > 0 {
			p.RegularPrice = req.RegularPrice
		}
		if req.Category != "" {
			p.Category = req.Category
		}
		p.SubCategory = req.SubCategory
		if req.CategoryID != "" {
			p.CategoryID = req.CategoryID
		}
		p.SubCategoryID = req.SubCategoryID
		if req.ImageURL != "" {
			p.ImageURL = req.ImageURL
		}
		if req.Images != nil {
			p.Images = req.Images
			if len(req.Images) > 0 && req.ImageURL == "" {
				p.ImageURL = req.Images[0]
			}
		}
		if req.Stock >= 0 {
			p.Stock = req.Stock
		}
		p.IsFeatured = req.IsFeatured

		if req.AgeGrading != "" {
			p.AgeGrading = req.AgeGrading
		}
		p.HasSmallParts = req.HasSmallParts
		p.HasSmallBall = req.HasSmallBall
		p.HasMarble = req.HasMarble
		p.HasBalloon = req.HasBalloon
		p.LegalWarningCopy = warningCopy

		if req.CountryOfOrigin != "" {
			p.CountryOfOrigin = req.CountryOfOrigin
		}
		if req.ManufacturerName != "" {
			p.ManufacturerName = req.ManufacturerName
		}
		if req.ManufacturerAddress != "" {
			p.ManufacturerAddress = req.ManufacturerAddress
		}
		if req.ImporterName != "" {
			p.ImporterName = req.ImporterName
		}
		if req.ImporterAddress != "" {
			p.ImporterAddress = req.ImporterAddress
		}
		if req.CustomerCareContact != "" {
			p.CustomerCareContact = req.CustomerCareContact
		}
		if req.NetQuantity != "" {
			p.NetQuantity = req.NetQuantity
		}
		if req.Materials != nil {
			p.Materials = req.Materials
		}
		p.IsNonToxic = req.IsNonToxic
		p.IsBpaFree = req.IsBpaFree
		p.IsWashable = req.IsWashable
		if req.BatteryRequirement != "" {
			p.BatteryRequirement = req.BatteryRequirement
		}
		if req.BatteryChemistry != "" {
			p.BatteryChemistry = req.BatteryChemistry
		}
		if req.Un383ReportUrl != "" {
			p.Un383ReportUrl = req.Un383ReportUrl
		}
		p.HasScrewLockBattery = req.HasScrewLockBattery
		p.NetWeightGrams = req.NetWeightGrams
		p.GrossWeightGrams = req.GrossWeightGrams
		p.GrossWidthCm = req.GrossWidthCm
		p.RequiresSuffocationTag = suffocationTag
		if req.BatchNumber != "" {
			p.BatchNumber = req.BatchNumber
		}
		if req.ListingStatus != "" {
			p.ListingStatus = req.ListingStatus
		}
	}

	if config.DB != nil {
		var dbProd models.Product
		if err := config.DB.First(&dbProd, "id = ?", id).Error; err == nil {
			updateFields(&dbProd)
			config.DB.Save(&dbProd)

			for i, p := range mockProducts {
				if p.ID == id {
					mockProducts[i] = dbProd
					break
				}
			}
			saveProductsToFile()

			c.JSON(http.StatusOK, gin.H{
				"message": "Product updated successfully",
				"product": dbProd,
			})
			return
		}
	}

	for i, p := range mockProducts {
		if p.ID == id {
			updateFields(&mockProducts[i])
			saveProductsToFile()

			c.JSON(http.StatusOK, gin.H{
				"message": "Product updated successfully",
				"product": mockProducts[i],
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}

func (pc *ProductController) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if config.DB != nil {
		config.DB.Where("id = ?", id).Delete(&models.Product{})
	}

	updated := make([]models.Product, 0)
	for _, p := range mockProducts {
		if p.ID != id {
			updated = append(updated, p)
		}
	}
	mockProducts = updated
	saveProductsToFile()

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (pc *ProductController) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		file, err = c.FormFile("image")
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided in request"})
		return
	}

	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	ext := filepath.Ext(file.Filename)
	baseName := strings.TrimSuffix(filepath.Base(file.Filename), ext)
	sanitizedName := strings.ReplaceAll(baseName, " ", "_")
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), sanitizedName, ext)
	dst := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
		return
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	host := c.Request.Host
	imageUrl := fmt.Sprintf("%s://%s/uploads/%s", scheme, host, filename)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Image uploaded successfully",
		"url":      imageUrl,
		"filename": filename,
	})
}

func computeLegalWarningCopy(smallParts, smallBall, marble, balloon bool) string {
	var parts []string
	if smallParts {
		parts = append(parts, "CHOKING HAZARD - Small parts.")
	}
	if smallBall {
		parts = append(parts, "CHOKING HAZARD - Toy contains a small ball.")
	}
	if marble {
		parts = append(parts, "CHOKING HAZARD - Toy contains a marble.")
	}
	if balloon {
		parts = append(parts, "CHOKING HAZARD - Children under 8 yrs can choke or suffocate on uninflated or broken balloons.")
	}
	if len(parts) == 0 {
		return "No choking hazard warnings required for this age grade."
	}
	return "WARNING: " + strings.Join(parts, " ") + " Not suitable for children under 3 years."
}

func (pc *ProductController) RecallBatch(c *gin.Context) {
	var body struct {
		BatchNumber     string `json:"batchNumber" binding:"required"`
		RecallReason    string `json:"recallReason" binding:"required"`
		NotifyCustomers bool   `json:"notifyCustomers"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recalledCount := 0
	for i, p := range mockProducts {
		if strings.EqualFold(p.BatchNumber, body.BatchNumber) {
			mockProducts[i].ListingStatus = "RECALLED"
			recalledCount++
		}
	}
	saveProductsToFile()

	c.JSON(http.StatusOK, gin.H{
		"success":             true,
		"recalledBatchNumber": body.BatchNumber,
		"recallReason":        body.RecallReason,
		"disabledSkusCount":   recalledCount,
		"message":             fmt.Sprintf("Batch %s recalled successfully. %d products disabled.", body.BatchNumber, recalledCount),
	})
}

func (pc *ProductController) GetExpiringCertifications(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"totalExpiring": 1,
		"expiringCertifications": []gin.H{
			{
				"certId":            "cert_102",
				"productId":         "prod_1",
				"productName":       "FunFillers Premium Wireless Gaming Headset",
				"certType":          "BIS_IS_9873",
				"certificateNumber": "BIS-CM/L-9812739",
				"expiryDate":        time.Now().AddDate(0, 0, 45).Format(time.RFC3339),
				"daysUntilExpiry":   45,
				"status":            "APPROVED",
			},
		},
	})
}

func (pc *ProductController) GetProductSharePreview(c *gin.Context) {
	id := c.Param("id")
	var prod models.Product
	found := false

	if config.DB != nil {
		if err := config.DB.First(&prod, "id = ?", id).Error; err == nil {
			found = true
		}
	}

	if !found {
		for _, p := range mockProducts {
			if p.ID == id {
				prod = p
				found = true
				break
			}
		}
	}

	if !found {
		c.String(http.StatusNotFound, "Product not found")
		return
	}

	host := c.Request.Host
	if host == "" {
		host = "ff-backend-klec.onrender.com"
	}
	scheme := "https"
	if strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.0.0.1") {
		scheme = "http"
	}

	fullURL := fmt.Sprintf("%s://%s/share/product/%s", scheme, host, prod.ID)
	priceStr := fmt.Sprintf("$%.2f", prod.Price)

	publicImageURL := prod.ImageURL
	if strings.Contains(publicImageURL, "localhost") || strings.Contains(publicImageURL, "127.0.0.1") {
		publicImageURL = regexp.MustCompile(`http://(localhost|127\.0\.0\.1):(5050|8080|5000)`).ReplaceAllString(publicImageURL, fmt.Sprintf("%s://%s", scheme, host))
	} else if strings.HasPrefix(publicImageURL, "/uploads") {
		publicImageURL = fmt.Sprintf("%s://%s%s", scheme, host, publicImageURL)
	}

	if publicImageURL == "" {
		publicImageURL = "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=600"
	}

	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en" prefix="og: http://ogp.me/ns#">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <!-- Primary Meta Tags -->
    <title>%s | FunFillers Toys & Games</title>
    <meta name="title" content="%s | FunFillers Toys">
    <meta name="description" content="Check out %s for %s on FunFillers! %s">

    <!-- Open Graph / Facebook / WhatsApp -->
    <meta property="og:type" content="website">
    <meta property="og:url" content="%s">
    <meta property="og:title" content="%s">
    <meta property="og:description" content="Check out %s for %s on FunFillers! %s">
    <meta property="og:image" content="%s">
    <meta property="og:image:secure_url" content="%s">
    <meta property="og:image:width" content="600">
    <meta property="og:image:height" content="600">
    <meta property="og:image:type" content="image/jpeg">
    <meta property="og:site_name" content="FunFillers Toys">
    <meta property="product:price:amount" content="%.2f">
    <meta property="product:price:currency" content="USD">

    <!-- Twitter Card -->
    <meta name="twitter:card" content="summary_large_image">
    <meta name="twitter:url" content="%s">
    <meta name="twitter:title" content="%s">
    <meta name="twitter:description" content="Check out %s for %s on FunFillers!">
    <meta name="twitter:image" content="%s">

    <style>
        * { box-sizing: border-box; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }
        body { background-color: #F8F9FD; margin: 0; padding: 24px; display: flex; justify-content: center; align-items: center; min-height: 100vh; }
        .card { background: white; border-radius: 24px; box-shadow: 0 10px 40px rgba(0,0,0,0.08); max-width: 420px; width: 100%%; overflow: hidden; text-align: center; padding: 32px 24px; border: 1px solid #EEF2F6; }
        .img-box { width: 100%%; height: 260px; border-radius: 16px; overflow: hidden; background: #F4F6F9; margin-bottom: 20px; display: flex; align-items: center; justify-content: center; }
        .img-box img { max-width: 100%%; max-height: 100%%; object-fit: contain; }
        h1 { font-size: 22px; color: #1E293B; margin: 0 0 8px 0; font-weight: 800; line-height: 1.3; }
        .category { font-size: 13px; color: #64748B; margin-bottom: 12px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.5px; }
        .price { font-size: 26px; color: #FF4785; font-weight: 900; margin-bottom: 16px; }
        .desc { font-size: 14px; color: #475569; line-height: 1.5; margin-bottom: 24px; text-align: left; }
        .btn { display: inline-block; width: 100%%; background: linear-gradient(135deg, #FF4785 0%%, #FF70A6 100%%); color: white; padding: 16px; border-radius: 14px; text-decoration: none; font-weight: 800; font-size: 16px; box-shadow: 0 6px 20px rgba(255, 71, 133, 0.3); transition: transform 0.2s; }
        .btn:hover { transform: translateY(-2px); }
    </style>
</head>
<body>
    <div class="card">
        <div class="img-box">
            <img src="%s" alt="%s" />
        </div>
        <div class="category">%s</div>
        <h1>%s</h1>
        <div class="price">%s</div>
        <div class="desc">%s</div>
        <a href="https://funfillers.netlify.app/product/%s" class="btn">View Product Page</a>
    </div>

    <script>
        // Redirect to Netlify user panel product page
        window.location.href = "https://funfillers.netlify.app/product/%s";
    </script>
</body>
</html>`,
		prod.Name, prod.Name, prod.Name, priceStr, prod.Description,
		fullURL, prod.Name, prod.Name, priceStr, prod.Description,
		publicImageURL, publicImageURL, prod.Price,
		fullURL, prod.Name, prod.Name, priceStr, publicImageURL,
		publicImageURL, prod.Name, prod.Category, prod.Name, priceStr, prod.Description,
		prod.ID, prod.ID,
	)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, htmlContent)
}


