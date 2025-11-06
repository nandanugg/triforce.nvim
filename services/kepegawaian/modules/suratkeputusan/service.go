package suratkeputusan

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/bsre"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/pdfcpu"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	WithTx(tx pgx.Tx) *repo.Queries
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
	ListKoreksiSuratKeputusanByPNSID(ctx context.Context, arg repo.ListKoreksiSuratKeputusanByPNSIDParams) ([]repo.ListKoreksiSuratKeputusanByPNSIDRow, error)
	CountKoreksiSuratKeputusanByPNSID(ctx context.Context, arg repo.CountKoreksiSuratKeputusanByPNSIDParams) (int64, error)
	GetPegawaiPNSIDByNIP(ctx context.Context, nip string) (string, error)
	ListAntreanKoreksiSuratKeputusanByNIP(ctx context.Context, arg repo.ListAntreanKoreksiSuratKeputusanByNIPParams) ([]repo.ListAntreanKoreksiSuratKeputusanByNIPRow, error)
	CountAntreanKoreksiSuratKeputusanByNIP(ctx context.Context, nipKorektor string) (int64, error)
	ListKorektorSuratKeputusanByID(ctx context.Context, id string) ([]repo.ListKorektorSuratKeputusanByIDRow, error)
	UpdateKorektorSuratKeputusanByID(ctx context.Context, arg repo.UpdateKorektorSuratKeputusanByIDParams) error
	UpdateStatusSuratKeputusanByID(ctx context.Context, arg repo.UpdateStatusSuratKeputusanByIDParams) error
	InsertRiwayatSuratKeputusan(ctx context.Context, arg repo.InsertRiwayatSuratKeputusanParams) error
	ListTandaTanganSuratKeputusanByPNSID(ctx context.Context, arg repo.ListTandaTanganSuratKeputusanByPNSIDParams) ([]repo.ListTandaTanganSuratKeputusanByPNSIDRow, error)
	CountTandaTanganSuratKeputusanByPNSID(ctx context.Context, arg repo.CountTandaTanganSuratKeputusanByPNSIDParams) (int64, error)
	ListTandaTanganSuratKeputusanAntreanByPNSID(ctx context.Context, arg repo.ListTandaTanganSuratKeputusanAntreanByPNSIDParams) ([]repo.ListTandaTanganSuratKeputusanAntreanByPNSIDRow, error)
	CountTandaTanganSuratKeputusanAntreanByPNSID(ctx context.Context, pnsID string) (int64, error)
	GetPegawaiTTDByNIP(ctx context.Context, nip string) (string, error)
	UpdateBerkasSuratKeputusanSignedByID(ctx context.Context, arg repo.UpdateBerkasSuratKeputusanSignedByIDParams) error
	InsertLogRequestSuratKeputusan(ctx context.Context, arg repo.InsertLogRequestSuratKeputusanParams) error
	GetPegawaiNIKByNIP(ctx context.Context, nip string) (string, error)
}

type service struct {
	repo repository
	db   *pgxpool.Pool
	bsre bsre.Client
}

func newService(r repository, db *pgxpool.Pool, bsre bsre.Client) *service {
	return &service{repo: r, db: db, bsre: bsre}
}

type listParams struct {
	limit        uint
	offset       uint
	nip          string
	noSK         string
	listStatusSK []int32
	kategoriSK   string
}

func (s *service) list(ctx context.Context, arg listParams) ([]suratKeputusan, uint, error) {
	data, err := s.repo.ListSuratKeputusanByNIP(ctx, repo.ListSuratKeputusanByNIPParams{
		Limit:        int32(arg.limit),
		Offset:       int32(arg.offset),
		Nip:          arg.nip,
		NoSk:         pgtype.Text{Valid: arg.noSK != "", String: arg.noSK},
		ListStatusSk: arg.listStatusSK,
		KategoriSk:   pgtype.Text{Valid: arg.kategoriSK != "", String: arg.kategoriSK},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-list] repo ListSuratKeputusanByNIP: %w", err)
	}

	count, err := s.repo.CountSuratKeputusanByNIP(ctx, repo.CountSuratKeputusanByNIPParams{
		NoSk:         pgtype.Text{Valid: arg.noSK != "", String: arg.noSK},
		ListStatusSk: arg.listStatusSK,
		KategoriSk:   pgtype.Text{Valid: arg.kategoriSK != "", String: arg.kategoriSK},
		Nip:          arg.nip,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-list] repo CountSuratKeputusanByNIP: %w", err)
	}

	listUnor, err := s.repo.ListUnitKerjaHierarchyByNIP(ctx, arg.nip)
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
	ListStatusSK   []int32
}

