# 项目名称

ChatGPT对接微信公众号/订阅号自动回复

特点
- 实现微信超时回复的消息回复（非客服实现）
- 自动回复超出500字时，内置http详情页查看
- 自动发布到订阅号草稿中

微信个人订阅号/公众号对接OpenAI时，会面临以下问题
- 微信自动回复5秒超时问题
- 微信自动回复500字超长问题

5秒超时问题，主要做了以下处理，首次提问会提交到OpenAI,如果OpenAI在5秒回复，将正常返回给客户端。如果超过5秒，微信第一次重复将从数据库中获取结果
如仍未有结果，则等待第二次尝试，如果第二次仍失败，提示用户重试
PS：如果是服务号建议直接对接客服信息，体验会更好

500字超长问题，超出500字，加入详情页，发客户端发送链接

不足之处：由于中途加入mongodb做了一次缓冲，因此在进行对话时，没有上下文关联

演示：微信搜索"虾客君"
源代码： https://github.com/yiuked/wx-aichat
## 安装

1. 安装依赖

```bash
version: '3'

services:
  mongo:
    image: 'bitnami/mongodb:latest'
    container_name: mongo
    environment:
      - MONGODB_ROOT_PASSWORD=123456
    volumes:
      - ./data:/bitnami/mongodb
    restart: always
  ghapi:
    image: wxgpt:1.0.0
    container_name: wxgpt
    privileged: true
    restart: always
    environment:
      WX_APPID: ${WX_APPID}
      WX_SECRET: ${WX_SECRET}
      WX_TOKEN: ${WX_TOKEN}
      WX_MEDIA_ID: ${WX_MEDIA_ID}
      WX_OPEN_ID: ${WX_OPEN_ID}
      WX_AUTHOR: ${WX_AUTHOR}
      AI_TOKEN: ${AI_TOKEN}
      DATA_SOURCE: ${DATA_SOURCE}
      HOST: ${HOST}
    ports:
      - "80:8089"
    depends_on:
      - mongo
```

在目录下创建`.env`文件
```
WX_APPID:     # 必填，用于自动回复验证的token,微信开发者设置中获取
AI_TOKEN:     # 必填，OpenAI,APIkey
DATA_SOURCE:  # 必填，mongo 数据源
HOST:         # 必填，接口链接地址，如`http://example.com/
# 以下内容仅需要自动发布到公众号的文章草稿时需要填写
WX_SECRET:    # 非必填
WX_TOKEN:     # 非必填
WX_MEDIA_ID:  # 非必填，发布到草稿默认的图片ID
WX_OPEN_ID:   # 非必填，发布者在订阅号下的openid
WX_AUTHOR:    # 非必填，自定义发布者名称
```

## 贡献者

- yiuked