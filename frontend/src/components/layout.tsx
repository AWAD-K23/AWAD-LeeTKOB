import { useEffect, useState, type ReactNode } from "react";
import { Outlet, useLocation } from "react-router-dom";
import { Link, NavLink } from "@/components/navigation";
import {
  BookOpen,
  Braces,
  CheckCircle2,
  ChevronRight,
  CircleAlert,
  Code2,
  FolderOpen,
  LoaderCircle,
  LogOut,
  Moon,
  Sun,
  Users,
} from "lucide-react";
import type { User } from "@/lib/api";

export function Brand() {
  return (
    <Link to="/" className="brand" aria-label="LeeTKOB home">
      <span className="brand-mark">
        <Braces size={22} />
      </span>
      <span>
        LeeT<span className="brand-accent">KOB</span>
      </span>
    </Link>
  );
}

export function Shell({ user, logout }: { user: User; logout: () => void }) {
  const location = useLocation();
  const [light, setLight] = useState(
    () => localStorage.getItem("leetkob-theme") === "light",
  );

  useEffect(() => {
    document.documentElement.dataset.theme = light ? "light" : "dark";
    localStorage.setItem("leetkob-theme", light ? "light" : "dark");
  }, [light]);

  return (
    <>
      <a className="skip-link" href="#content">
        Skip to content
      </a>
      <header className="topbar">
        <Brand />
        <nav className="main-nav" aria-label="Main navigation">
          <NavLink to="/classes">
            <FolderOpen size={16} />
            Classes
          </NavLink>
          {user.role !== "admin" && (
            <NavLink to="/labs">
              <BookOpen size={16} />
              Labs
            </NavLink>
          )}
          {user.role === "admin" && (
            <NavLink to="/people">
              <Users size={16} />
              People
            </NavLink>
          )}
          {user.role === "student" && (
            <NavLink to="/submissions">
              <Code2 size={16} />
              My submissions
            </NavLink>
          )}
        </nav>
        <div className="account-actions">
          <button
            className="icon-button"
            onClick={() => setLight(!light)}
            aria-label={light ? "Use dark theme" : "Use light theme"}
          >
            {light ? <Moon size={18} /> : <Sun size={18} />}
          </button>
          <span className="role-chip">
            {user.role === "ta" ? "Teaching assistant" : user.role}
          </span>
          <button
            className="icon-button"
            onClick={logout}
            aria-label="Sign out"
            title="Sign out"
          >
            <LogOut size={18} />
          </button>
        </div>
      </header>
      <main id="content" key={location.pathname} className="page-enter">
        <Outlet />
      </main>
      <footer className="site-footer">
        <span>LeeTKOB</span> A better place to learn C++.
      </footer>
    </>
  );
}

export function PageHeader({
  eyebrow,
  title,
  description,
  action,
}: {
  eyebrow?: string;
  title: string;
  description: string;
  action?: ReactNode;
}) {
  return (
    <div className="page-heading">
      <div>
        {eyebrow && <p className="eyebrow">{eyebrow}</p>}
        <h1>{title}</h1>
        <p className="muted">{description}</p>
      </div>
      {action}
    </div>
  );
}

export function Breadcrumb({
  parent,
  to,
  current,
}: {
  parent: string;
  to: string;
  current: string;
}) {
  return (
    <nav className="breadcrumb" aria-label="Breadcrumb">
      <Link to={to}>{parent}</Link>
      <ChevronRight size={14} />
      <span>{current}</span>
    </nav>
  );
}

export function EmptyState({
  title,
  children,
  action,
}: {
  title: string;
  children: ReactNode;
  action?: ReactNode;
}) {
  return (
    <div className="empty-state">
      <div className="empty-icon">
        <FolderOpen size={28} />
      </div>
      <h2>{title}</h2>
      <p className="muted">{children}</p>
      {action}
    </div>
  );
}

export function Feedback({
  error,
  retry,
}: {
  error: Error | null;
  retry?: () => void;
}) {
  if (!error) return null;

  return (
    <div className="error-banner" role="alert">
      <CircleAlert size={18} />
      <span>{error.message}</span>
      {retry && (
        <button className="btn btn-small" onClick={retry}>
          Retry
        </button>
      )}
    </div>
  );
}

export function Loading() {
  return (
    <div className="loading-state" role="status">
      <LoaderCircle className="spin" size={22} />
      <span>Loading your workspace…</span>
    </div>
  );
}

export function Verdict({
  status,
  verdict,
}: {
  status: string;
  verdict?: string;
}) {
  const accepted = verdict === "Accepted";
  const pending = !["completed", "failed"].includes(status);

  return (
    <span
      className={`verdict ${pending ? "pending" : accepted ? "accepted" : "failed"}`}
    >
      {pending ? (
        <LoaderCircle size={14} className="spin" />
      ) : accepted ? (
        <CheckCircle2 size={14} />
      ) : (
        <CircleAlert size={14} />
      )}
      {verdict || (status === "pending" ? "Queued" : status)}
    </span>
  );
}
