# How to start and stop a simulation

1. Log into [https://staging-app.ignitionrobotics.org](https://staging-app.ignitionrobotics.org)

2. Go to your [settings](https://staging-app.ignitionrobotics.org/settings),
   and select the `Access Tokens` tab.

3. Create a new token, and record it someplace safe. Use this token in
   subsequent REST calls where `YOUR_TOKEN` is specified. You can always
   create more tokens.   

4. Send a POST request with the docker image, and name. Docker image arguments are optional.

```
curl -X POST -H "Private-Token: YOUR_TOKEN" https://staging-cloudsim-nps.ignitionrobotics.org/1.0/start -F "image=DOCKER_IMAGE" -F "name=A_NAME" -F "args=OPTIONAL_ARGUMENTS"
```

Example:

```
curl -X POST -H "Private-Token: YOUR_TOKEN" https://staging-cloudsim-nps.ignitionrobotics.org/1.0/start -F "image=tfoote/test_novnc:latest" -F "name=npstest"
```

5. The POST command will return a JSON object that contains a "groupid". Use
   the groupid in subsequent REST calls in order to refer to a specific
   simulation instance.

6. Wait for the "running" status by periodically calling

```
curl -X GET -H "Private-Token: YOUR_TOKEN" https://staging-cloudsim-nps.ignitionrobotics.org/1.0/simulations/GROUP_ID
```

7. Copy the `uri` value from the JSON object returned in the previous step into
   a browser.

8. Stop the simulation using:

```
curl -X POST -H "Private-Token: YOUR_TOKEN" https://staging-cloudsim-nps.ignitionrobotics.org/1.0/stop/GROUP_ID`
```


## Filtering simulations

Query parameters can be passed to the `/simulations` route in order to
filter results. Each query parameter must be specified using
`q=**QUALIFIER**`, and can be combined with `&`.

When using `curl` it's important to put the URL with search parameters in
quotes. See the following example

```
curl -X GET -H "Private-Token: *TOKEN*" "https://staging-cloudsim-nps.ignitionrobotics.org/1.0/simulations?q=status:stopped&q=name:npstest"
```

### Filter by status

| Qualifier | Description|
|-----------|------------|
|`status:launching`| Return simulations that are in the process of launching.|
|`status:running`| Return running simulations.|
|`status:stopping`| Return simulations that are in the process of stopping.|
|`status:stopped`| Return stopped simulations.|


### Filter by name

| Qualifier | Description|
|-----------|------------|
|`name:*NAME*`| Return simulations with the given name.|

### Filter by group id

| Qualifier | Description|
|-----------|------------|
|`groupid:*GROUPID*`| Return simulations with the given group id.|