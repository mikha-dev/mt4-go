import dealer_api
import unittest

RET_LOGIN = None
RET_TICKET = None
LOGIN = 2003

class Dealer_ApiTestCase(unittest.TestCase):

    def test_001_create_user(self):
        global RET_LOGIN
        login = dealer_api.create_account("TestUser", "abcd1234",
                                 "tangwanwan@qq.com", "12322321121")
        RET_LOGIN = login
        print "RET_LOGIN: ", RET_LOGIN

    def test_002_check_password(self):
        global RET_LOGIN
        dealer_api.check_password(RET_LOGIN, "abcd1234")

    def test_003_reset_password(self):
        global RET_LOGIN
        dealer_api.reset_password(RET_LOGIN, "abcd2345")

    def test_004_recheck_password(self):
        global RET_LOGIN
        dealer_api.check_password(RET_LOGIN, "abcd2345")

    def test_005_get_asset(self):
        global RET_LOGIN
        ret = dealer_api.get_asset(RET_LOGIN)

        print ret

    def test_006_open_trade(self):
        global LOGIN, RET_TICKET
        ret = dealer_api.open_trade(LOGIN, dealer_api.OP_BUY, "GOLDx", 10)

        print ret

        RET_TICKET = ret["ticket"]
        print "RET_TICKET", RET_TICKET

    def test_007_get_trade(self):
        global LOGIN, RET_TICKET
        ret = dealer_api.get_trade(LOGIN, RET_TICKET)

        print ret

    def test_008_get_user_trades(self):
        global LOGIN
        ret = dealer_api.get_user_trades(LOGIN)

    def test_009_get_trades(self):
        global LOGIN

    def test_010_modify_trade(self):
        global LOGIN, RET_TICKET
        dealer_api.modify_trade(LOGIN, RET_TICKET, 0, 1350)

    def test_011_get_trade(self):
        global LOGIN, RET_TICKET
        ret = dealer_api.get_trade(LOGIN, RET_TICKET)
        print ret

    def test_012_close_trade(self):
        global LOGIN, RET_TICKET
        ret = dealer_api.close_trade(LOGIN, RET_TICKET, 10)

    def test_013_get_trade(self):
        global LOGIN, RET_TICKET
        ret = dealer_api.get_trade(LOGIN, RET_TICKET)
        print ret

if __name__ == "__main__":
    unittest.main()
