#!/bin/bash
#
# == Purpose ==
# This script will create an AWS ECR repository for a SubT team. An IAM 
# user is also created, along with a policy to control access to the team's ECR
# repository.
#
# == Usage ==
#
# Make sure you have proper AWS credentials established before running this
# script.
#
# $ ./subt_ecr_create.sh <team_name>
#
# The prefix "subt-" is prepended to the team name when creating an IAM user.
# The ECR repository will be called "subt/<team_name>".
#
# == AWS Management ==
#
# An admin can log into the AWS console and manage the IAM policy and ECS
# repositories

# Check for the existance of the team name on the command line
if [ $# -eq 0 ]; then
  echo "Specify the team name."
  echo "Usage:"
  echo "  ./subt_ecr_create.sh <unique_team_name>"
  exit
fi

# Make sure the team name is >2 in length
size=${#1}
if [ $size -lt 3 ]; then
  echo "Team name must be longer than 2 characters."
  exit
fi

# Check for the existance of your aws credentials
myAccessKeyId=`aws configure get aws_access_key_id`
mySecretAccessKey=`aws configure get aws_secret_access_key`
if [ -z "${myAccessKeyId}" ]; then
  echo "Your aws access key id is not set."
  exit
fi

if [ -z "${mySecretAccessKey}" ]; then
  echo "Your aws secret access key is not set."
  exit
fi

# This is the name of the ECR repository
ecrRepoName="subt/$1"

# This is the IAM username 
iAmName="subt-$1"

# Create a repository.
echo "Creating ECR repository"
aws ecr create-repository --repository-name $ecrRepoName

# This will make sure that you had valid AWS credentials. If not, exit early.
if [ $? -ne 0 ]; then
  echo "It looks like you have bad AWS credentials. Exiting"
  exit
fi
 
# Get the repo ARN.
repoArn=`aws ecr describe-repositories --repository-names $ecrRepoName | grep -oh 'arn:[a-zA-Z:0-9\/\_\-]*'`

# Get the repo ID.
repoId=`aws ecr describe-repositories --repository-name $ecrRepoName | grep -Po '"registryId":.*?".*?"' | cut -d: -f2 | sed 's/[ "\,]*//g'`

# Get the repo URI.
repoURI=`aws ecr describe-repositories --repository-name $ecrRepoName | grep -Po '"repositoryUri":.*?".*?"' | cut -d: -f2 | sed 's/[ "\,]*//g'`

# Create the IAM user
echo "Creating User"
aws iam create-user --user-name ${iAmName}

# Create an access key for the user, and get teh access key id and
# secret access key
userInfo=`aws iam create-access-key --user-name ${iAmName}`
userAccessKeyId=`echo ${userInfo} | grep -Po '"AccessKeyId":.*?".*?"' | cut -d: -f2 | sed 's/[ "\,]*//g'`
userSecretAccessKey=`echo ${userInfo} | grep -Po '"SecretAccessKey":.*?".*?"' | cut -d: -f2 | sed 's/[ "\,]*//g'`

# Create the main ECR policy for the user
echo "{
  \"Version\": \"2012-10-17\",
  \"Statement\":[{
    \"Effect\":\"Allow\",
    \"Action\":[
      \"ecr:BatchDeleteImage\",
      \"ecr:BatchCheckLayerAvailability\",
      \"ecr:GetDownloadUrlForLayer\",
      \"ecr:GetRepositoryPolicy\",
      \"ecr:DescribeRepositories\",
      \"ecr:ListImages\",
      \"ecr:DescribeImages\",
      \"ecr:BatchGetImage\",
      \"ecr:InitiateLayerUpload\",
      \"ecr:UploadLayerPart\",
      \"ecr:CompleteLayerUpload\",
      \"ecr:PutImage\"
    ],
    \"Resource\":\"${repoArn}\"
    }
  ]
}
" > /tmp/ecr_policy.json
 
echo "Setting ECR policies for the user"

# Attach the main ECR policy to the user
aws iam put-user-policy --user-name $iAmName --policy-name subt-$1-ecr-policy --policy-document file:///tmp/ecr_policy.json
 
# Create the authorization ECR policy for the user. This will allow the user
# to run: `aws ecr get-login --no-include-email`
echo "{
  \"Version\": \"2012-10-17\",
  \"Statement\":[{
    \"Effect\":\"Allow\",
    \"Action\":[
      \"ecr:GetAuthorizationToken\"
    ],
    \"Resource\":\"*\"
    }
  ]
}
" > /tmp/auth_policy.json

# Attach the ECR policy to the user
aws iam put-user-policy --user-name $iAmName --policy-name subt-$1-ecr-auth-policy --policy-document file:///tmp/auth_policy.json

# Create the team's S3 policy. This gives a team access to circuit log files.
echo "{
  \"Version\": \"2012-10-17\",
  \"Statement\": [{
    \"Sid\": \"AllowListingOfUserFolder\",
    \"Action\": [
      \"s3:ListBucket\"
    ],
    \"Effect\": \"Allow\",
    \"Resource\": [
      \"arn:aws:s3:::web-cloudsim-production-logs\"
    ],
    \"Condition\": {
      \"StringLike\": {
        \"s3:prefix\": [
          \"gz-logs/${1}/circuit_logs/tunnel/*\",
          \"gz-logs/${1}/circuit_logs/urban/*\",
          \"gz-logs/${1}/circuit_logs/cave/*\"
        ]
      }
    }
  },
  {
    \"Sid\": \"AllowGettingOfUserFolder\",
    \"Action\": [
      \"s3:GetObject\"
    ],
    \"Effect\": \"Allow\",
    \"Resource\": [
      \"arn:aws:s3:::web-cloudsim-production-logs/gz-logs/${1}/circuit_logs/tunnel/*\",
      \"arn:aws:s3:::web-cloudsim-production-logs/gz-logs/${1}/circuit_logs/urban/*\",
      \"arn:aws:s3:::web-cloudsim-production-logs/gz-logs/${1}/circuit_logs/cave/*\"
    ]
  }
]}
" > /tmp/s3_policy.json

# Attach the S3 policy to the user
aws iam put-user-policy --user-name $iAmName --policy-name subt-$1-s3-policy --policy-document file:///tmp/s3_policy.json

# Display useful information. This is the only time that you'll be able to
# access the secretkey. Make sure to store it someplace safe. If it becomes
# lost, we'll have to create new key for the user.
echo "Summary"
echo "======="
echo "User name:       ${iAmName}"
echo "AccessKeyId:     ${userAccessKeyId}"
echo "SecretAccessKey: ${userSecretAccessKey}"
echo "Region name:     us-east-1"
echo "Repository Name: ${ecrRepoName}"
echo "Repository ARN:  ${repoArn}"
echo "Repository ID:   ${repoId}"
echo "Repository URI:  ${repoURI}"

echo ""

# Instructions for the user
echo "How to create a docker image"
echo "============================"
echo "  1. Install AWS CLI"
echo "  2. Run: \"aws configure\", and use the above information"
echo "  3. Run: \"eval \`aws ecr get-login --no-include-email\`\""
echo "  4. Run: \"docker build -t <your_image> .\""
echo "  5. Run: \"docker tag <your_image> ${repoURI}:tunnel_circuit\""
echo "  6. Run: \"docker push ${repoURI}:tunnel_circuit\""
echo "  7. Submit to Portal: \"${repoURI}:tunnel_circuit\""

# == Appendix ==

# To List images in a user's repository
# aws ecr list-images --repository-name ${ecrRepoName}
