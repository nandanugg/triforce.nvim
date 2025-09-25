package suratkeputusan

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	CountSuratKeputusan(ctx context.Context, arg repo.CountSuratKeputusanParams) (int64, error)
	CountSuratKeputusanByNIP(ctx context.Context, arg repo.CountSuratKeputusanByNIPParams) (int64, error)
	GetBerkasSuratKeputusanByID(ctx context.Context, id string) (pgtype.Text, error)
	GetBerkasSuratKeputusanByNIPAndID(ctx context.Context, arg repo.GetBerkasSuratKeputusanByNIPAndIDParams) (pgtype.Text, error)
	GetBerkasSuratKeputusanSignedByID(ctx context.Context, id string) (pgtype.Text, error)
	GetBerkasSuratKeputusanSignedByNIPAndID(ctx context.Context, arg repo.GetBerkasSuratKeputusanSignedByNIPAndIDParams) (pgtype.Text, error)
	GetSuratKeputusanByID(ctx context.Context, id string) (repo.GetSuratKeputusanByIDRow, error)
	GetSuratKeputusanByNIPAndID(ctx context.Context, arg repo.GetSuratKeputusanByNIPAndIDParams) (repo.GetSuratKeputusanByNIPAndIDRow, error)
	ListSuratKeputusan(ctx context.Context, arg repo.ListSuratKeputusanParams) ([]repo.ListSuratKeputusanRow, error)
	ListSuratKeputusanByNIP(ctx context.Context, arg repo.ListSuratKeputusanByNIPParams) ([]repo.ListSuratKeputusanByNIPRow, error)
	ListUnitKerjaHierarchyByNIP(ctx context.Context, nip string) ([]repo.ListUnitKerjaHierarchyByNIPRow, error)
	ListLogSuratKeputusanByID(ctx context.Context, id string) ([]repo.ListLogSuratKeputusanByIDRow, error)
	ListUnitKerjaLengkapByIDs(ctx context.Context, ids []string) ([]repo.ListUnitKerjaLengkapByIDsRow, error)
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

