package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:255"`
	Email string `gorm:"uniqueIndex"`
}

type Company struct {
	CompanyID   string `gorm:"primaryKey;type:char(36)"`
	CompanyName string `gorm:"size:255"`
}

func (Company) TableName() string {
	return "company"
}

type Job struct {
	UUID       string  `gorm:"primaryKey;type:char(36)"`
	JobTitle   string  `gorm:"size:255"`
	CompanyID  string  `gorm:"type:char(36)"`
	Location   string  `gorm:"size:255"`
	Salary     float64 `gorm:"type:decimal(10,2)"`
	PostedDate string  `gorm:"type:date"`
}

type File struct {
	FileID   string `gorm:"primaryKey;type:char(36)"`
	FileName string `gorm:"size:255"`
	JobID    string `gorm:"type:char(36)"`
}

func main() {
	eagerLoad := flag.Bool("eager", false, "Run eager loading example")
	join := flag.Bool("join", false, "Run join example")
	lazyLoad := flag.Bool("lazy", false, "Run lazy loading example")
	jsonAggregate := flag.Bool("json", false, "Run JSON aggregate example")
	all := flag.Bool("all", false, "Run all examples and print metrics")

	flag.Parse()

	// Get database connection details from environment variables
	dbHost := getEnv("DB_HOST", "db")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "gorm")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// Add retry logic with better error handling
	var db *gorm.DB
	var err error
	maxRetries := 30 // Increase retries
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Printf("Successfully connected to database")
			break
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("failed to connect database after retries:", err)
	}

	if *eagerLoad {
		runEagerLoad(db)
	} else if *join {
		runJoin(db)
	} else if *lazyLoad {
		runLazyLoad(db)
	} else if *jsonAggregate {
		runJsonAggregate(db)
	} else if *all {
		runAllExamples(db)
	} else {
		runAllExamples(db)
	}
}

// Helper function to get environment variables with default values
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func runEagerLoad(db *gorm.DB) {
	start := time.Now()
	// Update the Job struct to include relationships
	type Job struct {
		UUID       string  `gorm:"primaryKey;type:char(36)"`
		JobTitle   string  `gorm:"size:255"`
		CompanyID  string  `gorm:"type:char(36)"`
		Location   string  `gorm:"size:255"`
		Salary     float64 `gorm:"type:decimal(10,2)"`
		PostedDate string  `gorm:"type:date"`
		Company    Company `gorm:"foreignKey:CompanyID;references:CompanyID"`
		Files      []File  `gorm:"foreignKey:JobID;references:UUID"`
	}

	var jobs []Job

	// Preload both Company and Files relationships
	result := db.Model(&Job{}).
		Preload("Company", func(db *gorm.DB) *gorm.DB {
			return db.Table("company")
		}).
		Preload("Files").
		Find(&jobs)

	if result.Error != nil {
		log.Fatal("failed to fetch jobs:", result.Error)
	}

	jsonData, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		log.Fatal("failed to marshal JSON")
	}

	err = os.WriteFile("eager_load_output.json", jsonData, 0644)
	if err != nil {
		log.Fatal("failed to write JSON to file")
	}

	elapsed := time.Since(start)
	log.Printf("runEagerLoad took %s", elapsed)
}

func runJoin(db *gorm.DB) {
	start := time.Now()

	type JobDetails struct {
		JobID       string  `json:"job_id"`
		JobTitle    string  `json:"job_title"`
		CompanyID   string  `json:"company_id"`
		CompanyName string  `json:"company_name"`
		Location    string  `json:"location"`
		Salary      float64 `json:"salary"`
		PostedDate  string  `json:"posted_date"`
		FileID      string  `json:"file_id"`
		FileName    string  `json:"file_name"`
	}

	var jobDetailsList []JobDetails

	db.Table("jobs").
		Select("jobs.uuid as job_id, jobs.job_title, jobs.company_id, company.company_name, jobs.location, jobs.salary, jobs.posted_date, files.file_id, files.file_name").
		Joins("left join company on jobs.company_id = company.company_id").
		Joins("left join files on jobs.uuid = files.job_id").
		Scan(&jobDetailsList)

	jsonData, err := json.MarshalIndent(jobDetailsList, "", "  ")
	if err != nil {
		log.Fatal("failed to marshal JSON")
	}

	err = os.WriteFile("join_output.json", jsonData, 0644)
	if err != nil {
		log.Fatal("failed to write JSON to file")
	}

	elapsed := time.Since(start)
	log.Printf("runJoin took %s", elapsed)
}

