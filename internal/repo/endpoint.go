package repo

import (
	"context"
	"time"

	"github.com/ItsNotGoodName/smtpbridge/internal/database"
	. "github.com/ItsNotGoodName/smtpbridge/internal/jet/table"
	"github.com/ItsNotGoodName/smtpbridge/internal/models"
	. "github.com/go-jet/jet/v2/sqlite"
)

var endpointPJ ProjectionList = ProjectionList{
	Endpoints.ID.AS("endpoint.id"),
	Endpoints.Internal.AS("endpoint.internal"),
	Endpoints.InternalID.AS("endpoint.internal_id"),
	Endpoints.Name.AS("endpoint.name"),
	Endpoints.AttachmentDisable.AS("endpoint.attachment_disable"),
	Endpoints.TextDisable.AS("endpoint.text_disable"),
	Endpoints.TitleTemplate.AS("endpoint.title_template"),
	Endpoints.BodyTemplate.AS("endpoint.body_template"),
	Endpoints.Kind.AS("endpoint.kind"),
	Endpoints.Config.AS("endpoint.config"),
}

func EndpointCreate(ctx context.Context, db database.Querier, end models.Endpoint) (int64, error) {
	now := models.NewTime(time.Now())
	m := struct {
		models.Endpoint
		CreatedAt models.Time
		UpdatedAt models.Time
	}{
		Endpoint:  end,
		CreatedAt: now,
		UpdatedAt: now,
	}
	res, err := Endpoints.
		INSERT(
			Endpoints.Internal,
			Endpoints.InternalID,
			Endpoints.Name,
			Endpoints.AttachmentDisable,
			Endpoints.TextDisable,
			Endpoints.TitleTemplate,
			Endpoints.BodyTemplate,
			Endpoints.Kind,
			Endpoints.Config,
			Endpoints.UpdatedAt,
			Endpoints.CreatedAt,
		).
		MODEL(m).
		ExecContext(ctx, db)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func EndpointGet(ctx context.Context, db database.Querier, id int64) (models.Endpoint, error) {
	var endpoint models.Endpoint
	err := Endpoints.
		SELECT(endpointPJ).
		WHERE(Endpoints.ID.EQ(Int64(id))).
		QueryContext(ctx, db, &endpoint)
	return endpoint, err
}

func EndpointList(ctx context.Context, db database.Querier) ([]models.Endpoint, error) {
	var endpoints []models.Endpoint
	err := Endpoints.
		SELECT(endpointPJ).
		WHERE(RawBool("1=1")).
		QueryContext(ctx, db, &endpoints)
	return endpoints, err
}

func EndpointUpdate(ctx context.Context, db database.Querier, end models.Endpoint) error {
	m := struct {
		models.Endpoint
		UpdatedAt models.Time
	}{
		Endpoint:  end,
		UpdatedAt: models.NewTime(time.Now()),
	}
	_, err := Endpoints.
		UPDATE(
			Endpoints.Internal,
			Endpoints.InternalID,
			Endpoints.Name,
			Endpoints.AttachmentDisable,
			Endpoints.TextDisable,
			Endpoints.TitleTemplate,
			Endpoints.BodyTemplate,
			Endpoints.Kind,
			Endpoints.Config,
			Endpoints.UpdatedAt,
		).
		MODEL(m).
		WHERE(Endpoints.ID.EQ(Int64(end.ID))).
		ExecContext(ctx, db)
	return err
}

func EndpointDelete(ctx context.Context, db database.Querier, id int64) error {
	res, err := Endpoints.
		DELETE().
		WHERE(Endpoints.ID.EQ(Int64(id))).
		ExecContext(ctx, db)
	if err != nil {
		return err
	}
	return one(res)
}
