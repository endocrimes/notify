# notify
A tiny little Go proxy for sending Slack Notifications

## Usage

To run the proxy:
```
$ godeps go build
$ ./notify -from :4000 -to https://slack.com/webhook/url
```

Then to send notificatons:
```
$ curl "localhost:4000?message='Hello%20World'&service='Test'"
```
