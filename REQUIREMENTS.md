# Requirements

- Rate limiter имеет один экземпляр и может перенаправлять на несколько сервисов путем маршрутизации
- Маршрутизация настраивается в конфиге

## Пример маршрутизации

Client -> `/api/users` -> rate limiter -> `service.v1.api/users` Client -> `/api/orders/create` -> rate limiter -> `service.v2/api/orders/create`