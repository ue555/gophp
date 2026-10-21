<?php
readonly class ReadonlyPoint {
    public function __construct(public int $x, public int $y) {}
}
$p = new ReadonlyPoint(2, 3);
echo $p->x, ':', $p->y, '|';
try { $p->x = 4; } catch (Error $e) { echo 'write|'; }
try { $p->extra = 5; } catch (Error $e) { echo 'dynamic|'; }
echo $p->x, "\n";
