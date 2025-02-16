# Avito Shop Service

Этот сервис позволяет сотрудникам Авито обмениваться монетами и покупать мерч.

## API Эндпоинты

### 1. **Аутентификация**
`POST /api/auth`
Авторизация и получение JWT-токена. При первой авторизации пользователь создаётся автоматически.

### 2. **Информация о пользователе**
`GET /api/info`
Возвращает баланс, инвентарь и историю переводов.

### 3. **Перевод монет**
`POST /api/sendCoin`
Позволяет отправить монеты другому пользователю.

### 4. **Покупка мерча**
`GET /api/buy/{item}`
Позволяет приобрести товар за монеты.

**Примечание:** Все защищённые эндпоинты требуют JWT-токен в заголовке `Authorization: Bearer <token>`.

---

## Запуск проекта

Требуется **Docker**, **Task** и **golang-migrate**.

1. **Установите `Task`:**
   ```sh
   go install github.com/go-task/task/v3/cmd/task@latest
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

2. **Установите `golang-migrate`:**
   ```sh
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   export PATH=$PATH:$(go env GOPATH)/bin
   ```

3. **Запустите сервис:**
   ```sh
   task run
   ```
Теперь сервис доступен по адресу [`http://localhost:8080`](http://localhost:8080)

4. **Остановите сервис:**
   ```sh
   task stop
   ```

---

## Тестирование

- **Юнит-тесты:** `task unit-test`
- **Покрытие:** `task coverage-test`
- **E2E-тесты:** `task e2e-test`

---

## Работа с базой

- **Применить миграции:** `task migrate-up`
- **Откатить миграцию:** `task migrate-down`
- **Создать новую миграцию:** `task migrate-create name=<название>`

---
