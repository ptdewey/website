const reducedMotion = () =>
  typeof matchMedia === "function" &&
  matchMedia("(prefers-reduced-motion: reduce)").matches;

export function fadeIn(el, duration = 240, delay = 0) {
  if (reducedMotion()) return;
  el.animate([{ opacity: 0 }, { opacity: 1 }], {
    duration,
    delay,
    easing: "ease-out",
    fill: "backwards",
  });
}

export function fadeInStagger(elements, step = 35, duration = 240) {
  let i = 0;
  for (const el of elements) {
    fadeIn(el, duration, i * step);
    i++;
  }
}

export function flipAnimate(elements, mutate, duration = 260) {
  if (reducedMotion()) {
    mutate();
    return;
  }
  const first = new Map();
  for (const el of elements) first.set(el, el.getBoundingClientRect());
  mutate();
  for (const [el, f] of first) {
    if (!el.isConnected) continue;
    const l = el.getBoundingClientRect();
    const dx = f.left - l.left;
    const dy = f.top - l.top;
    if (dx === 0 && dy === 0) continue;
    el.animate(
      [
        { transform: `translate(${dx}px, ${dy}px)` },
        { transform: "translate(0, 0)" },
      ],
      { duration, easing: "ease-out" },
    );
  }
}

export async function swap(el, update, out = 140, in_ = 200) {
  if (reducedMotion()) {
    update();
    return;
  }
  const outAnim = el.animate([{ opacity: 1 }, { opacity: 0 }], {
    duration: out,
    easing: "ease-out",
    fill: "forwards",
  });
  await outAnim.finished;
  update();
  outAnim.cancel();
  el.animate([{ opacity: 0 }, { opacity: 1 }], {
    duration: in_,
    easing: "ease-out",
  });
}
