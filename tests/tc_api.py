import time
import json
import requests

API_URL = "http://127.0.0.1:8082/api"
TOKEN = "8e3dac16-7d9f-48af-8266-c014d2b07dcb"

# Online
#API_URL = "http://13.112.90.111:8082/api"
#TOKEN = "8e5a5278-8f57-4b30-894e-c23eaa1e2534"

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

def set_trader(login):
    return request_api("set_trader", login=login)

def cancel_trader(login):
    return request_api("cancel_trader", login=login)

def force_cancel_trader(login):
    return request_api("force_cancel_trader", login=login)

def page_trader(page_index, page_size, **filter):
    return request_api("page_trader", page_index=page_index,
                       page_size=page_size, filter=filter)

def follow(trader, client, strategy, size, direction, exit):
    return request_api("follow", trader=trader, client=client,
                       strategy=strategy, size=size, direction=direction,
                       exit=exit)

def unfollow(trader, client):
    return request_api("unfollow", trader=trader, client=client)

def page_follow(page_index, page_size, **filter):
    return request_api("page_follow", page_index=page_index,
                       page_size=page_size, filter=filter)

def page_follow_order(page_index, page_size, **filter):
    return request_api("page_follow_order", page_index=page_index,
                       page_size=page_size, filter=filter)
