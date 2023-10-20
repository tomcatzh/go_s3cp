# go_s3cp

## Description

`go_s3cp` is a high-performance command-line utility written in Go for copying files from Amazon S3. It downloads S3 objects in parallel chunks to accelerate the download speed. The utility uses the AWS SDK for Go V2.

## Features

- Downloads S3 objects in parallel chunks for faster downloads.
- Automatically determines the local file path based on user input.
- Uses AWS SDK for Go V2 for efficient and reliable S3 operations.

## Installation

To install `go_s3cp`, you can clone the repository and build it using Go.

```bash
git clone https://github.com/your-github-username/go_s3cp.git
cd go_s3cp
go build
