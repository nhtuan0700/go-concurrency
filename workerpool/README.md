## Worker pool là gì?
- Cách hoạt động: Trong mô hình Worker Pool, một **Task Queue** sẽ chứa các công việc cần thực hiện. Nhiều Worker được khởi chạy song song để lấy công việc từ hàng đợi này và xử lý. Mỗi Worker sẽ thực hiện một tác vụ từ **Task Queue** và khi hoàn thành, Worker sẽ quay lại để lấy công việc tiếp theo

- Đặc điểm:
  - Synchronous: Các Worker thường lấy một task từ hàng đợi, xử lý nó rồi quay lại hàng đợi, điều này có thể là đồng bộ hoặc không đồng bộ.
  - Load Distribution: Mô hình này giúp phân bổ khối lượng công việc đồng đều giữa các Worker.
  - Control over Concurrency: Bạn có thể dễ dàng điều chỉnh số lượng Worker để quản lý độ song song của hệ thống.
