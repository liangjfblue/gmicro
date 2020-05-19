# go-micro

[![Build Status](https://travis-ci.org/liangjfblue/gmicro.svg?branch=master)](https://travis-ci.org/liangjfblue/gmicro)

## 服务
- 用户服务
- 发表文章服务
- 评论服务
- 鉴权服务


## 1、编译
./scripts/build.sh all


## 2、生成Dockerfile
./scripts/dockerfile.sh all


## 3、运行
创建deployments/db/mysql_data目录

## 3.0 打包
进入deployments目录: `sudo docker-compose build`

## 3.1、运行
进入deployments目录: `sudo docker-compose up`

## 3.2、停止
进入deployments目录: `sudo docker-compose down`


## 4、调用方法
- 注册用户

`http://172.16.7.16:7020/v1/user/register`

- 用户登录获取token

`http://172.16.7.16:7020/v1/user/login`


- 其余接口header带上Authorization（token）

- comment: `DELETE /v1/comment/:cid`
- comment: `GET    /v1/comment/list`
- user:    `POST   /v1/user/register`
- user:    `POST   /v1/user/login`
- user:    `GET    /v1/user/info`
- user:    `GET    /v1/coin/get`
- user:    `POST   /v1/coin/add`

## 5、分布式链路追踪
- opentracing
- jaeger

访问jaeger UI: `http://172.16.7.16:16686`
