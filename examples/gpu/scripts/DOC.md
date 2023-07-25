### Steps to Deploy

- Create a multi node kind cluster.
- Disable the native kube-scheduler in the kind cluster by removing the static manifest.
- Build the written plugin using TinyGo.
  ``tinygo build -o main.wasm -gc=custom -tags=custommalloc -scheduler=none --no-debug -target=wasi .``
  ``cp main.wasm ../../scheduler/cmd/scheduler/main.wasm ``
- Write the scheduler configuration.
- Build the wasm extender scheduler and pass the scheduler configuration.
  ``go build main.go``
- Run the wasm extender scheduler.
  ``./main --config ../../../examples/gpu/scripts/sched-config.yaml``
