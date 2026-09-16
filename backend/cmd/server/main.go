package main

import (
	"flag"
	"log"

	"personal-blog/internal/config"
	"personal-blog/internal/handler"
	"personal-blog/internal/router"
	"personal-blog/internal/storage"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := storage.Open(cfg)
	if err != nil {
		log.Fatalf("init db: %v", err)
	}

	authH := handler.NewAuthHandler(db, cfg)
	articleH := handler.NewArticleHandler(db)
	statsH := handler.NewStatsHandler(db)

	r := router.New(cfg, authH, articleH, statsH)
	log.Printf("server listening on %s", cfg.Server.Addr)
	if err := r.Run(cfg.Server.Addr); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
