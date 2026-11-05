// Package jobs đẩy việc chạy nền vào queue (River trên Postgres). Service và
// repository chỉ thấy Enqueuer, không import River.
//
// Args của job nằm trong package job/ của module sở hữu nó: chỗ enqueue
// (repository hoặc service) và worker đều import được, module khác thì không.
// Args là struct serialize được ra JSON, có Kind() không đổi qua các lần
// deploy (job cũ trong DB tìm worker theo kind). Job do người dùng yêu cầu
// nhúng ActorArgs:
//
//	// reports/job/generate_report.go
//	type GenerateReportArgs struct {
//		jobs.ActorArgs
//		ReportID uuid.UUID `json:"report_id"`
//	}
//
//	func (GenerateReportArgs) Kind() string { return "reports.generate" }
//
// Binary API dùng client chỉ insert, binary worker dùng client chạy job; bọc client
// rồi truyền Enqueuer cho repository/service:
//
//	client, err := jobs.NewInsertClient(pool, log)                   // binary API
//	client, err := jobs.NewWorkerClient(pool, log, cfg.Jobs, workers, periodic) // binary worker
//	enq := jobs.NewRiver(client)
//
// Job đi kèm một thay đổi dữ liệu thì repository enqueue bằng EnqueueTx, trong
// cùng database.WithTx với thay đổi đó: commit thì job mới tồn tại, rollback
// thì job mất theo (transactional outbox). Service quyết định có cần job hay
// không, nhưng không bao giờ thấy transaction:
//
//	// reports/repository/report_repository.go
//	func (r *ReportRepository) SaveRequested(ctx context.Context, rep *domain.Report) error {
//		return database.WithTx(ctx, r.pool, func(tx pgx.Tx) error {
//			if err := r.q.WithTx(tx).MarkReportRequested(ctx, rep.ID); err != nil {
//				return err
//			}
//			return r.jobs.EnqueueTx(ctx, tx, job.GenerateReportArgs{
//				ActorArgs: jobs.ActorArgsFrom(ctx),
//				ReportID:  rep.ID,
//			}, jobs.WithUniqueArgs())
//		})
//	}
//
// Worker (trong worker/ của module) đặt lại actor rồi gọi đúng một method của
// service:
//
//	func (w *GenerateReportWorker) Work(ctx context.Context, j *river.Job[job.GenerateReportArgs]) error {
//		ctx, err := jobs.RestoreActor(ctx, j.Args.ActorID, w.loadActor)
//		if err != nil {
//			return err
//		}
//		return w.svc.GenerateReport(ctx, j.Args.ReportID)
//	}
//
// Job không gắn với thay đổi nào thì service gọi thẳng Enqueue, không có đảm
// bảo trên. Scheduled job chạy bằng SystemActor.
//
// Worker có thể chạy lại một job nhiều lần (retry, crash giữa chừng), nên
// việc nó làm phải idempotent.
package jobs