func (s *service) listAdmin(ctx context.Context, arg listAdminParams) ([]suratKeputusan, uint, error) {
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
		ListStatusSk:   arg.ListStatusSK,
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
		ListStatusSk:   arg.ListStatusSK,
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
			NIPPemilik:  row.NipPemilikSk.String,
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

	listUnor, err := s.repo.ListUnitKerjaHierarchyByNIP(ctx, data.NipPemilikSk.String)
	if err != nil {
		return nil, fmt.Errorf("[suratkeputusan-getAdmin] repo ListUnitKerjaHierarchyByNIP: %w", err)
	}

	logs, err := s.repo.ListLogSuratKeputusanByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("[suratkeputusan-getAdmin] repo ListLogSuratKeputusanByID: %w", err)
	}

	logSk := typeutil.Map(logs, func(row repo.ListLogSuratKeputusanByIDRow) logSuratKeputusan {
		return logSuratKeputusan{
			Actor:     row.Actor.String,
			Log:       row.Log.String,
			Timestamp: row.WaktuTindakan.Time,
		}
	})

	return &suratKeputusan{
		IDSK:                 id,
		KategoriSK:           data.KategoriSk.String,
		NoSK:                 data.NoSk.String,
		TanggalSK:            db.Date(data.TanggalSk.Time),
		StatusSK:             statusSKText(data.StatusSk.Int16),
		UnitKerja:            s.getUnorLengkap(listUnor),
		NamaPemilik:          data.NamaPemilikSk.String,
		NIPPemilik:           data.NipPemilikSk.String,
		NamaPenandaTangan:    data.NamaPenandatangan.String,
		JabatanPenandaTangan: data.JabatanPenandatangan.String,
		Logs:                 &logSk,
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

type listKoreksiParams struct {
	limit       uint
	offset      uint
	unitKerjaID string
	namaPemilik string
	nipPemilik  string
	nip         string
	golonganID  int32
	jabatanID   string
	kategoriSK  string
	noSK        string
	status      string
}

func (s *service) listKoreksi(ctx context.Context, arg listKoreksiParams) ([]koreksiSuratKeputusan, uint, error) {
	pnsID, err := s.repo.GetPegawaiPNSIDByNIP(ctx, arg.nip)
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-listKoreksi] repo GetPegawaiPNSIDByNIP: %w", err)
	}
	statusKoreksi := statusKoreksiRequest(arg.status).value()
	data, err := s.repo.ListKoreksiSuratKeputusanByPNSID(ctx, repo.ListKoreksiSuratKeputusanByPNSIDParams{
		Limit:         int32(arg.limit),
		Offset:        int32(arg.offset),
		UnitKerjaID:   pgtype.Text{String: arg.unitKerjaID, Valid: arg.unitKerjaID != ""},
		NamaPemilik:   pgtype.Text{String: arg.namaPemilik, Valid: arg.namaPemilik != ""},
		NipPemilik:    pgtype.Text{String: arg.nipPemilik, Valid: arg.nipPemilik != ""},
		GolonganID:    pgtype.Int4{Int32: arg.golonganID, Valid: arg.golonganID != 0},
		JabatanID:     pgtype.Text{String: arg.jabatanID, Valid: arg.jabatanID != ""},
		KategoriSk:    pgtype.Text{String: arg.kategoriSK, Valid: arg.kategoriSK != ""},
		NoSk:          pgtype.Text{String: arg.noSK, Valid: arg.noSK != ""},
		StatusKoreksi: statusKoreksi,
		PnsID:         pnsID,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-listKoreksi] repo ListKoreksiSuratKeputusanByPNSID: %w", err)
	}

	count, err := s.repo.CountKoreksiSuratKeputusanByPNSID(ctx, repo.CountKoreksiSuratKeputusanByPNSIDParams{
		UnitKerjaID:   pgtype.Text{String: arg.unitKerjaID, Valid: arg.unitKerjaID != ""},
		NamaPemilik:   pgtype.Text{String: arg.namaPemilik, Valid: arg.namaPemilik != ""},
		NipPemilik:    pgtype.Text{String: arg.nipPemilik, Valid: arg.nipPemilik != ""},
		GolonganID:    pgtype.Int4{Int32: arg.golonganID, Valid: arg.golonganID != 0},
		JabatanID:     pgtype.Text{String: arg.jabatanID, Valid: arg.jabatanID != ""},
		KategoriSk:    pgtype.Text{String: arg.kategoriSK, Valid: arg.kategoriSK != ""},
		NoSk:          pgtype.Text{String: arg.noSK, Valid: arg.noSK != ""},
		StatusKoreksi: statusKoreksi,
		PnsID:         pnsID,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-listKoreksi] repo CountKoreksiSuratKeputusanByPNSID: %w", err)
	}

	uniqUnorIDs := typeutil.UniqMap(data, func(row repo.ListKoreksiSuratKeputusanByPNSIDRow, _ int) string {
		return row.UnorID.String
	})

	listUnorLengkap, err := s.repo.ListUnitKerjaLengkapByIDs(ctx, uniqUnorIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-listKoreksi] repo ListUnitKerjaLengkapByIDs: %w", err)
	}

	unorLengkapByID := typeutil.SliceToMap(listUnorLengkap, func(unorLengkap repo.ListUnitKerjaLengkapByIDsRow) (string, string) {
		return unorLengkap.ID, unorLengkap.NamaUnorLengkap
	})

	result := typeutil.Map(data, func(row repo.ListKoreksiSuratKeputusanByPNSIDRow) koreksiSuratKeputusan {
		return koreksiSuratKeputusan{
			IDSK:        row.FileID,
			NamaPemilik: row.NamaPemilikSk.String,
			NIPPemilik:  row.NipPemilikSk.String,
			KategoriSK:  row.KategoriSk.String,
			NoSK:        row.NoSk.String,
			TanggalSK:   db.Date(row.TanggalSk.Time),
			UnitKerja:   unorLengkapByID[row.UnorID.String],
		}
	})
	return result, uint(count), nil
}

