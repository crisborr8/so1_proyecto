# so1_proyecto_g1
Toda la configuración estarán en los archivos gcp_config y azure_config para gcp y aks respectivamente.
---
# AZ -- /azure_config
Antes que nada se debe de tener azure cli instalado
## Configuración previa
```
az login
```
```
az provider show -n Microsoft.OperationsManagement -o table
```
```
az provider show -n Microsoft.OperationalInsights -o table
```
```
az provider show -n Microsoft.Insights -o table
```
### Si no esta registrado correr el siguiente codigo
```
az provider register --namespace Microsoft.OperationsManagement
```
```
az provider register --namespace Microsoft.OperationalInsights
```
```
az provider register --namespace Microsoft.Insights
```
### Crear grupo
Se creará un grupo llamado k8azure
```
az group create --name k8azure --location eastus
```
## Obtener la version de kubernete disponible
```
az aks get-versions --location eastus --output table
```
## Creación del cluster
Se creará un cluster llamado cluster-azure
```
az aks create -g k8azure -n cluster-azure --enable-managed-identity --node-count 1 --enable-addons monitoring --enable-msi-auth-for-monitoring --enable-cluster-autoscaler --min-count 1 --max-count 4 --generate-ssh-keys --kubernetes-version 1.22.11
```
Instalar az aks cli si no esta instalado
```
az aks install-cli
```
Conectar al cluster
```
az aks get-credentials --resource-group k8azure --name cluster-azure
```
Creamos namespace usactar
```
kubectl create ns usactar
```
## Crear mongodb
Crear variables bash para almancenar información
```
export COSMOSDB_ACCOUNT_NAME=grupo1-$RANDOM
```
Crear cuenta en cosmos DB
```
az cosmosdb create --name $COSMOSDB_ACCOUNT_NAME --resource-group k8azure --kind MongoDB --enable-public-network true
```
Mostrar la base de datos
```
az cosmosdb mongodb database create --account-name $COSMOSDB_ACCOUNT_NAME --resource-group k8azure --name azureDB
```
Enumerar las bases de datos
```
az cosmosdb mongodb database list --account-name $COSMOSDB_ACCOUNT_NAME --resource-group k8azure -o table
```
Obtenemos el valor de DATABASE_MONGODB_URI a reemplazar en mongo.yaml
```
az cosmosdb keys list --type connection-strings -g k8azure -n $COSMOSDB_ACCOUNT_NAME --query "connectionStrings[0].connectionString" -o tsv
```
Aplicamos mongo.yaml
```
kubectl apply -f mongo.yaml
```
## Redis
Aplicamos redis.yaml
```
kubectl apply -f redis.yaml
```
## Go
Instalamos el controlador de ingress-nginx
```
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
```
```
helm install nginx-ingress ingress-nginx/ingress-nginx -n usactar
```
Aplicamos goapi.yaml
```
kubectl apply -f goapi.yaml
```
## gRPC server
Aplicar grpcserver.yaml
```
kubectl apply -f grpcserver.yaml
```













---
# GCP -- /gcp_config
## Creación del cluster
Verificamos las versiones son compatibles para el cluster y nodos (Estas solo pueden diferir por 1).
```
gcloud container get-server-config --zone=us-central1-a
```
Creamos el cluster llamado
cluster-google con cluster v1.21.x<br>
```
gcloud container clusters create cluster-google --num-nodes=1 --tags=allin,allout --machine-type=n1-standard-2 --no-enable-network-policy --zone=us-central1-a --cluster-version=1.21.14-gke.3000 --node-version=1.20.8-gke.2101 --enable-autoscaling --min-nodes=1 --max-nodes=4
```
Obtenemos las credenciales del cluster recien creado
```
gcloud container clusters get-credentials cluster-google --zone=us-central1-a
```
## Go Kafka
Creamos namespace usactar
```
kubectl create ns usactar
```
Realizamos el deploy de zookeeper.yaml
```
kubectl apply -f zookeeper.yaml 
```
Realizamos el deploy de kafka.yaml
```
kubectl apply -f kafka.yaml 
```
Instalamos ingress-nginx en gcp
```
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
```
```
helm install nginx-ingress ingress-nginx/ingress-nginx -n usactar
```
Realizamos el deploy de goapikafka.yaml
```
kubectl apply -f goapikafka.yaml 
```
Obtener el ADDRESS del ingress y agregarlo en web.yaml para correrlo nuevamente<br>
Instalar grpclient.yaml
```
kubectl apply -f grpcclient.yaml 
```











