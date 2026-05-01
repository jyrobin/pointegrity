# Letters

Monthly letters from the Pointegrity workshop, sent to subscribers
captured by `cmd/letterbox`.

## What's in this directory

- **`_template.md`** — checked into git. Skeleton with section
  guidance + placeholders (`{{UNSUBSCRIBE_URL}}` etc).
- **`YYYY-MM.md`** — actual filled-in letters. **NOT checked into
  git** (see `.gitignore`). Where they live long-term is your call;
  see "Storage options" below.

## Workflow

```sh
cp letters/_template.md letters/2026-05.md
# fill in. ~30-60 min.

cd ~/ws/jyws/pointegrity
poi-ops letterbox subscribers --csv > /tmp/subs.csv
awk -F, 'NR>1 {print $1}' /tmp/subs.csv     # email addresses
awk -F, 'NR>1 {print $4}' /tmp/subs.csv     # unsubscribe tokens

# Compose in your mail client:
#   From:    hello@pointegrity.com
#   To:      yourself
#   BCC:     paste the email list  (NEVER To: — exposes everyone)
#   Subject: line from frontmatter
#   Body:    everything below the frontmatter, with each
#            recipient's {{UNSUBSCRIBE_URL}} substituted to:
#            https://letterbox.pointegrity.com/unsubscribe?t=<token>

# After sending, fill in `sent_at:` in the frontmatter and move/copy
# the file to wherever you keep the archive.
```

The per-recipient `{{UNSUBSCRIBE_URL}}` substitution is the manual
tax. It's why `poi-ops letterbox send <markdown>` will pay for
itself by the second or third letter.

## Storage options for filled letters

These markdown files contain subscriber email count, occasional
ungrooomed wording, sometimes numbers (revenue, churn) you may not
want public. Three reasonable places:

1. **Local-only** (current default — gitignored).
   Pro: simplest. Con: lost on disk failure unless backed up.

2. **Private GitHub repo** (e.g. `pointegrity-private`).
   Pro: version control + GitHub backup. Free on personal plan.
   Con: another repo to maintain.

3. **Pouch** — drop the letter into your own data relay with a
   pinned/durable stream. Pro: dogfoods your own product; tags +
   timeline come for free. Con: pouch is optimized for transit,
   not archive — the read pattern for "show me the May 2024
   letter" isn't pouch's strongest move.

The format (markdown + frontmatter) is portable across all three.
Migrating later is a `cp` away.

## Frequency

Monthly. Less is fine; more is not. Subscribers leave from
"too frequent" much more than from "too infrequent."

## Subject line conventions

Always begin with `Letter from the workshop — `. Makes filtering /
threading trivial in subscribers' inboxes. Date suffix is the
month + year, not the send date — what the letter is *about* matters
more than when it went out.
