package local

// Store локальное хранилище ссылок в памяти приложения
var Store = store{store: make(map[string]linkEntity, 2)}
