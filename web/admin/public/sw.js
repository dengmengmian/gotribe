const ADMIN_CACHE_PREFIX = 'go-tribe-admin'

self.addEventListener('install', () => {
  self.skipWaiting()
})

self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches
      .keys()
      .then((cacheNames) =>
        Promise.all(
          cacheNames
            .filter((name) => name.startsWith(ADMIN_CACHE_PREFIX))
            .map((name) => caches.delete(name))
        )
      )
      .then(() => self.registration.unregister())
  )
})
