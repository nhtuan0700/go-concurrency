## Mô hình Pub/Sub

- Cách hoạt động: Mô hình này có hai thành phần chính là Publisher và Subscriber. Một Publisher phát ra các sự kiện hoặc thông điệp, và các Subscriber đăng ký để nhận các sự kiện này. Publisher và Subscriber hoạt động độc lập với nhau, và giao tiếp qua một hệ thống trung gian (message broker) như Kafka, NATS, RabbitMQ hoặc các channel nội bộ trong Golang.
- Đặc điểm:
  - Decoupled (Tách biệt): Publisher không cần biết về sự tồn tại của Subscriber, và ngược lại. Điều này giúp hệ thống mở rộng dễ dàng.
  - Broadcast: Một thông điệp từ Publisher có thể được gửi đến nhiều Subscriber cùng một lúc.
  - Asynchronous: Hệ thống thường chạy không đồng bộ, giúp giảm tải cho các thành phần chính của ứng dụng.
