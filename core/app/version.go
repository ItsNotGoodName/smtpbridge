package app

import (
	"github.com/ItsNotGoodName/smtpbridge/core/dto"
	"github.com/ItsNotGoodName/smtpbridge/core/version"
)

func (a *App) Version() *dto.VersionResponse {
	return &dto.VersionResponse{
		Version: version.CurrentVersion.Version,
		Commit:  version.CurrentVersion.Commit,
		Date:    version.CurrentVersion.Date,
		BuiltBy: version.CurrentVersion.BuiltBy,
	}
}