type listKoreksiAntreanParams struct {
	limit  uint
	offset uint
	nip    string
}

func (s *service) listKoreksiAntrean(ctx context.Context, arg listKoreksiAntreanParams) ([]antreanSK, uint, error) {
	data, err := s.repo.ListAntreanKoreksiSuratKeputusanByNIP(ctx, repo.ListAntreanKoreksiSuratKeputusanByNIPParams{
		Limit:       int32(arg.limit),
		Offset:      int32(arg.offset),
		NipKorektor: arg.nip,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-listKoreksiAntrean] repo ListKoreksiSuratKeputusanAntrean: %w", err)
	}

	count, err := s.repo.CountAntreanKoreksiSuratKeputusanByNIP(ctx, arg.nip)
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-listKoreksiAntrean] repo CountAntreanKoreksiSuratKeputusan: %w", err)
	}

	result := typeutil.Map(data, func(row repo.ListAntreanKoreksiSuratKeputusanByNIPRow) antreanSK {
		return antreanSK{
			KategoriSK: row.Kategori.String,
			Jumlah:     row.Jumlah,
		}
	})
	return result, uint(count), nil
}

func (s *service) getDetailSuratKeputusan(ctx context.Context, id, nip string) (*koreksiSuratKeputusan, error) {
	sk, err := s.repo.GetSuratKeputusanByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[suratkeputusan-getDetailSuratKeputusan] repo GetSuratKeputusanByID: %w", err)
	}

	listUnor, err := s.repo.ListUnitKerjaHierarchyByNIP(ctx, sk.NipPemilikSk.String)
	if err != nil {
		return nil, fmt.Errorf("[suratkeputusan-getDetailSuratKeputusan] repo ListUnitKerjaHierarchyByNIP: %w", err)
	}

	korektor, err := s.repo.ListKorektorSuratKeputusanByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("[suratkeputusan-getDetailSuratKeputusan] repo ListKorektorSuratKeputusanByID: %w", err)
	}

	aksi := s.getAksiSuratKeputusan(sk, korektor, nip)

	korektorKe := 0

	return &koreksiSuratKeputusan{
		IDSK:        id,
		KategoriSK:  sk.KategoriSk.String,
		NoSK:        sk.NoSk.String,
		TanggalSK:   db.Date(sk.TanggalSk.Time),
		UnitKerja:   s.getUnorLengkap(listUnor),
		NamaPemilik: sk.NamaPemilikSk.String,
		NIPPemilik:  sk.NipPemilikSk.String,
		ListKorektor: typeutil.Map(korektor, func(row repo.ListKorektorSuratKeputusanByIDRow) korektorSuratKeputusan {
			if nip == row.NipKorektor.String {
				korektorKe = int(row.KorektorKe.Int16)
			}
			return korektorSuratKeputusan{
				Nama:           row.NamaKorektor.String,
				NIP:            row.NipKorektor.String,
				GelarDepan:     row.GelarDepanKorektor.String,
				GelarBelakang:  row.GelarBelakangKorektor.String,
				StatusKoreksi:  statusKorektorSK(row.StatusKoreksi.Int16).String(),
				CatatanKoreksi: row.CatatanKoreksi.String,
				KorektorKe:     row.KorektorKe.Int16,
			}
		}),
		Aksi:       &aksi,
		StatusSK:   statusSKText(sk.StatusSk.Int16),
		KorektorKe: korektorKe,
	}, nil
}

