import { execSync } from "node:child_process";

function getRevision(): { full: string; short: string } {
  try {
    const full = execSync("jj log --no-graph -r @ -T change_id", {
      encoding: "utf8",
      stdio: ["ignore", "pipe", "ignore"],
    }).trim();
    return { full, short: full.slice(0, 12) };
  } catch {
    return { full: "", short: "" };
  }
}

export const revision = getRevision();
