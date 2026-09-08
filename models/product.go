package models

import "time"

type Product struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	SKU           string    `json:"sku"`
	Name          string    `json:"name" gorm:"not null"`
	Description   string    `json:"description"`
	Price         float64   `json:"price" gorm:"not null"`
	RegularPrice  float64   `json:"regularPrice"`
	ComparePrice  float64   `json:"comparePrice,omitempty"`
	CategoryID    string    `json:"categoryId"`
	SubCategoryID string    `json:"subCategoryId"`
	Category      string    `json:"category"`
	SubCategory   string    `json:"subCategory"`
	ImageURL      string    `json:"imageUrl"`
	Images        []string  `json:"images" gorm:"serializer:json"`
	Stock         int       `json:"stock"`
	Rating        float64   `json:"rating"`
	ReviewsCount  int       `json:"reviewsCount"`
	IsFeatured    bool      `json:"isFeatured"`
	CreatedAt     time.Time `json:"createdAt"`


	// Toy Safety & Compliance Metadata
	ListingStatus       string   `json:"listingStatus"` // DRAFT, PENDING_COMPLIANCE, ACTIVE, REJECTED, RECALLED
	AgeGrading          string   `json:"ageGrading"`    // 0-12m, 12-24m, 2-4y, 5-7y, 8-11y, 12+
	HasSmallParts       bool     `json:"hasSmallParts"`
	HasSmallBall        bool     `json:"hasSmallBall"`
	HasMarble           bool     `json:"hasMarble"`
	HasBalloon          bool     `json:"hasBalloon"`
	LegalWarningCopy    string   `json:"legalWarningCopy"`

	// Legal Metrology
	CountryOfOrigin     string   `json:"countryOfOrigin"`
	ManufacturerName    string   `json:"manufacturerName"`
	ManufacturerAddress string   `json:"manufacturerAddress"`
	ImporterName        string   `json:"importerName"`
	ImporterAddress     string   `json:"importerAddress"`
	CustomerCareContact string   `json:"customerCareContact"`
	NetQuantity         string   `json:"netQuantity"`

	// Logistics & Physical Specs
	Materials           []string `json:"materials" gorm:"serializer:json"`
	IsNonToxic          bool     `json:"isNonToxic"`
	IsBpaFree           bool     `json:"isBpaFree"`
	IsWashable          bool     `json:"isWashable"`
	BatteryRequirement  string   `json:"batteryRequirement"` // None, Included, Required
	BatteryChemistry    string   `json:"batteryChemistry"`   // Alkaline, Lithium-Ion, N/A
	Un383ReportUrl      string   `json:"un383ReportUrl"`
	HasScrewLockBattery bool     `json:"hasScrewLockBattery"`
	NetWeightGrams      float64  `json:"netWeightGrams"`
	GrossWeightGrams    float64  `json:"grossWeightGrams"`
	GrossWidthCm        float64  `json:"grossWidthCm"`
	RequiresSuffocationTag bool  `json:"requiresSuffocationTag"`

	// Traceability & Certifications
	BatchNumber         string                `json:"batchNumber"`
	Certifications      []SafetyCertification `json:"certifications" gorm:"serializer:json"`
}

type SafetyCertification struct {
	ID                 string    `json:"id"`
	CertType           string    `json:"certType"` // BIS_IS_9873, ASTM_F963_CPC, EN_71_CE
	CertificateNumber string    `json:"certificateNumber"`
	IssuingLab         string    `json:"issuingLab"`
	IssueDate          time.Time `json:"issueDate"`
	ExpiryDate         time.Time `json:"expiryDate"`
	ReportPdfUrl       string    `json:"reportPdfUrl"`
	VerificationStatus string    `json:"verificationStatus"` // PENDING, APPROVED, REJECTED
	RejectionReason    string    `json:"rejectionReason,omitempty"`
}

type ProductBatch struct {
	ID             string    `json:"id"`
	BatchNumber    string    `json:"batchNumber"`
	ProductID      string    `json:"productId"`
	ProductionDate time.Time `json:"productionDate"`
	TotalQuantity  int       `json:"totalQuantity"`
	Status         string    `json:"status"` // ACTIVE, RECALLED
	RecallReason   string    `json:"recallReason,omitempty"`
	RecalledAt     time.Time `json:"recalledAt,omitempty"`
}

