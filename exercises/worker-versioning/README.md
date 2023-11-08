# Exercise 2: Version the Change with Worker Versioning

During this exercise, you will

* Run a Workflow Execution and retrieve the Event History.
* Define Worker Build ID Version Sets and enable Versioning on your Worker.
* Make a change to your Workflow, and redeploy an updated Worker.
* Ensure that your Workflows have switched to using the new code path.
* Decommission an old Worker.

If you are running this Exercise using a local dev server, you will need to
enable Worker Versioning on server startup, as it is not enabled by default.
To do this, run `temporal server start dev` with additional parameters:

```shell
temporal server start-dev \
   --dynamic-config-value frontend.workerVersioningDataAPIs=true \
   --dynamic-config-value frontend.workerVersioningWorkflowAPIs=true \
   --dynamic-config-value worker.buildIdScavengerEnabled=true
```

You can also provide these parameters to a cluster in a YAML file. An
example is provided in `enable_worker_versioning.yaml`.

Make your changes to the code in the `practice` subdirectory (look for 
`TODO` comments that will guide you to where you should make changes to 
the code). If you need a hint or want to verify your changes, look at 
the complete version in the `solution` subdirectory.


## Part A: Run a Workflow to Completion

This exercise's `practice` directory contains a fully working
example to start, which does not yet use Worker Versioning.
You can run the Workflow without any changes to retrieve a copy of
the Event History.

1. Run `go run worker/main.go` in a terminal to start a Worker
2. Run `go run start/main.go` in another terminal. This will 
   process a pizza delivery order for a hypothetical customer
   whose address and details are in `start/main.go`.
3. Let this Workflow run to completion. Because this is an example
   and does not actually order any real pizzas, it should complete
   almost instantly.
4. As in Exercise 1, you will now download the history of this
   execution in JSON format. Open the Web UI (if you are running
   a local dev cluster, it will be running at http://localhost:8233),
   navigate to the detail page for this execution, and then click
   the **Download** button that appears on the right side of the page,
   just above the table showing the Event History. Save the file as
   `history_for_original_execution.json` in your `practice` directory.
   * NOTE: If you are running this exercise in GitPod, you may 
     be unable to download the file, due to the embedded browser
	 used in that environment. In this case, run the following 
	 command from the `practice`  directory `tctl wf show 
	 --workflow_id loan-processing-workflow-customer-a100 
	 --print_full > history_for_original_execution.json` to 
	 retrieve a copy. 

// change directories or terminating running workers?


## Part B: Assign a Build ID to your Worker and Task Queue

1. Edit the `worker/main.go` file to add a Build ID to your Worker
   and opt in to Worker Versioning. To do this, replace the empty
   `worker.Options{}` struct that's currently provided as an argument
   to `worker.New()` with a struct containing a `BuildID` and the
   `UseBuildIDForVersioning: true` parameter.
2. Edit the `start/main.go` file to call
   `client.UpdateWorkerBuildIdCompatibility()` before starting your
   Workflow.
3. Re-run the whole thing
4. See if it changed?
5. Promote to Default using CLI
6. Re-run again?



## Part C: Add Another Worker Using Version Sets

1. Add a more complicated compatible set and another worker
2. Do I need to keep a workflow running all this time? Does someone really love pizza?
3. Verify Polling via CLI (within 5 minutes): temporal task-queue describe --task-queue=MyTaskQueue --task-queue-type="workflow"
4. Click on the "Run Id" for your Workflow and then click on the Task Queue name to view active "Pollers" registered to handle these Tasks.


## Part D: Decommission Your Old Worker

1.  Verify whether Build ID is reachable: temporal task-queue get-build-id-reachability --build-id "2.0"
2.  Decommission old worker
3.  Does Event History get used at all here?


### This is the end of the exercise.

