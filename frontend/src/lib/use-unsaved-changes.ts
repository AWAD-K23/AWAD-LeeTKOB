import { useEffect, useRef } from "react";
import { useBlocker } from "react-router-dom";

// Cover both in-app navigation and closing a browser tab while authoring a lab.
export function useUnsavedChanges(dirty: boolean) {
  const allowNavigation = useRef(false);
  const blocker = useBlocker(
    ({ currentLocation, nextLocation }) =>
      dirty &&
      !allowNavigation.current &&
      currentLocation.pathname !== nextLocation.pathname,
  );

  useEffect(() => {
    if (blocker.state === "blocked") {
      if (window.confirm("Leave this page and discard your unsaved lab?"))
        blocker.proceed();
      else blocker.reset();
    }
  }, [blocker]);

  useEffect(() => {
    const warn = (event: BeforeUnloadEvent) => {
      if (dirty && !allowNavigation.current) event.preventDefault();
    };

    window.addEventListener("beforeunload", warn);
    return () => window.removeEventListener("beforeunload", warn);
  }, [dirty]);

  return () => {
    allowNavigation.current = true;
  };
}
