CREATE TABLE sky_auth_users (
                                id SERIAL PRIMARY KEY,                    -- 用户ID
                                username VARCHAR(100) NOT NULL UNIQUE,     -- 用户名
                                password VARCHAR(255) NOT NULL,            -- 密码（加密存储）
                                email VARCHAR(255),                       -- 邮箱
                                code INT,                                 -- 验证码
                                phone VARCHAR(20),                        -- 电话
                                user_type INT NOT NULL,                   -- 用户类型：1-管理员，2-客户
                                status BOOLEAN DEFAULT TRUE,              -- 用户状态：是否启用
                                created_by VARCHAR(255) NOT NULL,          -- 创建者
                                updated_by VARCHAR(255),                  -- 更新者
                                created_at TIMESTAMPTZ(6) DEFAULT CURRENT_TIMESTAMP, -- 创建时间
                                updated_at TIMESTAMPTZ(6) DEFAULT CURRENT_TIMESTAMP, -- 修改时间
                                is_deleted BOOLEAN DEFAULT FALSE,         -- 删除标志
                                notes TEXT                                -- 备注
);

-- 创建 sky_auth_tokens 表，并继承 common_base 表的字段
CREATE TABLE sky_auth_tokens (
                                 id SERIAL PRIMARY KEY,                    -- Token ID
                                 user_id INT NOT NULL REFERENCES sky_auth_users(id),  -- 关联用户表
                                 token VARCHAR(512) NOT NULL,              -- 认证 token
                                 created_at TIMESTAMPTZ(6) DEFAULT CURRENT_TIMESTAMP, -- 创建时间
                                 expires_at TIMESTAMPTZ(6),                -- 过期时间
                                 created_by VARCHAR(255) NOT NULL,          -- 创建者
                                 updated_by VARCHAR(255),                  -- 更新者
                                 is_deleted BOOLEAN DEFAULT FALSE,         -- 删除标志
                                 notes TEXT                                -- 备注
);