Simple serverless for K8S.

##Requirements
- Running K8S
- Config in home directory 

##Steps
1. Create namespace
2. Send POST request (endpoint /function)

Example path:
```console
localhost:10000/function
```

Example body:
```json
{
    "name": "fun",
    "namespace" : "default",
    "source" : "module.exports = {main: function(event, context) {return 'Hello World6!'}}",
    "deps": "{\n      \"name\": \"complaint-prev\",\n      \"version\": \"1.0.0\",\n      \"dependencies\": {\n          \"express\": \"^4.17.1\"\n        }\n    }"
}
```


3.  Create port-forward
```console
 kubectl port-forward service/fun1 3000:3000
```