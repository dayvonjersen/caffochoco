#!/bin/bash

output="../public/app.js"
output_min="../public/app.min.js"
echo "Building "$output"..."

rm -f $output
for file in $(LANG=C ls */*.js); do
    if [[ $file != $output && $file != $output_min ]]; then
        acorn --silent $file
        if [[ ! $? -eq 0 ]]; then
            echo $file
            echo -ne "\033[31mSyntax error.\033[0m\nAborted.\n"
            exit 1
        fi
        cat $file >> $output
    fi
done
sed 's/\r//g' $output > temp && mv temp $output
php -r "file_put_contents('$output', preg_replace('@^\s+return A;.*?\(function\(A\)\{@ms', '', file_get_contents('$output')));"

cat *.js >> $output

acorn --silent $output
if [[ ! $? -eq 0 ]]; then
    echo $output
    echo -ne "\033[31mSyntax error.\033[0m\nAborted.\n"
    exit 1
fi

echo "Minifying "$output" to "$output_min"..."
uglifyjs -c sequences,dead_code,conditionals,booleans,unused,if_return,join_vars,drop_console $output 2>/dev/null | uglifyjs -m --export-all > $output_min 2>/dev/null 
