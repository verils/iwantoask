# Iwantoask

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=verils_iwantoask&metric=alert_status)](https://sonarcloud.io/dashboard?id=verils_iwantoask)

### How to run

```shell script
docker pull verils/iwantoask

docker run --name iwantoask -p 8080:8080 -d verils/iwantoask
```

### Todos
- [x] 启动时显示版本号
- [x] 支持BaseUrl在环境变量配置（包括HTML和Go代码）
- [x] 使用CI进行构建（GitHub Actions）
- [x] 提交表单时进行参数校验和错误提示
- [x] 支持分页显示列表
- [x] 根据Cookie区分用户
- [x] 防止重复提交（禁用提交按钮）
- [ ] 针对请求url的测试（基于响应码的断言）
- [ ] 自动为用户生成用户名（支持中文）
- [ ] 支持对问题的Star
- [ ] 国际化支持
