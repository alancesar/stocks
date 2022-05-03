package main

import (
	"context"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"stocks/internal/mfinance"
	"stocks/internal/repository"
	"stocks/separator"
	"stocks/usecase"
)

var (
	lastPriceUseCase          *usecase.GetLastPrice
	createBuyOperationUseCase *usecase.BuyOperationUseCase
	listUseCase               *usecase.ListUseCase
	reportUseCase             *usecase.ReportUseCase
)

func init() {
	db, err := gorm.Open(sqlite.Open("stocks.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed to connect database")
	}

	database := repository.NewGormDatabase(db)
	provider := mfinance.NewProvider(http.DefaultClient)

	lastPriceUseCase = usecase.NewGetLastPrice(provider)
	createBuyOperationUseCase = usecase.NewBuyOperationUseCase(database)
	listUseCase = usecase.NewListUseCase(database)
	reportUseCase = usecase.NewReportUseCase(provider, database)
}

func main() {
	ctx := context.Background()

	switch os.Args[1] {
	case "buy":
		request, err := CreateBuyRequest(os.Args[2:]...)
		if err != nil {
			log.Fatalln(err)
		}

		operation, err := createBuyOperationUseCase.Execute(ctx, request)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("operation created succesffully: %v\n", operation)
	case "price":
		request, err := CreatePriceRequest(os.Args[2:]...)
		if err != nil {
			log.Fatalln(err)
		}

		lastPrice, err := lastPriceUseCase.Execute(ctx, request)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("%s: %.2f\n", request, lastPrice)
	case "list":
		operations, err := listUseCase.Execute(ctx)
		if err != nil {
			log.Fatalln(err)
		}

		if err := operations.Print(os.Stdout, separator.Tab); err != nil {
			log.Fatalln(err)
		}
	case "export":
		output := "stocks.csv"
		if len(os.Args) > 2 {
			output = os.Args[2]
		}

		file, err := os.Create(output)
		if err != nil {
			log.Fatalln(err)
		}

		defer func() {
			_ = file.Close()
		}()

		report, err := listUseCase.Execute(ctx)
		if err != nil {
			log.Fatalln(err)
		}

		if err := report.Print(file, separator.Comma); err != nil {
			log.Fatalln(err)
		}
	case "report":
		report, err := reportUseCase.Execute(ctx)
		if err != nil {
			log.Fatalln(err)
		}

		if err := report.Print(os.Stdout, separator.Tab); err != nil {
			log.Fatalln(err)
		}
	}
}
