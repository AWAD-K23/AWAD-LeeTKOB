import {
  Link as RouterLink,
  NavLink as RouterNavLink,
  type LinkProps,
  type NavLinkProps,
} from "react-router-dom";

// Router-managed snapshots animate both the departing and arriving page.
// Browsers without View Transitions still receive the CSS entrance animation.
export function Link(props: LinkProps) {
  return <RouterLink {...props} viewTransition />;
}

export function NavLink(props: NavLinkProps) {
  return <RouterNavLink {...props} viewTransition />;
}
