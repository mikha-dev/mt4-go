import time
import json
import requests

API_URL = "http://localhost:8080/api"
#API_URL = "http://58.64.189.209:1234/api"
#API_URL = "http://10.211.55.9:8080/api"
TOKEN = "8e3dac16-7d9f-48af-8266-c014d2b07dcb"

# Online
API_URL = "http://13.112.90.111:8080/api"
TOKEN = "8e5a5278-8f57-4b30-894e-c23eaa1e2534"

OP_BUY = 0
OP_SELL = 1
OP_BUY_LIMIT = 2
OP_SELL_LIMIT = 3
OP_BUY_STOP = 4
OP_SELL_STOP = 5

"""
Api request

POST API_URL
AgentToken: "your token"
{
    "func": "funcname"
    "args": {
       "argname": "argvalue",....
    }
}

Response:
{
    "result": true or false,
    "msg": "Error message",
    "data": ...
}

Object:

Asset:
{
    "login":        int,
    "balance":      float64,
    "credit":       float64,
    "margin":       float64,
    "free_margin":  float64,
    "margin_level":  float64
}

Trade(Unclosed Trades):
{
    "ticket":     int,
    "login":      int,
    "symbol":     string,
    "digits":     int,
    "cmd":        int,
    "volume":     int,
    "open_time":  int,
    "open_price": float64,
    "sl":         float64,
    "tp":         float64,
    "comment":    string,
    "expiration": int
}
"""

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
        raise Exception("Error: %s", ret["msg"])

    return ret["data"]

def open_trade(login, cmd, symbol, volume, **kwargs):
    """ Return ticket  """
    args = {
        "login": login,
        "cmd": cmd,
        "symbol": symbol,
        "volume": volume,
        "price": 0,
        "sl": 0,
        "tp": 0,
        "comment": "",
    }

    for k, v in kwargs.items():
        if k in args:
            args[k] = v

    return request_api("open_trade", **args)


def create_account(name, password, email, phone):
    """ Return login(int) """
    return request_api("create_account", name=name, password=password, email=email, phone=phone)

def check_password(login, password):
    """ Return bool """
    return request_api("check_password", login=login, password=password)

def reset_password(login, password):
    """ Return None """
    return request_api("reset_password", login=login, password=password)

def get_asset(login):
    """ Return Asset """

    return request_api("get_asset", login=login)

def close_trade(login, ticket, volume):
    """ Return None """
    return request_api("close_trade", login=login, ticket=ticket, volume=volume)

def modify_trade(login, ticket, sl, tp):
    """ Return None """
    return request_api("modify_trade", login=login, ticket=ticket, sl=sl, tp=tp)


# All belows trades means unclosed trades.
def get_trades():
    """ Return Trade Array """
    return request_api("get_trades")

def get_user_trades(login):
    """ Return Trade Array """
    return request_api("get_user_trades", login=login)

def get_trade(login, ticket):
    """ Return Trade or None """
    return request_api("get_trade", login=login, ticket=ticket)

def get_quote(symbol):
    """ Return Quote or None """
    return request_api("get_quote", symbol=symbol)

def get_quotes():
    """ Return Quote or None """
    return request_api("get_quotes")
