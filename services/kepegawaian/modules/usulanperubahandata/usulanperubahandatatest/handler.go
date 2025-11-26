package usulanperubahandatatest

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/usulanperubahandata"
)

type handler struct {
	db   *pgxpool.Pool
	repo *sqlc.Queries
}

type createRequest struct {
	Action string          `json:"action"`
	DataID string          `json:"data_id"`
	Data   json.RawMessage `json:"data"`
}

func (h *handler) createThenApprove(svc usulanperubahandata.ServiceInterface, jenisData string) func(echo.Context) error {
	return func(c echo.Context) error {
		var req createRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		ctx := c.Request().Context()
		nip := api.CurrentUser(c).NIP
		bytes, err := svc.GeneratePerubahanData(ctx, nip, req.Action, req.DataID, req.Data)
		if err != nil {
			var multiErr *api.MultiError
			if errors.As(err, &multiErr) {
				return echo.NewHTTPError(http.StatusBadRequest, multiErr.Error())
			}

			slog.ErrorContext(ctx, "[Test] Error generate usulan perubahan data.", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		id, err := h.repo.CreateUsulanPerubahanData(ctx, sqlc.CreateUsulanPerubahanDataParams{
			Nip:           nip,
			JenisData:     jenisData,
			DataID:        pgtype.Text{String: req.DataID, Valid: req.DataID != ""},
			PerubahanData: bytes,
			Action:        req.Action,
		})
		if err != nil {
			slog.ErrorContext(ctx, "[Test] Error create usulan perubahan data.", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		tx, err := h.db.Begin(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "[Test] Error begin tx usulan perubahan data.", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		defer func() {
			_ = tx.Rollback(ctx)
		}()

		sqlcTx := h.repo.WithTx(tx)
		if err := svc.SyncPerubahanData(ctx, sqlcTx, nip, req.Action, req.DataID, bytes); err != nil {
			slog.ErrorContext(ctx, "[Test] Error sync usulan perubahan data.", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		if err := sqlcTx.UpdateStatusUsulanPerubahanData(ctx, sqlc.UpdateStatusUsulanPerubahanDataParams{
			ID:        id,
			JenisData: jenisData,
			Status:    "Disetujui",
		}); err != nil {
			slog.ErrorContext(ctx, "[Test] Error approve usulan perubahan data.", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		if c.QueryParam("rollback") == "true" {
			if err := tx.Rollback(ctx); err != nil {
				slog.ErrorContext(ctx, "[Test] Error rollback usulan perubahan data.", "error", err)
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
		} else {
			if err := tx.Commit(ctx); err != nil {
				slog.ErrorContext(ctx, "[Test] Error commit usulan perubahan data.", "error", err)
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
		}

		return c.NoContent(http.StatusNoContent)
	}
}
