-- 定时任务测试用例
-- 插入到 job_center 表

-- 1. 简单测试任务 - 每5分钟执行一次
INSERT INTO `job_center` (`job_name`, `job_type`, `cron_expr`, `service_name`, `service_method`, `args`, `run_status`, `enabled`, `create_time`, `update_time`)
VALUES ('测试任务-Hello', 'cron_job', '0 */5 * * * *', 'testService', 'Hello', '{}', '待运行', 1, NOW(), NOW());

-- 2. 带参数测试任务 - 每10分钟执行一次
INSERT INTO `job_center` (`job_name`, `job_type`, `cron_expr`, `service_name`, `service_method`, `args`, `run_status`, `enabled`, `create_time`, `update_time`)
VALUES ('测试任务-带参数', 'cron_job', '0 */10 * * * *', 'testService', 'HelloWithArgs', '{"name": "张三", "count": 10}', '待运行', 1, NOW(), NOW());

-- 3. 模拟耗时工作 - 每30分钟执行一次
INSERT INTO `job_center` (`job_name`, `job_type`, `cron_expr`, `service_name`, `service_method`, `args`, `run_status`, `enabled`, `create_time`, `update_time`)
VALUES ('模拟耗时工作', 'cron_job', '0 */30 * * * *', 'testService', 'SimulateWork', '{"duration": 3}', '待运行', 1, NOW(), NOW());

-- 4. 模拟错误任务 - 每小时执行一次（已禁用，用于测试错误处理）
INSERT INTO `job_center` (`job_name`, `job_type`, `cron_expr`, `service_name`, `service_method`, `args`, `run_status`, `enabled`, `create_time`, `update_time`)
VALUES ('模拟错误任务', 'cron_job', '0 0 * * * *', 'testService', 'SimulateError', '{}', '待运行', 0, NOW(), NOW());

-- 5. 清理旧日志 - 每天凌晨2点执行
INSERT INTO `job_center` (`job_name`, `job_type`, `cron_expr`, `service_name`, `service_method`, `args`, `run_status`, `enabled`, `create_time`, `update_time`)
VALUES ('清理旧日志', 'cron_job', '0 0 2 * * *', 'logService', 'CleanOldLogs', '{"days": 30}', '待运行', 1, NOW(), NOW());

-- 6. 归档日志 - 每天凌晨3点执行
INSERT INTO `job_center` (`job_name`, `job_type`, `cron_expr`, `service_name`, `service_method`, `args`, `run_status`, `enabled`, `create_time`, `update_time`)
VALUES ('归档日志', 'cron_job', '0 0 3 * * *', 'logService', 'ArchiveLogs', '{}', '待运行', 1, NOW(), NOW());

-- 7. 发送每日报告 - 每天早上9点执行
INSERT INTO `job_center` (`job_name`, `job_type`, `cron_expr`, `service_name`, `service_method`, `args`, `run_status`, `enabled`, `create_time`, `update_time`)
VALUES ('发送每日报告', 'cron_job', '0 0 9 * * *', 'notifyService', 'SendDailyReport', '{}', '待运行', 1, NOW(), NOW());

-- 8. 计算每日统计 - 每天凌晨1点执行
INSERT INTO `job_center` (`job_name`, `job_type`, `cron_expr`, `service_name`, `service_method`, `args`, `run_status`, `enabled`, `create_time`, `update_time`)
VALUES ('计算每日统计', 'cron_job', '0 0 1 * * *', 'statsService', 'CalculateDailyStats', '{}', '待运行', 1, NOW(), NOW());

-- 9. 系统健康检查 - 每5分钟执行一次
INSERT INTO `job_center` (`job_name`, `job_type`, `cron_expr`, `service_name`, `service_method`, `args`, `run_status`, `enabled`, `create_time`, `update_time`)
VALUES ('系统健康检查', 'cron_job', '0 */5 * * * *', 'statsService', 'HealthCheck', '{}', '待运行', 1, NOW(), NOW());

-- 10. 数据备份 - 每天凌晨4点执行
INSERT INTO `job_center` (`job_name`, `job_type`, `cron_expr`, `service_name`, `service_method`, `args`, `run_status`, `enabled`, `create_time`, `update_time`)
VALUES ('数据备份', 'cron_job', '0 0 4 * * *', 'statsService', 'BackupData', '{}', '待运行', 1, NOW(), NOW());

-- 11. HTTP Hook 测试任务 - 每5分钟调用一次外部接口
INSERT INTO `job_center` (`job_name`, `job_type`, `cron_expr`, `service_name`, `service_method`, `hook_url`, `args`, `run_status`, `enabled`, `create_time`, `update_time`)
VALUES ('HTTP Hook测试', 'cron_hook', '0 */5 * * * *', '', '', 'http://localhost:8080/api/health', '{"test": true}', '待运行', 0, NOW(), NOW());