func runLazyLoad(db *gorm.DB) {
	start := time.Now()

	var jobs []Job
	db.Find(&jobs).Limit(50)

	type JobDetails struct {
		Job     Job     `json:"job"`
		Company Company `json:"company"`
		Files   []File  `json:"files"`
	}

	var jobDetailsList []JobDetails

	for _, job := range jobs {
		var company Company
		var files []File

		db.First(&company, "company_id = ?", job.CompanyID)
		db.Where("job_id = ?", job.UUID).Find(&files)

		jobDetails := JobDetails{
			Job:     job,
			Company: company,
			Files:   files,
		}

		jobDetailsList = append(jobDetailsList, jobDetails)
	}

	jsonData, err := json.MarshalIndent(jobDetailsList, "", "  ")
	if err != nil {
		log.Fatal("failed to marshal JSON")
	}

	err = os.WriteFile("lazy_load_output.json", jsonData, 0644)
	if err != nil {
		log.Fatal("failed to write JSON to file")
	}

	elapsed := time.Since(start)
	log.Printf("runLazyLoad took %s", elapsed)
}

func runJsonAggregate(db *gorm.DB) {
	start := time.Now()

	type JobDetails struct {
		JobID       string  `json:"job_id"`
		JobTitle    string  `json:"job_title"`
		CompanyID   string  `json:"company_id"`
		CompanyName string  `json:"company_name"`
		Location    string  `json:"location"`
		Salary      float64 `json:"salary"`
		PostedDate  string  `json:"posted_date"`
		FileID      string  `json:"file_id"`
		FileName    string  `json:"file_name"`
		Company     string  `json:"company"`
	}

	var jobDetailsList []JobDetails

	db.Raw(`
		SELECT 
			jobs.uuid as job_id, 
			jobs.job_title, 
			jobs.company_id, 
			(
				SELECT JSON_OBJECT('company_id', company.company_id, 'company_name', company.company_name) 
				FROM company 
				WHERE company.company_id = jobs.company_id
			) as company,
			jobs.location, 
			jobs.salary, 
			jobs.posted_date, 
			(
				SELECT JSON_ARRAYAGG(JSON_OBJECT('file_id', files.file_id, 'file_name', files.file_name)) 
				FROM files 
				WHERE files.job_id = jobs.uuid
			) as files
		FROM jobs
	`).Scan(&jobDetailsList)

	jsonData, err := json.MarshalIndent(jobDetailsList, "", "  ")
	if err != nil {
		log.Fatal("failed to marshal JSON")
	}

	err = os.WriteFile("json_aggregate_output.json", jsonData, 0644)
	if err != nil {
		log.Fatal("failed to write JSON to file")
	}

	elapsed := time.Since(start)
	log.Printf("runJsonAggregate took %s", elapsed)
}

func runAllExamples(db *gorm.DB) {
	start := time.Now()
	runEagerLoad(db)
	eagerLoadTime := time.Since(start)

	start = time.Now()
	runJoin(db)
	joinTime := time.Since(start)

	start = time.Now()
	runLazyLoad(db)
	lazyLoadTime := time.Since(start)

	start = time.Now()
	runJsonAggregate(db)
	jsonAggregateTime := time.Since(start)

	log.Printf("Eager Loading took %s", eagerLoadTime)
	log.Printf("Join took %s", joinTime)
	log.Printf("Lazy Loading took %s", lazyLoadTime)
	log.Printf("JSON Aggregate took %s", jsonAggregateTime)
}
