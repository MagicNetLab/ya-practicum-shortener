package storage

type LinkStore interface {
	PutLink(link string, short string) error

	PutBatchLinksArray(StoreBatchLicksArray map[string]string) error

	GetLink(short string) (string, error)

	HasShort(short string) (bool, error)

	GetShort(link string) (string, error)

	Init() error
}
