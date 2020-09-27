# Description

This is a tiny utility that opens [StackOverflow](https://stackoverflow.com) in a virtual headless browser and logs in there using provided credentials.

# Usage 

* run `go build ` to build an executable binary file
* run `main.exe --email=YOUR_EMAIL --password=YOUR_PASSWORD`

In case of successful login there will be an appropriate message in console.

In case of failure a file `output.html` will be generated with StackOverflow page content, and there will be a message with a failure reason in console.