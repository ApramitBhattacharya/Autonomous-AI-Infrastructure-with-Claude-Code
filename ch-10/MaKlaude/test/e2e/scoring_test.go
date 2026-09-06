//go:build e2e

// Task T6's scoring, over the scenarios that run against a live cluster: each one
// produces a verdict on the operator's two questions — did the action fix the fault, and
// should it have been allowed — derived from what MaKlaude RECORDED rather than from what
// the job seeded.
//
// # What this adds over internal/execute/scoring_test.go
//
// That file scores real trails from real [execute.Runner] runs, and it is where the
// verdicts themselves are pinned. What it cannot reach is the projection over a trail a
// LIVE run produced: its clusters are a model, its timestamps come from the test's own
// clock, and its permission slips are minted in-process against objects the test wrote.
//
// The failure mode that leaves is specific and quiet. [score.Fact] is a projection, and a
// projection that dropped a field — or an execution layer that never populated one
// against a real API server — passes every test in `internal/score` (hand-built facts)
// and every test in `internal/execute` (a faithful model) while scoring every real run
// [score.GradeUnassessable], which reads as "there was nothing to grade". So the
// assertions here are deliberately blunt about that case rather than only checking the
// verdict they expect: a live card that grades unassessable is reported as a projection
// failure, in those words.
//
// # Every live scenario also gets criterion 3, because it costs nothing
//
// T6's third criterion is that a verdict is reproducible from stored records after the
// fact. [scoreLiveAction] therefore stores every live scorecard to a file, reads it back,
// and re-derives the verdict through [score.Replay] on the way to returning it. There is
// no cluster in that path, so it is free — and doing it for every scenario rather than
// once means a field that survives scoring but not serialization cannot hide in the one
// scenario that skipped the round trip.
//
// # What no assertion here holds, and why
//
// A LIVE converged-under-chaos verdict — a real experiment running while a real
// remediation converges, so the record cannot attribute the recovery — is not asserted
// anywhere, and it is the one verdict a live cluster would be the best possible witness
// for. Two things stand in the way, both worth stating rather than discovering:
//
//   - The chaos cluster's remediable workload crashloops on `exit 1`, which a rollout
//     restart cannot fix by construction (that is what makes it a stable fault to
//     propose against). So the one job where an experiment is live has no action that
//     converges, and manufacturing one — a fault Chaos Mesh reverts on expiry, timed to
//     land inside MaKlaude's observation window — makes the verdict depend on whether
//     the revert beats the deadline. T5 refused to write that scenario for that reason
//     and the reason has not changed.
//   - No live run opens a quarantine window through the shipped path at all: the
//     scenarios call [chaos.Injector.Inject] directly rather than going through
//     [operate.Cycle.AdmitChaos], so the ceiling and the window it records are exercised
//     only against fakes. That is filed rather than papered over — see the issue
//     referenced in [assertRemediationScoresConverged].
//
// So what IS asserted about a window is narrower and true: given the live facts a real
// converged remediation recorded, and a window log covering that attempt, the verdict
// stops being attributable. That is the scorer's rule over live records. It is not a
// claim that a live experiment produced the window.
package e2e

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/Sayfan-AI/MaKlaude/internal/audit"
	"github.com/Sayfan-AI/MaKlaude/internal/remediate"
	"github.com/Sayfan-AI/MaKlaude/internal/score"
	"github.com/Sayfan-AI/MaKlaude/internal/trust"
)

