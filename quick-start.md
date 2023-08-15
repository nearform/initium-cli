# Quick Start Guide

#### initium-cli

This is a guide to help developers start using this repository. Even though it's based on a cluster spawned by k8s-kurated-addons, as long as the destination cluster has knative, the CLI will work.

### Prerequisites

- The project's dependencies as described in [the README](README.md)
- `tilt` (if using `initium-platform` to bring up the cluster)

### Using the software

Follow the steps:

1. Install kka-cli (run `make build` and copy the executable to PATH)
2. If you want to use `initium-platform`, clone it and run `make`
    * this will bring up a cluster using kind (kubernetes in Docker)
3. Wait for the cluster to stabilize
    * it's easier to spot checking `tilt` and `argocd` (they will show everything green)
4. Create a repo and enable read and write permissions for the GitHub Actions workflows as described in the [GitHub docs](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/enabling-features-for-your-repository/managing-github-actions-settings-for-a-repository#configuring-the-default-github_token-permissions)
5. Clone the new repo
6. Get argocd-password, enable and access the argocd port forward from Tilt
7. Create a new branch from main / master in the repo
8. Run `initium init config > .initium.yaml`
    * the app name needs to be unique since it will be used by knative to expose its service
    * it’s going to be on the domain GitHub Actions will output
    * if the organization name or GitHub account has uppercase characters, you will need to edit the `.kka` file and change the repo-name to a fully lowercase string
9. Run `initium init github`
10. Run `initium init service-account | kubectl apply -f -` to create the service account
11. Run the following script:
```
export KKA_LB_ENDPOINT="$(kubectl get service -n istio-ingress istio-ingressgateway -o go-template='{{(index .status.loadBalancer.ingress 0).ip}}'):80"
export KKA_CLUSTER_ENDPOINT=$(kubectl config view -o jsonpath='{.clusters[?(@.name == "kind-k8s-kurated-addons")].cluster.server}')
export KKA_CLUSTER_TOKEN=$(kubectl get secrets kka-cli-token -o jsonpath="{.data.token}" | base64 -d)
export KKA_CLUSTER_CA_CERT=$(kubectl get secrets kka-cli-token -o jsonpath="{.data.ca\.crt}" | base64 -d)
```
1.  `echo` the variables set above and set `CLUSTER_CA_CERT`, `CLUSTER_ENDPOINT` and `CLUSTER_TOKEN` as [GitHub Actions secrets](https://docs.github.com/en/actions/security-guides/encrypted-secrets#creating-encrypted-secrets-for-a-repository) in the new repo (or in the GitHub organization itself), using the values returned
    * The `CLUSTER_ENDPOINT` value should be in the URL:port format
    * to access the cluster from the github action you'll need to expose it via ngrok, grab the port from the `KKA_CLUSTER_ENDPOINT` variable and run `ngrok tcp --region=us <port>`
2.  Write the JS / Go code for your repository, exposing a port, and push it to GitHub
3.  Open a pull request
4.  A workflow should be running in GitHub Actions, building an image (even if there’s no Dockerfile in the repo), pushing the image to the registry, and deploying the service to the cluster
5.  Run `kubectl get ksvc -A` to see the application running in your cluster
6.  Run `curl -H “Host: <app URL shown in the action logs>” $KKA_LB_ENDPOINT`
18. You're all set!

The next steps are optional.

19. Make any change to the code and push it to the branch that has an open pull request
20. Run `curl -H “Host: <app URL shown in the action logs>” $KKA_LB_ENDPOINT` again to see the changes after the workflow finishes running
21. Merge the pull request so GitHub Actions can clean up the environment and deploy the version from main to the cluster
22. The docker image can be accessed in the GitHub repo packages session
23. Run `curl -H “Host: <app URL shown in the action logs>” $KKA_LB_ENDPOINT` again to see the changes after the workflow finishes running

Example usage:

One of the developers wants to change the app language from one to another. It’s as simple as:

- Change the code
- Push to a branch and open a pull request
- Wait for the action to run (it will build using a Dockerfile based on the new language)
- Ensure tests pass (or the build will fail)
- cURL the URL shown in the action logs
- Merge the PR
- Wait for the cleanup and new build
- cURL the new URL
- Done
