# Platform gateway session boundary

The Platform-private session capability is at `platform/internal/session`. Launcher and Platform services share one dependency-free security boundary for principals, cookies and tenant generation. The launcher root exposes composition contracts; there is no public launcher/session package.

Go restricts launcher/internal imports to launcher descendants. Platform services cannot import that path or import launcher root without cycles through bootstrap. The shared session capability therefore lives in the Platform internal directory. Launcher-only gateway wiring remains in launcher/internal/src.

Never derive principals or trusted tenant checks from HTTP input. Authentication and tenant-generation semantics are unchanged.
