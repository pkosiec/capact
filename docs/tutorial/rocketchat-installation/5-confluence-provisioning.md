# 5. Provision Confluence

To deploy Atlassian Confluence Data Center on Kubernetes, follow the steps:

1. Save your Confluence license in the `license.txt` file. 

1. Export Capact cluster domain name, license and PostgreSQL TypeInstance ID as environment variables:

    >**NOTE**: Use the PostgreSQL TypeInstance ID from the [Provision AWS RDS for PostgreSQL](./2-aws-rds-provisioning.md) tutorial.
   ```bash
   export DOMAIN_NAME={domain_name} # e.g. demo.cluster.capact.dev
   export LICENSE_KEY_BASE64=$(/bin/cat license.txt | base64 )
   export POSTGRESQL_TI_ID={ti_id} 
   ``` 

1. Create a file with installation parameters:

    ```bash
    cat > /tmp/confluence-params.yaml << ENDOFFILE
    ingress:
      host: confluence.${DOMAIN_NAME}
    
    confluence:
      # -- The Confluence license key.
      # If specified, the license is automatically populated during Confluence setup.
      # Otherwise, it will need to be provided via the browser after initial startup.
      licenseKeyInBase64: ${LICENSE_KEY_BASE64}
    
      resources:
        container:
          limits:
            cpu: "1"
            memory: "4G"
          requests:
            cpu: "1"
            memory: "2G"
	jvm:
          maxHeap: 2g
          minHeap: 512m
    ENDOFFILE
    ```
1. Create a file with TypeInstances IDs:
 
    ```bash
    cat > /tmp/confluence-ti.yaml << ENDOFFILE
    typeInstances:
      - name: "postgresql"
        id: "${POSTGRESQL_TI_ID}"
    ENDOFFILE
    ```
1. Create Kubernetes Namespace:
    ```bash
    kubectl create namespace confluence
    ```
1. Create Action:
    >**NOTE:** You need to have a proper cluster policy configuration as described in the [Configure Cluster Policy to prefer AWS solutions](./1-cluster-policy-configuration.md) tutorial.
 
    ```bash
    capact act cap.interface.productivity.confluence.install \
    --name confluence \
    --namespace confluence \
    --parameters-from-file /tmp/confluence-params.yaml
    --type-instances-from-file /tmp/confluence-ti.yaml
    ```
1. Run the Action:
    ```bash
    capact act run confluence
    ```
1. Watch the Action:
    ```bash
    capact act watch confluence
    ```
1. Once the Action is succeeded, view output TypeInstances:
   ```bash
   capact act status confluence
   ```
    
**Next steps:** Navigate back to the [main README](./README.md) and follow next steps.
