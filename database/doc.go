// Package database mở pool Postgres (pgxpool) dùng chung cho cả app.
//
// Mở một lần trong main, ping luôn nên sai URL hay DB chưa lên thì dừng ngay
// lúc khởi động, không đợi tới request đầu tiên:
//
//	type serverConfig struct {
//		DB database.Config // DB_MAX_CONNS, DB_CONNECT_TIMEOUT...
//	}
//
//	pool, err := database.Open(ctx, cfg.DB)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer database.Close(pool)
//
// Repository nhận *pgxpool.Pool, không tự mở kết nối. Method nào ghi thì làm
// hết thay đổi đó trong một WithTx, gồm cả ghi event vào outbox, để thay đổi
// và event cùng commit hoặc cùng mất. Service không bao giờ thấy transaction:
//
//	func (r *OrderRepository) Save(ctx context.Context, o *domain.Order) error {
//		return database.WithTx(ctx, r.pool, func(tx pgx.Tx) error {
//			n, err := r.q.WithTx(tx).UpdateOrderStatus(ctx, ...)
//			if err != nil {
//				return err
//			}
//			if n == 0 {
//				return domain.ErrOrderChanged // lỗi domain đi ra nguyên vẹn
//			}
//			return r.outbox.Append(ctx, tx, events...)
//		})
//	}
package database
