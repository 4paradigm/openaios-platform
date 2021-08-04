# Webhook组件

为了方便用户使用我们指定的算力规格，我们使用MutatingAdmissionWebhook对用户创建的的任务进行修改。
- 通过label`openaios.4paradigm.com/app: true`判断该pod为用户创建，然后读取其annotations，判断每个container所使用的算力规格，然后将每个container使用的算力规格作为list添加到annotation中。
- 对于某个container使用的算力规格，我们通过请求billing server获取到算力规格的详细信息，然后加入到yaml中
- 该模块参考了仓库 https://github.com/stackrox/admission-controller-webhook-demo
- 为了保证webhook server的服务证书不过期，我们通过了`generate-keys.sh`生成了有效期36500天的证书，但是如果需要改变service name或者namespace，应对脚本进行修改并且重新生成证书。