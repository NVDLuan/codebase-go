# CLAUDE.md — scripts/

Script tiện ích, gọi qua `Makefile` (`make local`, `make cleanup`).

| Script | Việc làm |
|---|---|
| `local-dev.sh` | Tạo cluster kind `my-app-local`, build image `my-app:local`, render `deploy/base` **bỏ Ingress + ExternalSecret** (cần AWS), tạo Secret giả, apply. |
| `cleanup.sh` | Mặc định xoá cluster kind. `--aws` → `terraform destroy` env **dev** (hỏi xác nhận, xoá Ingress trước để ALB được dọn). |

## Quy ước
- Bash, mở đầu `set -euo pipefail`, đường dẫn tính từ `ROOT` (chạy được từ bất kỳ thư mục nào).
- Script đụng tới AWS phải hỏi xác nhận và **chỉ được nhắm dev**. Không viết script destroy/sửa staging/prod.
- Không chứa credential; dùng AWS profile/env của người chạy.
- Thêm script mới → thêm target tương ứng trong `Makefile` và 1 dòng vào bảng trên.
- Giữ `chmod +x`.
