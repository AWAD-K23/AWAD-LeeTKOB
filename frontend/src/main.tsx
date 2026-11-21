// Mounts the placeholder page inside the HTML root element.
import React from "react";
import { createRoot } from "react-dom/client";

// Attach React to the root element defined in index.html.
createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <main>
      <h1>LeeTKOB</h1>
      <p>Project setup is in progress.</p>
    </main>
  </React.StrictMode>,
);
