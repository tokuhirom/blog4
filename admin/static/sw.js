// Service Worker for Blog4 Admin
// Required for Web Share Target API

const CACHE_NAME = 'blog4-admin-v2';
// Only cache static assets, not dynamic API endpoints
const urlsToCache = ['/admin/static/admin.css'];

// Install event - cache critical resources
self.addEventListener('install', (event) => {
    console.log('[SW] Installing service worker');
    event.waitUntil(
        caches
            .open(CACHE_NAME)
            .then((cache) => {
                console.log('[SW] Caching app shell');
                return cache.addAll(urlsToCache);
            })
            .then(() => self.skipWaiting()),
    );
});

// Activate event - clean up old caches
self.addEventListener('activate', (event) => {
    console.log('[SW] Activating service worker');
    event.waitUntil(
        caches
            .keys()
            .then((cacheNames) => {
                return Promise.all(
                    cacheNames.map((cacheName) => {
                        if (cacheName !== CACHE_NAME) {
                            console.log('[SW] Deleting old cache:', cacheName);
                            return caches.delete(cacheName);
                        }
                        return Promise.resolve();
                    }),
                );
            })
            .then(() => self.clients.claim()),
    );
});

// Fetch event - network-first for dynamic content, cache-first for static assets
self.addEventListener('fetch', (event) => {
    const url = new URL(event.request.url);

    // Never cache admin API endpoints - always use network
    if (
        url.pathname.startsWith('/admin/entries') ||
        url.pathname.startsWith('/admin/share-target') ||
        url.pathname.startsWith('/admin/upload')
    ) {
        event.respondWith(fetch(event.request));
        return;
    }

    // For static assets, use cache-first strategy
    event.respondWith(
        caches.match(event.request).then((response) => {
            // Cache hit - return response
            if (response) {
                return response;
            }

            // Clone the request
            const fetchRequest = event.request.clone();

            return fetch(fetchRequest).then((response) => {
                // Only cache static assets (CSS, images, etc.)
                if (
                    response &&
                    response.status === 200 &&
                    url.pathname.startsWith('/admin/static/')
                ) {
                    const responseToCache = response.clone();
                    caches.open(CACHE_NAME).then((cache) => {
                        cache.put(event.request, responseToCache);
                    });
                }

                return response;
            });
        }),
    );
});
