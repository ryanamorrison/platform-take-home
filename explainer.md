## Explainer

### Complaints:
> “The process in manual and error prone”
CI and CD/GitOps will solve this problem.  The CI process for building and packaging is automated and repeatable.  Without some of the configuration details from ArgoCD it would have been difficult to have created an end-to-end working example.  However, once CI has been run, and the code is pushed to a target branch such as *main*, *test*, *staging* or *prod*, ArgoCD will notice changes for the repo and branch it is configured for and reconcile those changes via the provided `kustomize` folder.  The advantage is that once ArgoCD is configured via its own git repo, engineers can fully deploy without having to engage SRE/DevOps/infra engineers.
> “Whenever a deployment happens, we suffer a tiny bit of downtime due to the server being down.”
Because Kubernetes is being used, this is solved using the `RollingUpdate` deployment strategy.  The settings provided in `kustomize/deployment.yaml` were guessed at for this application but they can be fine-tuned based on the deployment.  Usually this works best when used with Readiness and Liveness probes also provided.
> “there's no standardization for code formatting which leads to inconsistencies”
This was solved by adding a custom pre-commit hook (`.git/hooks/pre-commit`) which will execute `go fmt` on every commit.  A custom script is used rather than the popular pre-commit.com because there is another script that is being run on every commit also…
> “the proto-gen script is ran locally which leads to developers forgetting to do it before pushing code upstream”
Also solved with pre-commit since it does not seem to take that long to run.
> “downloading the tooling dependencies is a manual, undocumented process”
This is solved with a CI workflow and GitHub Actions.  The workflow will download dependencies, vet the code, lint the code, test the code, build the code and package it into a container.  Static testing was attempted but there are some issues with grpc and the golang tools so that job stage was commented out for further investigation.  Some of the exported golang functions and objects needed comments for the linter to approve them so they have been added.  Additionally, one test for the `store/` was added so that the testing stage had something to test.  
> “running the server locally is a bit of a pain,”
Yes, but Kubernetes is available and EKS assumed.  Testing can be done via a separate cluster if budget allows as the isolation provides an additional guardrail.  Alternatively, separate namespaces can be used in the same cluster.  
> “since it requires manually standing up a Postgres database (to replicate prod)”
The database could be solved a couple of different ways.  One would be to use the regular postgres container image and provide a migration.  The other is to simply pre-bake a test dataset into a postgres container.  The latter is what was chosen here.  Which is the most optimal really depends on the environment.  There are ways to have separate containers in a mono repo but the results looked messy unless one uses an external repo like Dockerhub (as opposed to the built-in one that is free with the GitHub repo).  The only drawback with the latter is that by using the native GH container registry, sensitive login information is hidden.  Using an external repo requires some additional configuration.
> “there's no easy way to share a feature update with our colleagues, since we only have a local and production environment.”
Since Kubernetes is being used but there was no way to access the environment for this take-home test, this particular point isn’t fully addressed in code.  What is assumed here is that ArgoCD would be configured to publish new features to a test deployment in an EKS cluster and that those changes would be viewable to others.
### Deliverables
> “In the README, please provide instructions on how to run your solution (whether locally or in cloud).”
Configure ArgoCD to watch a specific branch on the following repo and push the code to that branch.
https://github.com/ryanamorrison/platform-take-home/pkgs/container/platform-take-home
Nothing has been pushed to `main` yet so that branch could be used for ArgoCD.
The additional postgres repo, referenced in the `deployment.yaml` is here:
https://github.com/ryanamorrison/postgres-image

### Additional Note on Testing:
Without having some prior experience with gRPC, there wasn’t really a way to _completely_ test this.  The database image and go-lang application images were tested with basic `curl` and `psql` commands.  Kustomize was also tested to ensure it would work if ArgoCD were pointed at the repository.
```
ryan@enc1-dev-01:~/repos/go/src/github.com/skip-mev/platform-take-home$ kubectl apply -k kustomize/.
service/platform-take-home created
deployment.apps/platform-take-home created
ryan@enc1-dev-01:~/repos/go/src/github.com/skip-mev/platform-take-home$ kubectl get all -n platform-take-home-test
NAME                                      READY   STATUS    RESTARTS   AGE
pod/platform-take-home-6bcbd874cf-5dlxt   2/2     Running   0          54s
pod/platform-take-home-6bcbd874cf-d4dht   2/2     Running   0          54s
pod/platform-take-home-6bcbd874cf-ssz57   2/2     Running   0          54s

NAME                         TYPE           CLUSTER-IP    EXTERNAL-IP   PORT(S)                                        AGE
service/platform-take-home   LoadBalancer   10.8.44.177   10.8.48.7     9008:31201/TCP,8080:31068/TCP,8081:30017/TCP   54s

NAME                                 READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/platform-take-home   3/3     3            3           54s

NAME                                            DESIRED   CURRENT   READY   AGE
replicaset.apps/platform-take-home-6bcbd874cf   3         3         3       54s
ryan@enc1-dev-01:~/repos/go/src/github.com/skip-mev/platform-take-home$ kubectl exec -it -n platform-take-home-test platform-take-home-6bcbd874cf-5dlxt -- bash -c "psql -U postgres -c 'SELECT * FROM items;'"
Defaulted container "postgres" out of: postgres, platform-take-home
 id | created_at | updated_at | deleted_at | name | description
----+------------+------------+------------+------+-------------
  1 | 2024-11-19 | 2024-11-19 |            | test | ttt
(1 row)

ryan@enc1-dev-01:~/repos/go/src/github.com/skip-mev/platform-take-home$ curl 10.8.48.7:8081/health
# HELP go_gc_duration_seconds A summary of the wall-time pause (stop-the-world) duration in garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0.000264102
go_gc_duration_seconds{quantile="0.25"} 0.000264102
go_gc_duration_seconds{quantile="0.5"} 0.000501489
go_gc_duration_seconds{quantile="0.75"} 0.000600757
go_gc_duration_seconds{quantile="1"} 0.000600757
go_gc_duration_seconds_sum 0.001366348
go_gc_duration_seconds_count 3
# HELP go_gc_gogc_percent Heap size target percentage configured by the user, otherwise 100. This value is set by the GOGC environment variable, and the runtime/debug.SetGCPercent function. Sourced from /gc/gogc:percent
# TYPE go_gc_gogc_percent gauge
go_gc_gogc_percent 100
# HELP go_gc_gomemlimit_bytes Go runtime memory limit configured by the user, otherwise math.MaxInt64. This value is set by the GOMEMLIMIT environment variable, and the runtime/debug.SetMemoryLimit function. Sourced from /gc/gomemlimit:bytes
# TYPE go_gc_gomemlimit_bytes gauge
go_gc_gomemlimit_bytes 9.223372036854776e+18
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 20
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.23.2"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated in heap and currently in use. Equals to /memory/classes/heap/objects:bytes.
# TYPE go_memstats_alloc_bytes gauge
```

