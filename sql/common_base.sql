-- 公共表
CREATE TABLE common_base (
    id SERIAL PRIMARY KEY,  -- 使用 SERIAL 类型自动生成 id，等同于 DEFAULT nextval('common_base_id_seq')
    status INT4 DEFAULT 1,  -- 默认状态为 1
    is_deleted BOOLEAN DEFAULT false,  -- 默认未删除
    created_by VARCHAR(255) NOT NULL,  -- 创建者
    updated_by VARCHAR(255),  -- 更新者
    created_at TIMESTAMPTZ(6) DEFAULT CURRENT_TIMESTAMP,  -- 默认当前时间戳
    updated_at TIMESTAMPTZ(6) DEFAULT CURRENT_TIMESTAMP,  -- 默认当前时间戳
    notes TEXT  -- 备注
);