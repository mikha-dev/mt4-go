import websocket
import thread
import time

TOKEN = "8e3dac16-7d9f-48af-8266-c014d2b07dcb"
QUOTE_URL = "ws://localhost:8080/ws/quote"

#QUOTE_URL = "ws://58.64.189.209:1234/ws/quote"


# Online
QUOTE_URL = "ws://13.112.90.111:8080/ws/quote"
TOKEN = "8e5a5278-8f57-4b30-894e-c23eaa1e2534"


def on_message(ws, message):
    print message

def on_error(ws, error):
    print "### on open ###"
    print error

def on_close(ws):
    print "### closed ###"

def on_open(ws):
    print "### on open ###"

if __name__ == "__main__":
    websocket.enableTrace(True)
    ws = websocket.WebSocketApp(
        QUOTE_URL, header = ["AgentToken: %s" % TOKEN],
        on_message = on_message,
        on_error = on_error,
        on_close = on_close
    )
    ws.on_open = on_open
    ws.run_forever()
