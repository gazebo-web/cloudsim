# Overview

This directory contains tools and information related to AWS ECR for SubT.

1. `subt_ecr_create.sh`: Creates an ECR repository and IAM user for a team.
2. `Dockerfile`: An example dockerfile that can be used to test the results
   from the `subt_ecr_create.sh` script.

# subt ecr create

The `subt_ecr_create.sh` should be run for any team that needs an ECR
repository. Any additional management must be done through the AWS console.

See the contents of the script for more information.

# Dockerfile

This dockerfile will let you act as a team, and push an image to the team's
ECR repository.

1. Build the image. We are passing in the access key id and secrete key in
   as build args. This is not secure, but it is convenient for testing.
   Don't share the resulting image.

```
docker build -t ecr-test --build-arg access_key_id=<ACCESS_KEY_ID> --build-arg secret_access_key=<SECRET_ACCESS_KEY> .
```

2. Run the image.

```
docker run -it -v /var/run/docker.sock:/var/run/docker.sock ecr-test:latest /bin/bash
```

3. Inside the container, log into docker

```
eval `aws ecr get-login --no-include-email`
```

4. Tag an image that you want to push to ECR. Use `docker image list` to see
   what images you have available.

```
docker tag <image_to_upload>:<tag> <tag_name_from_subt_ecr_create_script>:<new_tag>
```

An real example might look like:

```
docker tag hello-world:latest 200670743174.dkr.ecr.us-east-1.amazonaws.com/subt/testing:tunnel-circuit-1
```

5. Push the image.

```
docker push <tag_name_from_subt_ecr_create_script>:<keyword>
```

Continuing the example from above, the push command would be:

```
docker push 200670743174.dkr.ecr.us-east-1.amazonaws.com/subt/testing:tunnel-circuit-1
```

6. List of images.

```
aws ecr describe-images --repository-name <team-repository-name>
```

For example:

```
aws ecr describe-images --repository-name subt/testing 
```

7. Submit an image to the SubT Portal. A team will choose one of their
   images to submit during circuit events. The team should submit the
   "repositoryName" and the "imageTag".
