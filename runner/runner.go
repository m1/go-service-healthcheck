package runner

import (
	"sync"
	"time"

	"github.com/m1/go-service-healthcheck/config"
	"github.com/m1/go-service-healthcheck/log"
	"github.com/m1/go-service-healthcheck/repositories"
	"github.com/m1/go-service-healthcheck/storage"
)

type Runner struct {
	Config       *config.Config
	DB           *storage.SQLiteDB
	Logger       *log.Logger
	Repositories *Repositories

	Worker *worker
}

type Repositories struct {
	Service repositories.ServiceRepository
}

// New ...
func New() *Runner {
	return &Runner{}
}

func (r *Runner) Run() error {
	var err error
	r.Config, err = config.LoadConfig()
	if err != nil {
		return err
	}

	r.Logger, err = log.NewLogger(r.Config)
	if err != nil {
		return err
	}

	r.DB, err = storage.NewSQLiteDB(r.Config.DBConfig)
	if err != nil {
		return err
	}

	// register repos
	repos := &Repositories{
		Service: repositories.NewServiceRepository(r.DB),
	}
	r.Repositories = repos

	r.Worker = newWorker(r.Config, r.Logger, r.Repositories)

	return r.tick()
}

func (r *Runner) tick() error {
	ticker := time.NewTicker(time.Duration(r.Config.RunnerConfig.ScrapeIntervalSeconds) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := r.fetchServiceStatuses()
			if err != nil {
				return err
			}
		}
	}
}

func (r *Runner) fetchServiceStatuses() error {
	services, err := r.Repositories.Service.GetAll()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(services))
	for _, service := range services {
		// implement work pool here
		go r.Worker.work(&wg, service)
	}
	wg.Wait()

	return nil
}