---
# Linkerd -- /linkerd
Si en algun momento llega a fallar la instalación, ir directamente con las credenciales del cluster que falla e instalar nuevamente.<br>
Instalar linkerd client
```
curl --proto '=https' --tlsv1.2 -sSfL https://run.linkerd.io/install | sh
```
Exportar la variable de entorno
```
export PATH=$PATH:$HOME/.linkerd2/bin
```
## Multicluster
Instalar step
```
wget https://dl.step.sm/gh-release/cli/docs-cli-install/v0.21.0/step-cli_0.21.0_amd64.deb
```
```
sudo dpkg -i step-cli_0.21.0_amd64.deb
```
```
sudo rm step-cli_0.21.0_amd64.deb
```
### Context
Si se desea se puede renombrar los contexts para un mejor manejo
```
kubectl config get-contexts
kubectl config rename-context oldName newName
```
En este caso se usaran los nombres cluster-google y cluster-azure
### Certificados
Generar certificados
```
step certificate create root.linkerd.cluster.local root.crt root.key \
  --profile root-ca --no-password --insecure
```
Generar las credenciales
```
step certificate create identity.linkerd.cluster.local issuer.crt issuer.key \
  --profile intermediate-ca --not-after 8760h --no-password --insecure \
  --ca root.crt --ca-key root.key
```
Instalar linkerdCRD en ambos clusters, el context será el nombre del cluster
```
linkerd install --crds | kubectl --context=cluster-google apply -f -

linkerd install --crds | kubectl --context=cluster-azure apply -f -
```
Instalar el control panel
```
linkerd install --identity-trust-anchors-file root.crt --identity-issuer-certificate-file issuer.crt --identity-issuer-key-file issuer.key | kubectl --context=cluster-google apply -f -

linkerd install --identity-trust-anchors-file root.crt --identity-issuer-certificate-file issuer.crt --identity-issuer-key-file issuer.key | kubectl --context=cluster-azure apply -f -
```
Instalar linkerd viz (Se tardara unos minutos)
```
linkerd --context=cluster-google viz install | kubectl --context=cluster-google apply -f -

linkerd --context=cluster-azure viz install | kubectl --context=cluster-azure apply -f -
```
Verificar la instalación (Tardará un poco)
```
for ctx in cluster-google cluster-azure; do
  echo "Checking cluster: ${ctx} ........."
  linkerd --context=${ctx} check || break
  echo "-------------"
done
```
### Preparacion del cluster
Instalar multiclusters
```
linkerd --context=cluster-google multicluster install | kubectl --context=cluster-google apply -f -

linkerd --context=cluster-azure multicluster install | kubectl --context=cluster-azure apply -f -
```
Verificar el gateway
```
for ctx in cluster-google cluster-azure; do
  echo "Checking gateway on cluster: ${ctx} ........."
  kubectl --context=${ctx} -n linkerd-multicluster \
    rollout status deploy/linkerd-gateway || break
  echo "-------------"
done
```
Verificar el balanceador de carga
```
for ctx in cluster-google cluster-azure; do
  printf "Checking cluster: ${ctx} ........."
  while [ "$(kubectl --context=${ctx} -n linkerd-multicluster get service -o 'custom-columns=:.status.loadBalancer.ingress[0].ip' --no-headers)" = "<none>" ]; do
      printf '.'
      sleep 1
  done
  printf "\n"
done
```
### Conección
Conectamos cluster-google con cluster-azure
```
linkerd --context=cluster-azure multicluster link --cluster-name cluster-azure |
  kubectl --context=cluster-google apply -f -
```
Verificammos que exista conexion en cluster-google
```
linkerd --context=cluster-google multicluster check
```
Verificamos que tenga conexion con el gateway de azure
```
linkerd --context=cluster-google multicluster gateways
```
### Exponer servicios
#### Mongo