func (s *service) getAksiSuratKeputusan(sk repo.GetSuratKeputusanByIDRow, korektor []repo.ListKorektorSuratKeputusanByIDRow, nip string) string {
	aksi := ""

	if s.cekTtd(sk) {
		aksi = "tandatangan"
	}

	if s.cekKoreksi(korektor, nip) {
		aksi = "koreksi"
	}

	return aksi
}

func (s *service) cekTtd(sk repo.GetSuratKeputusanByIDRow) bool {
	statusKoreksiSK := statusKoreksiSK(sk.StatusKoreksi.Int16)
	statusTtdSK := statusTtd(sk.StatusTtd.Int16)

	return statusTtdSK.belumTtd() && statusKoreksiSK.sudahDikoreksi()
}

func (s *service) cekKoreksi(korektor []repo.ListKorektorSuratKeputusanByIDRow, nip string) bool {
	var currentKorektor *repo.ListKorektorSuratKeputusanByIDRow
	for i := range korektor {
		if korektor[i].NipKorektor.String == nip {
			currentKorektor = &korektor[i]
			break
		}
	}

	if currentKorektor == nil {
		return false
	}

	korektorKe := currentKorektor.KorektorKe.Int16
	statusKoreksi := statusKorektorSK(currentKorektor.StatusKoreksi.Int16)

	if statusKoreksi.sudahDikoreksi() || statusKoreksi.dikembalikan() {
		return false
	}

	if korektorKe == 1 && statusKoreksi.belumDikoreksi() {
		return true
	}

	if korektorKe > 1 {
		for i := range korektor {
			if korektor[i].KorektorKe.Int16 < korektorKe {
				statusKoreksi := statusKorektorSK(korektor[i].StatusKoreksi.Int16)
				if statusKoreksi.belumDikoreksi() || statusKoreksi.dikembalikan() {
					return false
				}
			}
		}
	}

	return true
}

type koreksiSuratKeputusanParams struct {
	id             string
	statusKoreksi  string
	catatanKoreksi string
	nip            string
}

func (s *service) koreksiSuratKeputusan(ctx context.Context, arg koreksiSuratKeputusanParams) (bool, error) {
	pnsID, err := s.repo.GetPegawaiPNSIDByNIP(ctx, arg.nip)
	if err != nil {
		return false, fmt.Errorf("[suratkeputusan-koreksiSuratKeputusan] repo GetPegawaiPNSIDByNIP: %w", err)
	}
	sk, err := s.repo.GetSuratKeputusanByID(ctx, arg.id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("[suratkeputusan-koreksiSuratKeputusan] repo GetSuratKeputusanByID: %w", err)
	}

	korektor, err := s.repo.ListKorektorSuratKeputusanByID(ctx, arg.id)
	if err != nil {
		return false, fmt.Errorf("[suratkeputusan-koreksiSuratKeputusan] repo ListKorektorSuratKeputusanByID: %w", err)
	}

	aksi := s.getAksiSuratKeputusan(sk, korektor, arg.nip)
	if aksi != "koreksi" {
		return false, nil
	}

	if arg.statusKoreksi == "dikembalikan" {
		return s.rejectKoreksiSuratKeputusan(ctx, arg.id, arg.catatanKoreksi, arg.nip, pnsID)
	}

	return s.approveKoreksiSuratKeputusan(ctx, arg.id, arg.catatanKoreksi, arg.nip, pnsID, korektor)
}

