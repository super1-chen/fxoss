# Fxoss 

# FXOSS command-line tool (go version ) [![Build Status](https://travis-ci.org/super1-chen/fxoss.svg?branch=master)](https://travis-ci.org/super1-chen/fxoss)

This tool is used for get CDS asset list (with an option), show
CDS asset detail information, and login CDS by it's SN directly.
We assume that your system is Ubuntu 14.04 and
 
 go version 1.10


## Download and Install

goto [https://github.com/super1-chen/fxoss/releases](https://github.com/super1-chen/fxoss/releases) and download least version

and move execute files into  `/usr/local/bin`

## Setup FXOSS tool configuration

Add them in your environment variables.

```shell
export FXOSS_HOST=https://oss.fxdata.cn
export FXOSS_USER=admin
export FXOSS_PWD=admin
export FXOSS_SSH_USER=root
export FXOSS_SSH_PWD='xxxxxx'
```

__Notice:__
1. *_PWD must be included with quotes as 'password'

## Setup Email Configuration

if you want use `fxoss cds-report` you should setup email configuration first.

add json file into `/tmp/fx_email.json`

> /tmp/fx_email.json

```
{
    "address": "email@fxdata.cn",
	"password": "email password",
	"smtp_server": "smtp server"
}
```

## How to use the tool
### help information

fxoss -h \ --help

> help information 

```shell
fxoss is a command line tool for get cds list, show cds detail and ssh login cds server...

Usage:
  fxoss [command]

Available Commands:
  cds-list    Show cds list
  cds-login   SSH login remote server
  cds-port    Show cds port information
  cds-report  Make cds disk type report and send the report by email
  cds-show    Show cds detail info
  help        Help about any command
  version     Print the version number of fxoss

Flags:
  -h, --help      help for fxoss
  -v, --verbose   run fxoss in verbose mode

Use "fxoss [command] --help" for more information about a command.
```


### fxoss cds-list \[option\]


Show ALL cds list

```shell
$ fxoss cds-list
```

Show cds with option

```shell
$ fxoss cds-list 南京

Get cds list from api successfully
+-----+--------------------------+---------------+---------+---------------+-------------+------------------+---------------+-------------------+-----------------+-------------------+---------+---------------------+
| #   | company                  | sn            |  status | license_start | license_end | online_user(max) | hit_user(max) | service_kbps(max) | cache_kbps(max) | monitor_kbps(max) | version |          updated_at |
+-----+--------------------------+---------------+---------+---------------+-------------+------------------+---------------+-------------------+-----------------+-------------------+---------+---------------------+
| 14  | 南京农业大学工学院       | CAS0530000102 | healthy |          None |        None |       1143(1180) |       91(120) |      33792(78848) |     8192(45056) |    240640(349184) |   9.5.2 | 2017-06-15 14:20:28 |
| 15  | 南京航空航天大学江宁校区 | CAS0530000106 | healthy |          None |        None |       4303(4390) |      890(932) |    189440(366592) |   70656(171008) |  1687552(1917952) |   9.5.2 | 2017-06-15 14:16:21 |
| 29  | 南京中医药大学           | CAS0530000139 | healthy |          None |        None |       2811(2846) |      451(506) |     27648(107520) |    12288(83968) |    944128(955392) |   9.5.2 | 2017-06-15 14:19:36 |
| 33  | 南京大学                 | CAS0530000157 | healthy |          None |        None |     13627(13843) |    2716(2890) |   606357(2682450) | 144054(1065094) |  6183936(7887872) |   9.5.2 | 2017-06-15 14:18:06 |
| 60  | 南京理工大学             | CAS0530000216 | healthy |          None |        None |       5412(5699) |    1341(1341) |   724903(1185344) |   98092(683768) |  3734528(4180992) |   9.5.2 | 2017-06-15 14:16:28 |
| 71  | 南京航空航天大学新校区   | CAS0530000231 | healthy |          None |        None |       4642(4667) |    1279(1345) |    195719(301132) |   42467(215422) |  1984512(2298880) |   9.5.2 | 2017-06-15 14:16:56 |
| 104 | 南京工程学院             | CAS0530000281 | healthy |          None |        None |       2249(2406) |      311(405) |     63488(258048) |    6144(222208) |  1405952(1689600) |  9.2.02 | 2017-06-15 14:19:04 |
| 128 | 南京农业大学             | CAS0530000312 | healthy |          None |        None |       6811(6871) |    1048(1055) |    165724(267389) |   74446(286764) |    513024(618496) |   9.5.2 | 2017-06-15 14:16:32 |
+-----+--------------------------+---------------+---------+---------------+-------------+------------------+---------------+-------------------+-----------------+-------------------+---------+---------------------+

```
or

```shell
$ fxoss cds-list 147

Get cds list from api successfully
+-----+---------------+---------------+--------------------------------------------+---------------+-------------+------------------+---------------+-------------------+-----------------+-------------------+---------+---------------------+
| #   | company       | sn            |                                     status | license_start | license_end | online_user(max) | hit_user(max) | service_kbps(max) | cache_kbps(max) | monitor_kbps(max) | version |          updated_at |
+-----+---------------+---------------+--------------------------------------------+---------------+-------------+------------------+---------------+-------------------+-----------------+-------------------+---------+---------------------+
| 114 | 测试机-办公网 | CAS0510000147 | warn: cnc_http_2 offline, cnc_live offline |          None |        None |             2(2) |          0(0) |           0(1024) |      1024(2048) |      15360(54272) |   9.5.3 | 2017-06-15 14:17:23 |
+-----+---------------+---------------+--------------------------------------------+---------------+-------------+------------------+---------------+-------------------+-----------------+-------------------+---------+---------------------+
```



### fxoss cds-show <sn>

SHOW detail CDS information of CAS0510000147

```shell
$ fxoss cds-show CAS0510000147
GET cds detail information with 'CAS0510000147' success'
+---------------+---------------+--------------------------------------------+---------------------+---------------------+------------------+---------------+-------------------+-----------------+-------------------+---------+---------------------+
| company       | sn            | status                                     | license_start       | license_end         | online_user(max) | hit_user(max) | service_kbps(max) | cache_kbps(max) | monitor_kbps(max) | version | updated_at          |
+---------------+---------------+--------------------------------------------+---------------------+---------------------+------------------+---------------+-------------------+-----------------+-------------------+---------+---------------------+
| 测试机-办公网 | CAS0510000147 | warn: cnc_http_2 offline, cnc_live offline | 2017-02-15 00:00:00 | 2018-12-31 00:00:00 | 2(2)             | 0(0)          | 0(1024)           | 1024(2048)      | 15360(54272)      | 9.5.3   | 2017-06-15 14:07:22 |
+---------------+---------------+--------------------------------------------+---------------------+---------------------+------------------+---------------+-------------------+-----------------+-------------------+---------+---------------------+
CDS 'CAS0510000147' Nodes list
+---+---------------+------------+---------+---------------+-----------------+-------------------+
| # | sn            | type       | status  | hit_user(max) | cache_kbps(max) | service_kbps(max) |
+---+---------------+------------+---------+---------------+-----------------+-------------------+
| 1 | CAS0510000147 | icache     | healthy | 0(0)          | 1024(2048)      | 0(1024)           |
| 2 | VAS0510000147 | cnc_demand | healthy | 0(0)          | 0(0)            | 0(0)              |
| 3 | VBS0510000147 | cnc_live   | offline | 0(0)          | 0(0)            | 0(0)              |
| 4 | VCS0510000147 | cnc_http_2 | offline | 0(0)          | 0(0)            | 0(0)              |
| 5 | VDS0510000147 | xingyu     | healthy | 0(0)          | 0(0)            | 0(0)              |
+---+---------------+------------+---------+---------------+-----------------+-------------------+
```


### fxoss cds-login <sn> -u username -p password -t timeout -r retry

You can use default ssh username and password which stored in
your environ variables or input them by using command
 `-u username -p password` to login a CDS asset

Parameters:

| parameter         | description                            | example       |
| ----------------- | -------------------------------------- | ------------- |
| `-u` `--username` | SSH login username                     | `-u albert `  |
| `-p` `--password` | SSH login password (with quotes )      | `-p 'xxxxxx'` |
| `-t` `-timeout`   | SSH login timeout (default: 5secs)     | `-t 10`       |
| `-r` `-retry`     | SSH login retry times(default: 3times) | `-r 3`        |



Example:

login CDS use default username and password
```
$fxoss cds-login CAS0510000147 -t 10 -r 5  # login CDS assert CAS0510000147 by timeout 10 secs and 5 times retry.
Get icaches 'CAS0510000147' ports successfully
2017-06-15 14:41:30,544 WARNING  [fxoss] Unknown host key
Last login: Thu Jun 15 14:18:03 2017 from 192.168.2.21
[root@test94 ~ 14:41:31]#
```

OR login CDS with username and password

```
$ fxoss cds-login CAS0510000147 -u root -p 'xxxxxx'

```
__note:__

1. password must be included with quotes as 'password'
2. At first time you login the CDS, you should input password
   in interactive shell.

Use `exit` to quit the ssh session
