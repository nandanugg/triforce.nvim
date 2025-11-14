package pemberitahuan

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type pemberitahuan struct {
	ID                    int64              `json:"id"`
	JudulBerita           string             `json:"judul_berita"`
	DeskripsiBerita       string             `json:"deskripsi_berita"`
	IsCurrentPeriodPinned *bool              `json:"is_current_period_pinned,omitempty"`
	PinnedAt              pgtype.Timestamptz `json:"pinned_at"`
	Status                string             `json:"status"`
	DiterbitkanPada       time.Time          `json:"diterbitkan_pada"`
	DitarikPada           time.Time          `json:"ditarik_pada"`
	DiperbaruiOleh        string             `json:"diperbarui_oleh"`
	TerakhirDiperbarui    time.Time          `json:"terakhir_diperbarui"`
}

type status string

const (
	pemberitahuanStatusActive status = "ACTIVE"
	pemberitahuanStatusAll    status = "ALL"
)
