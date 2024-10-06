package local

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore_PutLink(t *testing.T) {
	t.Run("тест сохранения ссылки в локальном хранилище", func(t *testing.T) {
		putLink := linkEntity{shortLink: "dhfsjsj", originalURL: "http://rambler.ru", userID: 1, isDeleted: false}

		Store.store = make(map[string]linkEntity)
		assert.Equal(t, 0, len(Store.store))

		err := Store.PutLink(putLink.originalURL, putLink.shortLink, putLink.userID)
		assert.NoError(t, err)

		assert.Equal(t, 1, len(Store.store))

		dataInStore, ok := Store.store[putLink.shortLink]
		assert.True(t, ok)

		assert.Equal(t, putLink, dataInStore)
	})
}

func TestStore_PutBatchLinksArray(t *testing.T) {
	t.Run("тест пакетного сохранения ссылок в локальном хранилище", func(t *testing.T) {
		Store.store = make(map[string]linkEntity)
		assert.Equal(t, 0, len(Store.store))

		putData := map[string]string{"dshdgj": "http://rambler.ru", "xnmbsd": "http://yandex.ru"}

		err := Store.PutBatchLinksArray(putData, 1)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(Store.store))

		storeData := make(map[string]linkEntity)
		storeData["dshdgj"] = linkEntity{shortLink: "dshdgj", originalURL: "http://rambler.ru", userID: 1, isDeleted: false}
		storeData["xnmbsd"] = linkEntity{shortLink: "xnmbsd", originalURL: "http://yandex.ru", userID: 1, isDeleted: false}

		assert.Equal(t, storeData, Store.store)
	})
}

func TestStore_GetLink(t *testing.T) {
	t.Run("тест получения ссылки из store", func(t *testing.T) {
		Store.store = map[string]linkEntity{"dshdgj": {shortLink: "dshdgj", originalURL: "http://rambler.ru", userID: 1, isDeleted: true}}

		link, isDeleted, err := Store.GetLink("xnmbsd")
		assert.Error(t, err)
		assert.Empty(t, link)
		assert.False(t, isDeleted)

		link, isDeleted, err = Store.GetLink("dshdgj")
		assert.NoError(t, err)
		assert.Equal(t, "http://rambler.ru", link)
		assert.True(t, isDeleted)
	})
}

func TestStore_HasShort(t *testing.T) {
	t.Run("тест проверки наличия короткой ссылки в хранилище", func(t *testing.T) {
		Store.store = map[string]linkEntity{"dshdgj": {shortLink: "dshdgj", originalURL: "http://rambler.ru", userID: 1, isDeleted: false}}

		exists, err := Store.HasShort("jsdhjkj")
		assert.NoError(t, err)
		assert.False(t, exists)

		exists, err = Store.HasShort("dshdgj")
		assert.NoError(t, err)
		assert.True(t, exists)
	})
}

func TestStore_GetShort(t *testing.T) {
	t.Run("тест получения короткой ссылки из хранилища", func(t *testing.T) {
		Store.store = map[string]linkEntity{"dshdgj": {shortLink: "dshdgj", originalURL: "http://rambler.ru", userID: 1, isDeleted: false}}
		assert.Equal(t, 1, len(Store.store))

		short, err := Store.GetShort("http://yandex.ru")
		assert.Error(t, err)
		assert.Empty(t, short)

		short, err = Store.GetShort("http://rambler.ru")
		assert.NoError(t, err)
		assert.Equal(t, "dshdgj", short)
	})
}

func TestStore_GetUserLinks(t *testing.T) {
	t.Run("тест получения всех ссылок пользователя из хранилища", func(t *testing.T) {
		Store.store = map[string]linkEntity{
			"dshdgj":  {shortLink: "dshdgj", originalURL: "http://rambler.ru", userID: 1, isDeleted: false},
			"dfdjkjd": {shortLink: "dfdjkjd", originalURL: "http://yandex.ru", userID: 1, isDeleted: false},
			"mnbnvcx": {shortLink: "mnbnvcx", originalURL: "http://mail.ru", userID: 1, isDeleted: false},
			"hjhgsdf": {shortLink: "hjhgsdf", originalURL: "http://lenta.ru", userID: 7, isDeleted: false},
			"nbvcxz":  {shortLink: "nbvcxz", originalURL: "http://habr.com", userID: 7, isDeleted: false},
		}

		assert.Equal(t, 5, len(Store.store))

		res, err := Store.GetUserLinks(23)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(res))

		res, err = Store.GetUserLinks(1)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(res))

		res, err = Store.GetUserLinks(7)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(res))

	})
}

func TestStore_DeleteBatchLinksArray(t *testing.T) {
	t.Run("тест пакетного удаления ссылок пользователя", func(t *testing.T) {
		Store.store = map[string]linkEntity{
			"dshdgj":  {shortLink: "dshdgj", originalURL: "http://rambler.ru", userID: 1, isDeleted: false},
			"dfdjkjd": {shortLink: "dfdjkjd", originalURL: "http://yandex.ru", userID: 1, isDeleted: false},
			"mnbnvcx": {shortLink: "mnbnvcx", originalURL: "http://mail.ru", userID: 1, isDeleted: false},
			"hjhgsdf": {shortLink: "hjhgsdf", originalURL: "http://lenta.ru", userID: 7, isDeleted: false},
			"nbvcxz":  {shortLink: "nbvcxz", originalURL: "http://habr.com", userID: 7, isDeleted: false},
		}

		var shorts []string
		shorts = append(shorts, "mnbnvcx", "hjhgsdf", "nbvcxz")

		err := Store.DeleteBatchLinksArray(shorts, 7)
		assert.NoError(t, err)

		data, ok := Store.store["mnbnvcx"]
		assert.True(t, ok)
		assert.Equal(t, data.isDeleted, false)

		data, ok = Store.store["hjhgsdf"]
		assert.True(t, ok)
		assert.Equal(t, data.isDeleted, true)

		data, ok = Store.store["nbvcxz"]
		assert.True(t, ok)
		assert.Equal(t, data.isDeleted, true)

	})
}
