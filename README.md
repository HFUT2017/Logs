# Logs
This is a project to practice kafka, etcd, zookeeper, elasticsearch and docker.

首先，本项目主要是作为一个练习项目，可以用来入门学习kafka,etcd,zookeeper,elasticsearch和docker这几种组件。

我实现本项目的过程中的环境如下：

1.宿主机是Windows10的环境，我下载并VMware虚拟机，使用的是Centos7环境。

2.在虚拟机当中使用docker配置了各种各样的组件，这个过程也比较复杂和麻烦，经常因为虚拟机、端口映射、配置文件等原因导致组件无法使用。
但最后还是不断搜集资料完成了docker下的环境配置，在这个过程中也对docker命令有了比较深刻的理解.以下是我docker下的组件版本：

f5aa2953e79f   kibana:7.6.2          "/usr/local/bin/dumb…"   4 days ago   Up 3 hours   0.0.0.0:5601->5601/tcp, :::5601->5601/tcp                                              kibana

801e12a6d095   wurstmeister/kafka    "start-kafka.sh"         5 days ago   Up 3 hours   0.0.0.0:9092->9092/tcp, :::9092->9092/tcp                                              kafka

1ae107d29dfa   zookeeper             "/docker-entrypoint.…"   5 days ago   Up 3 hours   2888/tcp, 3888/tcp, 0.0.0.0:2181->2181/tcp, :::2181->2181/tcp, 8080/tcp                zookeeper

f2e63e7fce41   redis:4.0             "docker-entrypoint.s…"   6 days ago   Up 3 hours   0.0.0.0:6379->6379/tcp, :::6379->6379/tcp                                              redis

457e835f96b2   elasticsearch:7.6.2   "/usr/local/bin/dock…"   6 days ago   Up 2 hours   0.0.0.0:9200->9200/tcp, :::9200->9200/tcp, 0.0.0.0:9300->9300/tcp, :::9300->9300/tcp   elasticsearch

a5dbcab47862   mysql:5.7             "docker-entrypoint.s…"   7 days ago   Up 3 hours   0.0.0.0:3306->3306/tcp, :::3306->3306/tcp, 33060/tcp mysql


3.本项目主要包含两个部分，第一个部分是logagent，它主要的作用是不断的读取日志，并将日志发送到kafka当中，其中主要涉及的组件有tailf,kafka和etcd。tailf的作用是用来读取日志，etcd的作用类似于zookeeper
可以用来进行项目配置文件的管理，比如读取日志的路径，可以使得项目进行热更新。而kafka是核心组件，tailf将读到的日志交给kafka的producer。在这个模块当中，也可以学习到go语言当中比较多的channel和context的知识。

4.项目的第二个部分是logtransfer，它的主要作用是用kafka的consurmer来进行消费，将消费到的数据送入elasticsearch当中，之后再用kibana进行数据的可视化操作。

总结：
在这个项目当中，练习了go当中context和channel的使用，学习了kafka，etcd，elasticsearch，kibana和docker的使用，还算有所收获。
