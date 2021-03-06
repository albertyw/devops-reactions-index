const fs = require('fs');
const path = require('path');

const browserify = require('browserify');
require('dotenv').config();
const minifyStream = require('minify-stream');

const inputFile = path.join(__dirname, '..', 'server', 'static', 'js', 'app.js');
const outputFile = path.join(__dirname, '..', 'server', 'static', 'app.js');

browserify(inputFile, {debug: true})
  .transform('unassertify', {global: true})
  .transform('loose-envify')
  .plugin('common-shakeify')
  .plugin('browser-pack-flat/plugin')
  .transform('babelify',  {presets: ['@babel/preset-env']})
  .bundle()
  .pipe(minifyStream({
    mangle: false,
    toplevel: true,
    keep_fnames: true,
    keep_classnames: true,
  }))
  .pipe(fs.createWriteStream(outputFile));
