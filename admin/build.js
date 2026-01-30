import * as esbuild from 'esbuild';
import * as fs from 'fs';
import * as crypto from 'crypto';

const STATIC_DIR = 'static';

const bundles = [
    {
        entryPoint: 'src/codemirror-editor.js',
        prefix: 'codemirror-bundle',
        format: 'esm',
        jsx: false,
        templatePath: 'templates/htmx_entry_edit.html',
    },
    {
        entryPoint: 'src/entry-edit/index.jsx',
        prefix: 'entry-edit-app',
        format: 'esm',
        jsx: true,
        templatePath: 'templates/htmx_entry_edit.html',
    },
    {
        entryPoint: 'src/entry-list/index.jsx',
        prefix: 'entry-list-app',
        format: 'esm',
        jsx: true,
        templatePath: 'templates/htmx_entries.html',
    },
    {
        entryPoint: 'src/login/index.jsx',
        prefix: 'login-app',
        format: 'esm',
        jsx: true,
        templatePath: 'templates/htmx_login.html',
    },
];

for (const bundle of bundles) {
    // Clean up old bundles
    const files = fs.readdirSync(STATIC_DIR);
    for (const file of files) {
        if (file.startsWith(bundle.prefix) && file.endsWith('.js')) {
            fs.unlinkSync(`${STATIC_DIR}/${file}`);
        }
    }

    // Build options
    const buildOptions = {
        entryPoints: [bundle.entryPoint],
        bundle: true,
        minify: true,
        format: bundle.format,
        write: false,
    };

    if (bundle.jsx) {
        buildOptions.jsx = 'automatic';
        buildOptions.jsxImportSource = 'preact';
    }

    // Build the bundle
    const result = await esbuild.build(buildOptions);

    // Generate hash from content
    const content = result.outputFiles[0].contents;
    const hash = crypto.createHash('md5').update(content).digest('hex').slice(0, 8);
    const bundleName = `${bundle.prefix}.${hash}.js`;

    // Write the bundle
    fs.writeFileSync(`${STATIC_DIR}/${bundleName}`, content);
    console.log(`Built: ${STATIC_DIR}/${bundleName} (${(content.length / 1024).toFixed(1)}kb)`);

    // Update template with new hash
    const templatePath = bundle.templatePath;
    let template = fs.readFileSync(templatePath, 'utf8');
    const regex = new RegExp(`${bundle.prefix}\\.[a-zA-Z0-9]+\\.js`, 'g');
    template = template.replace(regex, bundleName);
    fs.writeFileSync(templatePath, template);
    console.log(`Updated: ${templatePath} with ${bundleName}`);
}
