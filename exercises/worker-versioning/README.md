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


## Part B: Assign a Build ID to your Worker and Task Queue

1. Edit the `worker/main.go` file to add a Build ID to your Worker
   and opt in to Worker Versioning. To do this, replace the empty
   `worker.Options{}` struct that's currently provided as an argument
   to `worker.New()` with a struct containing a `BuildID` and the
   `UseBuildIDForVersioning: true` parameter.
2. Edit the `start/main.go` file to call
   `client.UpdateWorkerBuildIdCompatibility()` before starting your
   Workflow.
3. Run a new version of your worker by running `go run worker/main.go`
   in a new terminal window. You do not need to terminate your old
   worker. If you restart your workflow by running
   `go run start/main.go` again, you should see some output from the
   terminal window of your new worker. This indicates that the workflow
   was chosen and run by the new, versioned worker, and the unversioned
   worker was ignored. You can also verify this by visiting the Web UI
   again. Nothing will have changed in the business logic of your
   Workflow, but you can examine the `pizza-tasks` Task Queue or the
   "Workers" UI tab to see the registered Workers and their versions.
   You may notice that there is a "Retirability" field listed in the
   Workers table of your Task Queue, and only your new, versioned
   Worker is listed as not retireable.
4. To retrieve the same information programmatically from the CLI, run
   `temporal task-queue get-build-ids --task-queue pizza-tasks`. You
   should receive output similar to this:
   ```output
         BuildIds        DefaultForSet   IsDefaultSet
     [revision-yymmdd]  revision-yymmdd  true
   ```


## Part C: Decommission Your Old Worker

1. Now that you have a new, versioned Worker running, you probably
   aren't planning on submitting any more Workflows to your old,
   unversioned Worker. However, it's important to make sure that your
   old Worker isn't still handling any Workflows -- and that it has
   a ready replacement -- before decommissioning it. There are a
   few ways of verifying this, using the CLI and the Web UI. The
   "Retirability" field you saw in the last step is one.
2. Run `temporal task-queue describe --task-queue pizza-tasks`. Both
   of your running Workers should report in, as they are continuing to
   poll despite not receiving any more Workflows:
   ```output
       Identity     LastAccessTime  RatePerSecond
     30404@Omelas@  39 seconds ago         100000
     8692@Omelas@   40 seconds ago         100000
   ```
3. Can you get Retirability from the CLI?


## Part D: Add Another Worker Using Version Sets

1. Lastly, you'll experiment with Version Sets using the CLI rather
   than the SDK. In Part B, you used
   `client.UpdateWorkerBuildIdCompatibility()` to update your Task
   Queue with new a Worker Version. You can accomplish the same
   thing by running `temporal task-queue update-build-ids` with
   matching parameters.
2. Assume you are adding another new Worker Version that is compatible with the
   Version that you added in Part B. To do this, try running:

   ```shell
   temporal task-queue update-build-ids add-new-compatible \
      --build-id="revision-yymmdd+1" \
      --existing-compatible-build-id="revision-yymmdd" \
      --task-queue="pizza-tasks"
   ```

   This will add `revision-yymmdd+1` as another compatible Build ID to the same
	Task Queue. This is equivalent to running the following SDK code:

   ```go
   c.UpdateWorkerBuildIdCompatibility(context.Background(), &client.UpdateWorkerBuildIdCompatibilityOptions{
		TaskQueue: pizza.TaskQueueName,
		Operation: &client.BuildIDOpAddNewCompatibleVersion{
			BuildID:                   "revision-yymmdd+1",
			ExistingCompatibleBuildId: "revision-yymmdd",
		},
	})
   ```

   (The `pizza.TaskQueueName` variable is stored in `shared.go`.)

   In practice, you should find this at least as useful as updating your
   Task Queues from the SDK.
3. Lastly, you can now deploy a new Worker, identified by this new Build ID,
   and it will be able to process the same Workflows. Although you would
   ordinarily do this when making code changes, you can do it without changing
   anything for the sake of this example. Update the `BuildID` field from your
   `worker.Options{}` declaration to reflect the new `revision-yymmdd+1`
   Versiom and restart your Worker once again, then re-run your Workflow
   starter. You should observe that both Worker Build IDs are compatible and
   able to process the Workflow.


### This is the end of the exercise.
