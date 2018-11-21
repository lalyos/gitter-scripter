[![](https://images.microbadger.com/badges/image/lalyos/gitter-scripter.svg)](https://microbadger.com/images/lalyos/gitter-scripter "Get your own image badge on microbadger.com")
[![Docker Automated build](https://img.shields.io/docker/automated/lalyos/gitter-scripter.svg)](https://hub.docker.com/r/lalyos/gitter-scripter)

This is a simple webapp authenticates with gitter oauth (github/gitlab/twitter).
Once authentication is done, it calles out a bash script. The default one:
- joins a central room
- sends a welcome message
- does some kubectl magic
This is all [k8s workshop specific](https://github.com/lalyos/k8s-workshop)

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
kubectl apply -f https://raw.githubusercontent.com/lalyos/gitter-scripter/master/gitter.yaml
```

## roadmap

muhaha

- [ ] custom script from web form
- [ ] several modular script, probably with [plugn](https://github.com/dokku/plugn) or [pluginhook](https://github.com/progrium/pluginhook)