# LeeTKOB frontend design plan

## 1. Direction and goals

Build a modern, focused coding platform inspired by LeetCode's compact navigation,
searchable lists, clear status indicators, and split problem-and-code workspace.
Keep LeeTKOB's own identity and prioritize its class, judging, and TA review workflows.

Reference: [LeetCode problem browsing](https://leetcode.com/problemset/).

The current frontend places nearly every workflow in `App.tsx`. Replace this with
dedicated pages, reusable components, clear navigation, and reliable request states.
Users must always understand their current context and next available action.

### Product principles

- Use compact, readable layouts with clear hierarchy and generous section spacing.
- Keep a single obvious primary action within each task.
- Prefer tables for browsing and panels for focused work.
- Make role boundaries visible and enforce permissions on the server.
- Display real data; do not ship fabricated statistics or inactive actions.
- Separate judge results, checker findings, TA decisions, and published marks.
- Use motion to explain changes and acknowledge actions without distracting from code.

## 2. Visual system

Dark mode is the default. Support light mode through semantic theme tokens and an
account preference. Persist the theme locally and respect system preference until
the user chooses a theme explicitly.

| Token or element | Design |
| --- | --- |
| Page background | `#18181B` |
| Navigation and panels | `#222225` |
| Elevated menus | `#2C2C30` |
| Borders | `#36363B` |
| Main text | `#F4F4F5` |
| Secondary text | `#A1A1AA` |
| Brand accent | Amber `#FFA116` |
| Success | Green, with an icon and text label |
| Pending / unsure | Amber, with an icon and text label |
| Failure / violation | Red, with an icon and text label |
| Body typography | System sans-serif, 14px with comfortable line height |
| Metadata | 12px; maintain readable contrast |
| Page titles | 20–28px, medium or semibold |
| Source and test output | Monospace, 13–14px |
| Spacing | 4, 8, 12, 16, 24, and 32px |
| Corners | 6–10px; restrained shadows |
| Icons | Lucide with consistent size and stroke |

Reserve amber for identity, selected states, and deliberate emphasis. Most controls
use neutral surfaces. Verify contrast in both themes, including status badges and
disabled controls. Never communicate status through color alone.

Replace broad global styles for `main`, `header`, and `form` with layout-specific
classes. Management pages use a centered content area; coding and review workspaces
use the available viewport width.

## 3. Navigation and application shell

A persistent top bar contains the logo, role-specific navigation, relevant class
context, and account menu. Detail pages include breadcrumbs.

| Role | Main destinations | Default destination |
| --- | --- | --- |
| Admin | Classes, People | Classes |
| TA | Classes, Labs, Review queue, Appeals | Labs |
| Student | Labs, My submissions, Grades | Labs |

Show review, appeals, and grade destinations only when their supporting APIs exist.
The account menu shows email, role, theme preference, and sign out. Signing out
clears account-specific cached data and returns to login.

Use a compact mobile navigation menu. Keep the current page and class identifiable.
Do not add a role switch that bypasses authentication or account permissions.

## 4. Pages and routes

| Route | Page and components | Main behavior |
| --- | --- | --- |
| `/login` | Centered login panel, logo, email, password visibility toggle, inline errors | Authenticate and redirect to the role's default page |
| `/classes` | Page header, search, class list, empty state | Browse classes; admin creates a class |
| `/classes/new` | Class name, TA combobox, student multi-select, selection summary | Create and open the class |
| `/classes/:id` | Class header, Labs and Members tabs | View assigned people and available labs |
| `/people` | Searchable account table and role filters | Admin finds accounts and assignments |
| `/labs` | Search, class and status filters, lab table | Open a lab; assigned TA creates a lab |
| `/classes/:id/labs/new` | Section navigation and lab authoring form | Define statement, tests, and constraints |
| `/labs/:id` | Split description and submission workspace | Read the task and submit C++ |
| `/submissions` | Submission history table and filters | Open an earlier submission |
| `/submissions/:id` | Judge summary, saved source, result details | Follow judging and inspect a submission |
| `/review` | Filterable TA review queue | Find required reviews |
| `/review/:id` | Source viewer, evidence panel, TA decision form | Review checker findings |
| `/labs/:id/publish` | Completion summary and publication confirmation | Publish after all required reviews |
| `/grades` | Published results grouped by class and lab | Inspect released marks and deductions |
| `/appeals` | Appeal list and state filters | Track or handle appeals |
| `/appeals/:id` | Deduction context, evidence, appeal, reasoned decision | Accept or reject an appeal |

Include a page-not-found screen, an access-denied screen, and a session-expired flow.
Each should provide a useful path back to the user's workspace.

## 5. Admin and class management

Replace UUID entry fields with account selectors displaying email and role.

- TA picker: searchable single selection restricted to TA accounts.
- Student picker: searchable multi-selection with removable chips.
- Summary: selected TA and student count before creation.
- Validation: field-level errors, with form data preserved after failures.
- Success: navigate to the created class and show a concise confirmation.

The class detail page includes its name, assigned TA, Labs tab, and Members tab.
Admin actions and TA actions remain distinct. An empty class explains the next step:
the assigned TA signs in and creates a lab.

Account search, member details, and assignment changes require appropriate backend
endpoints. Account creation remains the administrator CLI until an authorized
account-provisioning API is implemented.

## 6. TA lab authoring

Use a spacious authoring page with section navigation and four sections:

1. **Details:** class, title, statement, and statement preview.
2. **Test cases:** numbered input/output pairs and add/remove controls.
3. **Constraints:** requirement text and TA-selected deduction.
4. **Review:** complete summary with links back to each section.

Keep the primary Create lab action visible in a footer without covering form fields.
Preserve inputs after failed requests. Warn before leaving with unsaved changes.
On validation failure, focus the first invalid field and provide a linked summary.

Keep hidden judge tests separate from student-visible examples. The authoring UI may
need API additions for that distinction; do not expose hidden tests through the
student lab endpoint. Persistent drafts and editing require explicit backend support.

## 7. Student coding workspace

Desktop layout:

```text
┌─────────────────────────────────────────────────────────────┐
│ Logo / navigation       Class > Lab           Account menu   │
├────────────────────────────┬────────────────────────────────┤
│ Description                │ C++ source          Submit      │
│                            │                                │
│ Statement                  │ Source entry with line numbers │
│ Examples                   │                                │
│ Assignment constraints     │                                │
│                            ├────────────────────────────────┤
│                            │ Submission result              │
│                            │ Queued / judging / completed   │
└────────────────────────────┴────────────────────────────────┘
```

Use a resizable split with sensible minimum panel widths and a keyboard-operable
resize handle. On smaller screens, use Description, Code, and Result tabs.

Include a fixed C++ language label, copy source, reset with confirmation, and local
draft recovery scoped to account and lab. Explain draft recovery and clear sensitive
account state appropriately on sign out. Preserve source across request failures.

This is a submission workspace rather than a full online IDE, consistent with the
project brief. A Run button requires a separate test-run API. Do not add terminal,
debugger, dependency installation, or unsupported language controls.

After submitting, display progress, a link to the saved submission, and an explicit
retry action if status updates fail. Do not invent a queue position or percentage.

## 8. Results and submission history

Keep the following concepts distinct in components and copy:

| Concept | States or content |
| --- | --- |
| Job processing | Queued, processing, completed, failed |
| Judge result | Accepted, Wrong Answer, compilation failure, other supported verdicts |
| Constraint checking | Pending, pass, violation, unsure |
| TA review | Required, reviewed |
| Publication | Unpublished, published |

Accepted means the program passed judging; it does not imply a final mark.
Submission history includes lab, submission time, processing state, and verdict.
Show runtime, memory, source, and test details only when supported by the API.
Student views must exclude protected evidence and unpublished deductions.

## 9. TA review and publication

The review queue is a compact table with submission, lab, judge result, review
reason, and review state. Filter by class, lab, reason, and reviewed/unreviewed state.

The review workspace includes:

- Read-only source with stable line numbers.
- Constraint list showing checker source, verdict, and configured deduction.
- Evidence selection that scrolls to and highlights exact source lines.
- TA decision controls stored separately from checker verdicts.
- Previous/next navigation that warns about unsaved decisions.

Every violation, unsure verdict, and required sample of pass verdicts must receive
the review required by the project rules. Publication has a separate summary page
showing unresolved reviews. The server must block publication while required reviews
remain; disabling the frontend button alone is insufficient.

## 10. Grades and appeals

Grades display published results only, with clear deductions and associated reasons.
An appeal opens from a published deduction and retains its submission and constraint
context. Students can track appeal state. TAs review the original evidence and
provide a reason when accepting or rejecting an appeal.

Implement these pages after the supporting protected APIs and product rules exist.
Do not imply that marks have been released or an appeal saved without confirmation.

## 11. Reusable component inventory

| Group | Components |
| --- | --- |
| Foundations | Button, IconButton, Input, Textarea, Select, Checkbox, Label, Badge, Tooltip, Separator |
| Interaction | Dialog, AlertDialog, DropdownMenu, Tabs, Combobox, MultiSelect, ResizablePanels |
| Layout | AppShell, TopNav, MobileNav, AccountMenu, Breadcrumbs, PageHeader, SectionNav |
| Collections | FilterBar, DataTable, Pagination, EmptyState, SelectionSummary |
| Forms | FormField, FieldError, ErrorSummary, UnsavedChangesDialog |
| Feedback | ErrorBanner, LoadingSkeleton, ConnectionStatus, Toast |
| Classes | ClassPicker, AccountPicker, ClassList, MemberTable |
| Labs | LabTable, StatementPreview, TestCaseEditor, ConstraintEditor |
| Submissions | SourceInput, SourceViewer, VerdictBadge, JudgeResultPanel, SubmissionTable |
| Review | EvidencePanel, ReviewDecisionForm, ReviewQueueTable, PublishChecklist |
| Appeals | GradeSummary, DeductionList, AppealForm, AppealDecisionPanel |

Every reusable component defines default, hover, focus, disabled, loading, empty,
and error states where applicable. Continue using Radix-based primitives and Lucide.

## 12. Animation, transitions, and effects

### Motion principles

Motion should explain state changes, preserve orientation, and confirm interactions.
Prefer opacity and transform animations. Avoid animating layout on large tables,
source viewers, or every polling update. Keep content readable while animations run.

Use CSS transitions and keyframes for simple effects. Add a motion library only if
complex presence or layout transitions justify it. Centralize durations and easing
as design tokens instead of scattering arbitrary values through components.

| Motion token | Duration | Use |
| --- | --- | --- |
| Instant | 80–100ms | Press feedback and tiny state changes |
| Fast | 120–160ms | Hover, focus, selected state, tooltips |
| Standard | 180–220ms | Menus, tabs, toasts, panel content |
| Deliberate | 240–300ms | Dialogs, drawers, limited page entrances |
| Enter easing | `cubic-bezier(0.16, 1, 0.3, 1)` | Elements entering view |
| General easing | `cubic-bezier(0.2, 0, 0, 1)` | State and surface transitions |

### Component effects

| Component / interaction | Planned effect |
| --- | --- |
| Buttons | Subtle background and border transition; optional small press compression without layout movement |
| Inputs | Smooth border and focus-ring transition; maintain a clear keyboard outline |
| Navigation | Short active-indicator transition; current location remains clear without motion |
| Table rows | Neutral surface highlight on hover and keyboard focus; no row scaling |
| Class cards | Slight border emphasis on hover; avoid oversized shadows and excessive lift |
| Menus and tooltips | Fade with a 2–4px directional offset; respect viewport collision handling |
| Dialogs | Backdrop fade and small opacity/scale entrance; keyboard focus moves correctly |
| Mobile drawers | Short slide and backdrop fade; trap focus and restore it on close |
| Tabs | Brief content fade; avoid changing scroll position unexpectedly |
| Added tests / constraints | Short fade entrance that makes the new item easy to locate |
| Removed form sections | Brief exit only if it does not disrupt focus or lose unsaved data silently |
| Toasts | Small slide/fade entrance; readable timing and manual dismissal |
| Skeletons | Low-contrast shimmer during initial loading; match final content dimensions |
| Judging | Small spinner while processing; use honest state labels instead of artificial progress |
| Verdict changes | One brief transition when the verdict changes; never flash on every poll |
| Evidence selection | Scroll to the selected line and briefly emphasize its background |
| Resizable panels | Follow the pointer directly; do not ease the dragging interaction |
| Theme changes | Brief color transition on surfaces and text; avoid a bright initial flash |

### Page-specific motion

- Login: a single restrained panel entrance; keep error messages stable and readable.
- Class creation: navigate to the created class and show a success toast.
- Lab authoring: gently introduce new test cases; maintain focus on the added section.
- Coding workspace: preserve editor selection, source, and panel sizes during updates.
- Review: highlight selected evidence without shifting the entire layout.
- Publication: show a stable saved confirmation after server success; no celebratory
  effect that obscures the irreversible nature of publishing marks.

### Accessibility and performance

- Honor `prefers-reduced-motion: reduce`: remove spatial movement, shimmer, smooth
  scrolling, and repeated animation. Preserve static progress and status labels.
- Never depend on animation to convey a result or expose an action.
- Do not use flashing effects, parallax, autoplay decoration, or continuous glow.
- Use polite live announcements for meaningful asynchronous results. Avoid announcing
  every poll or decorative transition.
- Keep focus stable after updates; restore it after dialogs and menu dismissal.
- Reserve space for asynchronous content to avoid layout shifts.
- Avoid staggered entrances for long lists; keep navigation immediately usable.

## 13. Loading, errors, and connection behavior

- Resolve the session before requesting protected page data.
- Treat signed-out responses as a normal login state.
- Show initial loading skeletons and retryable connection errors.
- Preserve already loaded data when background refresh fails.
- Disable only the action currently running and guard against duplicate submissions.
- Preserve form input after request failures.
- Parse API problem responses safely; never expose an HTML parsing exception.
- Refresh affected queries after successful mutations.
- Distinguish an expired session from a server outage.
- Stop status polling when a submission reaches a terminal state.

## 14. Frontend structure

Use React Router for navigation and TanStack Query for server data and mutation
states. Keep authentication, API access, and presentation separate.

```text
frontend/src/
  app/                 # Router, providers, session, application shell
  components/
    ui/                # Reusable accessible primitives
    layout/            # Navigation, page headers, breadcrumbs
    feedback/          # Loading, error, connection, and toast patterns
  features/
    auth/
    classes/
    people/
    labs/
    submissions/
    review/
    grades/
    appeals/
  lib/                 # API client, query keys, shared utilities
  styles/              # Theme tokens, motion tokens, global foundation
```

When implementing code, separate logical parts with blank lines, comment important
behavior, and wrap long lines for readability.

## 15. API readiness and implementation boundaries

Current backend routes support sessions, class creation/listing, lab creation/listing,
submission creation, and protected submission status. Verify the contracts before
wiring each page.

Additional contracts are needed for account search, class details/members, submission
history and saved source, supported editing, TA review decisions, publication,
published grades, and appeals. Define role and ownership rules for every new endpoint.

Build the existing admin-to-TA-to-student journey first. Introduce future destinations
when the relevant API and permissions are ready. Never make real model calls in unit
tests, expose sandbox secrets, or publish marks with required reviews pending.

## 16. Delivery sequence

1. **Foundation:** semantic themes, motion tokens, primitives, routing, application
   shell, session handling, API error handling, and role-aware navigation.
2. **Existing workflows:** login, class list and creation, clear handoff to TA,
   TA lab authoring, student workspace, and submission status.
3. **Supporting APIs:** account search, class members/details, submission history
   and source, and explicitly supported editing actions.
4. **Review and release:** evidence review, TA decisions, publication gate, grades,
   and appeals with their protected backend contracts.
5. **Polish and verification:** responsive layouts, motion tuning, keyboard flows,
   reduced motion, errors, loading, and end-to-end journeys.

## 17. Acceptance criteria

- An admin can create a class and immediately see it and its assigned TA.
- A TA can find an assigned class and create a valid lab without raw UUID entry.
- An enrolled student can find a lab, preserve source, submit, and follow judging.
- Every network operation communicates loading, success, and actionable failure.
- Every visible action works with a real endpoint and server-side permissions.
- Future review and publication flows keep evidence, decisions, and marks separate.
- Desktop layouts work at 1280px and 1440px; tablet and mobile layouts remain usable
  at 768px and 375px without page-wide horizontal overflow.
- All primary journeys work with keyboard navigation and reduced motion enabled.
- Animation does not reset source, move focus unexpectedly, or delay interaction.
- Frontend lint, formatting, type-checking, and build pass using the documented check.
- Relevant browser journeys cover role navigation, form failure recovery, submission
  status, and permission boundaries. Complete backend checks when API changes occur.

## 18. Implementation status

The local implementation now includes the role-aware application shell, routed
login/class/people/lab/submission pages, searchable account selectors, class member
views, lab authoring, a responsive split coding workspace, account-scoped local
source drafts, saved-source views, polling, and student submission history.

The visual system includes light/dark themes, hover/focus/press transitions, page
and form-section entrances, progress indicators, and reduced-motion overrides.
Unsaved lab authoring is protected during in-app navigation and browser unload.
Browser API calls use a dedicated `/api` proxy prefix so page refreshes work.

Protected browsing endpoints were added for account search, class details, and
student history. Existing submission status includes saved source; lab lists include
public assignment requirement descriptions without exposing hidden tests or deductions.

Review queues, evidence decisions, publication, grades, and appeals remain future
phases requiring their own backend contracts and storage. They are intentionally
absent from the current navigation. Persistent lab drafts, lab editing, and a
separate test-run action also remain future work.
