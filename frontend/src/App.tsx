import { useState } from "react";
import {
  createBrowserRouter,
  RouterProvider,
  Navigate,
  Route,
  Routes,
  useNavigate,
} from "react-router-dom";
import {
  QueryClient,
  QueryCache,
  QueryClientProvider,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import { Eye, EyeOff, LockKeyhole, Code2 } from "lucide-react";
import { api, ApiError, type User } from "@/lib/api";
import {
  Brand,
  EmptyState,
  Feedback,
  Loading,
  Shell,
} from "@/components/layout";
import { ClassesPage, ClassPage, ClassForm, PeoplePage } from "@/pages/classes";
import { LabsPage, LabForm } from "@/pages/labs";
import {
  LabWorkspace,
  SubmissionPage,
  SubmissionsPage,
} from "@/pages/workspace";
import "./style.css";

const client = new QueryClient({
  // A protected request reporting an expired session returns to sign-in.
  queryCache: new QueryCache({
    onError: (error) => {
      if (error instanceof ApiError && error.status === 401) {
        client.setQueryData(["session"], null);
        client.removeQueries({
          predicate: (query) => query.queryKey[0] !== "session",
        });
      }
    },
  }),
  defaultOptions: { queries: { retry: false, staleTime: 15_000 } },
});

function Login({ complete }: { complete: () => Promise<unknown> }) {
  const [visible, setVisible] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  return (
    <div className="login-layout">
      <section className="login-story">
        <Brand />
        <div className="login-copy">
          <span className="eyebrow">YOUR C++ LABORATORY</span>
          <h1>
            Small steps.
            <br />
            Better solutions.
          </h1>
          <p>
            A focused space to write code, test your thinking, and keep
            learning.
          </p>
          <div className="code-preview" aria-hidden="true">
            <div className="preview-title">
              <Code2 size={16} />
              main.cpp <span>C++</span>
            </div>
            <pre>
              <span className="syntax-purple">#include</span>
              {" <iostream>\n\n"}
              <span className="syntax-amber">int</span>
              {" main() {\n"}
              {"    std::cout << "}
              <span className="syntax-green">"Hello, possibility."</span>
              {";\n"}
              {"    return 0;\n}"}
            </pre>
          </div>
        </div>
        <p className="muted">Write. Submit. Understand.</p>
      </section>
      <section className="login-form-area">
        <div className="login-card page-enter">
          <span className="login-lock">
            <LockKeyhole size={22} />
          </span>
          <h2>Welcome back</h2>
          <p className="muted">Sign in to your LeeTKOB workspace.</p>
          <Feedback error={error} />
          <form
            onSubmit={async (event) => {
              event.preventDefault();
              const data = new FormData(event.currentTarget);
              setBusy(true);
              setError(null);

              try {
                await api("/auth/login", {
                  email: data.get("email"),
                  password: data.get("password"),
                });
                await complete();
              } catch (cause) {
                setError(cause as Error);
              } finally {
                setBusy(false);
              }
            }}
          >
            <label>
              Email
              <input
                name="email"
                type="email"
                placeholder="you@example.edu"
                autoComplete="username"
                required
              />
            </label>
            <label htmlFor="password">Password</label>
            <div className="password-field">
              <input
                id="password"
                name="password"
                type={visible ? "text" : "password"}
                autoComplete="current-password"
                required
              />
              <button
                type="button"
                className="icon-button"
                onClick={() => setVisible(!visible)}
                aria-label={visible ? "Hide password" : "Show password"}
              >
                {visible ? <EyeOff size={18} /> : <Eye size={18} />}
              </button>
            </div>
            <button className="btn btn-primary" disabled={busy}>
              {busy ? "Signing in…" : "Sign in"}
            </button>
          </form>
          <p className="login-help">
            Need an account? Contact your class administrator.
          </p>
        </div>
      </section>
    </div>
  );
}

function Application() {
  const cache = useQueryClient();
  const navigate = useNavigate();
  const [logoutError, setLogoutError] = useState<Error | null>(null);
  const session = useQuery({
    queryKey: ["session"],
    queryFn: async () => {
      try {
        return await api<User>("/me");
      } catch (error) {
        if (error instanceof ApiError && error.status === 401) return null;
        throw error;
      }
    },
    staleTime: 0,
  });
  const user = session.data;

  if (session.isPending) return <Loading />;
  if (session.error)
    return (
      <div className="page-container">
        <Brand />
        <Feedback error={session.error} retry={() => void session.refetch()} />
      </div>
    );
  if (!user)
    return (
      <Login
        complete={async () => {
          cache.clear();
          await cache.fetchQuery({
            queryKey: ["session"],
            queryFn: () => api<User>("/me"),
          });
          navigate("/", { replace: true, viewTransition: true });
        }}
      />
    );

  const home = user.role === "admin" ? "/classes" : "/labs";
  const deny = (
    <div className="page-container">
      <EmptyState
        title="This page isn't available for your role"
        action={
          <a className="btn" href={home}>
            Return to workspace
          </a>
        }
      >
        Your account can access the destinations in the navigation above.
      </EmptyState>
    </div>
  );

  return (
    <>
      <Feedback error={logoutError} />
      <Routes>
        <Route
          element={
            <Shell
              user={user}
              logout={() => {
                void api("/auth/logout", {})
                  .then(() => {
                    cache.clear();
                    cache.setQueryData(["session"], null);
                    navigate("/login");
                  })
                  .catch(setLogoutError);
              }}
            />
          }
        >
          <Route index element={<Navigate to={home} replace />} />
          <Route path="login" element={<Navigate to={home} replace />} />
          <Route path="classes" element={<ClassesPage user={user} />} />
          <Route
            path="classes/new"
            element={user.role === "admin" ? <ClassForm /> : deny}
          />
          <Route path="classes/:id" element={<ClassPage user={user} />} />
          <Route
            path="people"
            element={user.role === "admin" ? <PeoplePage /> : deny}
          />
          <Route
            path="labs"
            element={user.role !== "admin" ? <LabsPage user={user} /> : deny}
          />
          <Route
            path="classes/:id/labs/new"
            element={user.role === "ta" ? <LabForm /> : deny}
          />
          <Route
            path="labs/new"
            element={user.role === "ta" ? <LabForm /> : deny}
          />
          <Route
            path="labs/:id"
            element={
              user.role !== "admin" ? <LabWorkspace user={user} /> : deny
            }
          />
          <Route
            path="submissions"
            element={user.role === "student" ? <SubmissionsPage /> : deny}
          />
          <Route
            path="submissions/:id"
            element={user.role !== "admin" ? <SubmissionPage /> : deny}
          />
          <Route
            path="*"
            element={
              <div className="page-container">
                <EmptyState
                  title="Page not found"
                  action={
                    <a className="btn" href={home}>
                      Return to workspace
                    </a>
                  }
                >
                  This address does not match a page in your workspace.
                </EmptyState>
              </div>
            }
          />
        </Route>
      </Routes>
    </>
  );
}

export default function App() {
  return (
    <QueryClientProvider client={client}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  );
}

const router = createBrowserRouter([{ path: "*", element: <Application /> }]);
