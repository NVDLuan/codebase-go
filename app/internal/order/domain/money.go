package domain

// Money lưu theo đơn vị nhỏ nhất (vd. đồng, cent) để tránh lỗi số thực.
type Money struct {
	Amount   int64
	Currency string
}

// Add cộng 2 khoản cùng loại tiền. Order đã kiểm tra currency khi tạo nên mismatch là lỗi lập trình.
func (m Money) Add(o Money) Money {
	if m.Currency != o.Currency {
		panic("currency mismatch")
	}
	return Money{Amount: m.Amount + o.Amount, Currency: m.Currency}
}

func (m Money) Mul(q int) Money { return Money{Amount: m.Amount * int64(q), Currency: m.Currency} }
