##Class setup and teardown

**Pre-requistes setup:**

1. login with valid authenticator application code, should succeed
2. check home page status and title, should succeed

**Post teardown:**

1. logout, should succeed


##Test Case 1: Profile full name

**Purpose: Test profile full name is working as expected**

**Steps:**

1. add full name and validate first name and last name are mandatory, should succeed
2. edit full name and validate first name and last name are mandatory, should succeed


##Test Case 2: Profile email addresses

**Purpose: Test profile email address is working as expected**

**Steps:**

1. add email address and validate label and email are mandatory, should succeed
2. edit email address and validate label and email are mandatory, should succeed


##Test Case 3: Profile phone numbers

**Purpose: Test profile phone number is working as expected**

**Steps:**

1. add phone number and validate label and phone number are mandatory, should succeed
2. edit phone number and validate label and phone number are mandatory, should succeed
3. remove phone number, should succeed


##Test Case 4: Profile facebook

**Purpose: Test profile facebook is working as expected**

**Steps:**

1. add facebook using valid account, should succeed
2. remove facebook, should succeed


##Test Case 5: Profile github

**Purpose: Test profile github is working as expected**

**Steps:**

1. add github using valid account, should succeed
2. remove github, should succeed


##Test Case 6: Profile addresses

**Purpose: Test profile address is working as expected**

**Steps:**

1. add address and validate all fields except remarks are mandatory, should succeed
2. edit address and validate all fields except remarks are mandatory, should succeed
3. remove address, should succeed


##Test Case 7: Profile bank accounts

**Purpose: Test profile bank account is working as expected**

**Steps:**

1. add bank account and validate all fields are mandatory, should succeed
2. edit bank account and validate all fields are mandatory, should succeed
3. remove bank account, should succeed


##Test Case 8: Profile digital wallet

**Purpose: Test profile digital wallet is working as expected**

**Steps:**

1. add digital wallet and validate all fields are mandatory, should succeed
2. edit digital wallet and validate all fields are mandatory, should succeed
3. remove digital wallet, should succeed


##Test Case 9: Notifications page

**Purpose: Test notifications page is working as expected**

**Steps:**

1. if no notifications should show "No unhandled notifications", should succeed


##Test Case 10: Settings api key

**Purpose: Test settings api key is working as expected**

**Steps:**

1. add api key and validate all fields are mandatory, should succeed
2. edit api key and validate all fields are mandatory, should succeed
3. remove api key, should succeed
