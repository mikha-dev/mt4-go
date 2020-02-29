#!/usr/bin/env python
# *_* coding=utf8 *_*

import time
import json
import requests

API_URL = "http://127.0.0.1:8081/api"
TOKEN = "8e3dac16-7d9f-48af-8266-c014d2b07dcb"

# Staging
#API_URL = "http://113.10.168.68:8081/api"
#TOKEN = "f6c0f0cc-7201-42fd-a934-f3a5c39e6a98"

# Online
#API_URL = "http://13.112.90.111:8081/api"
#TOKEN = "8e5a5278-8f57-4b30-894e-c23eaa1e2534"

OP_BUY = 0
OP_SELL = 1
OP_BUY_LIMIT = 2
OP_SELL_LIMIT = 3
OP_BUY_STOP = 4
OP_SELL_STOP = 5
OP_BALANCE = 6
OP_CREDIT = 7

def request_api(func, **kwargs):
    start = time.time()
    headers = {"AgentToken": TOKEN}
    data = json.dumps({"func": func, "args": kwargs})
    r = requests.post(API_URL, data=data, headers=headers)
    print "request_api:%s %s cost: %s" % (func, kwargs, (time.time() - start))
    print "resp code:", r.status_code, " body:", r.text

    if r.status_code != 200:
        raise Exception("Api request status:{0} error:{1}".format(
            r.status_code, r.text))

    ret = json.loads(r.text)
    if ret["result"] == False:
        raise Exception("Error: %s" % ret["msg"])

    return ret["data"]

"""
{
    "func": "page_trade",
    "args": {
        "page_index": 0,
        "page_size": 20,
        "filter": {}
    }
}
"""

def page_trade(page_index=0, page_size=20, **filter):
    """
    Parameter:
    filter: {
        "ticket": int, // 订单号
        "login": int, // 用户ID
        "symbol": string, // 产品
        "cmd": int, // 订单类型
        "open_time": [string, string], // 开仓时间 时间格式：["2016-12-11 23:00:01", "2016-12-12 23:00:01"]
        "close_time": [string, string], // 平仓时间 时间格式：["2016-12-11 23:00:01", "2016-12-12 23:00:01"]
    }
    Return:
    {
       "total": int(符合条件的数据总数)，
       "data": [{
           "ticket":814468, //int 订单号
           "login":2003, // int 用户ID
           "symbol":"GOLDx", // string 产品类型
           "digits":2,
           "cmd":1, // int 指令类型
           "volume":1000, // int volume / 1000 = lots(手数)
           "open_time":"2016-12-16T11:22:57Z", // 开仓时间
           "open_price":1132.77, // float64 开仓价格
           "close_price":1133.43, // float64 平仓价格
           "close_time":"2016-12-16T11:23:27Z", // 平仓时间 如等于 "1970-01-01T00:00:00Z" 表示仍未平仓
           "sl":0, // float64 stop loss
           "tp":0, // float64 take profit
           "expiration":"1970-01-01T00:00:00Z", // 过期时间 限价单过期时间
           "conv_rate1":1, // 外汇的转换率（可忽略）
           "conv_rate2":1, // 外汇的转换率（可忽略）
           "swaps":0,
           "profit":-660, // float64 利润
           "taxes":0,
           "comment":"abcf", // float64 备注
           "margin_rate":1,
           "timestamp":1481887407,
           "modify_time":"2016-12-16T17:55:27Z"
        }...]
    }
    """
    return request_api("page_trade", page_index=page_index,
                       page_size=page_size, filter=filter)

def page_user(page_index=0, page_size=20, **filter):
    """
    Parameter:
    filter: {
       "login": int // 用户ID
       "name": int // 用户名称
    }
    Return:
    {
        "total": int(符合条件的数据总数)，
        "data": [
            {
               "login":2003, // int 用户ID
               "name":"TEXT-TW3", // string 用户姓名
               "city":"rh-city2", // string 代理编号
               "phone":"",
               "email":"",
               "regdate":"2016-05-05T16:11:54Z", // string 注册时间
               "lastdate":"2016-12-13T10:10:15Z", // string 最后登录时间
               "balance":-18990.26, // float64 用户余额
               "credit":200000 // float64 用户信用
            }
        , ...]
    }
    """
    return request_api("page_user", page_index=page_index,
                       page_size=page_size, filter=filter)

def page_profit(page_index=0, page_size=20, **filter):
    """
    Parameter:
    filter: {
       "login": int // 用户ID
       "name": int // 用户名称
    }
    Return:
    {
        "total": int(符合条件的数据总数)，
        "data": [
            {
               "login":2003, // int 用户ID
               "history_profit": double ,
               "history_swaps": double ,
               "history_commission": double ,
               "hlistory_taxes": double ,
            }
        , ...]
    }
    """
    return request_api("page_profit", page_index=page_index,
                       page_size=page_size, filter=filter)