func (s *service) approveKoreksiSuratKeputusan(ctx context.Context, id, catatanKoreksi, nip, pnsID string, korektor []repo.ListKorektorSuratKeputusanByIDRow) (bool, error) {
	err := s.withTransaction(ctx, func(txRepo repository) error {
		tindakan := ""
		err := txRepo.UpdateKorektorSuratKeputusanByID(ctx, repo.UpdateKorektorSuratKeputusanByIDParams{
			ID:             id,
			PnsID:          pnsID,
			StatusKoreksi:  pgtype.Int2{Int16: int16(statusKorektorSK(statusKorektorSKSudahDikoreksi)), Valid: true},
			CatatanKoreksi: catatanKoreksi,
		})
		if err != nil {
			return fmt.Errorf("[approveKoreksiSuratKeputusan] UpdateKorektorSuratKeputusanByID: %w", err)
		}

		nextKorektor := s.cekKorektorSelanjutnya(korektor, nip)

		if nextKorektor == nil {
			err = txRepo.UpdateStatusSuratKeputusanByID(ctx, repo.UpdateStatusSuratKeputusanByIDParams{
				ID:            id,
				StatusSk:      pgtype.Int2{Int16: int16(statusSK(statusSKSudahDikoreksi)), Valid: true},
				StatusKoreksi: pgtype.Int2{Int16: int16(statusKoreksiSK(statusKoreksiSudahDikoreksi)), Valid: true},
			})
			if err != nil {
				return fmt.Errorf("[approveKoreksiSuratKeputusan] UpdateStatusSuratKeputusanByID: %w", err)
			}
			tindakan = string(diteruskanKePenandatangan)
		} else {
			err = txRepo.UpdateStatusSuratKeputusanByID(ctx, repo.UpdateStatusSuratKeputusanByIDParams{
				ID:       id,
				StatusSk: pgtype.Int2{Int16: int16(statusSK(statusSKSedangDikoreksi)), Valid: true},
			})
			if err != nil {
				return fmt.Errorf("[approveKoreksiSuratKeputusan] UpdateStatusSuratKeputusanByID: %w", err)
			}
			err = txRepo.UpdateKorektorSuratKeputusanByID(ctx, repo.UpdateKorektorSuratKeputusanByIDParams{
				ID:            id,
				PnsID:         nextKorektor.PegawaiKorektorID.String,
				StatusKoreksi: pgtype.Int2{Int16: int16(statusKorektorSK(statusKorektorSKBelumDikoreksi)), Valid: true},
			})
			if err != nil {
				return fmt.Errorf("[approveKoreksiSuratKeputusan] UpdateKorektorSuratKeputusanByID(next): %w", err)
			}

			tindakan = string(diteruskanKeKorektorSelanjutnya) + " " + strconv.Itoa(int(nextKorektor.KorektorKe.Int16))
		}

		err = txRepo.InsertRiwayatSuratKeputusan(ctx, repo.InsertRiwayatSuratKeputusanParams{
			FileID:          id,
			NipPemroses:     nip,
			Tindakan:        tindakan,
			CatatanTindakan: catatanKoreksi,
			AksesPengguna:   "web",
		})
		if err != nil {
			return fmt.Errorf("[approveKoreksiSuratKeputusan] InsertRiwayatSuratKeputusan: %w", err)
		}
		return nil
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *service) rejectKoreksiSuratKeputusan(ctx context.Context, id, catatanKoreksi, nip, pnsID string) (bool, error) {
	err := s.withTransaction(ctx, func(txRepo repository) error {
		err := txRepo.UpdateKorektorSuratKeputusanByID(ctx, repo.UpdateKorektorSuratKeputusanByIDParams{
			ID:             id,
			PnsID:          pnsID,
			StatusKoreksi:  pgtype.Int2{Int16: int16(statusKorektorSK(statusKorektorSKDikembalikan)), Valid: true},
			CatatanKoreksi: catatanKoreksi,
		})
		if err != nil {
			return fmt.Errorf("[rejectKoreksiSuratKeputusan] UpdateKorektorSuratKeputusanByID: %w", err)
		}

		err = txRepo.UpdateStatusSuratKeputusanByID(ctx, repo.UpdateStatusSuratKeputusanByIDParams{
			ID:            id,
			StatusKoreksi: pgtype.Int2{Int16: int16(statusKoreksiSK(statusKoreksiDikembalikan)), Valid: true},
			StatusTtd:     pgtype.Int2{Int16: int16(statusTtd(statusTtdDikembalikan)), Valid: true},
			StatusKembali: pgtype.Int2{Int16: 1, Valid: true},
			StatusSk:      pgtype.Int2{Int16: int16(statusSK(statusSKDikembalikan)), Valid: true},
			Catatan:       catatanKoreksi,
		})
		if err != nil {
			return fmt.Errorf("[rejectKoreksiSuratKeputusan] UpdateStatusSuratKeputusanByID: %w", err)
		}

		err = txRepo.InsertRiwayatSuratKeputusan(ctx, repo.InsertRiwayatSuratKeputusanParams{
			FileID:          id,
			NipPemroses:     nip,
			Tindakan:        string(suratKeputusanRiwayatMessage(dikembalikan)),
			CatatanTindakan: catatanKoreksi,
			AksesPengguna:   "web",
		})
		if err != nil {
			return fmt.Errorf("[rejectKoreksiSuratKeputusan] InsertRiwayatSuratKeputusan: %w", err)
		}
		return nil
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *service) cekKorektorSelanjutnya(korektor []repo.ListKorektorSuratKeputusanByIDRow, nip string) *repo.ListKorektorSuratKeputusanByIDRow {
	currentKorektorIndex := slices.IndexFunc(korektor, func(row repo.ListKorektorSuratKeputusanByIDRow) bool {
		return row.NipKorektor.String == nip
	})

	if currentKorektorIndex == len(korektor)-1 {
		return nil
	}

	return &korektor[currentKorektorIndex+1]
}

type listTandatanganParams struct {
	limit       uint
	offset      uint
	unitKerjaID string
	namaPemilik string
	nipPemilik  string
	golonganID  int32
	jabatanID   string
	kategoriSK  string
	noSK        string
	status      string
	nip         string
}

func (s *service) listTandatangan(ctx context.Context, arg listTandatanganParams) ([]koreksiSuratKeputusan, uint, error) {
	pnsID, err := s.repo.GetPegawaiPNSIDByNIP(ctx, arg.nip)
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-listTandatangan] repo GetPegawaiPNSIDByNIP: %w", err)
	}
	statusTandatangan := pgtype.Int4{Valid: false}
	if statusTandaTangan(arg.status).valid() {
		statusTandatangan = pgtype.Int4{Int32: statusTandaTangan(arg.status).value(), Valid: true}
	}
	data, err := s.repo.ListTandaTanganSuratKeputusanByPNSID(ctx, repo.ListTandaTanganSuratKeputusanByPNSIDParams{
		Limit:       int32(arg.limit),
		Offset:      int32(arg.offset),
		UnitKerjaID: pgtype.Text{String: arg.unitKerjaID, Valid: arg.unitKerjaID != ""},
		NamaPemilik: pgtype.Text{String: arg.namaPemilik, Valid: arg.namaPemilik != ""},
		NipPemilik:  pgtype.Text{String: arg.nipPemilik, Valid: arg.nipPemilik != ""},
		GolonganID:  pgtype.Int4{Int32: arg.golonganID, Valid: arg.golonganID != 0},
		JabatanID:   pgtype.Text{String: arg.jabatanID, Valid: arg.jabatanID != ""},
		KategoriSk:  pgtype.Text{String: arg.kategoriSK, Valid: arg.kategoriSK != ""},
		NoSk:        pgtype.Text{String: arg.noSK, Valid: arg.noSK != ""},
		StatusTtd:   statusTandatangan,
		PnsID:       pnsID,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-listTandatangan] repo ListTandaTanganSuratKeputusanByID: %w", err)
	}

	count, err := s.repo.CountTandaTanganSuratKeputusanByPNSID(ctx, repo.CountTandaTanganSuratKeputusanByPNSIDParams{
		UnitKerjaID: pgtype.Text{String: arg.unitKerjaID, Valid: arg.unitKerjaID != ""},
		NamaPemilik: pgtype.Text{String: arg.namaPemilik, Valid: arg.namaPemilik != ""},
		NipPemilik:  pgtype.Text{String: arg.nipPemilik, Valid: arg.nipPemilik != ""},
		GolonganID:  pgtype.Int4{Int32: arg.golonganID, Valid: arg.golonganID != 0},
		JabatanID:   pgtype.Text{String: arg.jabatanID, Valid: arg.jabatanID != ""},
		KategoriSk:  pgtype.Text{String: arg.kategoriSK, Valid: arg.kategoriSK != ""},
		NoSk:        pgtype.Text{String: arg.noSK, Valid: arg.noSK != ""},
		StatusTtd:   statusTandatangan,
		PnsID:       pnsID,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-listTandatangan] repo CountTandaTanganSuratKeputusanByID: %w", err)
	}

	uniqUnorIDs := typeutil.UniqMap(data, func(row repo.ListTandaTanganSuratKeputusanByPNSIDRow, _ int) string {
		return row.UnorID.String
	})
	listUnorLengkap, err := s.repo.ListUnitKerjaLengkapByIDs(ctx, uniqUnorIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-listTandatangan] repo ListUnitKerjaLengkapByIDs: %w", err)
	}

	unorLengkapByID := typeutil.SliceToMap(listUnorLengkap, func(unorLengkap repo.ListUnitKerjaLengkapByIDsRow) (string, string) {
		return unorLengkap.ID, unorLengkap.NamaUnorLengkap
	})

	result := typeutil.Map(data, func(row repo.ListTandaTanganSuratKeputusanByPNSIDRow) koreksiSuratKeputusan {
		return koreksiSuratKeputusan{
			IDSK:        row.FileID,
			NamaPemilik: row.NamaPemilikSk.String,
			NIPPemilik:  row.NipPemilikSk.String,
			KategoriSK:  row.KategoriSk.String,
			NoSK:        row.NoSk.String,
			TanggalSK:   db.Date(row.TanggalSk.Time),
			UnitKerja:   unorLengkapByID[row.UnorID.String],
		}
	})
	return result, uint(count), nil
}

func (s *service) listTandatanganAntrean(ctx context.Context, limit, offset uint, nip string) ([]antreanSK, uint, error) {
	pnsID, err := s.repo.GetPegawaiPNSIDByNIP(ctx, nip)
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-listTandatanganAntrean] repo GetPegawaiPNSIDByNIP: %w", err)
	}
	data, err := s.repo.ListTandaTanganSuratKeputusanAntreanByPNSID(ctx, repo.ListTandaTanganSuratKeputusanAntreanByPNSIDParams{
		Limit:  int32(limit),
		Offset: int32(offset),
		PnsID:  pnsID,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-listTandatanganAntrean] repo ListTandaTanganSuratKeputusanAntreanByID: %w", err)
	}

	count, err := s.repo.CountTandaTanganSuratKeputusanAntreanByPNSID(ctx, pnsID)
	if err != nil {
		return nil, 0, fmt.Errorf("[suratkeputusan-listTandatanganAntrean] repo CountTandaTanganSuratKeputusanAntreanByID: %w", err)
	}

	result := typeutil.Map(data, func(row repo.ListTandaTanganSuratKeputusanAntreanByPNSIDRow) antreanSK {
		return antreanSK{
			KategoriSK: row.Kategori.String,
			Jumlah:     row.Jumlah,
		}
	})
	return result, uint(count), nil
}

