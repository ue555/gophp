<?php
function callableApply(callable $callback, int $n): int { return $callback($n); }
function callableReturn(): callable { return 'callableTwice'; }
function callableTwice(int $n): int { return $n * 2; }
function closureOnly(Closure $callback): int { return $callback(6); }
$fn = callableReturn();
$closure = callableTwice(...);
echo callableApply($fn, 4), '|', callableApply($closure, 5), '|', closureOnly($closure), '|';
echo (int)is_callable($closure), '|', gettype($closure), '|', (int)(Closure::fromCallable($closure) === $closure);
