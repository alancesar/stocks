package main

import (
	"context"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"stocks/internal/okanebox"
	"stocks/internal/repositoty"
	"stocks/printer"
	"stocks/usecase"
)

var (
	lastPriceUseCase          *usecase.GetLastPrice
	createBuyOperationUseCase *usecase.BuyOperationUseCase
	reportUseCase             *usecase.ReportUseCase
)

func init() {
	db, err := gorm.Open(sqlite.Open("stocks.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed to connect database")
	}

	database := repositoty.NewGormDatabase(db)
	provider := okanebox.NewProvider(http.DefaultClient)

	lastPriceUseCase = usecase.NewGetLastPrice(provider)
	createBuyOperationUseCase = usecase.NewBuyOperationUseCase(database)
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
	case "report":
		report, err := reportUseCase.Execute(ctx)
		if err != nil {
			log.Fatalln(err)
		}

		if err := printer.Print(os.Stdout, report); err != nil {
			log.Fatalln(err)
		}
	}
}
