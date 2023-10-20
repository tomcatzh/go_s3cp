# go_s3cp

## Description

`go_s3cp` is a high-performance command-line utility written in Go for copying files from Amazon S3. It downloads S3 objects in parallel chunks to accelerate the download speed. The utility uses the AWS SDK for Go V2.

## Features

- Downloads S3 objects in parallel chunks for faster downloads.
- Automatically determines the local file path based on user input.
- Uses AWS SDK for Go V2 for efficient and reliable S3 operations.

## Performance

Based on internal tests focusing on the file "sd_xl_base_1.0.safetensors" (approximately 6.94 GB), `go_s3cp` has shown significant performance advantages over traditional methods like `aws s3 cp`. The tests were conducted on AWS EC2 instances with local NVMe SSD storage.

### Average Download Time for "sd_xl_base_1.0.safetensors" (seconds)

| Environment      | go_s3cp  | aws s3 cp | Improvement (%) |
| ---------------- | -------- | --------- | --------------- |
| ec2 g5.2xlarge   | 15.21    | 31.89     | 52.3            |
| ec2 i4i.2xlarge  | 8.31     | 20.65     | 59.78           |

## Installation

To install `go_s3cp`, you can clone the repository and build it using Go.

```bash
git clone https://github.com/tomcatzh/go_s3cp.git
cd go_s3cp
go build
