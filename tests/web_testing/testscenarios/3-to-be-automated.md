##Test Case 1: URL redirection integration with another website
**Purpose: Test URL redirection integration with another website is working as expected**

**Steps:**

1. so you start from website
2. press login
3. user is redirected to itsyouonline
4. user fills in login / pwd
5. user is authentoicated and approves access to the organisation of the website
6. user is redirected to the website with token
7. website uses token to get access_token from itsyou.online
8. website uses access token to get jwt from itsyou.online
