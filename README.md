# scaffold
this creates scaffold with gotmpl. be simple and customizable.

# usage
interactive mode
```
> scaffold create cleanArch targetDir

> AggregateRootName?: user
> MethodName?: create

created targetDir/domain/user.go
created targetDir/usecase/user/create.go
created targetDir/drivers/user.go
created targetDir/adapters/user/repository.go
...
```

cui mode
```
> cat user_create.yaml
AggregateRootName: user
MethodName: create
> scaffold create cleanArch targetDir < user_create.yaml

created targetDir/domain/user.go
created targetDir/usecase/user/create.go
created targetDir/drivers/user.go
created targetDir/adapters/user/repository.go
...
```