package scheduler

import (
	"log"
	"static-hosting-server/internal/config"
	"static-hosting-server/internal/services"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Scheduler struct {
	cron           *cron.Cron
	articleService *services.ArticleService
}

func Start(db *gorm.DB) *Scheduler {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config for scheduler: %v", err)
		return nil
	}

	articleService := services.NewArticleService(db, cfg)

	c := cron.New(cron.WithSeconds())

	scheduler := &Scheduler{
		cron:           c,
		articleService: articleService,
	}

	// 每小时检查一次过期文章
	c.AddFunc("0 0 * * * *", scheduler.cleanupExpiredArticles)

	// 启动定时任务
	c.Start()
	log.Println("Scheduler started")

	return scheduler
}

func (s *Scheduler) cleanupExpiredArticles() {
	log.Println("Starting cleanup of expired articles...")

	if err := s.articleService.CleanupExpiredArticles(); err != nil {
		log.Printf("Failed to cleanup expired articles: %v", err)
	} else {
		log.Println("Expired articles cleanup completed")
	}
}

func (s *Scheduler) Stop() {
	if s.cron != nil {
		s.cron.Stop()
		log.Println("Scheduler stopped")
	}
}