type tandatanganSKParams struct {
	id         string
	statusTtd  string
	nip        string
	catatanTtd string
	passphrase string
}

func (s *service) tandatanganSK(ctx context.Context, arg tandatanganSKParams) (string, error) {
	pnsID, err := s.repo.GetPegawaiPNSIDByNIP(ctx, arg.nip)
	if err != nil {
		return "", fmt.Errorf("[suratkeputusan-tandatanganSK] repo GetPegawaiPNSIDByNIP: %w", err)
	}

	sk, err := s.repo.GetSuratKeputusanByID(ctx, arg.id)
	if err != nil {
		return "", fmt.Errorf("[suratkeputusan-tandatanganSK] repo GetSuratKeputusanByID: %w", err)
	}

	if sk.TtdPegawaiID.String != pnsID {
		return errorMessageBukanPegawaiTTD.Error(), nil
	}

	if !s.cekTtd(sk) {
		return errorMessageStatusTtdInvalid.Error(), nil
	}

	statusTtd := statusTandatanganRequest(arg.statusTtd)
	if !statusTtd.tandaTangan() && !statusTtd.dikembalikan() {
		return errorMessageStatusTtdInvalid.Error(), nil
	}

	if statusTtd.dikembalikan() {
		return "", s.rejectTandatanganSK(ctx, arg.id, arg.catatanTtd, arg.nip)
	}

	return s.approveTandatanganSK(ctx, arg.id, arg.nip, arg.passphrase)
}

