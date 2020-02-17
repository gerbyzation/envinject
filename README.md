![](https://github.com/gerbyzation/envinject/workflows/Tests/badge.svg)

# Envinject

Run-time environment variable loading for static pages.

```html
<script>
  window.CONFIG = __ENV_INJECT__;
</script>
```

```
cat index.html | envinject --whitelist ENV,API_URL
```

```html
<script>
  window.CONFIG = { ENV: "production", API_URL: "http://api.my.app/" };
</script>
```

## Rational

According to [12 factor applicaitions](https://12factor.net/config) configuration should be stored in the environment.
For front-end applications this is not as easy though, as they only have access to environment variables if a build step is involved with for example webpack. The problem is that this approach requires rebuilding for each environment or configuration change.

A way to work around this is to have an API endpoint serving the configuration, but this increases the complexity considerably and hampers performance.
Another approach is to run a script on startup that generates a `env.js` or equivalent, based off the environment variables and loads this via a `script` tag in the html document, however this still adds an extra network request to retrieve the configuration. This project aims to improve on this solution by providing a utility script that injects the configuration directly into the html file on startup, requiring no extra network requests and making the configuration instantly available to the javascript application.
