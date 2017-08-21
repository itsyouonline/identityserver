##Test Case 1: Adding new organization

**Purpose: Test adding new organization is working as expected**

**Steps:**

1. select organization from the top menu
2. click new organization and enter new organization name, verify its mandatory, should succeed
3. check that you will be under people section as owner, should succeed
4. invite someone to this organization, should succeed
5. create sub-organization, should succeed

##Test Case 2: Set organization api key

**Purpose: Test organization api key generation is working as expected**

**Steps:**

1. click on organization.
2. press setting and add organization domain name, should succeed
3. press add api access key, enter lable and callbackurl and verify the format, should succeed
4. check api key is generated, should succeed


##Test Case 3: Organization cobranding 

**Purpose: Test cobranding is working as expected**

**Steps:**

### Show an organization logo on the login/register screen

When you use the authorization code flow to authenticate your users using Itsyou.Online, you can provide a better user experience by showing your logo on the login page.

To set an organization logo go to the settings page of an organization:

![Organization Settings](OrganizationSettingsTab.png)

and add a logo to your organization:

![Set organization logo](SetOrganizationLogo.png)

When a user is asked to login, this logo is added to the login/register page:

![Branded login page](BrandedLoginPage.png)


##Test Case 4: URL redirection integration with another website

**Purpose: Test URL redirection integration with another website is working as expected**

**Steps:**

1. go to website integrated with itsyouonline
2. click the link to sign-in using itsyou.online
3. login using your itsyou.online account and validate that successfully redirect you to the website page with your login information


##Test Case 5: Sharing of your personal information

**Purpose: Test sharing of your personal information with another website is working as expected**

**Steps:**

1. go to website integrated with itsyouonline
2. click the link to sign-in using itsyou.online
3. login using your itsyou.online account and validate that successfully redirect you to the website page with your login information
4. check that you got the required personal information from your itsyou.online account
