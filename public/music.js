const PDS = 'https://arabica.systems'
const DID = 'did:plc:hm5f3dnm6jdhrc55qp2npdja'
const NS = 'fm.teal.alpha'

function fmtMs(s) {
  const m = Math.floor(s / 60)
  return `${m}:${String(Math.floor(s % 60)).padStart(2, '0')}`
}

function timeAgo(date) {
  const s = Math.floor((Date.now() - new Date(date)) / 1000)
  if (s < 60) return 'just now'
  const m = Math.floor(s / 60)
  if (m < 60) return `${m}m ago`
  const h = Math.floor(m / 60)
  if (h < 24) return `${h}h ago`
  const d = Math.floor(h / 24)
  return `${d}d ago`
}

function renderPlay(tmpl, play) {
  const el = tmpl.content.cloneNode(true)
  const link = el.querySelector('.play-track')
  link.textContent = play.trackName
  if (play.originUrl) link.href = play.originUrl

  const timeEl = el.querySelector('time')
  timeEl.textContent = timeAgo(play.playedTime)
  timeEl.setAttribute('datetime', play.playedTime)

  el.querySelector('.play-artist').textContent = play.artists?.map(a => a.artistName).join(', ') || ''

  const releaseEl = el.querySelector('.play-release')
  if (play.releaseName) { releaseEl.textContent = play.releaseName } else { releaseEl.remove() }

  return el
}

let npTimer = null

async function loadStatus() {
  if (npTimer) { clearInterval(npTimer); npTimer = null }

  try {
    const rec = await ATProto.xrpc(PDS, 'com.atproto.repo.getRecord', {
      repo: DID, collection: `${NS}.actor.status`, rkey: 'self'
    })
    const status = rec?.value
    const container = document.getElementById('now-playing')
    if (!status?.item) { container.classList.add('hidden'); return }

    const startTime = Number(status.time)
    const expiryTime = Number(status.expiry)
    const now = Math.floor(Date.now() / 1000)
    if (expiryTime && expiryTime < now) { container.classList.add('hidden'); return }

    const total = status.item.duration || null
    document.getElementById('np-track').textContent = status.item.trackName
    document.getElementById('np-artist').textContent = status.item.artists?.map(a => a.artistName).join(', ') || ''
    const releaseEl = document.getElementById('np-release')
    if (status.item.releaseName) { releaseEl.textContent = status.item.releaseName }

    const timeEl = document.getElementById('np-time')
    const updateTime = () => {
      const elapsed = Math.floor(Date.now() / 1000) - startTime
      if (total) {
        const clamped = Math.min(elapsed, total)
        timeEl.textContent = `${fmtMs(clamped)} / ${fmtMs(total)}`
        if (clamped >= total) {
          clearInterval(npTimer)
          npTimer = null
          loadStatus()
          loadPlays()
        }
      } else {
        timeEl.textContent = fmtMs(elapsed)
      }
    }
    updateTime()
    npTimer = setInterval(updateTime, 1000)

    container.classList.remove('hidden')
  } catch (e) {
    console.error('Failed to load status:', e)
  }
}

async function loadPlays() {
  const container = document.getElementById('play-list')
  const tmpl = document.getElementById('play-tmpl')
  const setMessage = text => container.replaceChildren(
    Object.assign(document.createElement('p'), {className: 'opacity-60', textContent: text})
  )

  try {
    const records = await ATProto.listRecords(PDS, DID, `${NS}.feed.play`, 100)

    const recent = records
      .sort((a, b) => new Date(b.value.playedTime) - new Date(a.value.playedTime))
      .slice(0, 10)

    if (recent.length === 0) { setMessage('No plays yet.'); return }
    container.replaceChildren(...recent.map(r => renderPlay(tmpl, r.value)))
  } catch (e) {
    setMessage('Could not load plays.')
    console.error('Failed to load plays:', e)
  }
}

loadStatus()
loadPlays()
