import * as esbuild from 'esbuild';
import * as fs from 'fs';
import * as crypto from 'crypto';

const TEMPLATE_PATH = 'templates/htmx_entry_edit.html';
const STATIC_DIR = 'static';
const BUNDLE_PREFIX = 'codemirror-bundle';

// Clean up old bundles
const files = fs.readdirSync(STATIC_DIR);
for (const file of files) {
    if (file.startsWith(BUNDLE_PREFIX) && file.endsWith('.js')) {
        fs.unlinkSync(`${STATIC_DIR}/${file}`);
    }
}

// Build the bundle
const result = await esbuild.build({
    entryPoints: ['src/codemirror-editor.js'],
    bundle: true,
    minify: true,
    format: 'esm',
    write: false,
});

// Generate hash from content
const content = result.outputFiles[0].contents;
const hash = crypto.createHash('md5').update(content).digest('hex').slice(0, 8);
const bundleName = `${BUNDLE_PREFIX}.${hash}.js`;

// Write the bundle
fs.writeFileSync(`${STATIC_DIR}/${bundleName}`, content);
console.log(`Built: ${STATIC_DIR}/${bundleName} (${(content.length / 1024).toFixed(1)}kb)`);

// Update template
let template = fs.readFileSync(TEMPLATE_PATH, 'utf8');
template = template.replace(
    /codemirror-bundle\.[a-z0-9]+\.js/g,
    bundleName
);
fs.writeFileSync(TEMPLATE_PATH, template);
console.log(`Updated: ${TEMPLATE_PATH}`);
