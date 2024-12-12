package link

import (
	"errors"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/repo"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/repo/memory"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/repo/postgres"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/shortgen"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
	"golang.org/x/net/context"
	"net/http"
	"sync"
	"time"
)

// Shorten Сокращение ссылки присланной пользователем
func Shorten(url string, userID int) (short string, httpResponseStatus int) {
	short = shortgen.GetShortLink(7)
	httpResponseStatus = http.StatusCreated
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := repo.PutLink(ctx, url, short, userID)
	if err != nil {
		httpResponseStatus = http.StatusInternalServerError
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("error storing short link", args)
		notUniqueError := errors.Is(err, postgres.ErrLinkUniqueConflict) || errors.Is(err, memory.ErrorLinkNotUnique)
		if notUniqueError {
			short, err = repo.GetShort(ctx, url)
			if err == nil {
				httpResponseStatus = http.StatusConflict
			}
		}
	}

	return short, httpResponseStatus
}

// BatchDeleteLinks пакетное удаление ссылок пользователя
func BatchDeleteLinks(ctx context.Context, shorts []string, userID int) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	inputCh := make(chan string, 5)
	go func() {
		for _, short := range shorts {
			select {
			case <-doneCh:
				return
			case inputCh <- short:
			}
		}
		close(inputCh)
	}()

	numWorkers := 2
	wg := sync.WaitGroup{}
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(dataChan <-chan string) {
			defer wg.Done()

			var batch []string
			for {
				select {
				case short, ok := <-dataChan:
					if !ok {
						deleteLinks(batch, userID)
						return
					}
					batch = append(batch, short)
					if len(batch) == 5 {
						deleteLinks(batch, userID)
						batch = nil
					}
				case <-doneCh:
					return
				default:
					time.Sleep(100 * time.Millisecond)
				}
			}
		}(inputCh)
	}

	wg.Wait()
}

func deleteLinks(shorts []string, userID int) {
	if len(shorts) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err := repo.DeleteBatchLinksArray(ctx, shorts, userID)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("error deleting short links", args)
	}
}
