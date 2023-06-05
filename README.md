# Gaia

Gaia is a gRPC server that returns Mac OS and Nvidia based embedded systems' stats to client. 

## Installation

Make sure you have redis installed for logging data. 
```
brew install redis
```

For Bazel installation On MacOS : `brew install bazelisk` 
Next, in the project directory update repositories with 
```
bazel run //:gazelle

bazel run //:update-repos
```
Run the redis server:
```
redis-server --port 8002
```

Fianlly run Gaia
```
bazel run //:gaia
```

---

