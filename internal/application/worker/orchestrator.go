package worker

import (
	"fmt"

	"tsniper/pkg/logger"
)

type registeredWorker struct {
	Worker
	name string
}

type Orchestrator struct {
	workers     []*registeredWorker
	workerNames map[string]struct{}
	logger      logger.Logger
}

func NewOrchestrator(logger logger.Logger) *Orchestrator {
	return &Orchestrator{
		workers:     make([]*registeredWorker, 0),
		workerNames: make(map[string]struct{}),
		logger:      logger,
	}
}

func (s *Orchestrator) RegisterWorker(name string, worker Worker) {
	if _, ok := s.workerNames[name]; ok {
		s.logger.Warn("Worker %s already registered", name)
		return
	}
	s.workerNames[name] = struct{}{}

	registeredWorker := &registeredWorker{
		name:   name,
		Worker: worker,
	}
	s.workers = append(s.workers, registeredWorker)
	s.logger.Info("Registered worker %s", name)
}

func (s *Orchestrator) StartAll() error {
	for _, worker := range s.workers {
		if err := worker.Start(); err != nil {
			return fmt.Errorf("failed to start %s worker: %v", worker.name, err)
		}
		s.logger.Info("Started worker %s", worker.name)
	}
	return nil
}

func (s *Orchestrator) StopAll() error {
	for _, worker := range s.workers {
		if err := worker.Stop(); err != nil {
			return fmt.Errorf("failed to stop %s worker: %v", worker.name, err)
		}
		s.logger.Info("Stopped worker %s", worker.name)
	}
	return nil
}
