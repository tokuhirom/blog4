-- Insert dummy entries
INSERT INTO entry (path, title, body, visibility, format, published_at, last_edited_at, created_at, updated_at) VALUES
('getting-started', 'Getting Started with Blog4', 'Welcome to Blog4! This is a sample blog post to help you get started.\n\n## Features\n\n- **Markdown Support**: Write your posts in Markdown\n- **Wiki-style Links**: Use [[Entry Title]] to link between posts\n- **Amazon Product Links**: Use [asin:B00EXAMPLE] to embed product information\n- **Image Management**: Automatic image generation for entries\n\n## Creating Your First Post\n\n1. Click on "New Entry" in the admin panel\n2. Enter a title and path for your post\n3. Write your content in Markdown\n4. Set visibility to "public" when ready\n5. Save your post\n\n## Wiki Links\n\nYou can create links to other posts using wiki-style syntax. For example:\n- [[Docker Setup Guide]]\n- [[Markdown Cheatsheet]]\n- [[API Documentation]]\n\nThese links will automatically create new entries if they don\'t exist yet.', 'public', 'mkdn', NOW(), NOW(), NOW(), NOW()),

('docker-setup-guide', 'Docker Setup Guide', '# Docker Setup Guide\n\nThis guide will help you set up the development environment using Docker Compose.\n\n## Prerequisites\n\n- Docker and Docker Compose installed\n- Git\n\n## Quick Start\n\n```bash\n# Clone the repository\ngit clone https://github.com/tokuhirom/blog4.git\ncd blog4\n\n# Copy environment file\ncp .env.example .env\n\n# Start services\ndocker-compose up -d\n```\n\n## Services\n\nThe docker-compose setup includes:\n\n1. **MariaDB**: Database server\n2. **Go Backend**: API server with hot-reload\n3. **Frontend**: Vite dev server\n4. **MinIO**: S3-compatible storage\n\n## Troubleshooting\n\nIf you encounter issues, check the logs:\n\n```bash\ndocker-compose logs -f backend\n```\n\nFor more information, see [[Getting Started]].', 'public', 'mkdn', NOW(), NOW(), NOW(), NOW()),

('markdown-cheatsheet', 'Markdown Cheatsheet', '# Markdown Cheatsheet\n\nHere\'s a quick reference for Markdown syntax supported in Blog4.\n\n## Headers\n\n```markdown\n# H1 Header\n## H2 Header\n### H3 Header\n```\n\n## Text Formatting\n\n- **Bold**: `**text**`\n- *Italic*: `*text*`\n- ~~Strikethrough~~: `~~text~~`\n- `Code`: `` `code` ``\n\n## Lists\n\n### Unordered List\n- Item 1\n- Item 2\n  - Nested item\n\n### Ordered List\n1. First item\n2. Second item\n   1. Nested item\n\n## Links\n\n- External link: `[Google](https://google.com)`\n- Wiki link: `[[Another Entry]]`\n- Amazon product: `[asin:B08N5WRWNW]`\n\n## Code Blocks\n\n```javascript\nfunction hello() {\n  console.log("Hello, Blog4!");\n}\n```\n\n## Tables\n\n| Feature | Supported |\n|---------|----------|\n| Tables  | Yes      |\n| Images  | Yes      |\n| Videos  | No       |\n\n## Images\n\n```markdown\n![Alt text](https://example.com/image.jpg)\n```\n\nFor more formatting options, check the [[API Documentation]].', 'public', 'mkdn', NOW(), NOW(), NOW(), NOW()),

('api-documentation', 'API Documentation', '# API Documentation\n\n## Overview\n\nThe Blog4 API provides RESTful endpoints for managing blog entries.\n\n## Base URL\n\n```\nhttp://localhost:8181/admin/api\n```\n\n## Authentication\n\nIn production, the API uses Basic Authentication. For local development, CORS is enabled for `http://localhost:6173`.\n\n## Endpoints\n\n### List Entries\n\n```http\nGET /entries\n```\n\nReturns a list of all entries.\n\n### Get Entry\n\n```http\nGET /entries/{path}\n```\n\nReturns a specific entry by path.\n\n### Create Entry\n\n```http\nPOST /entries\nContent-Type: application/json\n\n{\n  "path": "my-new-post",\n  "title": "My New Post",\n  "body": "Post content in markdown"\n}\n```\n\n### Update Entry\n\n```http\nPUT /entries/{path}/body\nContent-Type: application/json\n\n{\n  "body": "Updated content"\n}\n```\n\n### Update Visibility\n\n```http\nPUT /entries/{path}/visibility\nContent-Type: application/json\n\n{\n  "visibility": "public"\n}\n```\n\n## File Upload\n\n```http\nPOST /upload\nContent-Type: multipart/form-data\n```\n\nUploads a file to S3-compatible storage.\n\n## Related Documentation\n\n- [[Getting Started]]\n- [[Docker Setup Guide]]\n- [[Markdown Cheatsheet]]', 'public', 'mkdn', NOW(), NOW(), NOW(), NOW()),

('private-draft', 'Private Draft Example', 'This is a private draft that won\'t be visible on the public site.\n\n## Work in Progress\n\nThis entry demonstrates:\n- Private visibility setting\n- Draft content\n- [[Links to Public Entries]]\n\n## TODO\n- [ ] Finish writing content\n- [ ] Add images\n- [ ] Review and publish', 'private', 'mkdn', NOW(), NOW(), NOW(), NOW()),

