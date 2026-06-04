-- Gaia Server 账户体系数据库初始化脚本
-- 在MySQL容器启动时自动执行
-- 数据库: account_system

USE account_system;

-- ==========================================
-- 1. 用户表 (users)
-- ==========================================
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    uuid VARCHAR(36) NOT NULL UNIQUE COMMENT '用户唯一标识UUID',
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    email VARCHAR(100) UNIQUE COMMENT '邮箱',
    phone VARCHAR(20) UNIQUE COMMENT '手机号',
    password_hash VARCHAR(255) NOT NULL COMMENT '密码哈希',
    salt VARCHAR(64) COMMENT '密码盐值',
    nickname VARCHAR(50) COMMENT '昵称',
    status TINYINT DEFAULT 1 COMMENT '状态: 0-禁用 1-正常 2-待验证 3-锁定',
    email_verified TINYINT DEFAULT 0 COMMENT '邮箱是否验证: 0-否 1-是',
    phone_verified TINYINT DEFAULT 0 COMMENT '手机是否验证: 0-否 1-是',
    last_login_at DATETIME COMMENT '最后登录时间',
    last_login_ip VARCHAR(45) COMMENT '最后登录IP',
    login_count INT UNSIGNED DEFAULT 0 COMMENT '登录次数',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME COMMENT '软删除时间',
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB COMMENT='用户表';

-- ==========================================
-- 2. 角色表 (roles)
-- ==========================================
CREATE TABLE IF NOT EXISTS roles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE COMMENT '角色名称',
    code VARCHAR(50) NOT NULL UNIQUE COMMENT '角色代码',
    description VARCHAR(255) COMMENT '角色描述',
    is_system TINYINT DEFAULT 0 COMMENT '是否系统角色: 0-否 1-是',
    status TINYINT DEFAULT 1 COMMENT '状态: 0-禁用 1-启用',
    sort_order INT DEFAULT 0 COMMENT '排序',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_code (code),
    INDEX idx_status (status)
) ENGINE=InnoDB COMMENT='角色表';

-- ==========================================
-- 3. 权限表 (permissions)
-- ==========================================
CREATE TABLE IF NOT EXISTS permissions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT '权限名称',
    code VARCHAR(100) NOT NULL UNIQUE COMMENT '权限代码',
    type VARCHAR(20) DEFAULT 'menu' COMMENT '类型: menu-菜单 button-按钮 api-接口',
    parent_id BIGINT UNSIGNED DEFAULT 0 COMMENT '父级ID',
    path VARCHAR(255) COMMENT '路由路径/API路径',
    icon VARCHAR(100) COMMENT '图标',
    sort_order INT DEFAULT 0 COMMENT '排序',
    status TINYINT DEFAULT 1 COMMENT '状态: 0-禁用 1-启用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_parent_id (parent_id),
    INDEX idx_type (type),
    INDEX idx_status (status)
) ENGINE=InnoDB COMMENT='权限表';

-- ==========================================
-- 4. 用户-角色关联表 (user_roles)
-- ==========================================
CREATE TABLE IF NOT EXISTS user_roles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    role_id BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_role (user_id, role_id),
    INDEX idx_user_id (user_id),
    INDEX idx_role_id (role_id)
) ENGINE=InnoDB COMMENT='用户角色关联表';

-- ==========================================
-- 5. 角色-权限关联表 (role_permissions)
-- ==========================================
CREATE TABLE IF NOT EXISTS role_permissions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    role_id BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    permission_id BIGINT UNSIGNED NOT NULL COMMENT '权限ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_role_permission (role_id, permission_id),
    INDEX idx_role_id (role_id),
    INDEX idx_permission_id (permission_id)
) ENGINE=InnoDB COMMENT='角色权限关联表';

