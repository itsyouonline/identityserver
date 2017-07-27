##Test Case 1: Access the registry

**Purpose: Access the registry using client_id and secret**

**Steps:**

1. add an entry in the registry depending if the client_id/secret combination is a user or an organization api key.
- ** user ** : Add/update an entry in a user registry using client credentials
- ** organization ** : Add/update an entry in an organization registry using client credentials

2. List the registry entries of the user/organization

3. An anonymous client is created and the registry entry created in step 1 is requested and printed.

4. Delete the entry created in step 1
