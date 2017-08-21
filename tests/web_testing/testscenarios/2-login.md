##Test Case 1: login with wrong username
**Purpose: Test login with wrong username is working as expected**

**Steps:**

1. login with invalid username(longString/numeric/specialChar), should fail with proper error message


##Test Case 2: login with wrong password
**Purpose: Test login with wrong password is working as expected**

**Steps:**

1. login with invalid password(longString/numeric/specialChar), should fail with proper error message


##Test Case 3: login with wrong username and password
**Purpose: Test login with wrong password is working as expected**

**Steps:**

1. login with invalid username and password(longString/numeric/specialChar), should fail with proper error message


##Test Case 4: login with authentication using wrong application code
**Purpose: Test login with wrong password is working as expected**

**Steps:**

1. login with valid username and password, should succeed
2. choose authentication using application, should succeed
3. login with invalid wrong authentication code(longString/numeric/specialChar), should fail with proper error message


##Test Case 5: login with authentication using wrong sms code
**Purpose: Test login with wrong sms code is working as expected**

**Steps:**

1. login with valid username and password, should succeed
2. choose authentication using sms code, should succeed
3. login with invalid wrong sms code(longString/numeric/specialChar), should fail with proper error message
