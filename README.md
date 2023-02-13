# Render LaTeX

将latex转换为`pdf`、`svg`、`png`、`jpg`，支持中文

可以通过此程序提供的接口，使博客或wiki支持数学公式

# API使用方法

通过GET请求获取，`?key={}&base={}&crop={}&type={}&transp={}&latex={}`。对应参数的取值请看下方的各参数介绍

## 参数

### key

若网站需要正确的key才能响应请求，则填入此项

### latex

期望渲染的latex代码

### crop

是否裁剪，`1`裁剪、`0`不裁剪

### type

latex代码的期望渲染格式，支持`pdf`、`svg`、`png`、`jpg`

### transp

当渲染格式为`png`时，是否将背景设为透明，`1`透明、`0`不透明

### base

latex嵌入的模板，分为`math`、`cjk`、`doc`、`empty`。

`{YOUR LATEX CODE}`将被替换为latex代码

1. `math`

```
\documentclass{article}
\usepackage[UTF8]{CJK}
\begin{document}
\thispagestyle{empty}
\begin{CJK*}{UTF8}{gbsn}
$
{YOUR LATEX CODE}
$
\end{CJK*}
\end{document}
```

2. `cjk`

```
\documentclass{article}
\usepackage[UTF8]{CJK}
\begin{document}
\thispagestyle{empty}
\begin{CJK*}{UTF8}{gbsn}
{YOUR LATEX CODE}
\end{CJK*}
\end{document}
```

3. `doc`

```
\documentclass{article}
\usepackage[UTF8]{CJK}
\begin{document}
{YOUR LATEX CODE}
\end{document}
```

4. `empty`

```
{YOUR LATEX CODE}
```

# 示例

## 请求信息

- `base`: `math`
- `key`: `zm8oj9R`。当不需要key访问时，此项将被忽略
- `crop`: `1`
- `type`: `svg`
- `transp`: `1`。由于输出并非png，此项将被忽略
- `latex`: 如下

```
\mbox{三角函数}
\left\{
\begin{array}{c}
    \cot{\alpha}=\frac{1}{\tan{\alpha}} \\\\
    \csc{\alpha}=\frac{1}{\sin{\alpha}} \\\\
    \sec{\alpha}=\frac{1}{\cos{\alpha}} \\
\end{array}
\right.
\\\\\\
\left\{
\begin{array}{r}
    \sin^{2}{\alpha}=1-\cos^{2}{\alpha} \\
    \tan^{2}{\alpha}=\sec^{2}{\alpha}-1 \\
    \cot^{2}{\alpha}=\csc^{2}{\alpha}-1 \\
\end{array}
\right.
```

## 请求方式

```go
http://172.16.1.206:8080/?base=math&key=zm8oj9R&crop=1&type=svg&transp=1&latex=\mbox{%E4%B8%89%E8%A7%92%E5%87%BD%E6%95%B0}%20\left\{%20\begin{array}{c}%20\cot{\alpha}=\frac{1}{\tan{\alpha}}%20\\\\%20\csc{\alpha}=\frac{1}{\sin{\alpha}}%20\\\\%20\sec{\alpha}=\frac{1}{\cos{\alpha}}%20\\%20\end{array}%20\right.%20\\\\\\%20\left\{%20\begin{array}{r}%20\sin^{2}{\alpha}=1-\cos^{2}{\alpha}%20\\%20\tan^{2}{\alpha}=\sec^{2}{\alpha}-1%20\\%20\cot^{2}{\alpha}=\csc^{2}{\alpha}-1%20\\%20\end{array}%20\right.
```

## 返回结果

![](res.svg)

## 清除缓存

渲染成功后的svg都会保留下来，以便下次请求相同公式时快速响应。若需要清理时，可以通过发起GET请求来清除。

GET方法，`/clear/cache?key=zm8oj9R`。请求时需要提供正确的key

## 清除临时目录

临时目录是渲染过程中使用的目录，正常情况下每次请求结束后，生成的文件都会自动删除。为防止未知情况下文件被保留，每次程序运行前都会自动清理，也可以手动清理。

GET方法，`/clear/temp?key=zm8oj9R`。请求时需要提供正确的key

# 编译运行

## Debian11

```shell
# latex渲染依赖
apt-get update
apt-get -y install software-properties-common
apt-get update
apt-get -y install --no-install-recommends apt-utils
apt-get -y install git
apt-get -y install texlive
apt-get -y install texlive-extra-utils
apt-get -y install latexmk
apt-get -y install latex-cjk-common
apt-get -y install latex-cjk-chinese
# go的编译环境，golang-1.18在bullseye-backports中
apt-get -y install golang-1.18
ln -s /usr/lib/go-1.18/bin/go /usr/bin/go
ln -s /usr/lib/go-1.18/bin/gofmt /usr/bin/gofmt
# 编译运行，使用git或手动下载zip并解压
git clone https://github.com/AkvicorEdwards/latexrender
cd latexrender
go build
# -c config.ini    指定配置文件路径，默认值config.ini。
# -k abc    指定key的值，默认值为random。若指定配置文件存在时，此选项不生效。
#                      设置为random则第一次运行时随机生成一个key写入配置文件。
#                      设置为empty则设置key为空，请求时不需要提交key。
./latexrender -c config.ini -k random
```

## Docker

### 运行

```shell
docker run -d -p 8080:8080 --name latexrender \
    --restart=always \
    akvicor/latexrender:latest
```

### 编译

```shell
./docker_build.sh
```
