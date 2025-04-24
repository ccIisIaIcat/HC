import requests
import os
import json
import mimetypes

# 服务器配置
BASE_URL = "https://jesuvukndxpo.sealoshzh.site"
API_BASE_URL = BASE_URL

# 管理员账户信息
ADMIN_USERNAME = "admin@yourdomain.com"
ADMIN_PASSWORD = "admin888"

# 图片配置
ITEM_IMAGE_PATH = "G:/cc/HC2/HC/doc/item/pic/美食探险背包.png"  # 物品大图路径
ITEM_ICON_PATH = "G:/cc/HC2/HC/doc/item/icon/美食探险背包.png"    # 物品图标路径

ITEM_NAME = "美食探险背包"
ITEM_DESCRIPTION = "什么什么什么？你居然尝试了20种新食物？我正想问你是不是饿昏头了...不过，这种冒险精神，还挺...特别的！这个探险背包给你，下次可以装更多奇怪食物，说不定能发现什么不得了的黑暗料理呢～"
ITEM_SOURCE = "饭团猫"

def get_file_mime_type(file_path):
    """自动检测文件的MIME类型"""
    mime_type, _ = mimetypes.guess_type(file_path)
    if mime_type is None:
        ext = os.path.splitext(file_path)[1].lower()
        mime_map = {
            '.png': 'image/png',
            '.jpg': 'image/jpeg',
            '.jpeg': 'image/jpeg',
            '.gif': 'image/gif'
        }
        mime_type = mime_map.get(ext, 'application/octet-stream')
    return mime_type

def login_admin():
    """登录管理员账户并获取token"""
    print("正在登录管理员账户...")
    response = requests.post(
        f"{API_BASE_URL}/api/login",
        json={
            "email": ADMIN_USERNAME,
            "password": ADMIN_PASSWORD
        }
    )
    if response.status_code == 200:
        token = response.json()["token"]
        print("管理员登录成功！")
        return token
    else:
        raise Exception(f"管理员登录失败: {response.text}")

def upload_image(file_path, directory, token):
    """上传图片并返回URL"""
    print(f"\n正在上传文件: {file_path}")
    print(f"目标目录: {directory}")
    
    if not os.path.exists(file_path):
        raise Exception(f"文件不存在: {file_path}")
    
    with open(file_path, 'rb') as f:
        mime_type = get_file_mime_type(file_path)
        files = {
            'file': (
                os.path.basename(file_path),
                f,
                mime_type
            )
        }
        data = {'directory': directory}
        headers = {'Authorization': f'Bearer {token}'}
        
        response = requests.post(
            f"{API_BASE_URL}/api/admin/upload/image",
            files=files,
            data=data,
            headers=headers
        )
        
        if response.status_code == 200:
            url = response.json()["url"]
            print(f"上传成功！URL: {url}")
            return url
        else:
            raise Exception(f"上传失败: {response.text}")

def create_item(token, item_data):
    """创建物品"""
    print("\n正在创建物品...")
    print(f"物品数据: {json.dumps(item_data, ensure_ascii=False, indent=2)}")
    
    headers = {
        'Authorization': f'Bearer {token}',
        'Content-Type': 'application/json'
    }
    
    response = requests.post(
        f"{API_BASE_URL}/api/admin/items",
        json=item_data,
        headers=headers
    )
    
    if response.status_code == 200:
        print("物品创建成功！")
        return response.json()
    else:
        raise Exception(f"创建物品失败: {response.text}")
    
