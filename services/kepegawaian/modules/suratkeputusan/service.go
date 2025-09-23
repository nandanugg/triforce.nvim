package suratkeputusan

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListSKByNIP(ctx context.Context, arg repo.ListSKByNIPParams) ([]repo.ListSKByNIPRow, error)
	GetSKByNIPAndID(ctx context.Context, arg repo.GetSKByNIPAndIDParams) (repo.GetSKByNIPAndIDRow, error)
	CountSKByNIP(ctx context.Context, arg repo.CountSKByNIPParams) (int64, error)
	ListUnitKerjaHierarchyByNIP(ctx context.Context, nip string) ([]repo.ListUnitKerjaHierarchyByNIPRow, error)
	GetBerkasSKByNIPAndID(ctx context.Context, arg repo.GetBerkasSKByNIPAndIDParams) (pgtype.Text, error)
	GetBerkasSKSignedByNIPAndID(ctx context.Context, arg repo.GetBerkasSKSignedByNIPAndIDParams) (pgtype.Text, error)
}

type service struct {
	repo repository
}

type listParams struct {
	Limit      uint
	Offset     uint
	Nip        string
	NoSK       string
	StatusSK   *int32
	KategoriSK string
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, arg listParams) ([]sk, uint, error) {
	statusSkParam := pgtype.Int4{Valid: false}
	if arg.StatusSK != nil {
		statusSkParam = pgtype.Int4{Valid: true, Int32: *arg.StatusSK}
	}

	data, err := s.repo.ListSKByNIP(ctx, repo.ListSKByNIPParams{
		Limit:      int32(arg.Limit),
		Offset:     int32(arg.Offset),
		Nip:        arg.Nip,
		NoSk:       pgtype.Text{Valid: arg.NoSK != "", String: arg.NoSK},
		StatusSk:   statusSkParam,
		KategoriSk: pgtype.Text{Valid: arg.KategoriSK != "", String: arg.KategoriSK},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-list] repo ListSKByNIP: %w", err)
	}

	count, err := s.repo.CountSKByNIP(ctx, repo.CountSKByNIPParams{
		NoSk:       pgtype.Text{Valid: arg.NoSK != "", String: arg.NoSK},
		StatusSk:   statusSkParam,
		KategoriSk: pgtype.Text{Valid: arg.KategoriSK != "", String: arg.KategoriSK},
		Nip:        arg.Nip,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-list] repo CountSKByNIP: %w", err)
	}

	listUnor, err := s.repo.ListUnitKerjaHierarchyByNIP(ctx, arg.Nip)
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-list] repo ListUnitKerjaHierarchyByNIP: %w", err)
	}

	unitKerja := s.getUnorLengkap(listUnor)
	result := typeutil.Map(data, func(row repo.ListSKByNIPRow) sk {
		return sk{
			IDSK:       row.FileID,
			KategoriSK: row.KategoriSk.String,
			NoSK:       row.NoSk.String,
			TanggalSK:  db.Date(row.TanggalSk.Time),
			StatusSK:   statusSKText(row.StatusSk.Int16),
			UnitKerja:  unitKerja,
		}
	})

	return result, uint(count), nil
}

func (s *service) get(ctx context.Context, nip, id string) (*sk, error) {
	data, err := s.repo.GetSKByNIPAndID(ctx, repo.GetSKByNIPAndIDParams{
		Nip: nip,
		ID:  id,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return &sk{}, fmt.Errorf("[suratkeputusan-get] repo GetSKByNIPAndIDParams: %w", err)
	}

	listUnor, err := s.repo.ListUnitKerjaHierarchyByNIP(ctx, nip)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return &sk{}, fmt.Errorf("[suratkeputusan-get] repo ListUnitKerjaHierarchyByNIP: %w", err)
	}

	return &sk{
		IDSK:              id,
		KategoriSK:        data.KategoriSk.String,
		NoSK:              data.NoSk.String,
		TanggalSK:         db.Date(data.TanggalSk.Time),
		StatusSK:          statusSKText(data.StatusSk.Int16),
		UnitKerja:         s.getUnorLengkap(listUnor),
		NamaPemilik:       data.NamaPemilikSk.String,
		NamaPenandaTangan: data.NamaPenandatangan.String,
	}, nil
}

func (s *service) getUnorLengkap(listUnor []repo.ListUnitKerjaHierarchyByNIPRow) string {
	if len(listUnor) == 0 {
		return ""
	}

	var unitKerja []string
	for _, unor := range listUnor {
		unitKerja = append(unitKerja, unor.NamaUnor.String)
	}

	return strings.Join(unitKerja, " - ")
}

func (s *service) getBerkas(ctx context.Context, nip, id string, signed bool) (string, []byte, error) {
	var res pgtype.Text
	var err error
	if signed {
		res, err = s.repo.GetBerkasSKSignedByNIPAndID(ctx, repo.GetBerkasSKSignedByNIPAndIDParams{
			Nip: nip,
			ID:  id,
		})
	} else {
		res, err = s.repo.GetBerkasSKByNIPAndID(ctx, repo.GetBerkasSKByNIPAndIDParams{
			Nip: nip,
			ID:  id,
		})
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", nil, fmt.Errorf("[suratkeputusan-getBerkas] repo get berkas: %w", err)
	}

	if len(res.String) == 0 {
		return "", nil, nil
	}

	return api.GetMimeTypeAndDecodedData(res.String)
}
