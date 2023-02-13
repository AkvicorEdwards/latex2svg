package main

import "fmt"

// 使用 \mbox{中文} 来插入中文
func latexMath(s string) string {
	return fmt.Sprintf(`\documentclass{article}
\usepackage[UTF8]{CJK}
\begin{document}
\thispagestyle{empty}
\begin{CJK*}{UTF8}{gbsn}
$
%s
$
\end{CJK*}
\end{document}`, s)
}

func latexCJK(s string) string {
	return fmt.Sprintf(`\documentclass{article}
\usepackage[UTF8]{CJK}
\begin{document}
\thispagestyle{empty}
\begin{CJK*}{UTF8}{gbsn}
%s
\end{CJK*}
\end{document}`, s)
}

func latexDoc(s string) string {
	return fmt.Sprintf(`\documentclass{article}
\usepackage[UTF8]{CJK}
\begin{document}
%s
\end{document}`, s)
}

// ?crop=[0/1] crop pdf
// ?type=[pdf/svg/png/jpg]
// ?transp=[0/1] transparent background instead of white (PNG)
func scriptBase() string {
	return `#!/bin/sh
cd ` + config.Path.Temp + `
# $1 latex file
# $2 output file type

tex=$1
crop=$2 # [0/1] crop pdf
tp=$3 # [svg/pdf/png/jpg]
tsp=$4 # [0/1] transp for png

echo "" | pdflatex "$tex" > /dev/null 2>&1

rm "${tex}.aux" > /dev/null 2>&1
rm "${tex}.log" > /dev/null 2>&1

if [ "$crop" = "1" ]; then
  pdfcrop "${tex}.pdf" "${tex}.pdf" > /dev/null 2>&1
fi

mv "${tex}.pdf" "${tex}" > /dev/null 2>&1

if [ "$tp" = "pdf" ]; then
  exit
elif [ "$tp" = "svg" ]; then
  pdftocairo -svg "${tex}" > /dev/null 2>&1
  mv "${tex}.svg" "${tex}" > /dev/null 2>&1
  exit
elif [ "$tp" = "jpg" ]; then
  pdftocairo -singlefile -jpeg "${tex}" > /dev/null 2>&1
  mv "${tex}.jpg" "${tex}" > /dev/null 2>&1
  exit
elif [ "$tp" = "png" ]; then
  if [ "$tsp" = "1" ]; then
    pdftocairo -singlefile -png -transp "${tex}" > /dev/null 2>&1
  else
    pdftocairo -singlefile -png "${tex}" > /dev/null 2>&1
  fi
  mv "${tex}.png" "${tex}" > /dev/null 2>&1
  exit
else
  exit
fi
`
}
