# Latex to SVG

将latex转换为svg，通过GET方式获取latex字符串并返回svg字符串

可以通过此程序提供的接口，使博客或wiki支持数学公式

# 使用方法

## LaTeX

```
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

GET方法，`/?latex=`

```go
http://172.16.1.206:8080/?latex=\left\{%20\begin{array}{c}%20\cot{\alpha}=\frac{1}{\tan{\alpha}}%20\\\\%20\csc{\alpha}=\frac{1}{\sin{\alpha}}%20\\\\%20\sec{\alpha}=\frac{1}{\cos{\alpha}}%20\\%20\end{array}%20\right.%20\\\\\\%20\left\{%20\begin{array}{r}%20\sin^{2}{\alpha}=1-\cos^{2}{\alpha}%20\\%20\tan^{2}{\alpha}=\sec^{2}{\alpha}-1%20\\%20\cot^{2}{\alpha}=\csc^{2}{\alpha}-1%20\\%20\end{array}%20\right.
```

## 返回结果

![](res.svg)

## 清除缓存

渲染成功后的svg都会保留下来，可以通过发起GET请求来清除。

GET方法，`/clear?key=clear_key`。请求时需要提供正确的key，key的值可以在`def.go`中设置

# 编译运行

以debian11为例

```shell
# latex渲染依赖
apt-get update
apt-get -y install software-properties-common
apt-get update
apt-get install -y --no-install-recommends apt-utils
apt-get -y install texlive
apt-get -y install texlive-extra-utils
apt-get -y install latexmk
apt-get -y install pdf2svg
# go的编译环境，golang-1.18在bullseye-backports中
apt-get -y install golang-1.18
ln -s /usr/lib/go-1.18/bin/go /usr/bin/go
ln -s /usr/lib/go-1.18/bin/gofmt /usr/bin/gofmt
```
