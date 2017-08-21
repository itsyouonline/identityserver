from api_testing.utils import BaseTest


class JWTBasicTests(BaseTest):

    def setUp(self):
        super(JWTBasicTests, self).setUp()
        self.response = self.Client_1.jwt.GetScope(self.user_1)
