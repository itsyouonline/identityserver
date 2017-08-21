import logging
import unittest
import time
#import os

from testconfig import config

from pytractor import webdriver
#from selenium.webdriver import FirefoxProfile
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.wait import WebDriverWait
from selenium.webdriver.common.by import By
from selenium.webdriver.common.action_chains import ActionChains

#from pyvirtualdisplay import Display

from testsuite.page_elements_xpath import login_page

class BaseTest(unittest.TestCase):
    def __init__(self, *args, **kwargs):
        super(BaseTest, self).__init__(*args, **kwargs)
        self.environment_url = config['main']['url']
        self.admin_username = config['main']['admin']
        self.admin_password = config['main']['passwd']
        self.browser = config['main']['browser']
        self.elements = login_page.elements.copy()

    def setUp(self):
        self._testID = self._testMethodName
        self._startTime = time.time()
        self._logger = logging.LoggerAdapter(logging.getLogger('itsyouonline_portal_testsuite'),
                                             {'testid': self.shortDescription() or self._testID})
        self.lg('Testcase %s Started at %s' % (self._testID, self._startTime))
#        display = Display(visible=0, size=(1600, 900))
#        display.start()
        self.set_browser()
        self.wait = WebDriverWait(self.driver, 30)
        self.driver.get(self.environment_url)
        self.driver.maximize_window()
        self.wait_until_element_located(self.elements["username_textbox"])

    def tearDown(self):
        """
        Environment cleanup and logs collection.
        """
        self.driver.quit()
        if hasattr(self, '_startTime'):
            executionTime = time.time() - self._startTime
        self.lg('Testcase %s ExecutionTime is %s sec.' % (self._testID, executionTime))

    def lg(self, msg):
        self._logger.info(msg)

    def login(self, username='', password=''):
        username = username or self.admin_username
        password = password or self.admin_password
        self.lg('Do login using username [%s] and passsword [%s]' % (username, password))
        self.set_text('username_textbox', username)
        self.set_text('password_textbox', password)
        self.click('login_button')
        self.lg('Login successfully using username [%s] and passsword [%s]' % (username, password))

    def logout(self):
        self.lg('Do logout')
        self.click('logout_button')
        self.lg('Logout done successfully')

    def set_browser(self):
        if self.browser == 'chrome':
            self.driver = webdriver.Chrome()
        elif self.browser == 'firefox':
#            fp = FirefoxProfile()
#            fp.set_preference("browser.download.folderList",2)
#            fp.set_preference("browser.download.manager.showWhenStarting",False)
#            fp.set_preference("browser.download.dir", os.path.expanduser("~")+"/Downloads/")
#            fp.set_preference("browser.helperApps.neverAsk.saveToDisk", "application/zip, application/octet-stream")
#            self.driver = webdriver.Firefox(firefox_profile=fp)
            self.driver = webdriver.Firefox()
        elif self.browser == 'ie':
            self.driver = webdriver.Ie()
        elif self.browser == 'opera':
            self.driver = webdriver.Opera()
        elif self.browser == 'safari':
            self.driver = webdriver.Safari
        else:
            raise AssertionError("Invalid broswer configuration [%s]" % self.browser)

    def element_is_enabled(self, element):
        return self.driver.find_element_by_xpath(self.elements[element]).is_enabled()

    def element_is_displayed(self, element):
        return self.driver.find_element_by_xpath(self.elements[element]).is_displayed()

    def element_background_color(self, element):
        return str(self.driver.find_element_by_xpath(self.elements[element])\
                   .value_of_css_property('background-color'))

    def wait_until_element_located(self, name):
        self.wait.until(EC.visibility_of_element_located((By.XPATH, name)))

    def wait_element(self, element):
        self.wait_until_element_located(self.elements[element])
        return True

    def wait_unti_element_clickable(self, name):
        self.wait.until(EC.element_to_be_clickable((By.XPATH, name)))

    def click(self, element):
        element = self.elements[element]
        self.wait_until_element_located(element)
        self.wait_unti_element_clickable(element)
        self.driver.find_element_by_xpath(element).click()
        time.sleep(1)

    def click_link(self, link):
        self.driver.get(link)

    def get_text(self, element):
        element = self.elements[element]
        self.wait_until_element_located(element)
        return self.driver.find_element_by_xpath(element).text

    def get_size(self,element):
        element = self.elements[element]
        self.wait_until_element_located(element)
        return self.driver.find_element_by_xpath(element).size

    def get_value(self, element):
        return self.get_attribute(element, "value")

    def element_is_readonly(self, element):
        return self.get_attribute(element, "readonly")

    def element_link(self, element):
        return self.get_attribute(element, "href")

    def get_attribute(self, element, attribute):
        element = self.elements[element]
        self.wait_until_element_located(element)
        return self.driver.find_element_by_xpath(element).get_attribute(attribute)

    def set_text(self, element, value):
        element = self.elements[element]
        self.wait_until_element_located(element)
        self.driver.find_element_by_xpath(element).clear()
        self.driver.find_element_by_xpath(element).send_keys(value)       

    def move_curser_to_element(self, element):
        element = self.elements[element]
        location = WebDriverWait(driver, 10).until(EC.visibility_of_element_located((By.XPATH, element)))
        ActionChains(self.driver).move_to_element(location).perform()
