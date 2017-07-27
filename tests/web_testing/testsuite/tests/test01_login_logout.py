from testsuite.utils import BaseTest


class LoginLogoutPortalTests(BaseTest):

    def test001_login_and_portal_title(self):
        """ PRTL-001
        *Test case for check user potal login and titles.*

        **Test Scenario:**

        #. check the login page title, should succeed
        #. do login using pre-created username/password, should succeed
        #. check the home page title, should succeed
        """
        self.lg('%s STARTED' % self._testID)
        self.lg('check the login page title, should succeed')
        self.assertEqual(self.driver.title, "It's You Online")
        self.lg('do login using admin username/password, should succeed')
        self.click('login_page_button')
        self.login()
        self.lg('check the home page title, should succeed')
        self.assertEqual(self.driver.title, "Home - It's You Online")
        self.lg('%s ENDED' % self._testID)