Exponiendo el servicio de mongodb (en azure)
```
kubectl get deploy mongodb -n usactar -o yaml | linkerd inject - | kubectl apply --context cluster-azure -n usactar -f -

kubectl expose deploy mongodb --target-port=27017 --port=27017 --context cluster-azure -n usactar

kubectl --context=cluster-azure label svc -n usactar mongodb mirror.linkerd.io/exported=true
```
Verificamos que el servicio se haya creado
```
kubectl --context=cluster-google -n usactar get svc mongodb-cluster-azure
```
Verificamos que este alcanzando el enpoint de mongodb
```
kubectl --context=cluster-google -n usactar get endpoints mongodb-cluster-azure \
  -o 'custom-columns=ENDPOINT_IP:.subsets[*].addresses[*].ip'
```
Las ips deben de ser iguales a el gateway
```
kubectl --context=cluster-azure -n linkerd-multicluster get svc linkerd-gateway \
-o "custom-columns=GATEWAY_IP:.status.loadBalancer.ingress[*].ip"
```
Inyectamos a grpcclient 
```
kubectl get deploy grpcclient -n usactar -o yaml --context cluster-google | linkerd inject - | kubectl apply --context cluster-google -n usactar -f -
```
#### grpcserver
Exponiendo el servicio de grpcserver (en azure)
```
kubectl get deploy grpcserver -n usactar -o yaml | linkerd inject - | kubectl apply --context cluster-azure -n usactar -f -

kubectl expose deploy grpcserver --target-port=50051 --port=50051 --context cluster-azure -n usactar

kubectl --context=cluster-azure label svc -n usactar grpcserver mirror.linkerd.io/exported=true
```
Verificamos que el servicio se haya creado
```
kubectl --context=cluster-google -n usactar get svc grpcserver-cluster-azure
```
Verificamos que este alcanzando el enpoint de grpcserver
```
kubectl --context=cluster-google -n usactar get endpoints grpcserver-cluster-azure \
  -o 'custom-columns=ENDPOINT_IP:.subsets[*].addresses[*].ip'
```
Las ips deben de ser iguales a el gateway
```
kubectl --context=cluster-azure -n linkerd-multicluster get svc linkerd-gateway \
-o "custom-columns=GATEWAY_IP:.status.loadBalancer.ingress[*].ip"
```
Inyectamos a grpcclient 
```
kubectl get deploy grpcclient -n usactar -o yaml --context cluster-google | linkerd inject - | kubectl apply --context cluster-google -n usactar -f -
```

## Observador
Inyectar al deployment (Si da fallos usar las credenciales de kubectl de cada cluster)
```
kubectl --context=cluster-google get -n usactar deploy -o yaml | linkerd inject - | kubectl apply -f -

kubectl --context=cluster-azure get -n usactar deploy -o yaml | linkerd inject - | kubectl apply -f -
```
## Dashboard
```
linkerd --context=cluster-google viz dashboard

linkerd --context=cluster-azure viz dashboard
```
Verificaciones extra
```
kubectl -n usactar get po -o jsonpath='{.items[0].spec.containers[*].name}'

linkerd --context=cluster-google -n usactar stat deploy

linkerd --context=cluster-azure -n usactar stat deploy
```











# Chaos mesh -- /chaos-mesh
## Configuración
Se debe tener helm instalado, esto se configuró con helm v3.10.1
```
helm repo add chaos-mesh https://charts.chaos-mesh.org
```
```
helm search repo chaos-mesh --version 2.3.2
```
Instalar en cluster-google y cluster-azure
```
helm install chaos-mesh chaos-mesh/chaos-mesh -n=usactar --set dashboard.create=true --set dashboard.securityMode=false --version 2.3.2 
```
Verificar que los pods esten corriendo
```
kubectl get pods --namespace usactar -l app.kubernetes.io/instance=chaos-mesh 
```
Abrir dashboard
```
kubectl port-forward -n usactar svc/chaos-dashboard 2333:2333
```
## Chaos experiments
Ir al dashboard de chaos-mesh (obtener id por el ingress) y generar experimentos, si pide token o autenticación desactivarla con el siguiente comando.
```
helm upgrade chaos-mesh chaos-mesh/chaos-mesh --namespace=chaos-mesh --version 2.3.2 --set dashboard.securityMode=false
```







# Curls
Crear pod para curl por si hay fallas y se desea verificar las conexiones
```
kubectl run mycurlpod --image=curlimages/curl -i --tty -- sh
```
Si ya se ha creado y se desea ingresar al sh
```
kubectl exec -i --tty mycurlpod -- sh
```
Ya se puede hacer curl al endpoint
---







# Bibliografía
* https://cloud.google.com/kubernetes-engine/docs/tutorials/http-balancer
* https://chaos-mesh.org/docs/production-installation-using-helm/
* https://chaos-mesh.org/docs/next/faqs/
* https://chaos-mesh.org/docs/2.3.2/manage-user-permissions/
* https://learn.microsoft.com/es-es/training/modules/aks-manage-application-state/3-exercise-create-resources
* https://saibharath005.medium.com/mongodb-in-azure-kubernetes-4b06c2bc152b
* https://azuredevopslabs.com/labs/vstsextend/kubernetes/?sa=X&ved=2ahUKEwjeoZ_n34noAhU8K6YKHduXAe0Q9QF6BAgKEAI
* https://artifacthub.io/packages/helm/chaos-mesh/chaos-mesh
* https://www.tutorialworks.com/kubernetes-curl/
* https://linkerd.io/2.11/getting-started/
* https://linkerd.io/2.12/tasks/multicluster/

