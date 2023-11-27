#!/bin/sh -eu
# script to turn single file into pdf with template

convert_dir() {
  input_dir="$1"
  if [ ! -d "$input_dir" ]; then
    echo "$input_dir not a dir :("
  fi

  find "$input_dir" -type f -name "*.md" | while read -r md_file; do
    dir_name=$(dirname "$md_file")
    file_base=$(basename "$md_file" ".md")

    pandoc --template=cmd/md2tex/template.tex -o "$dir_name"/"$file_base".pdf "$md_file"
    echo "$md_file -> $dir_name/$file_base.pdf"
  done
}

convert_dir "$1"
