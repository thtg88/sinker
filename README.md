# Sinker

Sinker is a Go executable that synchronises files from a given folder into an AWS S3 bucket.

**Please note**: this executable only watches for changed files if the executable is running. It will not update files if it was not executing when they have changed.

## Requirements

Go installed on your machine. On macOS you can install it via Homebrew with:

```bash
brew install go
```

Make sure to create a user in AWS IAM that has read and write access to your S3 bucket of choice. If you are uncertain what policy to give you can use the following JSON payload to give full rights (not only read/write) to your user on the bucket of choice:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "s3:*",
            "Resource": [
                "arn:aws:s3:::your-bucket-name-here",
                "arn:aws:s3:::your-bucket-name-here/*"
            ]
        }
    ]
}
```

## Setup

Clone the project on your machine:

```bash
git clone git@github.com:thtg88/sinker.git
```

Copy `.env.example` contents to a `.env` in the same directory.

Make sure to set the correct env variables:

- `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`: the access key ID and secret of an AWS user which has read and write access on your AWS S3 bucket
- `AWS_BUCKET`: the actual S3 bucket you want to write to
- `SINKER_BASE_PATH`: the base directory on your machine, which you want to keep synchronised

## Usage

```bash
go run cmd/sinker/main.go
```