achievements = [
    {
        "name": "饭团猫的礼物",
        "condition": "记录第一次食物",
        "description": "哼～总算开始记录了？我等得花儿都谢了！不过能迈出第一步已经很不错了啦～给你这只饭团当见面礼，它会陪你一起记录更多美食的。别误会，我才不是专门准备的礼物呢，只是...刚好做多了一个而已！",
        "icon_suggestion": "一只抱着饭团的傲娇猫咪",
        "image_suggestion": "一只白色猫咪，肚子是饭团形状，一边假装不屑地递出饭团，一边偷偷观察领奖人的反应"
    },
    {
        "name": "第七天的饭勺",
        "condition": "累计打卡7天",
        "description": "哼～累计记录七天就想要奖励？也不是不行啦～毕竟坚持这么久的也不多见，勉为其难奖励你这把闪亮勺子～",
        "icon_suggestion": "傲娇表情的金色小饭勺",
        "image_suggestion": "傲娇脸的金色饭勺，嘴巴一边撇着，一边闪着光"
    },
    {
        "name": "超人披风",
        "condition": "累计打卡21天",
        "description": "所以你真的用这个披风去拯救世界了？我以为你只是说说而已，没想到你真的做到了！",
        "icon_suggestion": "害羞的超人披风",
        "image_suggestion": "傲娇表情的披风，一边飘着一边假装不在意"
    },
    {
        "name": "不倒翁餐盘",
        "condition": "在打卡中断后重新打卡3天",
        "description": "哼～居然中断了打卡？真是让人失望！不过，看你又重新振作起来，倒是有点像我奶奶的不倒翁呢～给你这个不倒翁餐盘，算是对你坚持不懈的奖励吧！",
        "icon_suggestion": "翻白眼的不倒翁",
        "image_suggestion": "傲娇脸的不倒翁餐盘，一边晃一边嫌弃脸"
    },
    {
        "name": "太阳花",
        "condition": "记录一次早餐，且时间在上午9点前",
        "description": "噫～居然能在9点前起床？你是外星人吗？这朵太阳花送你了，它可是睡眼朦胧中做的，不要嫌弃它造型奇怪哦～对了，早餐记得吃点好的，你那些奇怪料理我都看不下去了！",
        "icon_suggestion": "困得要睡着的太阳花",
        "image_suggestion": "一朵打哈欠的太阳花，眼睛半睁半闭，看起来很困但很倔强"
    },
    {
        "name": "奇怪的鸡蛋",
        "condition": "记录一次高蛋白（超过20克）的餐食",
        "description": "20克蛋白质？你是要变成肌肉怪物吗？...虽然我也没说这样不好啦！这个蛋白质小人就送你吧，它比你强壮多了，或许能教你两招～",
        "icon_suggestion": "秀肌肉的鸡蛋人",
        "image_suggestion": "一个由鸡蛋变成的肌肉人，傲慢地秀着二头肌，同时偷偷看着领奖人"
    },
    {
        "name": "大小适中的碗",
        "condition": "记录一次热量适中（400-600卡）且三大营养素比例均衡的正餐",
        "description": "哎哟喂～看不出来你还挺会吃的嘛！这顿饭营养均衡、热量刚好，碗里放的居然都是人吃的东西！给你这个特制笑脸碗，别搞错了，我才不是夸你有品位，只是...刚好合格而已啦！",
        "icon_suggestion": "露出傲娇表情的饭碗",
        "image_suggestion": "一个碗，内部食物排列成傲娇表情，一边说'哼'一边脸红"
    },
    {
        "name": "几个糖分子",
        "condition": "记录一次低糖（碳水化合物少于30克）的餐食",
        "description": "什么？你居然能抵抗糖分的诱惑？啧啧，这不像正常人类能做到的事～给你这枚勇士勋章吧，从糖分大军中逃脱的英雄！...其实我也想尝试低糖饮食啦，只是暂时还没做好准备而已！",
        "icon_suggestion": "逃跑的糖分子",
        "image_suggestion": "一群拟人化的糖分子在追赶一个得意逃跑的勇士糖，逃兵背上背着'<30g'的旗帜"
    },
    {
        "name": "美食记录本",
        "condition": "连续3天每天记录2餐以上",
        "description": "哟～功课做得挺勤快嘛！连续三天记录这么多餐，你是要成为美食作家吗？给，这本日记本送你了～记得把你那些黑暗料理也好好记上，或许未来能出本《奇怪食谱大全》呢～",
        "icon_suggestion": "一脸嫌弃的日记本",
        "image_suggestion": "一本翻开的日记本，上面记录的食物照片旁边都有嫌弃的小评语和傲娇表情"
    },
    {
        "name": "正常的天秤",
        "condition": "一周内有3次餐食的三大营养素比例符合推荐标准",
        "description": "哎哟，看来你的天秤座属性暴露无遗啊～蛋白质、脂肪、碳水三兄弟被你安排得明明白白。给你这个天秤，我才不是因为你做得好才给的呢，只是...勉强算及格了吧！下次争取七次都平衡如何？",
        "icon_suggestion": "歪着头傲娇表情的天秤",
        "image_suggestion": "一个拟人化的天秤，三个盘子上放着不同营养素，天秤一边保持平衡一边傲娇地说'哼，刚好而已'"
    },
    {
        "name": "啄木鸟闹钟",
        "condition": "一周内记录至少4次早餐",
        "description": "不得不说，你这种能连续四天不赖床的生物还真是稀有～这枚勋章是给早起鸟专门定制的！喂，别误会，我这是在表扬你啦...",
        "icon_suggestion": "啄木鸟闹钟",
        "image_suggestion": "一个金色啄木鸟闹钟，背面刻着'早起的鸟儿也不一定有虫吃'"
    },
    {
        "name": "安眠小夜灯",
        "condition": "连续7天没有记录晚上9点后的进食",
        "description": "什么？你七天没有午夜觅食？是失恋了吗？还是冰箱被你吃空了？...咳咳，其实这样挺好的啦！你的胃应该感谢我，要不是我做的这个小夜灯提醒你，你肯定半夜还在跟零食约会！谢谢不用说，哼～",
        "icon_suggestion": "月亮下打呼噜的小夜灯",
        "image_suggestion": "一个月亮形状的小夜灯，旁边写着'别吵，朕在休息'"
    },
    {
        "name": "美食探险背包",
        "condition": "记录20种从未记录过的新食物",
        "description": "什么什么什么？你居然尝试了20种新食物？我正想问你是不是饿昏头了...不过，这种冒险精神，还挺...特别的！这个探险背包给你，下次可以装更多奇怪食物，说不定能发现什么不得了的黑暗料理呢～",
        "icon_suggestion": "被食物塞满显得很惊讶的背包",
        "image_suggestion": "一个探险背包，里面探出各种奇怪食物，背包表情很惊讶，写着'我的主人疯了，救命！'"
    }
]
token = login_admin()
# /home/devbox/project/backend/storage/apps/android/android_1.0.9.apk
v = achievements[7]
ITEM_IMAGE_PATH = f"G:/cc/HC2/HC/doc/item/pic/{v['name']}.png"
ITEM_ICON_PATH = f"G:/cc/HC2/HC/doc/item/icon/{v['name']}.png"
ITEM_NAME = v['name']
ITEM_DESCRIPTION = v['description'] + "【" + v['condition'] + "】"
ITEM_SOURCE = "饭团猫"
image_url = upload_image(ITEM_IMAGE_PATH, "item/pic", token)
icon_url = upload_image(ITEM_ICON_PATH, "item/icon", token)
item_data = {
        "name": ITEM_NAME,
        "description": ITEM_DESCRIPTION,
        "source": ITEM_SOURCE,
        "icon_url": icon_url,
        "image_url": image_url
    }
print(item_data)
result = create_item(token, item_data)
print("\n完整的物品信息:")
print(json.dumps(result, ensure_ascii=False, indent=2))

# for v in achievements:
#     ITEM_IMAGE_PATH = f"G:/cc/HC2/HC/doc/item/pic/{v['name']}.png"
#     ITEM_ICON_PATH = f"G:/cc/HC2/HC/doc/item/icon/{v['name']}.png"
#     ITEM_NAME = v['name']
#     ITEM_DESCRIPTION = v['description'] + "【" + v['condition'] + "】"
#     ITEM_SOURCE = "饭团猫"
#     image_url = upload_image(ITEM_IMAGE_PATH, "item/pic", token)
#     icon_url = upload_image(ITEM_ICON_PATH, "item/icon", token)
#     item_data = {
#             "name": ITEM_NAME,
#             "description": ITEM_DESCRIPTION,
#             "source": ITEM_SOURCE,
#             "icon_url": icon_url,
#             "image_url": image_url
#         }
#     print(item_data)
#     result = create_item(token, item_data)
#     print("\n完整的物品信息:")
#     print(json.dumps(result, ensure_ascii=False, indent=2))