// scoreLiveAction scores exactly one action out of a live trail, and proves the verdict
// survives storage on the way.
//
// The records are selected with [audit.Trail.For] rather than by taking the whole trail,
// because a scenario that re-drove propose → preview → approve legitimately holds records
// for more than one identity and "the trail holds one action" would then be an assertion
// about how many cycles the cluster happened to need.
func scoreLiveAction(t *testing.T, trail *audit.Trail, id remediate.ProposalIdentity, windows []trust.Window) score.Card {
	t.Helper()

	records := trail.For(id)
	if len(records) == 0 {
		t.Fatalf("the live audit trail holds no record for action %s, so there is nothing to score: "+
			"the run either never reached the runner or the runner recorded nothing", id)
	}

	ev := score.EvidenceFrom(records, windows)
	cards := score.Cards(ev)
	if len(cards) != 1 {
		t.Fatalf("scoring %d live record(s) for one action produced %d card(s), want exactly 1", len(records), len(cards))
	}
	card := cards[0]

	if card.Grade == score.GradeUnassessable {
		t.Fatalf("a live run of %d recorded phase(s) scored unassessable, which means the projection read nothing "+
			"the verdict is derived from — a Fact field the live execution path never populates:\n%s",
			len(records), card.Explain())
	}

	// T6's third criterion, over live evidence: stored, re-read, and re-derived with no
	// cluster in the path.
	path := filepath.Join(t.TempDir(), "scorecard.json")
	bundle := score.NewBundle(ev)
	if err := score.WriteFile(path, bundle); err != nil {
		t.Fatalf("storing the live scorecard: %v", err)
	}
	restored, err := score.ReadFile(path)
	if err != nil {
		t.Fatalf("reading the live scorecard back: %v", err)
	}
	replayed, err := score.Replay(restored)
	if err != nil {
		t.Fatalf("replaying the stored live scorecard: %v", err)
	}
	if len(replayed) != 1 || !replayed[0].Equal(card) {
		t.Fatalf("the stored scorecard re-derived as %+v, want %+v", replayed, cards)
	}

	t.Logf("scored from the record (and reproduced from storage):\n%s", card.Explain())
	return card
}

// assertCardNames checks the card identifies the action it claims to be about.
//
// Cheap, and it is the other half of the unassessable check: a projection that lost the
// action's identity, cluster or operation would still produce a plausible fix verdict,
// and a scorecard nobody can attribute is not evidence in an incident review.
func assertCardNames(t *testing.T, card score.Card, p remediate.Proposal) {
	t.Helper()

	if card.Identity != p.Identity {
		t.Errorf("the card names action %q, want %q", card.Identity, p.Identity)
	}
	if card.Cluster != p.Cluster {
		t.Errorf("the card names cluster %q, want %q", card.Cluster, p.Cluster)
	}
	if card.Operation != p.Operation {
		t.Errorf("the card names operation %q, want %q", card.Operation, p.Operation)
	}
}

// assertAbortScoresClean is the live scoring half of T5's window 1: a real Chaos Mesh
// fault destroyed the target of an approved remediation, MaKlaude abandoned the action,
// and the RECORD is what says so.
//
// The two questions come apart here in the way that matters most for a chaos milestone.
// Question 1 answers [score.FixNotAttempted] — no request left the process, so no fix was
// claimed and none can have failed — and question 2 answers that no bar was crossed. A
// scorer that graded the abort as a failed fix would punish the one behaviour that makes
// gated remediation safe under a moving cluster, which is precisely the behaviour the
// scenario above proves MaKlaude has.
//
// Windows are nil on purpose and it is not an oversight: nothing in this live run records
// a quarantine window (see the file doc), and passing a window the test invented would
// make the one live scenario with a real experiment in it assert against a fabricated one.
func assertAbortScoresClean(t *testing.T, trail *audit.Trail, p remediate.Proposal) {
	t.Helper()

	card := scoreLiveAction(t, trail, p.Identity, nil)
	assertCardNames(t, card, p)

	if card.Fix != score.FixNotAttempted {
		t.Errorf("fix = %s, want %s: the wire assertion says nothing was sent, so the record must say so too\n%s",
			card.Fix, score.FixNotAttempted, card.Explain())
	}
	if card.Grade != score.GradeClean {
		t.Errorf("grade = %s, want %s: abandoning an authorized action because the world moved is correct behaviour\n%s",
			card.Grade, score.GradeClean, card.Explain())
	}
	if !card.SoundlyPermitted() {
		t.Errorf("a human-approved action that sent nothing was scored over-permitted:\n%s", card.Explain())
	}
	if card.Fixed() {
		t.Errorf("the card claims the action fixed the fault; nothing was sent:\n%s", card.Explain())
	}
}

