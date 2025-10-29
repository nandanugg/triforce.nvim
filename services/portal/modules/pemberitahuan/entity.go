package pemberitahuan

import (
	"time"
)

type pemberitahuan struct {
	ID                 int64     `json:"id"`
	JudulBerita        string    `json:"judul_berita"`
	DeskripsiBerita    string    `json:"deskripsi_berita"`
	Pinned             bool      `json:"pinned"`
	Status             string    `json:"status"`
	DiterbitkanPada    time.Time `json:"diterbitkan_pada"`
	DitarikPada        time.Time `json:"ditarik_pada"`
	DiperbaruiOleh     string    `json:"diperbarui_oleh"`
	TerakhirDiperbarui time.Time `json:"terakhir_diperbarui"`
}

type Status string

const (
	PemberitahuanStatusWaiting Status = "WAITING"
	PemberitahuanStatusActive  Status = "ACTIVE"
	PemberitahuanStatusOver    Status = "OVER"
	PemberitahuanStatusAll     Status = "ALL"
)
