# 表示依赖 alpine 最新版
FROM alpine:latest
MAINTAINER Liu Jun<liujun@newcapec.net>
ENV VERSION 1.0

# 在容器根目录 创建一个 nct 目录
WORKDIR /nct

# 挂载容器目录
VOLUME ["/nct/configs"]
VOLUME ["/nct/data"]

# 拷贝当前目录下可以执行文件
COPY ShopService /nct/

# 拷贝配置文件到容器中
COPY configs /nct/configs

# 设置时区为上海
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone

# 设置编码
ENV LANG C.UTF-8

# 暴露端口
EXPOSE 7001
EXPOSE 7002

# 定义容器启动的主命令
ENTRYPOINT ["/nct/ShopService"]

# 定义默认参数
CMD ["start", "-c", "./configs"]
