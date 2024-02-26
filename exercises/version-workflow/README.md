# Exercise 1: Version the Change with the `GetVersion` API

During this exercise, you will 

* Run a Workflow Execution that completes successfully 
* Make and deploy a change that does not affect compatibility
* Make and deploy a change that breaks compatibility, causing a non-deterministic error
* Develop an automated test to check compatibility with previous executions
* Use the `GetVersion` API to implement versioning for the Workflow

Make your changes to the code in the `practice` subdirectory (look for 
`TODO` comments that will guide you to where you should make changes to 
the code). If you need a hint or want to verify your changes, look at 
the complete version in the `solution` subdirectory.


## Part A: Run a Workflow to Completion

1. Run `go run worker/main.go` in a terminal to start a Worker
2. Run `go run start/main.go a100` in another terminal. This will 
   start a Workflow that processes the loan for customer ID `a100`.
3. Let this Workflow run to completion. This customer has a loan 
   with 10 payments, and since the Workflow in this exercise uses 
   a Timer to add a three-second delay between each payment, it 
   should complete within 30 seconds.
4. You will now download the history of this execution in JSON
   format so that you can replay it in an automated test that you
   will develop later in this exercise. Open the Web UI (if you are
   running a local dev cluster, it will be running at
   http://localhost:8233), navigate to the **Event History** section
   of the detail page for this execution, and then click the
   **Download** button just above the table showing the Event History.
   In the **Download JSON** dialog, disable the
   **Decode Event History** toggle, and then click **Download**. Save
   the file as `history_for_original_execution.json` in your
   `practice` directory.
   * NOTE: If you are running this exercise in GitPod, you may 
     be unable to download the file, due to the embedded browser
	 used in that environment. In this case, run the following 
	 command from the `practice`  directory:
    
    ```shell
    temporal workflow show \
      --workflow-id loan-processing-workflow-customer-a100 \
      > history_for_original_execution.json
   ```
      
   to retrieve a copy. 
5. In the next section, you will make and deploy an incompatible 
   change, causing a non-deterministic error for an open execution.
   To allow time for you to do these things, edit the `workflow.go` 
   file and change the duration in the `workflow.Sleep` call from 
   3 seconds to 90 seconds.
6. Save your change to the `workflow.go` file and exit the editor
7. Restart the Worker by pressing Ctrl-C in the terminal window
   from step 1 and running the `go run worker/main.go` command again
8. Run the Workflow again: `go run start/main.go a100`
9. Use the Web UI to verify that the Workflow Execution from the 
   previous step is running before proceeding with the next part
   of this exercise.


## Part B: Deploy an Incompatible Change (without Versioning)

1. This Workflow uses the `SendThankYouToCustomer` Activity to 
   send a thank you message to the customer before charging 
   them with the first loan payment, but this was a mistake.
   This Activity should run after the last payment. To fix this, 
   edit the `workflow.go` file and move the five lines of code 
   (which begin with the `var notifyConfirmation string` statement)
   related to that Activity from just before the loop to just
   after it.
2. Save your change and exit the editor.
3. Restart the Worker by pressing Ctrl-C in the terminal 
   window where you started it and then running the 
   `go run worker/main.go` command again. Although the Worker
   immediately begins using the updated code after this restart, 
   the effects of your change to the Workflow logic may not be
   evident for up to 90 seconds due to the duration of the Timer.
4. Refresh the detail page for this execution in the Web UI. 
   Continue to refresh the page until the non-deterministic
   error is visible.

The non-deterministic error occurs because of your change to the 
Workflow logic. By moving the Activity from before the loop to after
it, the sequence of Commands generated during execution is different 
with the new code than it was prior to the change. 

Recall that you had an open Workflow Execution when you restarted the 
Worker during the deployment. The Worker used History Replay to 
recover the state of the open execution prior to the restart. Since 
the Commands generated when replaying it with the new code did not 
correspond to the Events that were generated when the Worker ran the 
original code before the restart, it is unable to recover the state 
and responds with the non-deterministic error you see.


## Part C: Use the Workflow Replayer to Test Compatibility

1. Edit the `workflow_test.go` file and uncomment the two import 
   statements near the top of the file, then implement the following
   in the `TestReplayWorkflowHistoryFromFile` function:
   * Create the Workflow Replayer
   * Register the Workflow Definition with the Replayer
   * Replay the Event History in the JSON file you downloaded
   * Use `assert.NoError` to verify that replaying the history 
     does not return an error
2. Save your changes
3. Run `go test`. You should find that this fails, which confirms 
   altering the execution order of the `SendThankYouToCustomer` 
   Activity) breaks compatibility. In the final part of this 
   exercise, you will use the `GetVersion` API to implement 
   versioning for your change, thereby making it compatible 
   with Workflow Executions started before or after the change.

## Part D: Version the Change with the `GetVersion` API

Just above the loop, where the `ExecuteActivity` call was prior to 
the change, add the following line:

```go
version := workflow.GetVersion(ctx, "MovedThankYouAfterLoop", workflow.DefaultVersion, 1)
```

This establishes a logical branch for code execution, identified 
by the user-defined Change ID `MovedThankYouAfterLoop`. Since there 
was no versioning in place prior to this change, the minimum supported 
version is `workflow.DefaultVersion` and the maximum supported version 
is `1`.

1. Add a conditional statement just after this new line: If the value
   of `version` is equal to `workflow.DefaultVersion`, meaning that it
   represents a Workflow Execution started when the Activity was called
   before the loop, then invoke `ExecuteActivity` call there.  In other
   words, copy the same lines you moved after the loop to inside the
   braces for this conditional statement, so that this Activity will be
   called if the condition evaluates to `true`.
3. Wrap the code you previously moved after the loop in a
   conditional statement that tests if `version` is equal to
   `1`. This will handle the Activity for Workflow
   Executions started after the change.
4. Change the duration of the `workflow.Sleep` statement at the
   bottom of the loop back to 3 seconds. This is unrelated to
   versioning, but will help you see the results more quickly.
5. Run `go test` again. You should find it succeeds this time,
   since you've used the `GetVersion` API to restore compatibility with
   the previous execution.
6. Restart the Worker by pressing Ctrl-C in the terminal
   window where you started it and then running the `go run
   worker/main.go` command again. You should find that the Workflow
   Execution completes successfully within the next 30 seconds.
   


### This is the end of the exercise.

