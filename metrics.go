package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metricsProcessor struct {
	em           entriesManager
	nbContainers *prometheus.GaugeVec
}

func newMetricsProcessor(em entriesManager) *metricsProcessor {
	nbContainers := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ghosts_nb_containers",
			Help: "The total number of created containers",
		},
		[]string{
			"category",
		},
	)
	return &metricsProcessor{em, nbContainers}
}

func (h *metricsProcessor) init() error {
	entries, err := h.em.get()
	if err != nil {
		return err
	}
	for _, entry := range entries {
		h.nbContainers.WithLabelValues(entry.Category[0]).Inc()
	}

	return nil
}

func (h *metricsProcessor) startEvent(id string) error {
	entries, err := h.em.get(id)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		h.nbContainers.WithLabelValues(entry.Category[0]).Inc()
	}
	return nil
}

func (h *metricsProcessor) dieEvent(id string) error {
	entries, err := h.em.get(id)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		h.nbContainers.WithLabelValues(entry.Category[0]).Dec()
	}
	return nil
}
