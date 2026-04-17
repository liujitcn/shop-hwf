##  🛸项目启动
```bash
# 进入服务端
cd server

# 初始化依赖
go mod tidy

# 启动项目
go run main.go
```
```bash
# 进入前端
cd web

# 安装依赖
pnpm install

# 运行项目
pnpm run dev

# 打包发布
pnpm run build
```

## Docker

```bash
# build
docker build -t shop-service.
docker save -o ./shop-service.tar shop-service:latest
docker load -i shop-service.tar

docker-compose up -d

docker compose up -d --force-recreate


docker run -d \
  --name nginx \
  --restart=always \
  -p 80:80 \
  -v /home/nginx/nginx.conf:/etc/nginx/nginx.conf \
  -v /home/nginx/conf.d:/etc/nginx/conf.d \
  -v /home/nginx/html:/usr/share/nginx/html \
  -v /home/nginx/logs:/var/log/nginx \
  -v /home/nginx/cert:/etc/nginx/cert \
  -v /home/data/shop:/home/data/shop \
  nginx
```



