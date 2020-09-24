[Insomnia](https://insomnia.rest) is a REST client used to setup and send requests using a graphical interface. 

There is an environment configuration file included in this directory that contains 
configurations for commonly used requests. This configuration file also contains a `local` environment that can be 
duplicated to create environments which can be switched between to send requests to different backends. 

## Configuration

### Setup

Certificate validation will block outgoing requests to Fuel and Cloudsim backends, so it must be disabled: 

1. Go to `Application → Preferences`.
2. Uncheck `Validate certificates`. 

### Environments

Environments allow you to setup named variables that can then be inserted into requests. There is a base environment 
that contains a set of variables. Any variable defined in this environment can be inserted in a request by 
pressing `Ctrl` + `Spacebar`. In order to be able to quickly change between different backends (`local`, `integration`, 
etc.), other environments can be defined to replace some of the values of the base environment. These environment
 variables have priority over base environment variables.
 
To view environments, press `Ctrl` + `E`. The base environment contains :
 
 * `org` User organization. Requests will typically use this value to fill in organization/owner fields.
 * `username` Account username. Some requests require that this username is the same as the one in the Auth0 bearer 
 token.
 * `solution_image` Used to start new simulations. This is the image used for solution containers.
 * `sim` Simulation GroupID. Used in some requests to target a specific simulation.
 * `token` Auth0 Security Bearer Token. All requests are configured to include this token in the header, so 
 configuring is required to send requests.
 
 These variables can be defined for specific environments (e.g. if you have different usernames in each backend, 
 etc.). All environments should have the corresponding Fuel/Cloudsim URLs. Make sure to include the API version in 
 the URLs. Changing the active environment will make requests go to the appropriate backend.
 
 Specific environments should contain variables to set backend URLs:
 
 * `fuel_url` Fuel backend URL.
 * `cloudsim_url` Cloudsim backend URL.

#### Bearer Token

Both Fuel and Cloudsim backends require and Auth0 bearer token to authenticate requests. The `token` environment 
variable must be setup to be able to comply with the backend authentication protocol. To setup the `token` variable:

1. Login to your target environment.
    1. For integration/staging, login to [staging.subtchallenge.world](https://staging.subtchallenge.world).
    2. For production, login to [subtchallenge.world](https://subtchallenge.world).
2. On your browser, access your local storage and copy the contents of `token`. This is your bearer token.
3. In your insomnia client, access your environments with `Ctrl` + `E`.
4. Select your target environment, this will depend on the environment you logged in during step 1.
5. Set the `token` environment variable.

The bearer token is valid for 2 hours, so you may need to repeat this process to update the `token` environment 
variable once it expires. 
