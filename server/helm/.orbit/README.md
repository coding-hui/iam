应用仓库中的目录格式以及意义如下所示：
```.
├── .orbit
│   ├── application-id.yaml #应用定义，启动顺序
│   ├── dashboards
│   │   └── app-dashboard-A.yaml #应用级自定义监控面板
│   ├── db-migrations
│   │   ├── ${db}
│   │   │   └── $timestamp.sql #数据库变更脚本
│   │   │── chagelogs.yaml #数据库变更快照
│   │   └── snapshot.yaml #数据库模型的结构化快照
│   ├── environments
│   │   ├── prod
│   │   │   ├── dashboards
│   │   │   │   └── dashboard-B.yaml #环境级自定义监控面板
│   │   │   └── env-prod.yaml #环境定义
│   │   └── testing
│   │       └── env-testing.yaml #环境定义
│   ├── pipelines
│   │   └── pipeline-A.yaml #部署流程
│   ├── strategies
│   │   └── canary.yaml #部署策略
│   └── versions
│       └── ${versionId}.yaml #版本定义，包含应用镜像变更，数据库变更文件和顺序
├── Charts.yaml
├── templates #helm 模版（helm）
│   ├── configmap-A.yaml
│   └── deployment-A.yaml
├── testing-values.yaml #环境 value（helm）
└── values.yaml #应用 value（helm）
```