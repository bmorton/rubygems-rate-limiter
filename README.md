# rubygems-rate-limiter

This is to repro the case outlined in bundler/bundler#5797 when Bundler receives
a 429 from index.rubygems.org and pulls every version of every gemspec individually.

This seems to occur when you get rate-limited somewhere in the middle of a bundle
install and you're using geminabox.  There may be other scenarios that this
occurs, but this scenario is reproducable.


## Steps

1. Run geminabox: `docker run -ti --rm -p 9292:9292 bmorton/geminabox`
1. Edit your /etc/hosts file (`sudo vim /etc/hosts`) and add:
```
127.0.0.1 rubygems.org index.rubygems.org
```
1. Pull down this repo: `go get github.com/bmorton/rubygems-rate-limiter`
1. `cd $GOPATH/src/github.com/bmorton/rubygems-rate-limiter`
1. Run the proxy: `go build && sudo ./rubygems-rate-limiter`
1. From the same directory, run bundler: `bundle install --deployment --verbose`

If you'd like to repro again, you'll need to restart the rate limiting proxy
since it'll just start rate limiting all requests after 5.
