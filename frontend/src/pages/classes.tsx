import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Link } from "@/components/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  BookOpen,
  Check,
  FolderOpen,
  Plus,
  Search,
  Users,
  X,
} from "lucide-react";
import {
  api,
  type Classroom,
  type Lab,
  type Person,
  type User,
} from "@/lib/api";
import {
  Breadcrumb,
  EmptyState,
  Feedback,
  Loading,
  PageHeader,
} from "@/components/layout";

export function ClassesPage({ user }: { user: User }) {
  const [search, setSearch] = useState("");
  const query = useQuery({
    queryKey: ["classes"],
    queryFn: () => api<Classroom[]>("/classes"),
  });

  return (
    <div className="page-container">
      <PageHeader
        eyebrow="YOUR WORKSPACE"
        title="Classes"
        description="People, problems, and progress. All in one place."
        action={
          user.role === "admin" && (
            <Link className="btn btn-primary" to="/classes/new">
              <Plus size={16} />
              Create class
            </Link>
          )
        }
      />
      <div className="section-toolbar">
        <label className="search-field">
          <Search size={17} />
          <input
            aria-label="Search classes"
            placeholder="Search classes…"
            value={search}
            onChange={(event) => setSearch(event.target.value)}
          />
        </label>
        <span className="muted">{query.data?.length ?? 0} classes</span>
      </div>
      <Feedback error={query.error} retry={() => void query.refetch()} />
      {query.isPending ? (
        <Loading />
      ) : (
        <div className="class-grid">
          {query.data
            ?.filter((item) =>
              item.name.toLowerCase().includes(search.toLowerCase()),
            )
            .map((item) => (
              <Link
                className="class-card"
                key={item.id}
                to={`/classes/${item.id}`}
              >
                <div className="class-card-top">
                  <span className="class-symbol">
                    <FolderOpen size={24} />
                  </span>
                  <span className="small-tag">Class</span>
                </div>
                <h2>{item.name}</h2>
                <p className="muted">Open class workspace</p>
                <div className="class-card-bottom">
                  <Users size={15} />
                  <span>Members & labs</span>
                  <span className="text-link">View class</span>
                </div>
              </Link>
            ))}
        </div>
      )}
      {query.data &&
        !query.data.some((item) =>
          item.name.toLowerCase().includes(search.toLowerCase()),
        ) && (
          <EmptyState
            title={
              search ? "No matching classes" : "Your next chapter starts here"
            }
          >
            {search
              ? "Try a different class name."
              : user.role === "admin"
                ? "Create a class and assign its TA and students to get started."
                : "Your administrator will assign you to a class. It will appear here."}
          </EmptyState>
        )}
      <div className="info-strip">
        <BookOpen size={20} />
        <div>
          <strong>A clear path from class to code</strong>
          <p>
            Administrators organize classes. Teaching assistants create labs.
            Students bring solutions.
          </p>
        </div>
      </div>
    </div>
  );
}

