package repo

import (
	"context"

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
