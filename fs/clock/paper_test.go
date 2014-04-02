package clock_test

import (
	"testing"

	"bazil.org/bazil/fs/clock"
)

// Tests based on the Vector Pair paper at
// http://publications.csail.mit.edu/tmp/MIT-CSAIL-TR-2005-014.pdf

func TestFigure2A(t *testing.T) {
	var b clock.Clock

	a := clock.Create(10, 1)
	if g, e := clock.Sync(a, &b), clock.Copy; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}
	b.ResolveTheirs(a)
	a.Update(10, 3)
	if g, e := clock.Sync(&b, a), clock.Nothing; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}
	a.ResolveOurs(&b)

	if g, e := a.String(), `{sync{10:3} mod{10:3} create{10:1}}`; g != e {
		t.Errorf("bad state A: %v != %v", g, e)
	}
	if g, e := b.String(), `{sync{10:1} mod{10:1} create{10:1}}`; g != e {
		t.Errorf("bad state B: %v != %v", g, e)
	}

	a.ValidateFile()
	b.ValidateFile()
}

func TestFigure2B(t *testing.T) {
	var b clock.Clock

	a := clock.Create(10, 1)
	if g, e := clock.Sync(a, &b), clock.Copy; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}
	b.ResolveTheirs(a)
	b.Update(11, 3)
	if g, e := clock.Sync(&b, a), clock.Copy; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}
	a.ResolveTheirs(&b)

	if g, e := a.String(), `{sync{10:1 11:3} mod{10:1 11:3} create{10:1}}`; g != e {
		t.Errorf("bad state A: %v != %v", g, e)
	}
	if g, e := b.String(), `{sync{10:1 11:3} mod{10:1 11:3} create{10:1}}`; g != e {
		t.Errorf("bad state B: %v != %v", g, e)
	}

	a.ValidateFile()
	b.ValidateFile()
}

func TestFigure2C(t *testing.T) {
	var b clock.Clock

	a := clock.Create(10, 1)
	if g, e := clock.Sync(a, &b), clock.Copy; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}
	b.ResolveTheirs(a)
	a.Update(10, 3)
	b.Update(11, 3)
	if g, e := clock.Sync(&b, a), clock.Conflict; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}

	if g, e := a.String(), `{sync{10:3} mod{10:3} create{10:1}}`; g != e {
		t.Errorf("bad state A: %v != %v", g, e)
	}
	if g, e := b.String(), `{sync{10:1 11:3} mod{10:1 11:3} create{10:1}}`; g != e {
		t.Errorf("bad state B: %v != %v", g, e)
	}

	a.ValidateFile()
	b.ValidateFile()
}

func TestFigure3B(t *testing.T) {
	var a clock.Clock
	var b clock.Clock
	var c clock.Clock

	b.Update(11, 1)

	if g, e := clock.Sync(&b, &a), clock.Copy; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}
	a.ResolveTheirs(&b)

	if g, e := clock.Sync(&b, &c), clock.Copy; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}
	c.ResolveTheirs(&b)

	a.Update(10, 3)
	b.Update(11, 3)

	if g, e := clock.Sync(&a, &b), clock.Conflict; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}
	// resolve conflict in favor of a
	b.ResolveTheirs(&a)

	if g, e := clock.Sync(&a, &b), clock.Nothing; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}

	if g, e := clock.Sync(&c, &b), clock.Nothing; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}

	a.Update(10, 6)
	if g, e := clock.Sync(&a, &b), clock.Copy; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}

	a.ValidateFile()
	b.ValidateFile()
	c.ValidateFile()
}

func TestFigure3C(t *testing.T) {
	var a clock.Clock
	var b clock.Clock
	var c clock.Clock

	b.Update(11, 1)

	if g, e := clock.Sync(&b, &a), clock.Copy; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}
	a.ResolveTheirs(&b)

	if g, e := clock.Sync(&b, &c), clock.Copy; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}
	c.ResolveTheirs(&b)

	a.Update(10, 3)
	b.Update(11, 3)

	if g, e := clock.Sync(&a, &b), clock.Conflict; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}
	// resolve conflict in favor of b
	b.ResolveOurs(&a)

	if g, e := clock.Sync(&a, &b), clock.Nothing; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}

	if g, e := clock.Sync(&c, &b), clock.Nothing; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}

	a.Update(10, 6)
	if g, e := clock.Sync(&a, &b), clock.Conflict; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}

	a.ValidateFile()
	b.ValidateFile()
	c.ValidateFile()
}

func TestFigure3D(t *testing.T) {
	var a clock.Clock
	var b clock.Clock
	var c clock.Clock

	b.Update(11, 1)

	if g, e := clock.Sync(&b, &a), clock.Copy; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}
	a.ResolveTheirs(&b)

	if g, e := clock.Sync(&b, &c), clock.Copy; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}
	c.ResolveTheirs(&b)

	a.Update(10, 3)
	b.Update(11, 3)

	if g, e := clock.Sync(&a, &b), clock.Conflict; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}
	// resolve conflict in favor of something new
	b.ResolveNew(&a)

	if g, e := clock.Sync(&a, &b), clock.Nothing; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}

	if g, e := clock.Sync(&c, &b), clock.Nothing; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}

	a.Update(10, 6)
	if g, e := clock.Sync(&a, &b), clock.Conflict; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}

	a.ValidateFile()
	b.ValidateFile()
	c.ValidateFile()
}

func TestFigure5A(t *testing.T) {
	a := clock.Create(10, 1)
	a.Update(10, 2) // delete

	b := &clock.Clock{}

	if g, e := clock.Sync(a, b), clock.Copy; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}
	b.ResolveTheirs(a)

	b = clock.Create(11, 4)
	if g, e := clock.SyncToMissing(b, a), clock.Copy; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}
	a.ResolveTheirs(b)

	a.ValidateFile()
	b.ValidateFile()
}

func TestFigure5C(t *testing.T) {
	a := clock.Create(10, 1)
	a.Update(10, 2) // delete

	b := clock.Create(11, 4)
	if g, e := clock.SyncToMissing(b, a), clock.Copy; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}
	a.ResolveTheirs(b)

	a.ValidateFile()
	b.ValidateFile()
}

func TestFigure5D(t *testing.T) {
	a := clock.Create(10, 1)
	b := clock.Create(11, 1)

	a.Update(10, 2) // delete
	b.Update(11, 2) // delete

	if g, e := clock.SyncToMissing(b, a), clock.Copy; g != e {
		t.Errorf("bad sync decision: %v is to %v -> %v != %v", a, b, g, e)
	}
	a.ResolveTheirs(b)

	a.ValidateFile()
	b.ValidateFile()
}
