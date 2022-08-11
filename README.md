# Harbor Scan Report Action
One of fantastic options of [Harbor](https://goharbor.io/) is integrations with different vulnerability scanners. 
They provide list of vulnerabilities found in image. 

This action gets this scan report and makes comment with scan results for given docker image.

Currently, it can comment Pull Requests (PR) and issues.

This action consists from 2 parts:
1. Getting scan report
2. Generating GitHub comment (optional)

## Why it might be useful?

* Your image is stored at [Harbor](https://goharbor.io/)
* Your might want to prevent unsecure images from been deployed

## Examples
* Clean image
![CleanImage](clean-image.png?raw=true)
* Vulnerable image
![VulnerableImage](vulnerable-image.png)

## Configuration
* Minimal valid example
```yaml
- name: Run Report
        uses: kyberorg/harbor-scan-report@v0.1
        with:
          harbor-host: my_harbor.tld
          image: my_harbor.tld/hub/redhat/ubi8:latest
```
* Full example
```yaml
- name: Run Report
        uses: kyberorg/harbor-scan-report@v0.1
        with:
          harbor-host: my_harbor.tld
          harbor-proto: http
          harbor-port: 8080
          harbor-robot: ${{ secrets.HARBOR_ROBOT }}
          harbor-token: ${{ secrets.HARBOR_TOKEN }}
          image: my_harbor.tld:8080/hub/redhat/ubi8:latest
          max-allowed-severity: high
          github-url: ${{ github.event.pull_request.comments_url }}
          github-token: ${{ secrets.GITHUB_TOKEN }}
```

## Inputs
### `harbor-host`
String with hostname, without protocol and port.

Required: `yes`

### `image`
Image to scan. Format: `registry.tld/project/repo:tag` or `project/repo:tag`. 

Tag is optional, if tag missing default tag `latest` will be used.  

Required: `yes`

### `harbor-robot`
Robot or Username to access Harbor with. Robot with limited privileges is preferred.

This parameter is optional, but without credentials action can access public repositories only. 

Required: `no`

### `harbor-token`
Token for Robot or password of User,  defined in `harbor-robot`.

This parameter is optional, but without credentials action can access public repositories only.

Required: `no`

### `max-allowed-severity`
Minimum Vulnerability severity after which action considered as failed. 

Valid values: `none`, `low`, `medium`, `high`, `critical`

Default value: `critical`

* `none` means zero-tolerance to any vulnerabilities i.e. action succeeds only if image hasn't any vulnerabilities.

* `critical` means that action never fails, even if image has critical vulnerabilities.

### `github-url`
GitHub API endpoint to use. Normally, you would use built-in var `github.event.issue.comments_url` for commenting issues
or `github.event.pull_request.comments_url` - for commenting pull requests.

If nothing defined - commenting mode is disabled.

Required: `no`
```yaml
github-url: ${{ github.event.pull_request.comments_url }}
```

### `github-token`
A GitHub personal access token used to comment on your behalf. 
Normally, it can is appended to action's secrets as `${{ secrets.GITHUB_TOKEN }}`
```yaml
github-token: ${{ secrets.GITHUB_TOKEN }}
```

### `harbor-proto`
Protocol of Harbor instance. Use it, if your Harbor instance can be accessed only by using `http`.

Valid values: `http` and `https`.

Default value is `https`. 

### `harbor-port`
Custom port of Harbor instance. Use it, if Harbor instance has custom port.

Any port within `0-65535` range.