func (s *service) list(ctx context.Context, arg listParams) ([]suratKeputusan, uint, error) {
	statusSkParam := pgtype.Int4{Valid: false}
	if arg.StatusSK != nil {
		statusSkParam = pgtype.Int4{Valid: true, Int32: *arg.StatusSK}
	}

	data, err := s.repo.ListSuratKeputusanByNIP(ctx, repo.ListSuratKeputusanByNIPParams{
		Limit:      int32(arg.Limit),
		Offset:     int32(arg.Offset),
		Nip:        arg.Nip,
		NoSk:       pgtype.Text{Valid: arg.NoSK != "", String: arg.NoSK},
		StatusSk:   statusSkParam,
		KategoriSk: pgtype.Text{Valid: arg.KategoriSK != "", String: arg.KategoriSK},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-list] repo ListSuratKeputusanByNIP: %w", err)
	}

	count, err := s.repo.CountSuratKeputusanByNIP(ctx, repo.CountSuratKeputusanByNIPParams{
		NoSk:       pgtype.Text{Valid: arg.NoSK != "", String: arg.NoSK},
		StatusSk:   statusSkParam,
		KategoriSk: pgtype.Text{Valid: arg.KategoriSK != "", String: arg.KategoriSK},
		Nip:        arg.Nip,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-list] repo CountSuratKeputusanByNIP: %w", err)
	}

	listUnor, err := s.repo.ListUnitKerjaHierarchyByNIP(ctx, arg.Nip)
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-list] repo ListUnitKerjaHierarchyByNIP: %w", err)
	}

	unitKerja := s.getUnorLengkap(listUnor)
	result := typeutil.Map(data, func(row repo.ListSuratKeputusanByNIPRow) suratKeputusan {
		return suratKeputusan{
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

func (s *service) get(ctx context.Context, nip, id string) (*suratKeputusan, error) {
	data, err := s.repo.GetSuratKeputusanByNIPAndID(ctx, repo.GetSuratKeputusanByNIPAndIDParams{
		Nip: nip,
		ID:  id,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[suratkeputusan-get] repo GetSuratKeputusanByNIPAndIDParams: %w", err)
	}

	listUnor, err := s.repo.ListUnitKerjaHierarchyByNIP(ctx, nip)
	if err != nil {
		return nil, fmt.Errorf("[suratkeputusan-get] repo ListUnitKerjaHierarchyByNIP: %w", err)
	}

	logs, err := s.repo.ListLogSuratKeputusanByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("[suratkeputusan-get] repo ListLogSKBySKID: %w", err)
	}

	logSk := typeutil.Map(logs, func(row repo.ListLogSuratKeputusanByIDRow) logSuratKeputusan {
		return logSuratKeputusan{
			Actor:     row.Actor.String,
			Log:       row.Log.String,
			Timestamp: row.WaktuTindakan.Time,
		}
	})

	return &suratKeputusan{
		IDSK:              id,
		KategoriSK:        data.KategoriSk.String,
		NoSK:              data.NoSk.String,
		TanggalSK:         db.Date(data.TanggalSk.Time),
		StatusSK:          statusSKText(data.StatusSk.Int16),
		UnitKerja:         s.getUnorLengkap(listUnor),
		NamaPemilik:       data.NamaPemilikSk.String,
		NamaPenandaTangan: data.NamaPenandatangan.String,
		Logs:              &logSk,
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
		res, err = s.repo.GetBerkasSuratKeputusanSignedByNIPAndID(ctx, repo.GetBerkasSuratKeputusanSignedByNIPAndIDParams{
			Nip: nip,
			ID:  id,
		})
	} else {
		res, err = s.repo.GetBerkasSuratKeputusanByNIPAndID(ctx, repo.GetBerkasSuratKeputusanByNIPAndIDParams{
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

type listAdminParams struct {
	Limit          uint
	Offset         uint
	UnitKerjaID    string
	NamaPemilik    string
	NIP            string
	GolonganID     int32
	JabatanID      string
	KategoriSK     string
	TanggalSKMulai db.Date
	TanggalSKAkhir db.Date
	StatusSK       *int32
}

func (s *service) listAdmin(ctx context.Context, arg listAdminParams) ([]suratKeputusan, uint, error) {
	statusSkParam := pgtype.Int4{Valid: false}
	if arg.StatusSK != nil {
		statusSkParam = pgtype.Int4{Valid: true, Int32: *arg.StatusSK}
	}

	data, err := s.repo.ListSuratKeputusan(ctx, repo.ListSuratKeputusanParams{
		Limit:          int32(arg.Limit),
		Offset:         int32(arg.Offset),
		UnitKerjaID:    pgtype.Text{String: arg.UnitKerjaID, Valid: arg.UnitKerjaID != ""},
		NamaPemilik:    pgtype.Text{String: arg.NamaPemilik, Valid: arg.NamaPemilik != ""},
		Nip:            pgtype.Text{String: arg.NIP, Valid: arg.NIP != ""},
		GolonganID:     pgtype.Int4{Int32: arg.GolonganID, Valid: arg.GolonganID != 0},
		JabatanID:      pgtype.Text{String: arg.JabatanID, Valid: arg.JabatanID != ""},
		KategoriSk:     pgtype.Text{String: arg.KategoriSK, Valid: arg.KategoriSK != ""},
		TanggalSkMulai: pgtype.Date{Time: time.Time(arg.TanggalSKMulai), Valid: !time.Time(arg.TanggalSKMulai).IsZero()},
		TanggalSkAkhir: pgtype.Date{Time: time.Time(arg.TanggalSKAkhir), Valid: !time.Time(arg.TanggalSKAkhir).IsZero()},
		StatusSk:       statusSkParam,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-listAdmin] repo ListSuratKeputusan: %w", err)
	}

	count, err := s.repo.CountSuratKeputusan(ctx, repo.CountSuratKeputusanParams{
		UnitKerjaID:    pgtype.Text{String: arg.UnitKerjaID, Valid: arg.UnitKerjaID != ""},
		NamaPemilik:    pgtype.Text{String: arg.NamaPemilik, Valid: arg.NamaPemilik != ""},
		Nip:            pgtype.Text{String: arg.NIP, Valid: arg.NIP != ""},
		GolonganID:     pgtype.Int4{Int32: arg.GolonganID, Valid: arg.GolonganID != 0},
		JabatanID:      pgtype.Text{String: arg.JabatanID, Valid: arg.JabatanID != ""},
		KategoriSk:     pgtype.Text{String: arg.KategoriSK, Valid: arg.KategoriSK != ""},
		TanggalSkMulai: pgtype.Date{Time: time.Time(arg.TanggalSKMulai), Valid: !time.Time(arg.TanggalSKMulai).IsZero()},
		TanggalSkAkhir: pgtype.Date{Time: time.Time(arg.TanggalSKAkhir), Valid: !time.Time(arg.TanggalSKAkhir).IsZero()},
		StatusSk:       statusSkParam,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-listAdmin] repo CountSuratKeputusan: %w", err)
	}

	uniqUnorIDs := typeutil.UniqMap(data, func(row repo.ListSuratKeputusanRow, _ int) string {
		return row.UnorID.String
	})

	listUnorLengkap, err := s.repo.ListUnitKerjaLengkapByIDs(ctx, uniqUnorIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-listAdmin] repo ListUnitKerjaLengkapByIDs: %w", err)
	}

	unorLengkapByID := typeutil.SliceToMap(listUnorLengkap, func(unorLengkap repo.ListUnitKerjaLengkapByIDsRow) (string, string) {
		return unorLengkap.ID, unorLengkap.NamaUnorLengkap
	})

	result := typeutil.Map(data, func(row repo.ListSuratKeputusanRow) suratKeputusan {
		return suratKeputusan{
			IDSK:        row.FileID,
			NamaPemilik: row.NamaPemilikSk.String,
			KategoriSK:  row.KategoriSk.String,
			NoSK:        row.NoSk.String,
			TanggalSK:   db.Date(row.TanggalSk.Time),
			UnitKerja:   unorLengkapByID[row.UnorID.String],
			StatusSK:    statusSKText(row.StatusSk.Int16),
		}
	})

	return result, uint(count), nil
}

func (s *service) getAdmin(ctx context.Context, id string) (*suratKeputusan, error) {
	data, err := s.repo.GetSuratKeputusanByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[suratkeputusan-getAdmin] repo GetSuratKeputusanByID: %w", err)
	}

	listUnor, err := s.repo.ListUnitKerjaHierarchyByNIP(ctx, data.NipSk.String)
	if err != nil {
		return nil, fmt.Errorf("[suratkeputusan-getAdmin] repo ListUnitKerjaHierarchyByNIP: %w", err)
	}

	return &suratKeputusan{
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

func (s *service) getBerkasAdmin(ctx context.Context, id string, signed bool) (string, []byte, error) {
	var res pgtype.Text
	var err error
	if signed {
		res, err = s.repo.GetBerkasSuratKeputusanSignedByID(ctx, id)
	} else {
		res, err = s.repo.GetBerkasSuratKeputusanByID(ctx, id)
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", nil, fmt.Errorf("[suratkeputusan-getBerkasAdmin] repo get berkas: %w", err)
	}

	if len(res.String) == 0 {
		return "", nil, nil
	}

	return api.GetMimeTypeAndDecodedData(res.String)
}
