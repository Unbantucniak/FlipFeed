// Existing browser profiles keep their saved session and public profiles on upgrade.
for (const suffix of ["accessToken", "user", "publicProfiles"]) {
  const previousKey = `gcfeed.${suffix}`;
  const currentKey = `flipfeed.${suffix}`;
  const value = localStorage.getItem(previousKey);
  if (value !== null) {
    if (localStorage.getItem(currentKey) === null) {
      try {
        localStorage.setItem(currentKey, value);
      } catch (error) {
        // Keep the original if the browser cannot store a second copy.
        if (error.name === "QuotaExceededError") continue;
        throw error;
      }
    }
    localStorage.removeItem(previousKey);
  }
}