type CreateProductRequest struct {
	SKU                 string   `json:"sku"`
	Name                string   `json:"name" binding:"required"`
	Description         string   `json:"description"`
	Price               float64  `json:"price"`
	RegularPrice        float64  `json:"regularPrice"`
	CategoryID          string   `json:"categoryId"`
	SubCategoryID       string   `json:"subCategoryId"`
	Category            string   `json:"category"`
	SubCategory         string   `json:"subCategory"`
	ImageURL            string   `json:"imageUrl"`
	Images              []string `json:"images"`
	Stock               int      `json:"stock"`
	IsFeatured          bool     `json:"isFeatured"`

	// Safety & Metrology
	AgeGrading          string   `json:"ageGrading"`
	HasSmallParts       bool     `json:"hasSmallParts"`
	HasSmallBall        bool     `json:"hasSmallBall"`
	HasMarble           bool     `json:"hasMarble"`
	HasBalloon          bool     `json:"hasBalloon"`
	CountryOfOrigin     string   `json:"countryOfOrigin"`
	ManufacturerName    string   `json:"manufacturerName"`
	ManufacturerAddress string   `json:"manufacturerAddress"`
	ImporterName        string   `json:"importerName"`
	ImporterAddress     string   `json:"importerAddress"`
	CustomerCareContact string   `json:"customerCareContact"`
	NetQuantity         string   `json:"netQuantity"`
	Materials           []string `json:"materials"`
	IsNonToxic          bool     `json:"isNonToxic"`
	IsBpaFree           bool     `json:"isBpaFree"`
	IsWashable          bool     `json:"isWashable"`
	BatteryRequirement  string   `json:"batteryRequirement"`
	BatteryChemistry    string   `json:"batteryChemistry"`
	Un383ReportUrl      string   `json:"un383ReportUrl"`
	HasScrewLockBattery bool     `json:"hasScrewLockBattery"`
	NetWeightGrams      float64  `json:"netWeightGrams"`
	GrossWeightGrams    float64  `json:"grossWeightGrams"`
	GrossWidthCm        float64  `json:"grossWidthCm"`
	BatchNumber         string   `json:"batchNumber"`
}

type UpdateProductRequest struct {
	SKU                 string   `json:"sku"`
	Name                string   `json:"name"`
	Description         string   `json:"description"`
	Price               float64  `json:"price"`
	RegularPrice        float64  `json:"regularPrice"`
	CategoryID          string   `json:"categoryId"`
	SubCategoryID       string   `json:"subCategoryId"`
	Category            string   `json:"category"`
	SubCategory         string   `json:"subCategory"`
	ImageURL            string   `json:"imageUrl"`
	Images              []string `json:"images"`
	Stock               int      `json:"stock"`
	IsFeatured          bool     `json:"isFeatured"`

	// Safety & Metrology
	AgeGrading          string   `json:"ageGrading"`
	HasSmallParts       bool     `json:"hasSmallParts"`
	HasSmallBall        bool     `json:"hasSmallBall"`
	HasMarble           bool     `json:"hasMarble"`
	HasBalloon          bool     `json:"hasBalloon"`
	CountryOfOrigin     string   `json:"countryOfOrigin"`
	ManufacturerName    string   `json:"manufacturerName"`
	ManufacturerAddress string   `json:"manufacturerAddress"`
	ImporterName        string   `json:"importerName"`
	ImporterAddress     string   `json:"importerAddress"`
	CustomerCareContact string   `json:"customerCareContact"`
	NetQuantity         string   `json:"netQuantity"`
	Materials           []string `json:"materials"`
	IsNonToxic          bool     `json:"isNonToxic"`
	IsBpaFree           bool     `json:"isBpaFree"`
	IsWashable          bool     `json:"isWashable"`
	BatteryRequirement  string   `json:"batteryRequirement"`
	BatteryChemistry    string   `json:"batteryChemistry"`
	Un383ReportUrl      string   `json:"un383ReportUrl"`
	HasScrewLockBattery bool     `json:"hasScrewLockBattery"`
	NetWeightGrams      float64  `json:"netWeightGrams"`
	GrossWeightGrams    float64  `json:"grossWeightGrams"`
	GrossWidthCm        float64  `json:"grossWidthCm"`
	BatchNumber         string   `json:"batchNumber"`
	ListingStatus       string   `json:"listingStatus"`
}


