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

6. Wait for the "Address has been acquired" status by periodically calling

```
curl -X GET -H "Private-Token: YOUR_TOKEN" https://staging-cloudsim-nps.ignitionrobotics.org/1.0/simulations/GROUP_ID
```

7. Copy the `uri` value from the JSON object returned in the previous step into
   a browser.

8. Stop the simulation using:

```
curl -X POST -H "Private-Token: YOUR_TOKEN" https://staging-cloudsim-nps.ignitionrobotics.org/1.0/stop/GROUP_ID`
```
