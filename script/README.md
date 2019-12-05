## Hcloud-classic

## vnc 관련 자동 실행 스크립트

총 두가지의 파일이 필요하다.
vnc를 컨트롤 하는 스크립트 그리고 특정 프로세스가 실행중인지 확인해주는 스크립트

프로젝트의 `script` 디렉토리를 보면 `hcc_init`, `isprocess` 두가지 스크립트를
각각 적절한 위치에 옮겨야 한다.

hcc_init은 현재 vnc에 관하여 start,stop,restart,status 기능을 지원한다.





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