-- ==========================================
-- 6. 登录日志表 (login_logs)
-- ==========================================
CREATE TABLE IF NOT EXISTS login_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED COMMENT '用户ID',
    username VARCHAR(50) COMMENT '用户名',
    login_type VARCHAR(20) DEFAULT 'password' COMMENT '登录方式: password-密码 sms-短信 oauth-第三方',
    ip_address VARCHAR(45) COMMENT 'IP地址',
    user_agent VARCHAR(500) COMMENT '用户代理',
    device_type VARCHAR(20) COMMENT '设备类型: pc/mobile/tablet',
    os VARCHAR(50) COMMENT '操作系统',
    browser VARCHAR(50) COMMENT '浏览器',
    location VARCHAR(100) COMMENT '登录地点',
    status TINYINT DEFAULT 1 COMMENT '状态: 0-失败 1-成功',
    fail_reason VARCHAR(255) COMMENT '失败原因',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at),
    INDEX idx_ip_address (ip_address)
) ENGINE=InnoDB COMMENT='登录日志表';

-- ==========================================
-- 7. 用户Token表 (user_tokens)
-- ==========================================
CREATE TABLE IF NOT EXISTS user_tokens (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    token_type VARCHAR(20) DEFAULT 'access' COMMENT 'Token类型: access-访问 refresh-刷新',
    token VARCHAR(500) NOT NULL COMMENT 'Token值',
    device_id VARCHAR(100) COMMENT '设备ID',
    expires_at DATETIME NOT NULL COMMENT '过期时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_token (token(100)),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB COMMENT='用户Token表';

-- ==========================================
-- 8. 验证码表 (verification_codes)
-- ==========================================
CREATE TABLE IF NOT EXISTS verification_codes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    target VARCHAR(100) NOT NULL COMMENT '目标(邮箱/手机号)',
    code VARCHAR(10) NOT NULL COMMENT '验证码',
    type VARCHAR(20) NOT NULL COMMENT '类型: register-注册 login-登录 reset_password-重置密码 bind-绑定',
    used TINYINT DEFAULT 0 COMMENT '是否已使用: 0-否 1-是',
    expires_at DATETIME NOT NULL COMMENT '过期时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_target_type (target, type),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB COMMENT='验证码表';

-- ==========================================
-- 9. 第三方登录绑定表 (oauth_bindings)
-- ==========================================
CREATE TABLE IF NOT EXISTS oauth_bindings (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    provider VARCHAR(30) NOT NULL COMMENT '第三方平台: wechat/qq/weibo/github/google',
    open_id VARCHAR(100) NOT NULL COMMENT '第三方OpenID',
    union_id VARCHAR(100) COMMENT '第三方UnionID',
    access_token VARCHAR(500) COMMENT '访问Token',
    refresh_token VARCHAR(500) COMMENT '刷新Token',
    expires_at DATETIME COMMENT 'Token过期时间',
    nickname VARCHAR(50) COMMENT '第三方昵称',
    avatar_url VARCHAR(500) COMMENT '第三方头像',
    raw_data JSON COMMENT '原始数据',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_provider_openid (provider, open_id),
    INDEX idx_user_id (user_id)
) ENGINE=InnoDB COMMENT='第三方登录绑定表';

-- ==========================================
-- 10. 操作日志表 (operation_logs)
-- ==========================================
CREATE TABLE IF NOT EXISTS operation_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED COMMENT '用户ID',
    username VARCHAR(50) COMMENT '用户名',
    module VARCHAR(50) COMMENT '模块',
    action VARCHAR(50) COMMENT '操作',
    method VARCHAR(10) COMMENT '请求方法',
    path VARCHAR(255) COMMENT '请求路径',
    params JSON COMMENT '请求参数',
    result JSON COMMENT '返回结果',
    ip_address VARCHAR(45) COMMENT 'IP地址',
    user_agent VARCHAR(500) COMMENT '用户代理',
    duration INT COMMENT '耗时(ms)',
    status TINYINT DEFAULT 1 COMMENT '状态: 0-失败 1-成功',
    error_msg TEXT COMMENT '错误信息',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_module (module),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB COMMENT='操作日志表';

-- ==========================================
-- 初始化默认角色
-- ==========================================
INSERT INTO roles (name, code, description, is_system, sort_order) VALUES
('超级管理员', 'super_admin', '拥有系统所有权限', 1, 1),
('管理员', 'admin', '系统管理员', 1, 2),
('普通用户', 'user', '普通注册用户', 1, 3),
('访客', 'guest', '游客/未登录用户', 1, 4);

-- ==========================================
-- 初始化默认权限
-- ==========================================
INSERT INTO permissions (name, code, type, parent_id, path, sort_order) VALUES
('系统管理', 'system', 'menu', 0, '/system', 1),
('用户管理', 'system:user', 'menu', 1, '/system/user', 1),
('用户查看', 'system:user:view', 'button', 2, NULL, 1),
('用户新增', 'system:user:add', 'button', 2, NULL, 2),
('用户编辑', 'system:user:edit', 'button', 2, NULL, 3),
('用户删除', 'system:user:delete', 'button', 2, NULL, 4),
('角色管理', 'system:role', 'menu', 1, '/system/role', 2),
('角色查看', 'system:role:view', 'button', 7, NULL, 1),
('角色新增', 'system:role:add', 'button', 7, NULL, 2),
('角色编辑', 'system:role:edit', 'button', 7, NULL, 3),
('角色删除', 'system:role:delete', 'button', 7, NULL, 4),
('权限管理', 'system:permission', 'menu', 1, '/system/permission', 3),
('日志管理', 'system:log', 'menu', 1, '/system/log', 4),
('登录日志', 'system:log:login', 'menu', 13, '/system/log/login', 1),
('操作日志', 'system:log:operation', 'menu', 13, '/system/log/operation', 2);

-- ==========================================
-- 为角色分配权限
-- ==========================================
-- 超级管理员拥有所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions;

-- 管理员拥有除权限管理外的所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 2, id FROM permissions WHERE code NOT LIKE '%permission%';

-- 普通用户只有查看权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 3, id FROM permissions WHERE code LIKE '%:view';

-- ==========================================
-- 11. 任务中心表 (job_center) - Gaia框架表
-- ==========================================
CREATE TABLE IF NOT EXISTS job_center (
    id bigint(20) NOT NULL AUTO_INCREMENT,
    job_name varchar(255) NOT NULL COMMENT '任务名称',
    job_type varchar(50) NOT NULL COMMENT '任务类型: cron_job/cron_hook',
    cron_expr varchar(100) NOT NULL COMMENT 'Cron表达式',
    hook_url varchar(500) DEFAULT NULL COMMENT 'Hook URL（cron_hook类型）',
    service_name varchar(255) DEFAULT NULL COMMENT '服务名（cron_job类型）',
    service_method varchar(255) DEFAULT NULL COMMENT '方法名（cron_job类型）',
    args text COMMENT '参数JSON',
    timeout int(11) DEFAULT '300' COMMENT '超时时间（秒）',
    enabled tinyint(1) DEFAULT '1' COMMENT '是否启用',
    run_status varchar(50) DEFAULT '待运行' COMMENT '运行状态',
    last_run_time datetime DEFAULT NULL COMMENT '最后运行时间',
    create_time datetime DEFAULT CURRENT_TIMESTAMP,
    update_time datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_enabled (enabled),
    KEY idx_job_type (job_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='定时任务中心';

-- ==========================================
-- 12. 任务执行记录表 (job_record) - Gaia框架表
-- ==========================================
CREATE TABLE IF NOT EXISTS job_record (
    id bigint(20) NOT NULL AUTO_INCREMENT,
    job_id bigint(20) NOT NULL COMMENT '任务ID',
    run_time datetime DEFAULT NULL COMMENT '运行时间',
    end_time datetime DEFAULT NULL COMMENT '结束时间',
    status varchar(50) DEFAULT NULL COMMENT '执行状态',
    result text COMMENT '执行结果',
    create_time datetime DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_job_id (job_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务执行记录';
