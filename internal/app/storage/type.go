package storage

type LinkStore interface {
	PutLink(link string, short string, userID int) error

	PutBatchLinksArray(StoreBatchLinksArray map[string]string, userID int) error

	GetLink(short string) (string, error)

	HasShort(short string) (bool, error)

	GetShort(link string) (string, error)

	GetUserLinks(userID int) (map[string]string, error)

	Init() error
}