('sample-with-images', 'Sample Post with Images', '# Beautiful Images in Blog4\n\nThis post demonstrates image handling in Blog4.\n\n## Uploading Images\n\nYou can upload images through the admin interface. They are stored in MinIO (S3-compatible storage) and automatically optimized.\n\n![Sample Image](https://picsum.photos/800/400)\n\n## Image Features\n\n1. **Automatic Optimization**: Images are processed for web delivery\n2. **CDN Support**: Works with CloudFront or other CDNs\n3. **Responsive**: Images adapt to different screen sizes\n\n## Gallery Example\n\n![Nature](https://picsum.photos/400/300?nature)\n![Architecture](https://picsum.photos/400/300?architecture)\n![Technology](https://picsum.photos/400/300?technology)\n\n## Best Practices\n\n- Use descriptive alt text\n- Optimize images before uploading\n- Consider image placement for readability\n\nFor more tips, see [[Markdown Cheatsheet]].', 'public', 'mkdn', NOW() - INTERVAL 1 DAY, NOW() - INTERVAL 1 DAY, NOW() - INTERVAL 1 DAY, NOW()),

('japanese-content', '日本語コンテンツのサンプル', '# 日本語対応について\n\nBlog4は日本語コンテンツを完全にサポートしています。\n\n## 機能\n\n### 全文検索\nMariaDBのn-gramパーサーを使用して、日本語の全文検索に対応しています。\n\n### 文字エンコーディング\n- UTF-8mb4を使用\n- 絵文字も使用可能\n\n### サンプルテキスト\n\n吾輩は猫である。名前はまだ無い。どこで生れたかとんと見当がつかぬ。何でも薄暗いじめじめした所でニャーニャー泣いていた事だけは記憶している。\n\n### リンクの例\n\n他の記事へのリンク：\n- [[Getting Started]]\n- [[Markdown Cheatsheet]]\n\n### コードブロック\n\n```go\nfunc main() {\n    fmt.Println("こんにちは、世界！")\n}\n```\n\n### 表の例\n\n| 機能 | 対応状況 |\n|------|----------|\n| 日本語入力 | 対応 |\n| 全文検索 | 対応 |\n| 絵文字 | 対応 |\n\n詳しくは[[API Documentation]]をご覧ください。', 'public', 'mkdn', NOW() - INTERVAL 2 DAY, NOW() - INTERVAL 2 DAY, NOW() - INTERVAL 2 DAY, NOW()),

('tech-stack', 'Technology Stack', '# Blog4 Technology Stack\n\n## Backend\n\n### Go\n- **Framework**: Chi router\n- **ORM**: SQLC for type-safe SQL\n- **API**: OpenAPI 3.0 with Ogen code generation\n\n### Database\n- **MariaDB 10.11.13**\n- Full-text search with n-gram parser\n- Separate schemas for admin and public access\n\n## Frontend\n\n### Admin Panel\n- **Framework**: React with TypeScript\n- **UI Library**: Material-UI\n- **Build Tool**: Vite\n- **API Client**: Auto-generated from OpenAPI spec\n\n### Features\n- Hot-reload in development\n- Type-safe API calls\n- Responsive design\n\n## Infrastructure\n\n### Storage\n- **MinIO**: S3-compatible object storage\n- Automatic image optimization\n- Backup support\n\n### Development\n- **Docker Compose**: One-command setup\n- **Air**: Hot-reload for Go\n- **TypeSpec**: API specification\n\n## Code Generation Pipeline\n\n1. TypeSpec → OpenAPI spec\n2. OpenAPI → Go server (Ogen)\n3. OpenAPI → TypeScript client (Orval)\n4. SQL → Go code (SQLC)\n\nFor setup instructions, see [[Docker Setup Guide]].', 'public', 'mkdn', NOW() - INTERVAL 3 DAY, NOW() - INTERVAL 3 DAY, NOW() - INTERVAL 3 DAY, NOW());

-- Insert entry links based on wiki-style links in the content
INSERT INTO entry_link (src_path, dst_title) VALUES
('getting-started', 'Docker Setup Guide'),
('getting-started', 'Markdown Cheatsheet'),
('getting-started', 'API Documentation'),
('docker-setup-guide', 'Getting Started'),
('markdown-cheatsheet', 'API Documentation'),
('api-documentation', 'Getting Started'),
('api-documentation', 'Docker Setup Guide'),
('api-documentation', 'Markdown Cheatsheet'),
('private-draft', 'Links to Public Entries'),
('sample-with-images', 'Markdown Cheatsheet'),
('japanese-content', 'Getting Started'),
('japanese-content', 'Markdown Cheatsheet'),
('japanese-content', 'API Documentation'),
('tech-stack', 'Docker Setup Guide');

-- Insert sample entry images (these would normally be generated by the worker)
INSERT INTO entry_image (path, url) VALUES
('getting-started', 'https://picsum.photos/1200/630?random=1'),
('docker-setup-guide', 'https://picsum.photos/1200/630?random=2'),
('markdown-cheatsheet', 'https://picsum.photos/1200/630?random=3'),
('api-documentation', 'https://picsum.photos/1200/630?random=4'),
('sample-with-images', 'https://picsum.photos/1200/630?random=5'),
('japanese-content', 'https://picsum.photos/1200/630?random=6'),
('tech-stack', 'https://picsum.photos/1200/630?random=7');

-- Insert sample Amazon cache data
INSERT INTO amazon_cache (asin, title, image_medium_url, link) VALUES
('B08N5WRWNW', 'Echo (4th Gen) | With premium sound, smart home hub, and Alexa', 'https://m.media-amazon.com/images/I/71JB6hM6Z6L._AC_SL1000_.jpg', 'https://www.amazon.com/dp/B08N5WRWNW'),
('B07FZ8S74R', 'Echo Dot (3rd Gen) - Smart speaker with Alexa', 'https://m.media-amazon.com/images/I/6182S7MYC2L._AC_SL1000_.jpg', 'https://www.amazon.com/dp/B07FZ8S74R');