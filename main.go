package main

import (
	"fmt"
	"time"

	db "github.com/dihanto/cli/config"
	"github.com/dihanto/cli/controller"
	"github.com/dihanto/cli/repository"
	"github.com/dihanto/cli/usecase"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "A brief description of your application",
	Long:  `A longer description of your application`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, World!")
	},
}

func main() {

	productRepository := repository.NewMysqlProductRepostitory(db.GetDB())

	// Create a new instance of productUsecase with a timeout of 10 seconds
	productUsecase := usecase.NewProductUsecase(productRepository, 10*time.Second)

	// Create a new instance of the ProductController and add a subcommand to the root command
	productController := controller.NewProductController(productUsecase, rootCmd)
	productController.InsertCmd()
	productController.ShowCmd()
	productController.UpdateCmd()
	productController.DeleteCmd()
	productController.SelectCmd()
	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
