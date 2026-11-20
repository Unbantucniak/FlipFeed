# GCFeed

GCFeed 是一个短视频 Feed 服务。用户可以上传视频、关注作者，在最新、关注、热门和推荐四个视图中浏览作品。React 播放页会预加载接下来的视频，并上报播放质量。

## 视频发布后如何进入 Feed

Go 后端分为 API 和 Worker 两个进程。API 使用 Gin 处理登录、上传、互动和分页请求，视频、用户和关注关系保存在 MySQL。上传 MP4 时，服务使用 FFmpeg 将视频索引移到文件头，方便浏览器边下载边播放。

发布视频后，Worker 从 RabbitMQ 接收事件，预热 Feed 缓存并生成视频向量。普通作者的新作品会写入关注者的收件箱。粉丝数达到阈值的作者使用发件箱，读关注流时再合并，减少热门作者发布时的写入量。

最新流按发布时间和视频 ID 做游标分页。Redis 缓存列表页和视频卡片，也用于维护热榜和互动计数。

## 启动

安装带 Compose 插件的 Docker。在 Windows 上可以直接运行 [start-gcfeed.cmd](start-gcfeed.cmd)。脚本会启动容器，检查首页和最新流接口，然后打开浏览器。运行 [stop-gcfeed.cmd](stop-gcfeed.cmd) 可停止服务。

使用默认 Compose 配置时，从仓库根目录启动服务。

```sh
docker compose -f apps/docker-compose.yml up -d --build
```

| 服务 | Windows 脚本 | 默认 Compose |
| --- | --- | --- |
| Web | `http://127.0.0.1:5174` | `http://127.0.0.1:5173` |
| API 健康检查 | `http://127.0.0.1:8080/health` | `http://127.0.0.1:8080/health` |
| Prometheus | `http://127.0.0.1:9092` | `http://127.0.0.1:9090` |
| Grafana | `http://127.0.0.1:3001` | `http://127.0.0.1:3000` |

注册账号并发布第一条视频后，可以用下面的请求在最新流中找到它。

```sh
curl "http://127.0.0.1:8080/api/feed-items?scene=timeline&limit=10"
```

响应中的 `items` 是当前页的视频。若 `has_more` 为真，把 `next_cursor` 作为下一次请求的 `cursor` 参数。Grafana 的 **GCFeed Overview** 面板可以查看 Feed 的 P95 延迟、缓存命中率和 Worker 任务成功率。

使用默认 Compose 配置时，在仓库根目录查看日志或停止容器。

```sh
docker compose -f apps/docker-compose.yml logs -f api worker
docker compose -f apps/docker-compose.yml stop
```

## 服务配置

默认 Compose 端口绑定到本机。API 启动时生成 JWT 签名密钥。要让登录状态跨重启生效，请设置固定的 `GCFEED_JWT_SECRET`。调用需要内部鉴权的接口时，通过 `GCFEED_INTERNAL_TOKEN` 配置调用方和 API 共用的令牌。

使用 [Kubernetes 部署清单](apps/deploy.yaml) 时，先在 `gcfeed` 命名空间创建 `gcfeed-secret`，包含 `mysql-root-password` 和 `api-config.yaml` 两个键。配置文件的结构见 [容器配置](apps/api/configs/config.docker.yaml)，其中的数据库密码应与 Secret 一致，并为各个 API 副本设置相同的 JWT 签名密钥。

## 开发与测试

HTTP 路由在 [router.go](apps/api/internal/interfaces/http/router/router.go)，后台任务在 [Worker 入口](apps/api/cmd/worker/main.go) 注册。前端位于 [apps/web](apps/web)，使用 React 和 Vite。

运行后端测试需要 Go 1.26.1 或更高版本。

```sh
cd apps/api
go test ./...
```

从仓库根目录安装前端依赖并构建播放页。

```sh
npm --prefix apps/web ci
npm --prefix apps/web run build
```
