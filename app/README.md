# base-go — app

Service Go theo Clean Architecture (feature `order` đầy đủ, `customer` mới có domain).
Chi tiết quy ước: [CLAUDE.md](CLAUDE.md).

```bash
make run                 # chạy local, adapter in-memory
make test                # unit test
make test-integration    # repository + Postgres (Docker)
make test-e2e            # E2E (Docker)
make lint                # golangci-lint + depguard
```

API:
```bash
curl -XPOST localhost:8080/api/orders -d '{"customer_id":"<uuid>","items":[{"product_id":"<uuid>","unit_price":1000,"quantity":2}]}'
curl localhost:8080/api/orders/<id>
curl -XPOST localhost:8080/api/orders/<id>/cancel
```
