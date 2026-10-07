# ADR-0001: Bounded HTTP receipts

- Дата: 2026-10-07
- Статус: принято

## Контекст

`io.ReadAll(resp.Body)` без лимита позволяет SDK выделять память пропорционально внешнему потоку. Проверка длины в приложении происходит слишком поздно. Raw HTTP provenance требует точных complete bytes и безопасной метаинформации. Существующие конструкторы, map methods и экспортируемые RawResponse поля уже публичны.

## Варианты

1. Проверять после чтения в приложении: не ограничивает SDK allocation.
2. Удалить mutable поля и заменить весь raw API: ломает source compatibility.
3. Ограничить единый body reader и добавить неизменяемые accessors: сохраняет используемые публичные сигнатуры.

## Решение

Использовать вариант 3. Положительный лимит по умолчанию 16 MiB применяется ко всем body, `limit+1` выявляет overflow. Полный body сохраняется на raw error; oversized/incomplete body не становится receipt. Опции HTTP client и лимита задаются при создании; глобальный transport не меняется. Время receipt фиксируется сразу после complete read, до decode. Legacy Body — независимая копия. Header/selector allowlists не включают signature, ключи, cookies или redirect URLs.

## Последствия

Память ограничена configured body size с дополнительными bounded копиями/JSON decode. Ответы крупнее default требуют явно поднять лимит. Accessors нельзя изменить извне; legacy поля не считаются authoritative evidence. Map provider errors сохраняют типы; raw errors поддерживают errors.As через wrapper. SDK не добавляет retries. Historical source semantics проверяются отдельно и могут оставаться неподдержанными даже при working receipts.

## Проверка

- `http/receipt_test.go`, `services/market_receipt_test.go`, `client_http_options_test.go`.
- `docs/HTTP_RECEIPTS.md`, `services/testdata/contracts/`.
