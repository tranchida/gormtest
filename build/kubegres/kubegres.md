Install kubegres operator

kubectl apply -f https://raw.githubusercontent.com/reactive-tech/kubegres/v1.19/kubegres.yaml

see complete doc https://www.kubegres.io/doc/getting-started.html

before create a new cluster, create a secret

```
apiVersion: v1
kind: Secret
metadata:
  name: mypostgres-secret
  namespace: default
type: Opaque
stringData:
  superUserPassword: postgresSuperUserPsw
  replicationUserPassword: postgresReplicaPsw
```

create a cluster with the operator

```
apiVersion: kubegres.reactive-tech.io/v1
kind: Kubegres
metadata:
  name: mypostgres
  namespace: default

spec:

   replicas: 3
   image: postgres:17.2

   database:
      size: 200Mi

   env:
      - name: POSTGRES_PASSWORD
        valueFrom:
           secretKeyRef:
              name: mypostgres-secret
              key: superUserPassword

      - name: POSTGRES_REPLICATION_PASSWORD
        valueFrom:
           secretKeyRef:
              name: mypostgres-secret
              key: replicationUserPassword
```

