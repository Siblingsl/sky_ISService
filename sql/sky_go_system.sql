-- ----------------------------
-- Table structure for sky_system_admins
-- ----------------------------
DROP TABLE IF EXISTS "public"."sky_system_admins";
CREATE TABLE "public"."sky_system_admins" (
                                              "id" int4 NOT NULL DEFAULT nextval('sky_auth_admins_id_seq'::regclass),
                                              "username" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
                                              "password" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
                                              "full_name" varchar(255) COLLATE "pg_catalog"."default",
                                              "email" varchar(255) COLLATE "pg_catalog"."default",
                                              "phone" varchar(20) COLLATE "pg_catalog"."default",
                                              "status" bool DEFAULT true,
                                              "created_by" int4,
                                              "updated_by" int4,
                                              "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
                                              "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
                                              "deleted" bool DEFAULT false,
                                              "token" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;
ALTER TABLE "public"."sky_system_admins" OWNER TO "postgres";
COMMENT ON COLUMN "public"."sky_system_admins"."id" IS '主键';
COMMENT ON COLUMN "public"."sky_system_admins"."username" IS '员工用户名';
COMMENT ON COLUMN "public"."sky_system_admins"."password" IS '员工密码';
COMMENT ON COLUMN "public"."sky_system_admins"."full_name" IS '员工全名';
COMMENT ON COLUMN "public"."sky_system_admins"."email" IS '员工电子邮件';
COMMENT ON COLUMN "public"."sky_system_admins"."phone" IS '员工电话';
COMMENT ON COLUMN "public"."sky_system_admins"."status" IS '员工状态：是否启用';
COMMENT ON COLUMN "public"."sky_system_admins"."created_by" IS '创建人ID';
COMMENT ON COLUMN "public"."sky_system_admins"."updated_by" IS '更新人ID';
COMMENT ON COLUMN "public"."sky_system_admins"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."sky_system_admins"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."sky_system_admins"."deleted" IS '删除标志';
COMMENT ON COLUMN "public"."sky_system_admins"."token" IS 'token';

-- ----------------------------
-- Table structure for sky_system_menus
-- ----------------------------
DROP TABLE IF EXISTS "public"."sky_system_menus";
CREATE TABLE "public"."sky_system_menus" (
                                             "id" int4 NOT NULL DEFAULT nextval('sky_system_menus_id_seq'::regclass),
                                             "menu_name" varchar(100) COLLATE "pg_catalog"."default" NOT NULL,
                                             "menu_url" varchar(255) COLLATE "pg_catalog"."default",
                                             "parent_id" int4,
                                             "menu_type" int4,
                                             "menu_icon" varchar(100) COLLATE "pg_catalog"."default",
                                             "description" text COLLATE "pg_catalog"."default",
                                             "status" bool DEFAULT true,
                                             "created_by" int4,
                                             "updated_by" int4,
                                             "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
                                             "updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
                                             "deleted" bool DEFAULT false
)
;
ALTER TABLE "public"."sky_system_menus" OWNER TO "postgres";
COMMENT ON COLUMN "public"."sky_system_menus"."id" IS '主键';
COMMENT ON COLUMN "public"."sky_system_menus"."menu_name" IS '菜单名称';
COMMENT ON COLUMN "public"."sky_system_menus"."menu_url" IS '菜单链接（可以为空）';
COMMENT ON COLUMN "public"."sky_system_menus"."parent_id" IS '父菜单ID';
COMMENT ON COLUMN "public"."sky_system_menus"."menu_type" IS '菜单类型: 1-目录，2-菜单，3-按钮''';
COMMENT ON COLUMN "public"."sky_system_menus"."menu_icon" IS '菜单图标';
COMMENT ON COLUMN "public"."sky_system_menus"."description" IS '描述';
COMMENT ON COLUMN "public"."sky_system_menus"."status" IS '状态字段，默认启用';
COMMENT ON COLUMN "public"."sky_system_menus"."created_by" IS '创建人ID';
COMMENT ON COLUMN "public"."sky_system_menus"."updated_by" IS '更新人ID';
COMMENT ON COLUMN "public"."sky_system_menus"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."sky_system_menus"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."sky_system_menus"."deleted" IS '删除标志';

-- ----------------------------
-- Table structure for sky_system_roles
-- ----------------------------
DROP TABLE IF EXISTS sky_system_roles;
CREATE TABLE sky_system_roles (
                                  id SERIAL PRIMARY KEY,              -- 主键，自增
                                  role_name VARCHAR(100) NOT NULL,     -- 角色名称
                                  role_key VARCHAR(100) NOT NULL,      -- 角色唯一标识
                                  role_sort VARCHAR(100) NOT NULL,     -- 角色排序
                                  description TEXT,                    -- 角色描述
                                  status BOOLEAN DEFAULT TRUE,         -- 状态字段，默认为 true
                                  created_by INT,                      -- 创建者
                                  updated_by INT,                      -- 更新者
                                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
                                  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
                                  is_deleted BOOLEAN DEFAULT FALSE        -- 软删除标志，默认为 false
);

ALTER TABLE "public"."sky_system_roles" OWNER TO "postgres";
COMMENT ON COLUMN "public"."sky_system_roles"."id" IS '主键';
COMMENT ON COLUMN "public"."sky_system_roles"."role_name" IS '角色名称';
COMMENT ON COLUMN "public"."sky_system_roles"."description" IS '角色描述';
COMMENT ON COLUMN "public"."sky_system_roles"."status" IS '角色状态：是否启用';
COMMENT ON COLUMN "public"."sky_system_roles"."created_by" IS '创建人ID';
COMMENT ON COLUMN "public"."sky_system_roles"."updated_by" IS '更新人ID';
COMMENT ON COLUMN "public"."sky_system_roles"."created_at" IS '创建时间';
COMMENT ON COLUMN "public"."sky_system_roles"."updated_at" IS '更新时间';
COMMENT ON COLUMN "public"."sky_system_roles"."deleted" IS '删除标志';


-- ----------------------------
-- Table structure for admins_roles
-- ----------------------------
DROP TABLE IF EXISTS "public"."admins_roles";
CREATE TABLE "public"."admins_roles" (
                                         "admin_id" int4 NOT NULL,
                                         "role_id" int4 NOT NULL,
                                         "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
;
ALTER TABLE "public"."admins_roles" OWNER TO "postgres";

-- ----------------------------
-- Table structure for role_menu
-- ----------------------------
DROP TABLE IF EXISTS "public"."role_menu";
CREATE TABLE "public"."role_menu" (
                                      "role_id" int4 NOT NULL,
                                      "menu_id" int4 NOT NULL,
                                      "created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
;
ALTER TABLE "public"."role_menu" OWNER TO "postgres";