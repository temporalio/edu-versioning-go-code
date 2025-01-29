# Code Repository for Temporal Versioning (Go)
This repository provides code used for exercises and demonstrations
included in the Go version of the 
[Versioning](https://learn.temporal.io/courses/versioning) 
training course.

It's important to remember that the example code used in this course was designed to support learning a specific aspect of Temporal, not to serve as a ready-to-use template for implementing a production system.

For the exercises, make sure to run `temporal server start-dev --ui-port 8080 --db-filename clusterdata.db` in one terminal to start the Temporal server. For more details on this command, please refer to the `Setting up a Local Development Environment` chapter in the course. Note: If you're using the Gitpod environment to run this exercise, you can skip this step.

## Hands-On Exercises

Directory Name                     | Exercise
:--------------------------------- | :-------------------------------------------------------
`exercises/version-workflow`       | [Exercise 1](exercises/version-workflow/README.md)


## Reference
The following links provide additional information that you may find helpful as you work through this course.
* [General Temporal Documentation](https://docs.temporal.io/)
* [Temporal Go SDK Documentation](https://pkg.go.dev/go.temporal.io/sdk)
* [Go Language Documentation](https://go.dev/doc/)
* [GitPod Documentation: Troubleshooting](https://www.gitpod.io/docs/troubleshooting)


## Exercise Environment for this Course
You can launch an exercise environment for this course in GitPod by 
clicking the button below:

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/temporalio/edu-versioning-go-code)

Alternatively, you can follow 
[these instructions](https://learn.temporal.io/getting_started/go/dev_environment/) to 
set up your own Temporal Cluster with Docker Compose, which you can use as an 
exercise environment.
