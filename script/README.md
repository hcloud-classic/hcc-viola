## Hcloud-classic

## vnc 관련 자동 실행 스크립트

총 두가지의 파일이 필요하다.
vnc를 컨트롤 하는 스크립트 그리고 특정 프로세스가 실행중인지 확인해주는 스크립트

프로젝트의 `script` 디렉토리를 보면 `hcc_init`, ~~`isprocess`~~ 두가지 스크립트를
각각 적절한 위치에 옮겨야 한다.

hcc_init은 현재 vnc에 관하여 start,stop,restart,status 기능을 지원한다.

### 클러스터 준비

```shell
$ mkdir -p /etc/hcc/viola
$ cp viola.conf /etc/hcc/viola/viola.conf
# interface  name 은 현재 시스템의 pxe boot 때 Dhcp를 받는 nic 의 이름으로 해주면 된다.
$ cp hcc_init hcc_viola /etc/init.d
$ cp viola /etc/hcc/viola/
$ chkconfig hcc_init on
$ chkconfig hcc_viola on
```

### Basic Install

#### rabbitmq
http://www.rabbitmq.com/install-rpm.html 을 참고 하자
https://packagecloud.io/rabbitmq/erlang/install#bash-rpm
```shell
$ yum install erlang -y
# rabbitmq install
$ curl -s https://packagecloud.io/install/repositories/rabbitmq/erlang/script.rpm.sh | sudo bash

$ yum install rabbitmq-server
service rabbitmq-server start
```

#### telegraf

```shell
$ wget https://repos.influxdata.com/rhel/6/x86_64/stable/telegraf-1.15.1-1.x86_64.rpm
$ rpm -ivh telegraf-1.15.1-1.x86_64.rpm
$ chown telegraf:telegraf /var/log/telegraf
$ service telegraf start
$ chkconfig telegraf on
```


### hcc_init

```shell
$ cp hcc_init /etc/init.d
$ cd /etc/init.d
$ chmod 755 hcc_init
$ update-rc.d hcc_init defaults 99
```


## isprocess

```shell
$ cp isprocess /usr/local/sbin/
$ chmod 755 /usr/local/sbin/isprocess
```

