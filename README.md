cfyndns
=======

A dynamic DNS client for Cloudflare.

```
$ cfdyndns -name dyn.example.com -email user@example.com -key 6a2d0f96158761d85f2334e9bf5b604212343
2015/11/26 21:30:42 Current external IP is 192.0.2.90
2015/11/26 21:30:44 Current DNS IP is 192.0.2.90
2015/11/26 21:30:45 Updated record dyn.example.com -> 192.0.2.90
...
```

Uses `curl http://icanhazip.com/` to figure out the external address by default, but any command can be used. See `-help`.

License
-------

MIT