export function ClassForm() {
  const navigate = useNavigate();
  const cache = useQueryClient();
  const [search, setSearch] = useState("");
  const [ta, setTa] = useState<Person | null>(null);
  const [students, setStudents] = useState<Person[]>([]);
  const people = useQuery({
    queryKey: ["people", search],
    queryFn: () => api<Person[]>(`/people?q=${encodeURIComponent(search)}`),
  });
  const save = useMutation({
    mutationFn: (name: string) =>
      api<{ id: string }>("/classes", {
        name,
        taId: ta?.id,
        studentIds: students.map((student) => student.id),
      }),
    onSuccess: async (data) => {
      await cache.invalidateQueries({ queryKey: ["classes"] });
      navigate(`/classes/${data.id}`, { viewTransition: true });
    },
  });

  return (
    <div className="page-container narrow-page">
      <Breadcrumb parent="Classes" to="/classes" current="New class" />
      <PageHeader
        title="Create a class"
        description="Build a home for your next lab."
      />
      <Feedback error={save.error} />
      <form
        className="form-layout"
        onSubmit={(event) => {
          event.preventDefault();
          save.mutate(String(new FormData(event.currentTarget).get("name")));
        }}
      >
        <section className="panel form-section">
          <div className="section-label">Class details</div>
          <label>
            Class name
            <input
              name="name"
              placeholder="e.g. Introduction to C++ · Group A"
              required
            />
          </label>
        </section>
        <section className="panel form-section">
          <div className="section-label">Assign people</div>
          <p className="muted">
            Choose existing accounts. Search by email to find your TA and
            students.
          </p>
          <label className="search-field">
            <Search size={17} />
            <input
              value={search}
              aria-label="Search accounts"
              placeholder="Search accounts by email…"
              onChange={(event) => setSearch(event.target.value)}
            />
          </label>
          <Feedback error={people.error} retry={() => void people.refetch()} />
          <div className="picker-columns">
            <div>
              <h3>
                Teaching assistant <span className="required">*</span>
              </h3>
              <div className="account-list">
                {people.data
                  ?.filter((p) => p.role === "ta")
                  .map((person) => (
                    <label className="account-option" key={person.id}>
                      <input
                        type="radio"
                        name="ta"
                        checked={ta?.id === person.id}
                        onChange={() => setTa(person)}
                      />
                      <span>
                        {person.email}
                        <small>Teaching assistant</small>
                      </span>
                    </label>
                  ))}
              </div>
            </div>
            <div>
              <h3>Students</h3>
              <div className="account-list">
                {people.data
                  ?.filter((p) => p.role === "student")
                  .map((person) => (
                    <label className="account-option" key={person.id}>
                      <input
                        type="checkbox"
                        checked={students.some((p) => p.id === person.id)}
                        onChange={(event) =>
                          setStudents(
                            event.target.checked
                              ? [...students, person]
                              : students.filter((p) => p.id !== person.id),
                          )
                        }
                      />
                      <span>
                        {person.email}
                        <small>Student</small>
                      </span>
                    </label>
                  ))}
              </div>
            </div>
          </div>
          {people.data?.length === 0 && (
            <p className="muted">No accounts found. Try another email.</p>
          )}
          <p className="help-text">
            Showing up to 200 matching accounts. Account provisioning is handled
            by your administrator.
          </p>
        </section>
        <div className="selection-summary">
          <Check size={17} />
          <span>
            {ta ? `TA: ${ta.email}` : "Select a teaching assistant"} ·{" "}
            {students.length} students
          </span>
        </div>
        {students.length > 0 && (
          <div className="chips">
            {students.map((person) => (
              <button
                type="button"
                className="small-tag"
                key={person.id}
                onClick={() =>
                  setStudents(students.filter((p) => p.id !== person.id))
                }
              >
                {person.email}
                <X size={12} aria-label="Remove" />
              </button>
            ))}
          </div>
        )}
        <div className="form-footer">
          <Link className="btn" to="/classes">
            Cancel
          </Link>
          <button className="btn btn-primary" disabled={!ta || save.isPending}>
            {save.isPending ? "Creating…" : "Create class"}
          </button>
        </div>
      </form>
    </div>
  );
}

