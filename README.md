[![](https://images.microbadger.com/badges/image/lalyos/gitter-scripter.svg)](https://microbadger.com/images/lalyos/gitter-scripter "Get your own image badge on microbadger.com")
[![Docker Automated build](https://img.shields.io/docker/automated/lalyos/gitter-scripter.svg)](https://hub.docker.com/r/lalyos/gitter-scripter)

## create a gitter oauth application

On gitter's developer site [create new app](https://developer.gitter.im/apps/new)
Use redirect url as: `http://gitter.yuorcustom.domain/login/callback`

## store oauth credentials

Store the newly generated **oauth key** and **oauth secret** in a secret:
```
kubectl create secret generic gitter \
  --from-literal=GITTER_OAUTH_KEY=$GITTER_OAUTH_KEY \
  --from-literal=GITTER_OAUTH_SECRET=$GITTER_OAUTH_SECRET
```

## deploy

Now you can create all the k8s resources:
- deployment
- service
- ingress

```
kubectl apply -f http://...
```