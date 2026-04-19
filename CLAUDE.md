# CLAUDE.md

Guidance for Claude when working in this repository.

## Project

Patrick Dewey's personal site. Astro 6 static build, Tailwind v4 (CSS-first via `@tailwindcss/vite`), markdown content collections for the blog, vanilla JS islands for the atproto-backed `/music` and `/coffee` pages. Built with `pnpm build`; Nix flake provides `nix run` and `nix develop`.

Content directories:
- `src/content/blog/*.md` — blog posts, schema in `src/content.config.ts`
- `src/pages/` — top-level pages (index, projects, coffee, music, now, blog/*, rss.xml.js)
- `src/layouts/Layout.astro` — shared shell
- `src/lib/revision.ts` — jj change ID read at build time, shown in footer
- `public/` — static assets, per-page vanilla JS scripts (`coffee.js`, `music.js`, `atproto.js`, `bluesky-comments.js`)

## Design Context

### Users

A personal corner of the internet for Patrick Dewey — a backend engineer and atproto enthusiast. Visitors: Patrick himself using it as a logbook; friends/strangers from Bluesky/GitHub/project pages; atproto, coffee, and music niche visitors arriving via `/coffee`, `/music`, and future atproto-integrated surfaces. Not a résumé. Not a lead funnel. A personal research notebook that happens to be public.

### Brand Personality

**Lived-in, mechanical, honest.** The voice of a well-worn engineering notebook — precise where it matters, loose where it doesn't, no marketing veneer. Quiet confidence; minor delight on small details; nothing performative. No hero section, no "about me." The page is the content.

### Aesthetic Direction

**Terminal-first / hacker almanac.** A well-maintained `.plan` file crossed with a printed technical manual. Dense, not airy. Rewards scanning. Every pixel earns its place.

- Monospace is structural (column alignment, ASCII dividers, tables), not decorative sprinkles.
- Earthy dark palette stays: `--background: #24211e`, warm cream primary, green/orange/teal/yellow/red semantic accents mapped to heading levels. Refine with OKLCH; do not replace.
- Pair Iosevka Patrick (body, custom, already installed) with one printed-manual display face (IM Fell English) for the wordmark and post titles — one precious display moment per page, not sprinkled everywhere.
- Surface atproto identity visually: `at://` URIs, DID, rkey, PDS links.

**Anti-references:** generic SaaS landing pages; portfolio-site-as-résumé; over-designed agency sites with scroll-jacking/awwwards energy.

### Design Principles

1. **The page is the manifest.** No hero, no intro paragraph, no "welcome."
2. **Density is a feature.** Homepage reads like `cat ~/.plan` — recent brew, recent play, latest post, `/now` headline.
3. **Monospace alignment is structural.** Columns should actually align. Dates, rkeys, and revision IDs sit in fixed-width slots.
4. **Semantic color over decorative color.** The h1→h5 mapping to teal/orange/red/yellow is a system, not variety-for-its-own-sake.
5. **Small honest details beat flashy gestures.** jj revision in the footer is the right size of flex.

## Working Conventions

- Prefer editing existing files to creating new ones.
- No marketing copy, no emoji sprinkles, no "Welcome to..." openers.
- New build-time data fetches should live in `src/lib/` with graceful fallbacks — the homepage rendering must never fail because a PDS request did.
- Match spacing to the 4pt scale already in `src/styles/app.css`. Use `gap` over margins.
- If adding a third-party component or integration, first ask: could this be a plain `<script is:inline>` or a 20-line helper instead?
