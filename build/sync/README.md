sync
====

The sync tool is a proof-of-concept implementation of a tool for synchronizing mattermost plugin
repositories with the mattermost-plugin-starter-template repo.

Overview
--------

At its core the tool is just a collection of checks and actions that are executed according to a
synchronization plan (see `./build/sync/plan.yml` for an example). The plan defines a set of files
and/or directories that need to be kept in sync between the plugin repository and the template (this
repo).

For each path, a set of actions to be performed is outlined. No more than one action of that set
will be executed - the first one whose checks pass. Other actions are meant to act as fallbacks.
The idea is to be able to e.g. overwrite a file if it has no local changes or apply a format-specific
merge algorithm otherwise.

Before running each action, the tool will check if any checks are defined for that action. If there
are any, they will be executed and their results examined. If all checks pass, the action will be executed.
If there is a check failure, the tool will locate the next applicable action according to the plan and
start over with it.

The synchronization plan can also run checks before running any actions, e.g. to check if the plugin and
template worktrees are clean.

Running
-------

The tool can be executed from the root of this repository with a command:
```
$ go run ./build/sync/main.go ./build/sync/plan.yml ../mattermost-plugin-github
```

(assuming `mattermost-plugin-github` is the plugin repository we want to synchronize with the template).

Caveat emptor
-------------

This is a very basic proof-of-concept and there are many things that should be improved/implemented:
(in no specific order)

   1. Format-specific merge actions for `go.mod`, `go.sum`, `webapp/package.json` and other files should
       be implemented.
   2. Better logging should be implemented.
   3. Handling action dependencies should be investigated.
      e.g. if the `build` directory is overwritten, that will in some cases mean that the go.mod file also needs
      to be updated.
   4. Storing the tree-hash of the template repository that the plugin was synchronized with would allow
      improving the performance of the tool by restricting the search space when examining if a file
      has been altered in the plugin repository.
