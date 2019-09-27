---
title: Expose and secure a service
type: Tutorials
---

This tutorial shows you how to expose and secure services or lambda functions using the API Gateway Controller. The controller reacts to an instance of the APIRule custom resource (CR) and creates an Istio Virtual Service and [Oathkeeper Access Rules](https://www.ory.sh/docs/oathkeeper/api-access-rules) according to the details specified in the CR. To interact with the secured services, the tutorial uses an OAuth2 client registered through the Hydra Maester controller.

The tutorial comes with a sample HttpBin service deployment and a sample lambda function.

## Register an OAuth2 client and get tokens

1. Export these values as environment variables:

  - The name of your client and the Secret which stores the client credentials:

    ```
    export CLIENT_NAME={YOUR_CLIENT_NAME}
    ```

  - The Namespace in which you want to create the client and the Secret that stores its credentials:

    ```
    export CLIENT_NAMESPACE={YOUR_CLIENT_NAMESPACE}
    ```

  - The domain of your cluster:

    ```
    export DOMAIN={CLUSTER_DOMAIN}
    ```

2. Create an OAuth2 client with `read` and `write` scopes. Run:

  ```
  cat <<EOF | kubectl apply -f -
  apiVersion: hydra.ory.sh/v1alpha1
  kind: OAuth2Client
  metadata:
    name: $CLIENT_NAME
    namespace: $CLIENT_NAMESPACE
  spec:
    grantTypes:
      - "client_credentials"
    scope: "read write"
    secretName: $CLIENT_NAME
  EOF
  ```

3. Export the credentials of the created client as environment variables. Run:

  ```
  export CLIENT_ID="$(kubectl get secret -n $CLIENT_NAMESPACE $CLIENT_NAME -o jsonpath='{.data.client_id}' | base64 --decode)"
  export CLIENT_SECRET="$(kubectl get secret -n $CLIENT_NAMESPACE $CLIENT_NAME -o jsonpath='{.data.client_secret}' | base64 --decode)"
  ```

4. Encode your client credentials and export them as an environment variable:

  ```
  export ENCODED_CREDENTIALS=$(echo -n "$CLIENT_ID:$CLIENT_SECRET" | base64)
  ```

5. Get tokens to interact with secured resources:

<div tabs>
  <details>
  <summary>
  Token with "read" scope
  </summary>

  1. Get the token:

      ```
      curl -ik -X POST "https://oauth2.$DOMAIN/oauth2/token" -H "Authorization: Basic $ENCODED_CREDENTIALS" -F "grant_type=client_credentials" -F "scope=read"
      ```

  2. Export the issued token as an environment variable:

      ```
      export ACCESS_TOKEN_READ={ISSUED_READ_TOKEN}
      ```

  </details>
  <details>
  <summary>
  Token with "write" scope
  </summary>

  1. Get the token:

      ```
      curl -ik -X POST "https://oauth2.$DOMAIN/oauth2/token" -H "Authorization: Basic $ENCODED_CREDENTIALS" -F "grant_type=client_credentials" -F "scope=write"
      ```

  2. Export the issued token as an environment variable:

      ```
      export ACCESS_TOKEN_WRITE={ISSUED_WRITE_TOKEN}
      ```

   </details>
</div>

## Deploy, expose, and secure the sample resources

Follow the instructions in the tabs to deploy an instance of the HttpBin service or a sample lambda function, expose them, and secure them with Oathkeeper rules.

<div tabs>

  <details>
  <summary>
  HttpBin - secure endpoints of a service
  </summary>

1. Deploy an instance of the HttpBin service:

  ```
  kubectl apply -f https://raw.githubusercontent.com/istio/istio/master/samples/httpbin/httpbin.yaml
  ```

2. Expose the service and secure it by creating an APIRule CR:

  ```
  cat <<EOF | kubectl apply -f -
  apiVersion: gateway.kyma-project.io/v1alpha1
  kind: APIRule
  metadata:
    name: httpbin
  spec:
    gateway: kyma-gateway.kyma-system.svc.cluster.local
    service:
      name: httpbin-proxy
      port: 8000
      host: httpbin-proxy.kyma.local
    rules:
      - path: /.*
        methods: ["GET"]
        accessStrategies:
          - handler: oauth2_introspection
            config:
              required_scope: ["read"]
      - path: /post
        methods: ["POST"]
        accessStrategies:
          - handler: oauth2_introspection
            config:
              required_scope: ["write"]
  EOF
  ```

>**NOTE:** If you are running Kyma on Minikube, add `httpbin-proxy.kyma.local` to the entry with Minikube IP in your system's `/etc/hosts` file.

The exposed service requires tokens with "read" scope for `GET` requests in the entire service and tokens with "write" scope for `POST` requests to the `/post` endpoint of the service.

  </details>

  <details>
  <summary>
  Secure a lambda function
  </summary>

1. Create a lambda function using the [supplied code](./assets/lambda.yaml):

  ```
  kubectl apply -f ./assets/lambda.yaml
  ```

2. Expose the lambda function and secure it by creating an APIRule CR:

  ```
  cat <<EOF | kubectl apply -f -
  apiVersion: gateway.kyma-project.io/v1alpha1
  kind: APIRule
  metadata:
    name: lambda
  spec:
    gateway: kyma-gateway.kyma-system.svc.cluster.local
    service:
      name: labmda-proxy
      port: 8080
      host: labmda-proxy.kyma.local
    rules:
      - path: /lambda
        methods: ["GET"]
        accessStrategies:
          - handler: oauth2_introspection
            config:
              required_scope: ["read"]
  EOF
  ```

>**NOTE:** If you are running Kyma on Minikube, add `lambda-proxy.kyma.local` to the entry with Minikube IP in your system's `/etc/hosts` file.  

The exposed lambda function requires all `GET` requests to have a valid token with the "read" scope.

  </details>
</div>

## Access the secured resources

Follow the instructions in the tabs to call the secured service or lambda functions using the tokens issued for the client you registered.

<div tabs>

  <details>
  <summary>
  Call secured endpoints of a service
  </summary>

1. Send a `GET` request with a token with the "read" scope to the HttpBin service:

  ```
  curl -ik -X GET https://httpbin-proxy.$DOMAIN/headers -H "Authorization: Bearer $ACCESS_TOKEN_READ"
  ```

2. Send a `POST` request with a token with the "write" scope to the HttpBin's `/post` endpoint:

  ```
  curl -ik -X POST https://httpbin-proxy.$DOMAIN/post -d "test data" -H "Authorization: bearer $ACCESS_TOKEN_WRITE"
  ```

These calls return a code `200` response. If you call the service without a token, you get a code `401` response. If you call the service or its secured endpoint with a token with the wrong scope, you get the code `403` response.

  </details>

  <details>
  <summary>
  Call the secured lambda function
  </summary>

Send a `GET` request with a token with the "read" scope to the lambda function:

  ```
  curl -ik https://lambda-proxy.$DOMAIN/lambda -H "Authorization: bearer $ACCESS_TOKEN_READ"
  ```

This call returns a code `200` response. If you call the lambda function without a token, you get a code `401` response. If you call the lambda function with a token with the wrong scope, you get the code `403` response.

  </details>
</div>