# CLAUDE.md — docs/

Tài liệu dự án, viết bằng tiếng Việt, Markdown.

| File | Nội dung | Cập nhật khi |
|---|---|---|
| `architecture.md` | Sơ đồ, bảng thành phần, bảng môi trường | Thêm/bỏ thành phần hạ tầng, đổi thông số env |
| `runbook.md` | Kết nối cluster, rollback, xử lý sự cố | Gặp sự cố mới hoặc đổi quy trình vận hành |
| `adr/NNNN-*.md` | Architecture Decision Records | Có quyết định kiến trúc mới |

## ADR
- Đặt tên `NNNN-tieu-de-khong-dau.md`, số tăng dần (hiện có 0001–0004).
- Mẫu: tiêu đề `# N. ...`, `Trạng thái` (Đề xuất / Chấp nhận / Thay thế bởi ADR-X), `Ngày`, rồi `Bối cảnh` → `Quyết định` → `Hệ quả`.
- Không sửa nội dung ADR đã chấp nhận; muốn đổi thì viết ADR mới và đánh dấu ADR cũ "Thay thế bởi".

## Quy ước
- Lệnh trong runbook phải chạy được nguyên văn (chỉ thay `<env>`, `<pod>`...).
- Thông số trong `architecture.md` phải khớp `infra/envs/*/terraform.tfvars` và `deploy/overlays/*`; đổi cấu hình thì sửa docs trong cùng PR.
