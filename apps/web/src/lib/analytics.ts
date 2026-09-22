/**
 * Minimal analytics interface so pages can call `track(...)` without the
 * template committing to a specific provider. The default implementation
 * just logs to the console — no external network calls, so it's safe in a
 * local/offline dev environment and adds no third-party script to the CSP.
 *
 * To wire up a real provider later:
 *   - PostHog: swap the body of `track`/`identify` for `posthog.capture(...)`
 *     and `posthog.identify(...)`, and load the PostHog snippet in index.html.
 *   - Plausible: Plausible is mostly automatic (script tag) but custom events
 *     go through `window.plausible(eventName, { props })`.
 * Keep the same `track`/`identify` call sites in your pages either way — only
 * this file should need to change.
 */

type EventProperties = Record<string, string | number | boolean>;

export interface Analytics {
  track: (event: string, properties?: EventProperties) => void;
  identify: (userId: string, traits?: EventProperties) => void;
}

const consoleAnalytics: Analytics = {
  track(event, properties) {
    // eslint-disable-next-line no-console
    console.info("[analytics] track", event, properties ?? {});
  },
  identify(userId, traits) {
    // eslint-disable-next-line no-console
    console.info("[analytics] identify", userId, traits ?? {});
  },
};

export const analytics: Analytics = consoleAnalytics;
