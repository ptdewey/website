import { listRecords, DID } from "./atproto.js";
import { fadeInStagger } from "./anim.js";

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

loadPlays();
