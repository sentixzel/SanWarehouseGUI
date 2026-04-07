-- Дополнительные настройки базы данных
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Примечание: Таблицы будут созданы автоматически GORM при первом запуске сервера
-- Индексы будут созданы автоматически через GORM
-- Если нужно создать индексы вручную, раскомментируйте строки ниже после первого запуска

-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_sku ON products(sku);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_category ON products(category);
-- CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_products_status ON products(status);doc