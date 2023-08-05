package db

import (
	"github.com/ItsNotGoodName/smtpbridge/internal/core"
	"github.com/ItsNotGoodName/smtpbridge/internal/endpoints"
)

func EndpointGet(cc core.Context, id int64) (endpoints.Endpoint, error) {
	var end endpoints.Endpoint
	err := cc.DB.NewSelect().Model(&end).Where("id = ?", id).Scan(cc)
	return end, err
}

func EndpointList(cc core.Context) ([]endpoints.Endpoint, error) {
	var ends []endpoints.Endpoint
	err := cc.DB.NewSelect().Model(&ends).Scan(cc)
	return ends, err
}

func EndpointListByRule(cc core.Context, id int64) ([]endpoints.Endpoint, error) {
	var ends []endpoints.Endpoint
	err := cc.DB.
		NewSelect().
		Model(&ends).
		Join("JOIN rules_to_endpoints ON endpoint.id = rules_to_endpoints.endpoint_id").
		Where("rules_to_endpoints.rule_id = ?", id).
		Scan(cc)
	return ends, err
}
