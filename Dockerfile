FROM golang:1.18.10-bullseye

RUN apt-get update
RUN apt-get -y install texlive
RUN apt-get -y install texlive-extra-utils
RUN apt-get -y install latexmk
RUN apt-get -y install latex-cjk-common
RUN apt-get -y install latex-cjk-chinese
#RUN apt-get -y install git

WORKDIR /wp

#RUN mkdir -p /wp/dl
#RUN cd /wp/dl && git clone https://github.com/AkvicorEdwards/latexrender

RUN mkdir -p /wp/dl/latexrender
COPY . /wp/dl/latexrender

RUN cd /wp/dl/latexrender && go build && mv latexrender /wp/latexrender

CMD ["./latexrender"]
