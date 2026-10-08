# CLAUDE.md — app/ (Go)

Service Go 1.26 theo **Clean Architecture**, tổ chức theo **feature** trong `internal/`.
Tài liệu gốc: `backend-clean-architecture-go.md`. Stack: `net/http` (router Go 1.22+), pgx/v5, sqlc, goose, kafka-go, log/slog.

## Lớp & quy tắc phụ thuộc
| Package | Chứa gì | Được import |
|---|---|---|
| `internal/<feature>/domain` | Entity (field unexported), value type, sentinel error | Chỉ stdlib + uuid |
| `internal/<feature>/app` | Use case + **interface port** (khai báo ở phía dùng) | `domain` |
| `internal/<feature>/transport/http` | Handler, request/response, map lỗi → HTTP | `app`, `domain`, `platform/httpserver` |
| `internal/<feature>/adapter/*` | postgres, stripe, kafka, local | `app`, `domain`, `platform` |
| `internal/platform/*` | config, database, httpserver, observability | **Không** import feature |
| `cmd/server/main.go` | Composition root: chọn adapter, ráp use case | Tất cả |

- Kiểm tra bằng **depguard** trong `.golangci.yml` (CI chạy `golangci-lint`). Đổi module path → sửa luôn các rule.
- `transport` và `adapter` không import lẫn nhau. Feature không import feature khác.

## Cấu trúc feature `order` (mẫu cho feature mới)
```
domain/     order.go money.go status.go events.go errors.go (+ _test)
app/        ports.go dto.go place_order.go cancel_order.go get_order.go (+ _test, fakes_test.go)
transport/http/  handler.go request.go response.go errors.go (+ handler_test.go)
adapter/postgres/  repository.go queries.sql db/ (sqlc, KHÔNG sửa tay) repository_test.go (tag integration)
adapter/stripe/  adapter/kafka/  adapter/local/ (in-memory + log)
```
`customer` mới có `domain`; các lớp khác là `doc.go` giữ chỗ.

## Lệnh
```bash
make run | test | test-integration | test-e2e | lint | sqlc
```
`test-integration`, `test-e2e` dùng testcontainers (cần Docker).

## Quy ước
- `context.Context` là tham số đầu tiên từ handler xuống adapter.
- Lỗi: domain trả sentinel (`errors.Is`), app bọc bằng `fmt.Errorf("...: %w")`. Map lỗi → status ở **một chỗ mỗi feature**: `transport/http/errors.go`.
- Interface nhỏ, đặt ở nơi dùng; "accept interfaces, return structs". Handler định nghĩa interface 1 method cho use case.
- Test use case bằng fake tự viết (`app/fakes_test.go`), không dùng mock framework.
- Thêm/đổi SQL: sửa `queries.sql` hoặc migration → `make sqlc` → commit cả thư mục `db/` (CI chạy `sqlc diff`).
- Migration: thêm file `migrations/0000N_*.sql` (goose, có `-- +goose Up/Down`), được nhúng vào binary.
- Tránh package `models`, `utils`, `common`.

## Hợp đồng với Kubernetes (đừng phá vỡ)
- **Port 8080**, **`/health/live`** (không kiểm tra phụ thuộc) và **`/health/ready`** (ping DB nếu có) — khớp `deploy/base/deployment.yaml` và `ingress.yaml`.
- **UID 10001, non-root, root filesystem read-only** (chỉ ghi được `/tmp`).
- **Graceful shutdown khi nhận SIGTERM.**
- **Migration tự chạy khi khởi động** nếu có `DATABASE_URL`.
- Log JSON ra stdout.

## Cấu hình (env)
| Biến | Nguồn | Bỏ trống thì |
|---|---|---|
| `DATABASE_URL` | Secret (Secrets Manager `my-app/<env>`) | Repository in-memory |
| `STRIPE_API_KEY` | Secret | Payment gateway local (luôn thành công) |
| `KAFKA_BROKERS` | ConfigMap | Publisher ghi event ra log |
| `KAFKA_TOPIC`, `LOG_LEVEL`, `APP_ENV` | ConfigMap | Mặc định `orders`, `info`, `local` |

Adapter `local` chỉ để chạy local / dev chưa có hạ tầng — **production phải cấu hình đủ 3 biến đầu**.
Hạ tầng RDS/MSK chưa có trong `infra/`; NetworkPolicy đã mở egress 5432 và 9092/9094/9096.

## API
`POST /api/orders` (201) · `GET /api/orders/{id}` (200/404) · `POST /api/orders/{id}/cancel` (204/404/409)
