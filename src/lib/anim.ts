const reducedMotion = () =>
  typeof matchMedia === "function" &&
  matchMedia("(prefers-reduced-motion: reduce)").matches;

export function fadeIn(el: Element, duration = 240, delay = 0): void {
  if (reducedMotion()) return;
  (el as HTMLElement).animate([{ opacity: 0 }, { opacity: 1 }], {
    duration,
    delay,
    easing: "ease-out",
    fill: "backwards",
  });
}

// Fade in a list of elements with a small stagger between each.
export function fadeInStagger(
  elements: Iterable<Element>,
  step = 35,
  duration = 240,
): void {
  let i = 0;
  for (const el of elements) {
    fadeIn(el, duration, i * step);
    i++;
  }
}

// FLIP animation: measure positions before mutate(), run it, animate each
// surviving element from its old position back to its new one. Elements
// removed during mutate() are skipped silently. Good for the slide-down
// effect when inserting into a sorted list.
export function flipAnimate(
  elements: Element[],
  mutate: () => void,
  duration = 260,
): void {
  if (reducedMotion()) { mutate(); return; }
  const first = new Map<Element, DOMRect>();
  for (const el of elements) first.set(el, el.getBoundingClientRect());
  mutate();
  for (const [el, f] of first) {
    if (!el.isConnected) continue;
    const l = el.getBoundingClientRect();
    const dx = f.left - l.left;
    const dy = f.top - l.top;
    if (dx === 0 && dy === 0) continue;
    (el as HTMLElement).animate(
      [
        { transform: `translate(${dx}px, ${dy}px)` },
        { transform: "translate(0, 0)" },
      ],
      { duration, easing: "ease-out" },
    );
  }
}

export async function swap(
  el: Element,
  update: () => void,
  out = 140,
  in_ = 200,
): Promise<void> {
  if (reducedMotion()) {
    update();
    return;
  }
  const node = el as HTMLElement;
  const outAnim = node.animate([{ opacity: 1 }, { opacity: 0 }], {
    duration: out,
    easing: "ease-out",
    fill: "forwards",
  });
  await outAnim.finished;
  update();
  outAnim.cancel();
  node.animate([{ opacity: 0 }, { opacity: 1 }], {
    duration: in_,
    easing: "ease-out",
  });
}
