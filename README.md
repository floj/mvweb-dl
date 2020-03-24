# MvWeb-DL

Query [mediathekviewweb.de](https://mediathekviewweb.de/) for your shows and download new ones.

## Motivation
I had a bunch of scripts I triggered via Cron. Some of them 
- queried the [ARD Mediathek GraphQL endpoints](https://api.ardmediathek.de/public-gateway) and downloaded new _Die Sendung mit der Maus_ episodes.
- queries the [ZDF RSS feeds](https://www.zdf.de/rss/podcast/video/zdf/comedy/die-anstalt) and downloaded new _Die Anstalt_ episodes (via `youtube-dl`).
- scaped parts of the [ZDF website](https://www.zdf.de/) and downloaded _Robin Hood_ episodes

While I kept adding more and more scripts this got a bit messy. And since all shows and episodes are covered by [mediathekviewweb.de](https://mediathekviewweb.de/) I decided to consolidate the scripts and make a decent app out of it.

# Usage
tbd

# Configuration
tbd

# Todos
- [ ] Add tests
- [ ] Add more filters
- [ ] Maybe a CI pipeline?