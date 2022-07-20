package memdb

import (
	"bytes"
	"io/fs"
	"time"
)

type dataFile struct {
	dataBlock dataBlock
	reader    *bytes.Reader
}

func newDataFile(dataBlock dataBlock) *dataFile {
	return &dataFile{
		dataBlock: dataBlock,
		reader:    bytes.NewReader(dataBlock.data),
	}
}

// fs.File

func (df *dataFile) Read(data []byte) (int, error) { return df.reader.Read(data) }
func (df *dataFile) Close() error                  { return nil }
func (df *dataFile) Stat() (fs.FileInfo, error)    { return df, nil }
func (df *dataFile) Seek(offset int64, whence int) (int64, error) {
	return df.reader.Seek(offset, whence)
}

// fs.FileInfo

func (df *dataFile) Name() string       { return df.dataBlock.fileName } // base name of the file
func (df *dataFile) Size() int64        { return df.dataBlock.length }   // length in bytes for regular files; system-dependent for others
func (df *dataFile) Mode() fs.FileMode  { return 0777 }                  // file mode bits
func (df *dataFile) ModTime() time.Time { return df.dataBlock.modtime }  // modification time
func (df *dataFile) IsDir() bool        { return false }                 // abbreviation for Mode().IsDir()
func (df *dataFile) Sys() interface{}   { return nil }                   // underlying data source (can return nil)
