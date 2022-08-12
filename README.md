Sometimes Syscalls Get Stuck
============================

When running a few hundred million os.Open or os.Lstat calls on a
large NFS filer, some tiny percentage of them would never (at least in
a few days) return, breaking my "is the work all done" logic.

This package wraps these calls in a form that has timeouts as well.

e.g.:

```
	import (
	   timeout "github.com/jayalane/go-syscalls-timeout"
	)

....

	des, err := timeout.ReadDir(".")

```
