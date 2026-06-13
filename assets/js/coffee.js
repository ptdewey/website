import { resolveHandle, listRecords, getRecord, refToUri } from "./atproto.js";
import { fadeInStagger } from "./anim.js";

const NS = "social.arabica.alpha";
const DOT = " \u00b7 ";

function formatTime(s) {
  if (!s) return "";
  const m = Math.floor(s / 60);
  return m > 0 ? `${m}:${String(s % 60).padStart(2, "0")}` : `${s}s`;
}

function setOrRemove(el, text) {
  if (!el) return;
  if (text) el.textContent = text;
  else el.remove();
}

function renderRating(container, rating) {
  const n = Math.min(Math.max(rating, 0), 10);
  container.textContent = `★ ${n}/10`;
}

function setParam(el, value) {
  if (!el) return;
  if (!value) {
    el.remove();
    return;
  }
  el.querySelector("dd").textContent = value;
}

function renderPours(container, pours) {
  if (!container || !pours?.length) {
    container?.remove();
    return;
  }

  const label = document.createElement("span");
  label.className = "brew-label";
  label.textContent = "Pours:";
  container.append(label);

  pours.forEach((pour, i) => {
    const pill = document.createElement("span");
    pill.className = "brew-pour";
    const index = document.createElement("span");
    index.className = "brew-pour__index";
    index.textContent = `${i + 1}`;
    pill.append(index);
    if (pour.waterAmount) {
      const water = document.createElement("span");
      water.className = "brew-pour__water";
      water.textContent = `${pour.waterAmount}g`;
      pill.append(water);
    }
    if (pour.timeSeconds) {
      const dot = document.createElement("span");
      dot.className = "brew-pour__dot";
      dot.textContent = "·";
      const time = document.createElement("span");
      time.className = "brew-pour__time";
      time.textContent = formatTime(pour.timeSeconds);
      pill.append(dot, time);
    }
    container.append(pill);
  });
}

function renderBrew(tmpl, brew) {
  const el = tmpl.content.cloneNode(true);

  const bean = el.querySelector(".brew-bean");
  bean.textContent = brew.beanName || "Unknown bean";
  bean.href = brew.brewURL;

  const timeEl = el.querySelector("time");
  timeEl.textContent = new Date(brew.createdAt).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
  timeEl.setAttribute("datetime", brew.createdAt);

  setOrRemove(el.querySelector(".brew-roaster"), brew.roasterName);
  setOrRemove(
    el.querySelector(".brew-sub"),
    [
      brew.beanOrigin,
      brew.beanRoastLevel,
      brew.beanVariety,
      brew.beanProcess,
      brew.coffeeAmount ? `${brew.coffeeAmount}g` : null,
    ]
      .filter(Boolean)
      .join(DOT),
  );

  const water =
    brew.waterAmount ||
    brew.pours?.reduce((s, p) => s + (p.waterAmount ?? 0), 0) ||
    0;
  const timeStr = brew.timeSeconds ? formatTime(brew.timeSeconds) : null;

  const grinderStr = brew.grinderName
    ? `${brew.grinderName}${brew.grindSize ? ` (${brew.grindSize})` : ""}`
    : brew.grindSize
      ? `grind ${brew.grindSize}`
      : null;
  setParam(el.querySelector(".brew-param-grinder"), grinderStr);
  setParam(el.querySelector(".brew-param-water"), water ? `${water}g` : null);
  setParam(
    el.querySelector(".brew-param-temp"),
    brew.temperature ? `${(brew.temperature / 10).toFixed(1)}\u00b0C` : null,
  );
  setParam(el.querySelector(".brew-param-time"), timeStr);
  if (!el.querySelector(".brew-param")) el.querySelector(".brew-params")?.remove();

  const brewerText = brew.brewerName || brew.method || "";
  if (!brewerText) el.querySelector(".brew-brewer-row")?.remove();
  const brewerEl = el.querySelector(".brew-brewer");
  if (brewerEl) brewerEl.textContent = brewerText;

  renderPours(el.querySelector(".brew-pours"), brew.pours);

  setOrRemove(el.querySelector(".brew-notes"), brew.tastingNotes);

  const ratingEl = el.querySelector(".brew-rating");
  if (brew.rating) renderRating(ratingEl, brew.rating);
  else ratingEl.remove();

  return el;
}

async function loadBrews() {
  const container = document.getElementById("brew-list");
  const tmpl = document.getElementById("brew-tmpl");
  const setMessage = (text) => {
    const p = document.createElement("p");
    p.className = "log-msg";
    p.textContent = text;
    container.replaceChildren(p);
  };

  try {
    const did = await resolveHandle();
    const records = await listRecords(`${NS}.brew`, 15, { repo: did });
    const sorted = records.sort(
      (a, b) =>
        new Date(b.value.createdAt).getTime() -
        new Date(a.value.createdAt).getTime(),
    );
    const recent = sorted.slice(0, 10);

    const brews = await Promise.all(
      recent.map(async ({ uri, value }) => {
        const [beanRec, brewerRec, grinderRec] = await Promise.all([
          getRecord(`${NS}.bean`, refToUri(value.beanRef), { repo: did }),
          getRecord(`${NS}.brewer`, refToUri(value.brewerRef), { repo: did }),
          getRecord(`${NS}.grinder`, refToUri(value.grinderRef), { repo: did }),
        ]);
        const roasterRec = await getRecord(
          `${NS}.roaster`,
          refToUri(beanRec?.roasterRef),
          { repo: did },
        );
        const rkey = uri.split("/").pop() ?? "";
        return {
          ...value,
          rkey,
          brewURL: `https://alpha.arabica.social/brews/pdewey.com/${rkey}`,
          beanName: beanRec?.name ?? "",
          beanOrigin: beanRec?.origin ?? "",
          beanRoastLevel: beanRec?.roastLevel ?? "",
          beanProcess: beanRec?.process ?? "",
          beanVariety: beanRec?.variety ?? "",
          roasterName: roasterRec?.name ?? "",
          brewerName: brewerRec?.name ?? "",
          grinderName: grinderRec?.name ?? "",
        };
      }),
    );

    if (brews.length === 0) {
      setMessage("No brews yet.");
      return;
    }

    const grid = document.createElement("div");
    grid.className = "brew-grid";
    grid.replaceChildren(...brews.map((b) => renderBrew(tmpl, b)));
    container.replaceChildren(grid);
    fadeInStagger(grid.children);
  } catch (e) {
    setMessage("Could not load brews.");
    console.error("Failed to load brews:", e);
  }
}

loadBrews();
