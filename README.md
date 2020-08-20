Sometimes namespaces get stuck in terminating state.  This small program/container/app finds namespaces in terminating state and removes them.

 - By default we read kubeconfig at $HOME or $USERPROFILE if on Windows
   - You can override this by supplying the following param: `thisbinary kubeconfig=./path/to/dir`