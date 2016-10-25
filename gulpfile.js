"use strict";
const gulp   = require('gulp');
const glob   = require('glob');
const zalgo  = require('./zalgo.js/zalgo.js');

gulp.task('webcomponents', () => {
    let files = glob.sync("webcomponents/*.html");
    let output = [];
    files.forEach((file) => {
        output[output.length] = zalgo.Devour(file);
    });
    zalgo.Bundle(output, "public/components.js", "public/components.css");
});


gulp.task('watch', () => {
    gulp.watch('webcomponents/*.html', ['webcomponents']);
});

gulp.task('default', ['watch']);
