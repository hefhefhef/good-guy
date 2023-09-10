package jobs

import (
	"time"

	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/external/ctms"
	"github.com/openuniland/good-guy/external/facebook"
	"github.com/openuniland/good-guy/internal/articles"
	"github.com/openuniland/good-guy/internal/users"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

type Jobs struct {
	cfg        *configs.Configs
	articleUC  articles.UseCase
	userUC     users.UseCase
	facebookUC facebook.UseCase
	ctmsUS     ctms.UseCase
}

func NewJobs(cfg *configs.Configs, articleUC articles.UseCase, userUC users.UseCase, facebookUC facebook.UseCase, ctmsUS ctms.UseCase) *Jobs {
	return &Jobs{cfg: cfg, articleUC: articleUC, userUC: userUC, facebookUC: facebookUC, ctmsUS: ctmsUS}
}

func (j *Jobs) Run() {

	c := cron.New(cron.WithSeconds())

	// every 5 seconds
	// c.AddFunc("*/5 * * * * *", func() {
	// 	log.Info().Msgf("[JOBS]:[TEST]:[TIME=%v]", time.Now())
	// })

	c.AddFunc("0 0 23 * * *", func() {
		log.Info().Msgf("[JOBS]:[0 0 23 * * *]:[TEST]:[TIME=%v]", time.Now())
	})

	c.AddFunc("0 20 22 * * *", func() {
		log.Info().Msgf("[JOBS]:[0 20 22 * * *]:[TEST]:[TIME=%v]", time.Now())
	})

	// TEST

	//every 10 minutes
	c.AddFunc("*/600 * * * * *", func() {
		log.Info().Msgf("[JOBS]:[Start sync articles]:[TIME=%v]", time.Now())
		go j.syncArticles()
	})

	//20h every day
	c.AddFunc("0 0 20 * * *", func() {
		log.Info().Msg("Running getUpcomingExamSchedule")
		go j.getUpcomingExamSchedule()
	})

	// 6h45 am every day
	c.AddFunc("0 45 6 * * *", func() {
		log.Info().Msg("Running morningClassSchedule")
		go j.morningClassSchedule()
	})

	// 12h00 pm every day
	c.AddFunc("0 0 12 * * *", func() {
		log.Info().Msg("Running afternoonClassSchedule")
		go j.afternoonClassSchedule()
	})

	// 16h30 pm every day
	c.AddFunc("0 30 16 * * *", func() {
		log.Info().Msg("Running eveningClassSchedule")
		go j.eveningClassSchedule()
	})

	c.Start()

}
