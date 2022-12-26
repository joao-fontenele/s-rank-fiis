package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joao-fontenele/s-rank-fiis/internal/model"
	"github.com/joao-fontenele/s-rank-fiis/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

const (
	DBURL = "postgresql://usr:pwd@localhost:5432/ranks" // TODO: get from env variable
)

func parseOperations(purchasesPath string) ([]model.Operation, error) {
	var ops []model.Operation
	f, err := os.Open(purchasesPath)
	if err != nil {
		return ops, err
	}

	defer func() {
		_ = f.Close()
	}()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for i, line := range data {
		if i <= 0 {
			continue
		}
		d, err := time.Parse("2006-01-02", line[1])
		if err != nil {
			log.Printf("failed to parse date: %v", err)
		}
		amount, _ := strconv.Atoi(line[2])
		price, _ := strconv.ParseFloat(line[3], 64)
		ops = append(ops, model.Operation{
			Code:         line[0],
			PurchaseDate: d,
			Amount:       amount,
			Price:        price,
		})
	}

	return ops, nil
}

func main() {
	ops, err := parseOperations("./purchases.csv")
	if err != nil {
		log.Fatalf("failed to parse operations: %v", err)
	}

	dbGorm, err := gorm.Open(postgres.Open(DBURL), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}

	opRepo := repository.Operation{Conn: dbGorm}
	for _, op := range ops {
		err := opRepo.Save(&op)
		if err != nil {
			log.Fatalf("failed to save operation: %v", err)
		}
	}
}
