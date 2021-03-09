# goidx(google shared drive index)

Google Shared Drive index and play audio/video via cloudflare worker

more details in [Build a CloudFlare Workers Running GDIndex](https://github.com/maple3142/GDIndex#manual-way)

Supported Commands:

1. dd: mount a shared drive (no args)
2. cd: without args is select mode
3. ls: list
4. play: select mode, use mpv play anything playable (`.ass` Subtitle files will be loaded automatically)

TODO

- [ ] Cache
- [ ] Play Folder as A play list
- [ ] Upload File(size limited by cloudflare worker)
