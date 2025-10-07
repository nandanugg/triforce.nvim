package jabatan

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
	ListRefJabatan(ctx context.Context, arg sqlc.ListRefJabatanParams) ([]sqlc.ListRefJabatanRow, error)
	ListRefJabatanWithKeyword(ctx context.Context, arg sqlc.ListRefJabatanWithKeywordParams) ([]sqlc.ListRefJabatanWithKeywordRow, error)
	CountRefJabatan(ctx context.Context, nama pgtype.Text) (int64, error)
	CountRefJabatanWithKeyword(ctx context.Context, keyword pgtype.Text) (int64, error)
	GetRefJabatan(ctx context.Context, id int32) (sqlc.GetRefJabatanRow, error)
	CreateRefJabatan(ctx context.Context, arg sqlc.CreateRefJabatanParams) (sqlc.CreateRefJabatanRow, error)
	UpdateRefJabatan(ctx context.Context, arg sqlc.UpdateRefJabatanParams) (sqlc.UpdateRefJabatanRow, error)
	DeleteRefJabatan(ctx context.Context, id int32) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) listJabatan(ctx context.Context, nama string, limit, offset uint) ([]jabatanPublic, int64, error) {
	pgNama := pgtype.Text{Valid: nama != "", String: nama}
	data, err := s.repo.ListRefJabatan(ctx, sqlc.ListRefJabatanParams{
		Nama:   pgNama,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRefJabatan(ctx, pgNama)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	result := typeutil.Map(data, func(row sqlc.ListRefJabatanRow) jabatanPublic {
		return jabatanPublic{
			ID:   row.KodeJabatan,
			Nama: row.NamaJabatan.String,
		}
	})

	return result, count, nil
}

func (s *service) listAdmin(ctx context.Context, keyword string, limit, offset uint) ([]jabatan, int64, error) {
	pgKeyword := s.stringToPgtypeText(keyword)
	data, err := s.repo.ListRefJabatanWithKeyword(ctx, sqlc.ListRefJabatanWithKeywordParams{
		Keyword: pgKeyword,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[listAdmin] error ListRefJabatan: %w", err)
	}

	count, err := s.repo.CountRefJabatanWithKeyword(ctx, pgKeyword)
	if err != nil {
		return nil, 0, fmt.Errorf("[listAdmin] error CountRefJabatan: %w", err)
	}

	result := typeutil.Map(data, func(row sqlc.ListRefJabatanWithKeywordRow) jabatan {
		return jabatan{
			Kode:      row.KodeJabatan,
			ID:        row.ID,
			Nama:      row.NamaJabatan.String,
			NamaFull:  row.NamaJabatanFull.String,
			Jenis:     typeutil.ValueOrNil(row.JenisJabatan.Int16, row.JenisJabatan.Valid),
			Kelas:     typeutil.ValueOrNil(row.Kelas.Int16, row.Kelas.Valid),
			Pensiun:   typeutil.ValueOrNil(row.Pensiun.Int16, row.Pensiun.Valid),
			KodeBkn:   row.KodeBkn.String,
			NamaBkn:   row.NamaJabatanBkn.String,
			Kategori:  row.KategoriJabatan.String,
			BknID:     row.BknID.String,
			Tunjangan: row.TunjanganJabatan.Int64,
		}
	})

	return result, count, nil
}

func (s *service) get(ctx context.Context, id int32) (*jabatan, error) {
	row, err := s.repo.GetRefJabatan(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[get] error GetRefJabatan: %w", err)
	}
	result := &jabatan{
		Kode:      row.KodeJabatan,
		ID:        row.ID,
		Nama:      row.NamaJabatan.String,
		NamaFull:  row.NamaJabatanFull.String,
		Jenis:     typeutil.ValueOrNil(row.JenisJabatan.Int16, row.JenisJabatan.Valid),
		Kelas:     typeutil.ValueOrNil(row.Kelas.Int16, row.Kelas.Valid),
		Pensiun:   typeutil.ValueOrNil(row.Pensiun.Int16, row.Pensiun.Valid),
		KodeBkn:   row.KodeBkn.String,
		NamaBkn:   row.NamaJabatanBkn.String,
		Kategori:  row.KategoriJabatan.String,
		BknID:     row.BknID.String,
		Tunjangan: row.TunjanganJabatan.Int64,
	}

	return result, nil
}

type createParams struct {
	kodeJabatan      string
	namaJabatan      string
	namaJabatanFull  string
	jenisJabatan     *int16
	kelas            *int16
	pensiun          *int16
	kodeBkn          string
	namaJabatanBkn   string
	kategoriJabatan  string
	bknID            string
	tunjanganJabatan int64
}

func (s *service) create(ctx context.Context, params createParams) (*jabatan, error) {
	row, err := s.repo.CreateRefJabatan(ctx, sqlc.CreateRefJabatanParams{
		KodeJabatan:      params.kodeJabatan,
		NamaJabatan:      s.stringToPgtypeText(params.namaJabatan),
		NamaJabatanFull:  s.stringToPgtypeText(params.namaJabatanFull),
		JenisJabatan:     s.intToPgtypeInt(params.jenisJabatan),
		Kelas:            s.intToPgtypeInt(params.kelas),
		Pensiun:          s.intToPgtypeInt(params.pensiun),
		KodeBkn:          s.stringToPgtypeText(params.kodeBkn),
		NamaJabatanBkn:   s.stringToPgtypeText(params.namaJabatanBkn),
		KategoriJabatan:  s.stringToPgtypeText(params.kategoriJabatan),
		BknID:            s.stringToPgtypeText(params.bknID),
		TunjanganJabatan: pgtype.Int8{Valid: true, Int64: params.tunjanganJabatan},
	})
	if err != nil {
		return nil, fmt.Errorf("[create] error CreateRefJabatan: %w", err)
	}

	result := &jabatan{
		Kode:      row.KodeJabatan,
		Nama:      row.NamaJabatan.String,
		NamaFull:  row.NamaJabatanFull.String,
		Jenis:     typeutil.ValueOrNil(row.JenisJabatan.Int16, row.JenisJabatan.Valid),
		Kelas:     typeutil.ValueOrNil(row.Kelas.Int16, row.Kelas.Valid),
		Pensiun:   typeutil.ValueOrNil(row.Pensiun.Int16, row.Pensiun.Valid),
		KodeBkn:   row.KodeBkn.String,
		NamaBkn:   row.NamaJabatanBkn.String,
		Kategori:  row.KategoriJabatan.String,
		BknID:     row.BknID.String,
		ID:        row.ID,
		Tunjangan: row.TunjanganJabatan.Int64,
	}

	return result, nil
}

type updateParams struct {
	kodeJabatan      string
	namaJabatan      string
	namaJabatanFull  string
	jenisJabatan     *int16
	kelas            *int16
	pensiun          *int16
	kodeBkn          string
	namaJabatanBkn   string
	kategoriJabatan  string
	bknID            string
	tunjanganJabatan int64
}

func (s *service) update(ctx context.Context, id int32, params updateParams) (*jabatan, error) {
	row, err := s.repo.UpdateRefJabatan(ctx, sqlc.UpdateRefJabatanParams{
		ID:               id,
		KodeJabatan:      params.kodeJabatan,
		NamaJabatan:      s.stringToPgtypeText(params.namaJabatan),
		NamaJabatanFull:  s.stringToPgtypeText(params.namaJabatanFull),
		JenisJabatan:     s.intToPgtypeInt(params.jenisJabatan),
		Kelas:            s.intToPgtypeInt(params.kelas),
		Pensiun:          s.intToPgtypeInt(params.pensiun),
		KodeBkn:          s.stringToPgtypeText(params.kodeBkn),
		NamaJabatanBkn:   s.stringToPgtypeText(params.namaJabatanBkn),
		KategoriJabatan:  s.stringToPgtypeText(params.kategoriJabatan),
		BknID:            s.stringToPgtypeText(params.bknID),
		TunjanganJabatan: pgtype.Int8{Valid: true, Int64: params.tunjanganJabatan},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[update] error UpdateRefJabatan: %w", err)
	}

	return &jabatan{
		Kode:      row.KodeJabatan,
		ID:        row.ID,
		Nama:      row.NamaJabatan.String,
		NamaFull:  row.NamaJabatanFull.String,
		Jenis:     typeutil.ValueOrNil(row.JenisJabatan.Int16, row.JenisJabatan.Valid),
		Kelas:     typeutil.ValueOrNil(row.Kelas.Int16, row.Kelas.Valid),
		Pensiun:   typeutil.ValueOrNil(row.Pensiun.Int16, row.Pensiun.Valid),
		KodeBkn:   row.KodeBkn.String,
		NamaBkn:   row.NamaJabatanBkn.String,
		Kategori:  row.KategoriJabatan.String,
		BknID:     row.BknID.String,
		Tunjangan: row.TunjanganJabatan.Int64,
	}, nil
}

func (s *service) delete(ctx context.Context, id int32) (bool, error) {
	rowsAffected, err := s.repo.DeleteRefJabatan(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("[delete] error DeleteRefJabatan: %w", err)
	}
	return rowsAffected > 0, nil
}

func (s *service) stringToPgtypeText(params string) pgtype.Text {
	return pgtype.Text{Valid: params != "", String: params}
}

func (s *service) intToPgtypeInt(params *int16) pgtype.Int2 {
	if params == nil {
		return pgtype.Int2{Valid: false, Int16: 0}
	}

	return pgtype.Int2{Valid: true, Int16: *params}
}
