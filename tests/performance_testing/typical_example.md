## Typical example of performance testing for itsyou.online

** Add Test Plane
![JMeter1](1.png)

** Add Thread Group
![JMeter2](2.png)

** Add HTTP Request Defaults
![JMeter3](3.png)

** Add HTTP Cache Manager
![JMeter4](4.png)

** Add HTTP Cookie Manager
![JMeter5](5.png)

** Add HTTP Header Manager
![JMeter6](6.png)

** Add HTTP Request for home page
![JMeter7](7.png)

** Add HTTP Request for api login
![JMeter8](8.png)

** Add Regular Expression Extractor to get the access token after login
![JMeter9](9.png)

** Add Response Assertion to validate the response
![JMeter10](10.png)

** Add HTTP Request for api call (get info)
![JMeter11](11.png)

** Download this example from here:
[itsyouonline.jmx](itsyouonline.jmx)

** Performance results for 100 thread:
[100 thread](test_results/100_threads/README.md)
