import random
from datetime import datetime, timedelta

# 配置
START_ID = 1115
PROJECT_ID = 35
BATCH_SIZE = 10
TOTAL_COUNT = 100

# PS5测试相关数据
modules = ["Dashboard UI", "Game Library", "PlayStation Store", "Settings", "Network", 
           "Controller", "Audio System", "Video Output", "Storage Management", "User Account",
           "Trophy System", "Share/Capture", "Party Chat", "Remote Play", "VR Support"]

types = ["Functional", "UI", "UIInteraction", "Compatibility", "Performance", "Security", "Environment"]
severities = ["Critical", "Major", "Minor", "Trivial"]
priorities = ["A", "B", "C", "D"]
statuses = ["New", "InProgress", "Confirmed", "Resolved", "Closed", "Reopened"]
phases = ["UT", "IT", "ST", "UAT", "Regression"]
versions = ["v1.0.0", "v1.0.1", "v1.1.0", "v1.2.0", "v2.0.0-beta", "v2.0.0"]
components = ["Frontend", "Backend", "Firmware", "Network Stack", "Graphics Engine", "Audio Engine", "Storage Driver"]
teams = ["QA-Tokyo", "QA-Shanghai", "QA-California", "Dev-Test", "Beta Tester"]
users = ["Tanaka", "Yamamoto", "Zhang Wei", "Li Ming", "John Smith", "Emily Chen", "Sato", "Kim"]
assignees = ["dev_tanaka", "dev_suzuki", "dev_wang", "dev_chen", "dev_johnson", "dev_kim"]
models = ["PS5 Standard Edition", "PS5 Digital Edition", "PS5 Slim", "PS5 Pro"]
frequencies = ["Always", "Often", "Sometimes", "Rarely", "Once"]

actions = ["loading game", "starting application", "connecting to server", "downloading content",
           "updating system", "switching user", "entering rest mode", "resuming from suspend",
           "capturing screenshot", "streaming gameplay", "voice chat active", "installing game"]

objects = ["Home screen", "Store page", "Game menu", "Settings panel", "Trophy list",
           "Friends list", "Download queue", "Capture gallery", "User profile", "Party interface"]

# 生成时间范围
start_date = datetime(2025, 11, 1)
end_date = datetime(2026, 2, 6)

def generate_bug(seq_id):
    defect_id = f"{seq_id:06d}"
    module = random.choice(modules)
    action = random.choice(actions)
    obj = random.choice(objects)
    
    title = f"[{module}] {obj} issue when {action}"
    
    description = f"""## 问题描述
在使用 {module} 功能时，{obj} 发生异常。

## 复现步骤
1. 启动 PS5 主机
2. 进入 {module} 菜单
3. 执行 {action} 操作
4. 观察到问题发生

## 测试环境
- 机型: {random.choice(models)}
- 固件版本: {random.choice(versions)}"""
    
    status = random.choice(statuses)
    fix_version = random.choice(versions) if status in ["Resolved", "Closed"] else ""
    resolution = random.choice(["Fixed in firmware update", "Code patch applied", "Configuration issue resolved"]) if status in ["Resolved", "Closed"] else ""
    
    # 生成随机时间
    days_diff = (end_date - start_date).days
    created_at = start_date + timedelta(days=random.randint(0, days_diff), hours=random.randint(0, 23))
    updated_at = created_at + timedelta(hours=random.randint(1, 72))
    
    return f"""('{defect_id}-' || hex(randomblob(8)), '{defect_id}', {PROJECT_ID}, '{title}', '{module}', '{description.replace("'", "''")}', '重启主机或重新进入{module}功能', '{random.choice(priorities)}', '{random.choice(severities)}', '{random.choice(types)}', '{random.choice(frequencies)}', '{random.choice(versions)}', '{random.choice(phases)}', 'TC-PS5-{random.randint(1000, 9999)}', '{random.choice(assignees)}', '{random.choice(["S", "A", "B", "C"])}', '{random.choice(teams)}', '{module} > {obj}', '{fix_version}', '', '{random.choice(components)}', '{resolution}', '{random.choice(models)}', '{random.choice(users)}', '{status}', 1, 1, '{created_at.strftime("%Y-%m-%d %H:%M:%S")}', '{updated_at.strftime("%Y-%m-%d %H:%M:%S")}')"""

# 生成10批SQL文件
for batch in range(10):
    batch_num = batch + 1
    start_seq = START_ID + (batch * BATCH_SIZE)
    
    sql_lines = [f"-- Batch {batch_num}: Defects {start_seq} to {start_seq + BATCH_SIZE - 1}"]
    sql_lines.append("INSERT INTO defects (id, defect_id, project_id, title, subject, description, recovery_method, priority, severity, type, frequency, detected_version, phase, case_id, assignee, recovery_rank, detection_team, location, fix_version, sqa_memo, component, resolution, models, detected_by, status, created_by, updated_by, created_at, updated_at) VALUES")
    
    values = []
    for i in range(BATCH_SIZE):
        seq_id = start_seq + i
        values.append(generate_bug(seq_id))
    
    sql_lines.append(",\n".join(values) + ";")
    
    filename = f"../scripts/ps5_bugs_batch{batch_num}.sql"
    with open(filename, 'w', encoding='utf-8') as f:
        f.write("\n".join(sql_lines))
    
    print(f"Generated {filename}")

print(f"\nGenerated {TOTAL_COUNT} PS5 bugs in 10 batches")
print(f"Defect IDs: {START_ID:06d} to {START_ID + TOTAL_COUNT - 1:06d}")
