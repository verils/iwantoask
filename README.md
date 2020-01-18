# Iwantoask

### How to run

```shell script
docker pull verils/iwantoask

docker run --name iwantoask -p 8080:8080 -d verils/iwantoask
```

### TODO
- [x] 启动时显示版本号
- [x] 支持BaseUrl在环境变量配置（包括HTML和Go代码）
- [x] 使用CI进行构建（GitHub Actions）
- [x] 提交表单时进行参数校验和错误提示
- [x] 支持分页显示列表
- [ ] 根据Cookie区分用户
- [ ] 自动为用户生成用户名（支持中文）
- [ ] 支持对问题的Star
- [ ] 国际化支持
