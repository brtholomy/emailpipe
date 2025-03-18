This is my emailpipe. There are many like it, but this one is mine.

It parses a given RSS XML file, matches an `<item><link>` against a given `--slug`, assembles the HTML email using the provided template, and sends the email via the Buttondown API.

---

How to send to production:

```sh
./emailpipe --slug glossolalia --prod
```

To read from local RSS rather than the default bartholomy.ooo live URL:

```sh
./emailpipe --slug glossolalia --source "/home/bth/y/dev/bartholomy.ooo/public/posts/index.xml"
```
