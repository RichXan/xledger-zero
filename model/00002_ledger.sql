-- 分类表（简化命名）
CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(50),
    color VARCHAR(20),
    type INT NOT NULL,           -- 1: 收入, 2: 支出
    sort_order INT DEFAULT 0,
    is_system BOOLEAN DEFAULT FALSE,
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_categories_user ON categories(user_id, status);

-- 子分类表
CREATE TABLE sub_categories (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL REFERENCES categories(id),
    user_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(50),
    color VARCHAR(20),
    sort_order INT DEFAULT 0,
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sub_categories_category ON sub_categories(category_id, status);

-- 交易记录表
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL REFERENCES categories(id),
    sub_category_id BIGINT REFERENCES sub_categories(id),
    amount DECIMAL(15, 2) NOT NULL,
    type INT NOT NULL,           -- 1: 收入, 2: 支出
    description TEXT,
    note TEXT,
    transaction_date TIMESTAMP NOT NULL,
    tags TEXT[],
    location VARCHAR(200),
    images TEXT[],
    status INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_transactions_user ON transactions(user_id, transaction_date DESC);
CREATE INDEX idx_transactions_category ON transactions(category_id);
CREATE INDEX idx_transactions_date ON transactions(transaction_date);

-- 插入默认分类数据
-- 支出分类
INSERT INTO categories (user_id, name, icon, color, type, sort_order, is_system, status) VALUES
(0, '餐饮', '🍔', '#FF6B6B', 2, 1, true, 1),
(0, '交通', '🚗', '#4ECDC4', 2, 2, true, 1),
(0, '购物', '🛍️', '#95E1D3', 2, 3, true, 1),
(0, '娱乐', '🎮', '#F38181', 2, 4, true, 1),
(0, '住房', '🏠', '#AA96DA', 2, 5, true, 1),
(0, '医疗', '⚕️', '#FCBAD3', 2, 6, true, 1),
(0, '教育', '📚', '#A8D8EA', 2, 7, true, 1),
(0, '其他支出', '💸', '#FFFFD2', 2, 99, true, 1);

-- 收入分类
INSERT INTO categories (user_id, name, icon, color, type, sort_order, is_system, status) VALUES
(0, '工资', '💼', '#51CF66', 1, 1, true, 1),
(0, '奖金', '🎁', '#69DB7C', 1, 2, true, 1),
(0, '投资收益', '📈', '#8CE99A', 1, 3, true, 1),
(0, '兼职', '💻', '#B2F2BB', 1, 4, true, 1),
(0, '其他收入', '💰', '#D3F9D8', 1, 99, true, 1);
