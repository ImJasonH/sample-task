# Sample Tekton Custom Task Controller

## One-Time Setup

1. Choose a name for your task. This is usually a single word, and should
   describe what the task does, or the service it integrates with. Replace the
   string `replaceme` everywhere in this repo with that word, all lowercase.

For instance, to call your new task `pineapple`:

```
find ./ -type f -exec sed -i 's/replaceme/pineapple/g' {} \;
```

1. Choose an `apiVersion` for your task. This should be unique across all
   tasks, so usually it's a domain name you control, and a version string like
   `v1alpha1`, `v2beta3`, or `v4`, e.g., `mydomain.net/v5beta6`. Replace the
   string `my.api-group.dev/v0alpha0` in `main.go` with your `apiVersion`.

1. Choose a `kind` name for your task. This will normally be the name you chose
   in step 1, capitalized (e.g., `Pineapple`). Replace the string `ReplaceMe`
   in `main.go` with your `kind`.

1. Build and deploy your new controller in your cluster, using
   [`ko`](https://github.com/google/ko):

```
ko apply -f config/
```

This will run your controller as a single-replica deployment,
called `pineapple-task-controller`, in the namespace `pineapple-task`,
authorized as the ServiceAccount `pineapple-task-controller`, with permissions
to watch and update Tekton `Run` objects.

1. See the logs for your controller:

```
kubectl logs $(kubectl get pods -n pineapple-task -ojsonpath={.items[0].metadata.name}) -n pineapple-task
```

You should see logs like:

```
{"level":"info","caller":"logging/config.go:111","msg":"Successfully created the logger."}
{"level":"info","caller":"logging/config.go:112","msg":"Logging level set to info"}
{"level":"info","caller":"logging/config.go:79","msg":"Fetch GitHub commit ID from kodata failed","error":"open /var/run/ko/HEAD: no such file or directory"}
{"level":"info","logger":"pineapple-task-controller","caller":"profiling/server.go:59","msg":"Profiling enabled: false"}
{"level":"info","logger":"pineapple-task-controller","caller":"sharedmain/main.go:197","msg":"pineapple-task-controller will not run in leader-elected mode"}
{"level":"info","logger":"pineapple-task-controller","caller":"sharedmain/main.go:175","msg":"Starting configuration manager..."}
{"level":"info","logger":"pineapple-task-controller","caller":"sharedmain/main.go:179","msg":"Starting informers..."}
{"level":"info","logger":"pineapple-task-controller","caller":"sharedmain/main.go:183","msg":"Starting controllers..."}
{"level":"info","logger":"pineapple-task-controller","caller":"controller/controller.go:367","msg":"Starting controller and workers"}
{"level":"info","logger":"pineapple-task-controller","caller":"controller/controller.go:377","msg":"Started workers"}
```

## Implement Your Controller

The `ReconcileKind` method in `main.go` will be called every time a `Run` that
references your specified `apiVersion` and `kind`. Any updates to the `Run`
that is passed in will be persisted and any downstream clients watching the
resource will be notified.

Normally, your reconciler should take some action in one of a few situations:

* When the `Run` being reconciled is newly created; that is, when it has no status. In this case, your reconciler should take some action (e.g., schedule a remote execution), then update the `Run`'s status to indicate that it's ongoing:

```
r.Status.SetConditions([]apis.Conditions{{
	Type:    apis.ConditionSucceeded,
	Status:  corev1.ConditionUnknown, // Ongoing
	Reason:  "MyReason",
	Message: "Human-readable explanation of the current status",
}})
```

The `Reason` (`"MyReason"`) is intended to be a machine-consumable camel-case
string describing one of a few possible states. The `Message` is intended to
describe the state of the `Run` to a human reader, in as much detail as
is reasonable.

After updating the status to signal that the `Run` is ongoing, it can be
useful to use the sample `enqueueAfter` method to schedule a reconciliation
check at some point in the future. This can be useful to poll the state of an
external resource and update the status of the `Run`.

* To signal that the `Run` has succeeded:

```
r.Status.SetConditions([]apis.Conditions{{
	Type:    apis.ConditionSucceeded,
	Status:  corev1.ConditionTrue, // Successful!
	Reason:  "YayItPassed",
	Message: "Human-readable explanation of the successful Run",
}})
```

As above, `Reason` should be a short camel-case string intended to describe
one of a handful of possible states, and `Message` should be a human-readable
sentence providing more detail.

* To signal that the `Run` failed:

```
r.Status.SetConditions([]apis.Conditions{{
	Type:    apis.ConditionSucceeded,
	Status:  corev1.ConditionFalse, // Failed
	Reason:  "BooItFailed",
	Message: "Human-readable explanation of the unsuccessful Run",
}})
```

This signals that the `Run` failed.

# TODO: Params
# TODO: UnnamedaTasks
# TODO: Timeout, cancellation
