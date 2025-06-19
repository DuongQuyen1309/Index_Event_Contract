package main

import (
	"context"
	"fmt"

	"github.com/DuongQuyen1309/indexevent/internal/datastore"
	"github.com/DuongQuyen1309/indexevent/internal/db"
	"github.com/DuongQuyen1309/indexevent/internal/router"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("error load env file", err)
	}
}
func main() {
	ctx := context.Background()
	db.ConnectDB()
	if err := datastore.CreateRequestCreatedEvent(ctx); err != nil {
		fmt.Println("error create request created event table", err)
		return
	}
	if err := datastore.CreateResponseCreatedEvent(ctx); err != nil {
		fmt.Println("error create response created event table", err)
		return
	}
	// if err := service.IndexEvent(ctx); err != nil {
	// 	fmt.Println("error index event", err)
	// 	return
	// }
	datastoreInDB := &datastore.DataStoreInDB{}
	router := router.SetupRouter(datastoreInDB)
	router.Run(":8080")
	// router := gin.Default()
	// handler.NewHandler(&handler.Config{
	// 	R: router,
	// })
	// srv := &http.Server{
	// 	Addr:    ":8080",
	// 	Handler: router,
	// }
	// srv.ListenAndServe()
}
