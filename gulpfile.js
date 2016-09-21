"use strict";
const fs     = require('fs');
const acorn  = require('acorn');
const gulp   = require('gulp');
const uglify = require('uglify-js-harmony');
const glob   = require('glob');
const zalgo  = require('./zalgo.js/zalgo.js');

function lintJS(file) {
    let code = fs.readFileSync(file);
    try {
        acorn.parse(code);
    } catch(e) {
        console.error(file, e.message);
        return false;
    }
    return true;
}

gulp.task('app', () => {
    let files = glob.sync("app/*/*.js");
    files[files.length] = "app/app.js";
    let output = "";
    for(let i = 0; i < files.length; i++) {
        if(!lintJS(files[i])) {
            console.error("Compilation aborted due to errors.");
        }
        output += fs.readFileSync(files[i]).toString();
    }
    fs.writeFileSync("./public/app.js", output);
    try {
        var result = uglify.minify(["./public/app.js"], {
            compress: {
                sequences: true,
                dead_code: true,
                conditionals: true,
                booleans: true,
                unused: true,
                if_return: true,
                join_vars: true,
                drop_console: true,
            }
        });
        fs.writeFileSync("./public/app.min.js", result.code);
    } catch(e) {
        console.error(e.message);
    }
});

gulp.task('webcomponents', () => {
    let files = glob.sync("webcomponents/*.html");
    let output = [];
    files.forEach((file) => {
        output[output.length] = zalgo.Devour(file);
    });
    zalgo.Bundle(output, "public/components.js", "public/components.css");
});


gulp.task('watch', () => {
    gulp.watch('app/**/*.js', ['app']);
    gulp.watch('webcomponents/*.html', ['webcomponents']);
});

gulp.task('default', ['watch']);
