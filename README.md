# 场景
用于将configmap定时备份到阿里云oss

## 使用示例
```shell
apiVersion: xinyu.com/v1
kind: PvMonitor
metadata:
  name: pvmonitor-sample
spec:
  regex: ".*"
  email:
    host: "smtp server"
    port: "smtp server port"
    user: "your email"
    password: "your email password"
    subject: "容量报告"
    to:
      - your email
    cc:
      - your email
```

### 安装
```
cd pvmonitor
make output
kubectl -f output/deploy.yaml
```

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

