package attachment

import (
	"context"
	"io/fs"
)

type DataService struct {
	repository RepositoryData
	uri        string
	url        func(att *Attachment) string
}

func localURL(repository RepositoryData, prefix string) func(att *Attachment) string {
	return func(att *Attachment) string {
		return prefix + repository.File(att)
	}
}

func remoteURL(repository RepositoryData) func(att *Attachment) string {
	return func(att *Attachment) string {
		return repository.URL(att)
	}
}

func NewDataService(repository RepositoryData, host, uri string) *DataService {
	ds := DataService{
		repository: repository,
		uri:        uri,
	}

	if repository.Remote() {
		ds.url = remoteURL(repository)
	} else {
		ds.url = localURL(repository, host+uri)
	}

	return &ds
}

func (ds *DataService) Remote() bool {
	return ds.repository.Remote()
}

func (ds *DataService) FS() fs.FS {
	return ds.repository.FS()
}

func (ds *DataService) URI() string {
	return ds.uri
}

func (ds *DataService) URL(att *Attachment) string {
	return ds.url(att)
}

func (ds *DataService) File(att *Attachment) string {
	return ds.repository.File(att)
}

func (ds *DataService) Create(ctx context.Context, att *Attachment) error {
	data, err := att.GetData()
	if err != nil {
		return err
	}

	return ds.repository.Create(ctx, att, data)
}

func (ds *DataService) Delete(ctx context.Context, att *Attachment) error {
	return ds.repository.Delete(ctx, att)
}

func (ds *DataService) Load(ctx context.Context, att *Attachment) error {
	data, err := ds.repository.Get(ctx, att)
	if err != nil {
		return err
	}

	return att.SetData(data)
}

func (ds *DataService) Size(ctx context.Context) (int64, error) {
	return ds.repository.Size(ctx)
}