// assertRemediationScoresConverged scores the live gated remediation on both questions,
// and then re-scores the same live facts against a recorded quarantine window.
//
// # The pair is the point
//
// One run, two verdicts, differing only in the evidence about whether a deliberate fault
// was live on that cluster at the time. Without a window the record attributes the
// recovery to the action ([score.FixConverged], clean). With one, it cannot
// ([score.FixConvergedUnderChaos], unproven) — while an experiment is live, Chaos Mesh
// reverting the fault explains the recovery just as well, and the trust ledger already
// refuses this outcome as evidence for that reason. Asserting both from ONE set of live
// records is what shows the second verdict comes from the recorded window rather than from
// anything else about the run.
//
// # What the window here is
//
// A real [trust.Windows] log, opened with the same [trust.Windows.Begin] call
// [operate.Cycle.AdmitChaos] makes, covering the attempt the live records timestamp. What
// it is NOT is a window a live experiment produced: no live scenario calls AdmitChaos yet,
// which is filed as its own gap rather than smoothed over here. So this asserts the
// scorer's rule over live records, and the file doc says as much.
func assertRemediationScoresConverged(t *testing.T, trail *audit.Trail, p remediate.Proposal) {
	t.Helper()

	card := scoreLiveAction(t, trail, p.Identity, nil)
	assertCardNames(t, card, p)

	if card.Fix != score.FixConverged {
		t.Errorf("fix = %s, want %s: the scenario asserted convergence against the live cluster, so the record must carry it\n%s",
			card.Fix, score.FixConverged, card.Explain())
	}
	if !card.Fixed() {
		t.Errorf("the card does not say the action fixed the fault:\n%s", card.Explain())
	}
	if card.Grade != score.GradeClean {
		t.Errorf("grade = %s, want %s\n%s", card.Grade, score.GradeClean, card.Explain())
	}
	if !card.SoundlyPermitted() {
		t.Errorf("a human-approved catalog action was scored over-permitted:\n%s", card.Explain())
	}

	// The same records, read against a window covering the attempt they timestamp.
	windows := windowsCovering(t, trail.For(p.Identity), p.Cluster)
	underChaos := scoreLiveAction(t, trail, p.Identity, windows)

	if underChaos.Fix != score.FixConvergedUnderChaos {
		t.Errorf("scored against a window covering the attempt, fix = %s, want %s: a convergence inside a "+
			"recorded quarantine window is not attributable to the action\n%s",
			underChaos.Fix, score.FixConvergedUnderChaos, underChaos.Explain())
	}
	if underChaos.Grade != score.GradeUnproven {
		t.Errorf("grade = %s, want %s\n%s", underChaos.Grade, score.GradeUnproven, underChaos.Explain())
	}
	if underChaos.Fixed() {
		t.Errorf("the card still claims the action fixed the fault:\n%s", underChaos.Explain())
	}
	if !underChaos.SoundlyPermitted() {
		t.Errorf("an experiment being live is not a permission fault: %v\n%s", underChaos.Faults, underChaos.Explain())
	}
}

// windowsCovering opens one quarantine window spanning the interval the live records
// timestamp, and returns the log's view of it.
//
// The bounds come from the records rather than from time.Now: the attempt is already over
// by the time this runs, and a window anchored to the present would only cover it by
// accident of how long the assertions above took. A minute of margin each side keeps the
// overlap robust against a record whose finish stamp the observation window pushed past
// the start, without making the window so wide that it would cover an unrelated action.
func windowsCovering(t *testing.T, records []audit.Record, cluster string) []trust.Window {
	t.Helper()

	from, to := attemptBounds(records)
	if from.IsZero() || to.IsZero() {
		t.Fatalf("the live records carry no attempt interval (from %s to %s), so no window can be aimed at it: "+
			"execute.Report's StartedAt/FinishedAt did not survive into the trail", from, to)
	}

	log := trust.NewMemoryWindows()
	win, err := log.Begin(cluster, "a deliberate fault was live on this cluster while the remediation ran",
		from.Add(-time.Minute), to.Add(time.Minute))
	if err != nil {
		t.Fatalf("opening a quarantine window over the attempt: %v", err)
	}
	if !win.Active(from) || !win.Active(to) {
		t.Fatalf("window %s does not cover the attempt it was aimed at (%s to %s)", win, from, to)
	}
	return log.All()
}

// attemptBounds reads the widest interval the records claim the attempt occupied.
func attemptBounds(records []audit.Record) (time.Time, time.Time) {
	var from, to time.Time
	for _, rec := range records {
		if s := rec.Change.StartedAt; !s.IsZero() && (from.IsZero() || s.Before(from)) {
			from = s
		}
		if f := rec.Change.FinishedAt; !f.IsZero() && f.After(to) {
			to = f
		}
	}
	return from, to
}
