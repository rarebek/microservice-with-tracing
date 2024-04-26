package app

import (
	"log"
	"os"
	"os/signal"
	"product_service/internal/app"
	"product_service/internal/pkg/config"
	"syscall"

	"github.com/spf13/cobra"
)

const (
	PRODUCT_CREATE_CONSUMER = "user_create_consumer"
)

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "To run consumer give the name followed by arguments consumer",
	Long: `Example : 
	go run cmd/main.go consumer name_of_consumer`,
	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		consumerName := args[0]

		switch consumerName {
		case PRODUCT_CREATE_CONSUMER:
			UserCreateConsumerRun()
		default:
			log.Fatalf("No consumer with name: '%s'", consumerName)
		}
	},
}

func init() {
	rootCmd.AddCommand(consumerCmd)
}

func UserCreateConsumerRun() {
	config := config.New()

	app, err := app.NewProductConsumer(config)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		app.Run()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	app.Logger.Info("user service stops")

	// stop app
	app.Close()
}
