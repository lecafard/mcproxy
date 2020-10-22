mcproxy
=======

start a minecraft server on-demand when people are playing to keep costs down. servers
will automatically shut down if unused.

costs
-----
you just need to pay for the time that the minecraft server is on + a micro or nano
instance to run the proxy.

howto
-----
1. compile using `go build .`
2. copy config-example.json to config.json and modify to your liking.
3. run `./mcproxy`

other setup
-----------
- install some plugin on the server to shut it down when noone's on.
- find documentation on how to start up VM instances on GCP/Azure/AWS/Whatever Cloud Provider

documentation
-------------
:shrug:
