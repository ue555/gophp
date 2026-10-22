<?php
function side($x) { echo 'eval:', $x, "\n"; return $x; }
foreach ([[1], [null, false], ['unknown' => 3]] as $args) {
    try { getmypid(...$args); }
    catch (Throwable $e) { echo get_class($e), ': ', $e->getMessage(), "\n"; }
}
try { getmypid(side(1), side(2)); }
catch (Throwable $e) { echo get_class($e), ': ', $e->getMessage(), "\n"; }
try { getmypid(unknown: side(3)); }
catch (Throwable $e) { echo get_class($e), ': ', $e->getMessage(), "\n"; }
echo getmypid(...[]) === getmypid() ? "unchanged\n" : "invalid\n";
