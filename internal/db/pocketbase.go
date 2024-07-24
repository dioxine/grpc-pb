package pocketbase

import (
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func PbInit() *pocketbase.PocketBase {
	cfg := &pocketbase.Config{
		DefaultDataDir:  "./internal/db/pb_data",
		DefaultDev:      false,
		HideStartBanner: false,
	}

	db := pocketbase.NewWithConfig(*cfg)

	// serves static files from the provided public dir (if exists)
	db.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
		return nil
	})

	err := db.Bootstrap()
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func PbStart(db *pocketbase.PocketBase) {
	// Start server
	db.RootCmd.SetArgs([]string{"serve"})
	log.Println("Starting database")
	if err := db.Start(); err != nil {
		log.Fatal(err)
	}
}
