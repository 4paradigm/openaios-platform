# 如何创立自定义镜像供环境使用

目前环境支持以ssh或者jupyterlab的交互方式进入环境内部进行操作。

在创建环境的时候，可以任意选择public分区中的环境镜像，这些镜像预先安装了环境交互需要的软件包，所以在环境开始后，你就可以通过选择的交互方式操作环境。

但如果你需要用自定义镜像来创建环境，可能需要在镜像中安装对应交互方式需要的软件包。
你在环境创建界面的交互方式选择，会直接影响到镜像中entrypoint中启动的服务。*如果使用了未安装软件包的交互方式，将会导致环境启动失败。*

## ssh
目前，ssh服务需要在镜像中安装sshd软件包

- Ubuntu
  ```shell
  sudo apt-get install openssh-server
  ```
- Centos
  ```shell
  sudo yum install openssh-server
  ```

装完可以确认一下是否`/usr/sbin/sshd`二进制文件存在

## jupyterlab
可以通过以下命令通过pip安装
```shell
pip install jupyterlab
```
安装完后确认jupyterlab二进制文件存在，并且位于`PATH`环境变量目录下。