export function ClassPage({ user }: { user: User }) {
  const { id } = useParams();
  const [tab, setTab] = useState("labs");
  const query = useQuery({
    queryKey: ["class", id],
    queryFn: () => api<Classroom>(`/classes/${id}`),
  });
  const labs = useQuery({
    queryKey: ["labs"],
    queryFn: () => api<Lab[]>("/labs"),
    enabled: user.role !== "admin",
  });

  if (query.isPending) return <Loading />;

  return (
    <div className="page-container">
      <Breadcrumb
        parent="Classes"
        to="/classes"
        current={query.data?.name ?? "Class"}
      />
      <Feedback error={query.error} retry={() => void query.refetch()} />
      {query.data && (
        <>
          <PageHeader
            eyebrow="CLASS WORKSPACE"
            title={query.data.name}
            description="A shared space for your people and their next challenge."
            action={
              user.role === "ta" && (
                <Link
                  className="btn btn-primary"
                  to={`/classes/${id}/labs/new`}
                >
                  <Plus size={16} />
                  Create lab
                </Link>
              )
            }
          />
          <div className="tabbar" role="tablist" aria-label="Class sections">
            {(user.role === "admin"
              ? ["members"]
              : user.role === "ta"
                ? ["labs", "members"]
                : ["labs"]
            ).map((item) => (
              <button
                key={item}
                role="tab"
                aria-selected={
                  (user.role === "admin" ? "members" : tab) === item
                }
                className={
                  (user.role === "admin" ? "members" : tab) === item
                    ? "active"
                    : ""
                }
                onClick={() => setTab(item)}
              >
                {item === "labs" ? "Labs" : "Members"}
              </button>
            ))}
          </div>
          {tab === "members" || user.role === "admin" ? (
            <div className="panel table-wrap">
              <table>
                <thead>
                  <tr>
                    <th>Account</th>
                    <th>Role</th>
                  </tr>
                </thead>
                <tbody>
                  {query.data.members?.map((person) => (
                    <tr key={person.id}>
                      <td>
                        <span className="person-cell">
                          <span className="avatar">
                            {person.email[0].toUpperCase()}
                          </span>
                          {person.email}
                        </span>
                      </td>
                      <td>
                        <span className="small-tag">
                          {person.role === "ta"
                            ? "Teaching assistant"
                            : "Student"}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <div className="panel">
              <Feedback error={labs.error} retry={() => void labs.refetch()} />
              {labs.isPending ? (
                <Loading />
              ) : labs.data?.some((lab) => lab.classId === id) ? (
                labs.data
                  .filter((lab) => lab.classId === id)
                  .map((lab) => (
                    <Link
                      className="list-row"
                      to={`/labs/${lab.id}`}
                      key={lab.id}
                    >
                      <BookOpen size={18} />
                      <strong>{lab.title}</strong>
                      <span className="text-link">Open lab</span>
                    </Link>
                  ))
              ) : (
                <EmptyState title="Room for your next challenge">
                  {user.role === "ta"
                    ? "Create the first lab for this class."
                    : "Your TA hasn't created a lab for this class yet."}
                </EmptyState>
              )}
            </div>
          )}
          {user.role === "admin" && (
            <div className="info-strip">
              <BookOpen size={20} />
              <div>
                <strong>Ready for your teaching assistant</strong>
                <p>
                  The assigned TA can sign in to create labs for this class.
                </p>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
}

export function PeoplePage() {
  const [search, setSearch] = useState("");
  const [role, setRole] = useState("");
  const query = useQuery({
    queryKey: ["people", search],
    queryFn: () => api<Person[]>(`/people?q=${encodeURIComponent(search)}`),
  });

  return (
    <div className="page-container">
      <PageHeader
        eyebrow="ADMINISTRATION"
        title="People"
        description="The people behind every solution."
      />
      <div className="section-toolbar">
        <label className="search-field">
          <Search size={17} />
          <input
            aria-label="Search people"
            placeholder="Search by email…"
            value={search}
            onChange={(event) => setSearch(event.target.value)}
          />
        </label>
        <select
          aria-label="Filter by role"
          value={role}
          onChange={(event) => setRole(event.target.value)}
        >
          <option value="">All roles</option>
          <option value="admin">Admin</option>
          <option value="ta">Teaching assistant</option>
          <option value="student">Student</option>
        </select>
      </div>
      <Feedback error={query.error} retry={() => void query.refetch()} />
      {query.isPending ? (
        <Loading />
      ) : (
        <div className="panel table-wrap">
          <table>
            <thead>
              <tr>
                <th>Email</th>
                <th>Role</th>
                <th>Account ID</th>
              </tr>
            </thead>
            <tbody>
              {query.data
                ?.filter((p) => !role || p.role === role)
                .map((person) => (
                  <tr key={person.id}>
                    <td>
                      <span className="person-cell">
                        <span className="avatar">
                          {person.email[0].toUpperCase()}
                        </span>
                        {person.email}
                      </span>
                    </td>
                    <td>
                      <span className="small-tag">{person.role}</span>
                    </td>
                    <td className="mono muted">{person.id}</td>
                  </tr>
                ))}
            </tbody>
          </table>
          {query.data?.filter((p) => !role || p.role === role).length === 0 && (
            <EmptyState title="No matching people">
              Try another email or role.
            </EmptyState>
          )}
        </div>
      )}
      <p className="help-text">
        Up to 200 matching accounts. New accounts are provisioned with the
        administrator CLI.
      </p>
    </div>
  );
}
