// Package storage lưu file upload và file export. Hiện chỉ có bản ổ đĩa.
//
// Service nhận Store (interface), main chọn bản cài đặt:
//
//	store, err := storage.NewLocal(cfg.Storage) // STORAGE_DIR
//
// Key do code tạo từ ID, không lấy tên file người dùng gửi lên:
//
//	key := "uploads/" + rep.ID.String() + "/source.pdf"
//	err := store.Put(ctx, key, file)
package storage
