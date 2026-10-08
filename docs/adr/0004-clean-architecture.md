# 4. Clean Architecture cho source code ứng dụng

- Trạng thái: Chấp nhận
- Ngày: 2026-10-07

## Bối cảnh
App cần dễ test, đổi được hạ tầng (DB, payment, message broker) mà không sửa nghiệp vụ, và giữ cấu trúc nhất quán khi thêm feature.

## Quyết định
- Tổ chức theo **feature** (`order`, `customer`); mỗi feature có 4 lớp: domain → application → (presentation | infrastructure).
- Phụ thuộc chỉ hướng vào trong; port (interface/trait/Protocol) khai báo ở application, adapter hiện thực ở infrastructure.
- Package `domain/app/transport/adapter` theo feature trong `internal/`; kiểm tra bằng **depguard** (golangci-lint).
- Composition root duy nhất chọn adapter theo env. Có adapter `local` (in-memory/log) để chạy khi chưa có DB/Stripe/Kafka.
- Map lỗi nghiệp vụ → HTTP status ở một chỗ.

## Hệ quả
- Unit test domain/use case không cần DB hay mock framework (dùng fake).
- Nhiều file/boilerplate hơn cho feature nhỏ.
- Production phải cấu hình `DATABASE_URL`, `STRIPE_API_KEY`, `KAFKA_BROKERS`; RDS/MSK chưa có trong `infra/`.