func (s *service) rejectTandatanganSK(ctx context.Context, id, catatanTtd, nip string) error {
	err := s.withTransaction(ctx, func(txRepo repository) error {
		err := txRepo.UpdateStatusSuratKeputusanByID(ctx, repo.UpdateStatusSuratKeputusanByIDParams{
			ID:            id,
			StatusTtd:     pgtype.Int2{Int16: int16(statusTtdDikembalikan), Valid: true},
			StatusSk:      pgtype.Int2{Int16: int16(statusSK(statusSKDikembalikan)), Valid: true},
			StatusKoreksi: pgtype.Int2{Int16: int16(statusKoreksiDikembalikan), Valid: true},
			Catatan:       catatanTtd,
		})
		if err != nil {
			return fmt.Errorf("[suratkeputusan-rejectTandatanganSK] UpdateStatusSuratKeputusanByID: %w", err)
		}

		err = txRepo.InsertRiwayatSuratKeputusan(ctx, repo.InsertRiwayatSuratKeputusanParams{
			FileID:          id,
			NipPemroses:     nip,
			Tindakan:        string(suratKeputusanRiwayatMessage(dikembalikanOlehPenandatangan)),
			CatatanTindakan: catatanTtd,
		})
		if err != nil {
			return fmt.Errorf("[suratkeputusan-rejectTandatanganSK] InsertRiwayatSuratKeputusan: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("[suratkeputusan-rejectTandatanganSK] %w", err)
	}
	return err
}

func (s *service) approveTandatanganSK(ctx context.Context, id, nip, passphrase string) (string, error) {
	message := ""
	err := s.withTransaction(ctx, func(txRepo repository) error {
		err := txRepo.UpdateStatusSuratKeputusanByID(ctx, repo.UpdateStatusSuratKeputusanByIDParams{
			ID:        id,
			StatusTtd: pgtype.Int2{Int16: int16(statusTtdSudahTtd), Valid: true},
			StatusSk:  pgtype.Int2{Int16: int16(statusSK(statusSKSudahTtd)), Valid: true},
		})
		if err != nil {
			return fmt.Errorf("[suratkeputusan-approveTandatanganSK] UpdateStatusSuratKeputusanByID: %w", err)
		}

		err = txRepo.InsertRiwayatSuratKeputusan(ctx, repo.InsertRiwayatSuratKeputusanParams{
			FileID:          id,
			NipPemroses:     nip,
			Tindakan:        string(suratKeputusanRiwayatMessage(ditandatangani)),
			CatatanTindakan: "",
			AksesPengguna:   "web",
		})
		if err != nil {
			return fmt.Errorf("[suratkeputusan-approveTandatanganSK] InsertRiwayatSuratKeputusan: %w", err)
		}

		message, err = s.signSK(ctx, txRepo, id, nip, passphrase)
		if message != "" {
			return errorMessage(message)
		}
		if err != nil {
			return fmt.Errorf("[suratkeputusan-approveTandatanganSK] %w", err)
		}

		return err
	})
	if err != nil {
		if validationErr, ok := err.(errorMessage); ok {
			return validationErr.Error(), nil
		}
		return "", fmt.Errorf("[suratkeputusan-approveTandatanganSK] %w", err)
	}
	return message, nil
}

func (s *service) signSK(
	ctx context.Context,
	txRepo repository,
	id, nip, passphrase string,
) (string, error) {
	sk, err := txRepo.GetSuratKeputusanByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("[suratkeputusan-signSK] repo GetSuratKeputusanByID: %w", err)
	}

	berkas, err := txRepo.GetBerkasSuratKeputusanByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("[suratkeputusan-signSK] repo GetBerkasSuratKeputusanByID: %w", err)
	}

	sigBase64, err := txRepo.GetPegawaiTTDByNIP(ctx, nip)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorMessageBelumAdaTTD.Error(), nil
		}
		return "", fmt.Errorf("[suratkeputusan-signSK] repo GetPegawaiTTDByNIP: %w", err)
	}

	posisiTtd := letakTTDsk(sk.LetakTtd.Int16)
	x, y := posisiTtd.koordinat()
	page := strconv.Itoa(int(sk.HalamanTtd.Int16))
	signed, err := pdfcpu.AddSignatureToPDF(
		berkas.String, sigBase64, x, y, 0.1,
		page,
	)
	if err != nil {
		return "", fmt.Errorf("[suratkeputusan-signSK] pdfcpu.AddSignatureToPDF: %w", err)
	}

	nik, err := txRepo.GetPegawaiNIKByNIP(ctx, nip)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errorMessageNIKNotFound.Error(), nil
		}
		return "", fmt.Errorf("[suratkeputusan-signSK] repo GetPegawaiNIKByNIP: %w", err)
	}

	signedBytes, statusCode, err := s.bsre.Sign(
		bsre.SignParams{
			NIK:        nik,
			Passphrase: passphrase,
			Tampilan:   bsre.InvisibleMode,
		},
		[]bsre.UploadFile{
			{
				Field:         "file",
				ContentBase64: signed,
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("[suratkeputusan-signSK] bsre.Send: %w", err)
	}

	if statusCode != http.StatusOK {
		decoded, _ := base64.StdEncoding.DecodeString(signedBytes)

		var (
			responseBody bsre.ErrorResponse
			message      string
		)

		if err := json.Unmarshal(decoded, &responseBody); err == nil {
			message = responseBody.StatusCode.Message()
		}

		logErr := s.repo.InsertLogRequestSuratKeputusan(
			ctx,
			repo.InsertLogRequestSuratKeputusanParams{
				FileID:     id,
				Nik:        nik,
				Keterangan: string(decoded),
				Status:     1,
				ProsesCron: true,
			},
		)
		if logErr != nil {
			return "", fmt.Errorf("[suratkeputusan-signSK] repo InsertLogRequestSuratKeputusan: %w", logErr)
		}

		return message, fmt.Errorf("[suratkeputusan-signSK] bsre.Send: %w", err)
	}

	err = txRepo.UpdateBerkasSuratKeputusanSignedByID(ctx, repo.UpdateBerkasSuratKeputusanSignedByIDParams{
		ID:             id,
		FileBase64Sign: signedBytes,
	})
	if err != nil {
		return "", fmt.Errorf("[suratkeputusan-signSK] repo UpdateBerkasSuratKeputusanSignedByID: %w", err)
	}

	err = txRepo.InsertLogRequestSuratKeputusan(ctx, repo.InsertLogRequestSuratKeputusanParams{
		FileID:     id,
		Nik:        nik,
		Keterangan: "Berhasil tandatangan.",
		Status:     2,
		ProsesCron: true,
	})
	if err != nil {
		return "", fmt.Errorf("[suratkeputusan-signSK] repo InsertLogRequestSuratKeputusan: %w", err)
	}

	return "", nil
}
