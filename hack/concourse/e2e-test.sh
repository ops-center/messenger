#!/bin/sh

set -x -e

# start docker and log-in to docker-hub
entrypoint.sh
docker login --username=$DOCKER_USER --password=$DOCKER_PASS
docker run hello-world

apt-get update &> /dev/null
apt-get install -y git python python-pip &> /dev/null

# install kubectl
curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl &> /dev/null
chmod +x ./kubectl
mv ./kubectl /bin/kubectl

# install onessl
curl -fsSL -o onessl https://github.com/kubepack/onessl/releases/download/0.3.0/onessl-linux-amd64 \
  && chmod +x onessl \
  && mv onessl /usr/local/bin/

# copy voyager to $GOPATH
mkdir -p $GOPATH/src/github.com/appscode
cp -r messenger $GOPATH/src/github.com/appscode

pushd $GOPATH/src/github.com/appscode/messenger

# build and push docker image
./hack/builddeps.sh
export APPSCODE_ENV=dev
export DOCKER_REGISTRY=shudipta
./hack/docker/make.sh build
./hack/docker/make.sh push


# create config/.env file that have all necessary creds
cat > hack/config/.env <<EOF
AUTH_TOKEN_TO_SEND_MSG=TOKEN_TO_SEND
AUTH_TOKEN_TO_SEE_HIST=TOKEN_TO_SEE
EOF

# run tests
source ./hack/deploy/messenger.sh --docker-registry=shudipta
./hack/make.py test e2e --kubeconfig=/home/kube/config --webhook=false

popd

