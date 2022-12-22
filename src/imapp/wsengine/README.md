# 程序启动的环境变量

```bash
NACOS_IP=172.21.112.1;                                                         #nacos
NACOS_PORT=8848;                                                               #nacos
NACOS_GROUP=DEFAULT_GROUP;                                                     #nacos
NACOS_NAMESPACE_ID=18f2260c-ffde-4e4e-a526-be1c4e5ee8b2;                       #nacos
NACOS_DATA_ID=cluster_whatsapp_engine;                                         #nacos
NACOS_USERNAME=nacos;
NACOS_PASSWORD=nacos;
```

# Nacos配置

```yaml
MysqlDataBase:     #数据库
  Username: ykappuser2
  Password: LCLi7FaZskxnKo5MW5
  Host: 172.21.112.1
  Port: 3306
  DataBaseName: whatsapp
  CharSet: utf8
  LogLevel: info
  MaxIdleConnection: 8
  MaxOpenConnection: 20
  MaxLifeTime: 60
```

# 其他依赖库

## FFmpeg

> *因媒体文件需要转码使用了FFmpeg, 如果没有使用私信媒体文件功能则不影响程序使用*

```bash
apt install ffmpeg
```

## libvips

> *因图片文件需要进行特定的算法，并且不以代码依赖的方式使用了libvips, 如果没有使用私信图片文件功能则不影响程序使用*

#### Debian

```bash
apt install libvips-dev
```

#### Centos

```bash
yum install https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm
yum install yum-utils
yum install http://rpms.remirepo.net/enterprise/remi-release-7.rpm
yum-config-manager --enable remi
yum install vips vips-devel vips-tools
```


# 解析protobuf序列化的内容

#### 具备linux系统环境和protoc执行文件 可以这么操作去解析hex protobuf

    echo hex内容 | xxd -r -p | ./protoc --decode_raw

> [下载链接](https://github.com/protocolbuffers/protobuf/releases/download/v21.5/protoc-21.5-linux-x86_64.zip)
