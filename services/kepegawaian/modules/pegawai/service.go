package pegawai

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	GetProfilePegawaiByPNSID(ctx context.Context, pnsID string) (sqlc.GetProfilePegawaiByPNSIDRow, error)
	ListUnitKerjaHierarchy(ctx context.Context, id string) ([]sqlc.ListUnitKerjaHierarchyRow, error)
	ListPegawaiAktif(ctx context.Context, arg sqlc.ListPegawaiAktifParams) ([]sqlc.ListPegawaiAktifRow, error)
	ListUnitKerjaLengkapByIDs(ctx context.Context, ids []string) ([]sqlc.ListUnitKerjaLengkapByIDsRow, error)
	CountPegawaiAktif(ctx context.Context, arg sqlc.CountPegawaiAktifParams) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) getProfileByPNSID(ctx context.Context, pnsID string) (*profile, error) {
	data, err := s.repo.GetProfilePegawaiByPNSID(ctx, pnsID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repo get profile: %w", err)
	}

	unitOrganisasi := make([]string, 0)
	if data.UnorID.Valid {
		rows, err := s.repo.ListUnitKerjaHierarchy(ctx, data.UnorID.String)
		if err != nil {
			return nil, fmt.Errorf("repo list unit kerja hierarchy: %w", err)
		}

		unitOrganisasi = typeutil.FilterMap(rows, func(row sqlc.ListUnitKerjaHierarchyRow) (string, bool) {
			return row.NamaUnor.String, row.NamaUnor.String != ""
		})
	}

	return &profile{
		NIPLama:        data.NipLama.String,
		NIPBaru:        data.NipBaru.String,
		GelarDepan:     data.GelarDepan.String,
		GelarBelakang:  data.GelarBelakang.String,
		Nama:           data.Nama.String,
		Pangkat:        data.Pangkat.String,
		Golongan:       typeutil.Cast[string](data.Golongan),
		Jabatan:        data.Jabatan.String,
		UnitOrganisasi: unitOrganisasi,
	}, nil
}

type adminListPegawaiParams struct {
	limit      uint
	offset     uint
	keyword    string
	unitID     string
	golonganID int32
	jabatanID  string
	status     string
}

func (s *service) adminListPegawai(ctx context.Context, arg adminListPegawaiParams) ([]pegawai, uint, error) {
	data, err := s.repo.ListPegawaiAktif(ctx, sqlc.ListPegawaiAktifParams{
		Limit:       int32(arg.limit),
		Offset:      int32(arg.offset),
		Keyword:     pgtype.Text{Valid: arg.keyword != "", String: arg.keyword},
		UnitKerjaID: pgtype.Text{Valid: arg.unitID != "", String: arg.unitID},
		GolonganID:  pgtype.Int4{Valid: arg.golonganID != 0, Int32: arg.golonganID},
		JabatanID:   pgtype.Text{Valid: arg.jabatanID != "", String: arg.jabatanID},
		StatusHukum: getStatusHukum(arg.status),
		Mpp:         statusPNSMPP,
		StatusPns:   pgtype.Text{Valid: arg.status != "", String: arg.status},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[pegawai-adminListPegawai] repo ListPegawaiAktif: %w", err)
	}

	uniqUnorIDs := typeutil.UniqMap(data, func(row sqlc.ListPegawaiAktifRow, _ int) string {
		return row.UnorID.String
	})

	listUnorLengkap, err := s.repo.ListUnitKerjaLengkapByIDs(ctx, uniqUnorIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("[pegawai-adminListPegawai] repo ListUnitKerjaLengkapByIDs: %w", err)
	}
	unorLengkapByID := typeutil.SliceToMap(listUnorLengkap, func(unorLengkap sqlc.ListUnitKerjaLengkapByIDsRow) (string, string) {
		return unorLengkap.ID, unorLengkap.NamaUnorLengkap
	})

	count, err := s.repo.CountPegawaiAktif(ctx, sqlc.CountPegawaiAktifParams{
		Keyword:     pgtype.Text{Valid: arg.keyword != "", String: arg.keyword},
		UnitKerjaID: pgtype.Text{Valid: arg.unitID != "", String: arg.unitID},
		GolonganID:  pgtype.Int4{Valid: arg.golonganID != 0, Int32: arg.golonganID},
		JabatanID:   pgtype.Text{Valid: arg.jabatanID != "", String: arg.jabatanID},
		StatusHukum: getStatusHukum(arg.status),
		StatusPns:   pgtype.Text{Valid: arg.status != "", String: arg.status},
		Mpp:         statusPNSMPP,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[pegawai-adminListPegawai] repo CountPegawaiAktif: %w", err)
	}

	result := typeutil.Map(data, func(row sqlc.ListPegawaiAktifRow) pegawai {
		return pegawai{
			NIP:           row.Nip.String,
			GelarDepan:    row.GelarDepan.String,
			Nama:          row.Nama.String,
			GelarBelakang: row.GelarBelakang.String,
			Golongan:      row.Golongan.String,
			Jabatan:       row.Jabatan.String,
			UnitKerja:     unorLengkapByID[row.UnorID.String],
			Status:        row.Status,
		}
	})

	return result, uint(count), nil
}
