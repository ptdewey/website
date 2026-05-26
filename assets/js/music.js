import { xrpc, listRecords, DID } from "./atproto.js";
import { fadeIn, fadeInStagger } from "./anim.js";

const NS = "fm.teal.alpha";
const PDSLS = "https://pdsls.dev/at/";

function fmtStamp(iso) {
  const d = new Date(iso);
  const now = new Date();
  const sameDay =
    d.getFullYear() === now.getFullYear() &&
    d.getMonth() === now.getMonth() &&
    d.getDate() === now.getDate();
  const time = d.toLocaleTimeString("en-GB", {
    hour: "2-digit",
    minute: "2-digit",
  });
  if (sameDay) return time;
  const date = d.toLocaleDateString("en-CA", {
    month: "2-digit",
    day: "2-digit",
  });
  return `${date} ${time}`;
}

function renderPlay(tmpl, uri, play) {
  const el = tmpl.content.cloneNode(true);

  const link = el.querySelector(".play-track");
  link.textContent = play.trackName;
  if (play.originUrl) link.href = play.originUrl;
  else link.removeAttribute("target");

  const timeEl = el.querySelector("time");
  timeEl.textContent = fmtStamp(play.playedTime);
  timeEl.setAttribute("datetime", play.playedTime);

  el.querySelector(".play-artist").textContent =
    play.artists?.map((a) => a.artistName).join(", ") ?? "";

  const releaseEl = el.querySelector(".play-release");
  if (play.releaseName) releaseEl.textContent = play.releaseName;
  else releaseEl.remove();

  const rkeyEl = el.querySelector(".play-rkey");
  const rkey = uri.split("/").pop() ?? "";
  const label = document.createElement("span");
  label.textContent = "rec  ";
  const a = document.createElement("a");
  a.href = `${PDSLS}${DID}/${NS}.feed.play/${rkey}`;
  a.target = "_blank";
  a.rel = "noopener";
  a.textContent = rkey;
  rkeyEl.append(label, a);

  return el;
}

let npTimer = null;

function hasNowPlaying(status) {
  const item = status?.item;
  if (!item) return false;
  return Boolean(
    item.trackName ||
    item.releaseName ||
    item.artists?.some((artist) => artist.artistName),
  );
}

function hideNowPlaying() {
  document.getElementById("now-playing")?.classList.add("hidden");
}

async function loadStatus() {
  if (npTimer) {
    clearTimeout(npTimer);
    npTimer = null;
  }

  try {
    const rec = await xrpc("com.atproto.repo.getRecord", {
      repo: DID,
      collection: `${NS}.actor.status`,
      rkey: "self",
    }).catch(() => null);
    const status =
      rec && typeof rec === "object" && "value" in rec ? rec.value : rec;
    const container = document.getElementById("now-playing");
    if (!hasNowPlaying(status)) {
      hideNowPlaying();
      return;
    }

    const expiryTime = Number(status.expiry);
    const now = Math.floor(Date.now() / 1000);
    if (expiryTime && expiryTime < now) {
      hideNowPlaying();
      return;
    }

    document.getElementById("np-track").textContent = status.item.trackName;
    document.getElementById("np-artist").textContent =
      status.item.artists?.map((a) => a.artistName).join(", ") ?? "";
    const releaseEl = document.getElementById("np-release");
    if (status.item.releaseName) {
      releaseEl.textContent = status.item.releaseName;
      releaseEl.classList.remove("hidden");
    } else {
      releaseEl.textContent = "";
    }

    if (expiryTime) {
      const refreshInMs = Math.max(1000, (expiryTime - now + 1) * 1000);
      npTimer = setTimeout(() => {
        npTimer = null;
        loadStatus();
        loadPlays();
      }, refreshInMs);
    }

    const wasHidden = container.classList.contains("hidden");
    container.classList.remove("hidden");
    if (wasHidden) fadeIn(container);
  } catch (e) {
    hideNowPlaying();
    console.error("Failed to load status:", e);
  }
}

async function loadPlays() {
  const container = document.getElementById("play-list");
  const tmpl = document.getElementById("play-tmpl");
  const setMessage = (text) => {
    const p = document.createElement("p");
    p.className = "log-msg";
    p.textContent = text;
    container.replaceChildren(p);
  };

  try {
    const records = await listRecords(`${NS}.feed.play`, 15);
    const sorted = records.sort(
      (a, b) =>
        new Date(b.value.playedTime).getTime() -
        new Date(a.value.playedTime).getTime(),
    );
    const recent = sorted.slice(0, 10);

    if (recent.length === 0) {
      setMessage("No plays yet.");
      return;
    }

    container.replaceChildren(
      ...recent.map((r) => renderPlay(tmpl, r.uri, r.value)),
    );
    fadeInStagger(container.children);
  } catch (e) {
    setMessage("Could not load plays.");
    console.error("Failed to load plays:", e);
  }
}

loadStatus();
loadPlays();
