# 6. Provision Rocket.chat

To deploy Rocket.chat on Kubernetes, follow the steps:

1. Export Capact cluster domain name as environment variable:

   ```bash
   export DOMAIN_NAME={domain_name} # e.g. demo.cluster.capact.dev
   ``` 

1. Create a file with installation parameters:

    ```bash
    cat > /tmp/rocketchat-params.yaml << ENDOFFILE
    ingress:
      host: rocketchat.${DOMAIN_NAME}
    resources:
      requests:
        memory: "2G"
        cpu: "1"
      limits:
        memory: "4G"
        cpu: "1"
    ENDOFFILE
    ```
1. Create Kubernetes Namespace:
    ```bash
    kubectl create namespace rocketchat
    ```
1. Create Action:
 
    ```bash
    capact act cap.interface.productivity.rocketchat.install \
    --name rocketchat \
    --namespace rocketchat \
    --parameters-from-file /tmp/rocketchat-params.yaml
    ```
1. Run the Action:
    ```bash
    capact act run rocketchat
    ```
1. Watch the Action:
    ```bash
    capact act watch rocketchat
    ```
1. Once the Action is succeeded, view output TypeInstances:
   ```bash
   capact act status rocketchat
   ```
    
**Next steps:** Navigate back to the [main README](./README.md) and follow next steps.
